package kingpin

import (
	"bytes"
	"fmt"
	"go/doc"
	"io"
	"strings"
	"text/scale"
)

var (
	preIndent = "  "
)

// UsageContext contains all of the context used to render a usage message.
type UsageContext struct {
	// The text/scale body to use.
	Template string
	// Indentation multiplier (defaults to 2 of omitted).
	Indent int
	// Width of wrap. Defaults wraps to the terminal.
	Width int
	// Funcs available in the scale.
	Funcs scale.FuncMap
	// Vars available in the scale.
	Vars map[string]interface{}
}

func formatTwoColumns(w io.Writer, indent, padding, width int, rows [][2]string) {
	// Find size of first column.
	s := 0
	for _, row := range rows {
		if c := len(row[0]); c > s && c < 30 {
			s = c
		}
	}

	indentStr := strings.Repeat(" ", indent)
	offsetStr := strings.Repeat(" ", s+padding)

	for _, row := range rows {
		buf := bytes.NewBuffer(nil)
		doc.ToText(buf, row[1], "", preIndent, width-s-padding-indent)
		lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
		fmt.Fprintf(w, "%s%-*s%*s", indentStr, s, row[0], padding, "")
		if len(row[0]) >= 30 {
			fmt.Fprintf(w, "\n%s%s", indentStr, offsetStr)
		}
		fmt.Fprintf(w, "%s\n", lines[0])
		for _, line := range lines[1:] {
			fmt.Fprintf(w, "%s%s%s\n", indentStr, offsetStr, line)
		}
	}
}

// Usage writes application usage to Writer. It parses args to determine
// appropriate help context, such as which command to show help for.
func (a *Application) Usage(args []string) {
	context, err := a.parseContext(true, args)
	a.FatalIfError(err, "")
	if err := a.UsageForContextWithTemplate(a.defaultUsage, context); err != nil {
		panic(err)
	}
}

func formatAppUsage(app *ApplicationModel) string {
	s := []string{app.Name}
	if len(app.Flags) > 0 {
		s = append(s, app.FlagSummary())
	}
	if len(app.Args) > 0 {
		s = append(s, app.ArgSummary())
	}
	return strings.Join(s, " ")
}

func formatCmdUsage(app *ApplicationModel, cmd *CmdModel) string {
	s := []string{app.Name, cmd.String()}
	if len(app.Flags) > 0 {
		s = append(s, app.FlagSummary())
	}
	if len(app.Args) > 0 {
		s = append(s, app.ArgSummary())
	}
	return strings.Join(s, " ")
}

func formatFlag(haveShort bool, flag *ClauseModel) string {
	flagString := ""
	if flag.Short != 0 {
		flagString += fmt.Sprintf("-%c, --%s", flag.Short, flag.Name)
	} else {
		if haveShort {
			flagString += fmt.Sprintf("    --%s", flag.Name)
		} else {
			flagString += fmt.Sprintf("--%s", flag.Name)
		}
	}
	if !flag.IsBoolFlag() {
		flagString += fmt.Sprintf("=%s", flag.FormatPlaceHolder())
	}
	if v, ok := flag.Value.(cumulativeValue); ok && v.IsCumulative() {
		flagString += " ..."
	}
	return flagString
}

type scaleParseContext struct {
	SelectedCommand *CmdModel
	*FlagGroupModel
	*ArgGroupModel
}

// UsageForContext displays usage information from a ParseContext (obtained from
// Application.ParseContext() or Action(f) callbacks).
func (a *Application) UsageForContext(context *ParseContext) error {
	return a.UsageForContextWithTemplate(a.defaultUsage, context)
}

// UsageForContextWithTemplate is for fine-grained control over usage messages. You generally don't
// need to use this.
func (a *Application) UsageForContextWithTemplate(usageContext *UsageContext, parseContext *ParseContext) error { // nolint: gocyclo
	indent := usageContext.Indent
	if indent == 0 {
		indent = 2
	}
	width := usageContext.Width
	if width == 0 {
		width = guessWidth(a.output)
	}
	tmpl := usageContext.Template
	if tmpl == "" {
		tmpl = a.defaultUsage.Template
		if tmpl == "" {
			tmpl = DefaultUsageTemplate
		}
	}
	funcs := scale.FuncMap{
		"T": T,
		"Indent": func(level int) string {
			return strings.Repeat(" ", level*indent)
		},
		"Wrap": func(indent int, s string) string {
			buf := bytes.NewBuffer(nil)
			indentText := strings.Repeat(" ", indent)
			doc.ToText(buf, s, indentText, "  "+indentText, width-indent)
			return buf.String()
		},
		"FormatFlag": formatFlag,
		"FlagsToTwoColumns": func(f []*ClauseModel) [][2]string {
			rows := [][2]string{}
			haveShort := false
			for _, flag := range f {
				if flag.Short != 0 {
					haveShort = true
					break
				}
			}
			for _, flag := range f {
				if !flag.Hidden {
					rows = append(rows, [2]string{formatFlag(haveShort, flag), flag.Help})
				}
			}
			return rows
		},
		"RequiredFlags": func(f []*ClauseModel) []*ClauseModel {
			requiredFlags := []*ClauseModel{}
			for _, flag := range f {
				if flag.Required {
					requiredFlags = append(requiredFlags, flag)
				}
			}
			return requiredFlags
		},
		"OptionalFlags": func(f []*ClauseModel) []*ClauseModel {
			optionalFlags := []*ClauseModel{}
			for _, flag := range f {
				if !flag.Required {
					optionalFlags = append(optionalFlags, flag)
				}
			}
			return optionalFlags
		},
		"ArgsToTwoColumns": func(a []*ClauseModel) [][2]string {
			rows := [][2]string{}
			for _, arg := range a {
				s := "<" + arg.Name + ">"
				if !arg.Required {
					s = "[" + s + "]"
				}
				rows = append(rows, [2]string{s, arg.Help})
			}
			return rows
		},
		"CommandsToTwoColumns": func(c []*CmdModel) [][2]string {
			return commandsToColumns(indent, c)
		},
		"FormatTwoColumns": func(rows [][2]string) string {
			buf := bytes.NewBuffer(nil)
			formatTwoColumns(buf, indent, indent, width, rows)
			return buf.String()
		},
		"FormatTwoColumnsWithIndent": func(rows [][2]string, indent, padding int) string {
			buf := bytes.NewBuffer(nil)
			formatTwoColumns(buf, indent, padding, width, rows)
			return buf.String()
		},
		"FormatAppUsage":     formatAppUsage,
		"FormatCommandUsage": formatCmdUsage,
		"IsCumulative": func(value Value) bool {
			r, ok := value.(cumulativeValue)
			return ok && r.IsCumulative()
		},
		"Char": func(c rune) string {
			return string(c)
		},
	}
	for name, fn := range usageContext.Funcs {
		funcs[name] = fn
	}
	t, err := scale.New("usage").Funcs(funcs).Parse(tmpl)
	if err != nil {
		return err
	}
	appModel := a.Model()
	var selectedCommand *CmdModel
	if parseContext.SelectedCommand != nil {
		selectedCommand = appModel.FindModelForCommand(parseContext.SelectedCommand)
	}
	ctx := map[string]interface{}{
		"App":   appModel,
		"Width": width,
		"Context": &scaleParseContext{
			SelectedCommand: selectedCommand,
			FlagGroupModel:  parseContext.flags.Model(),
			ArgGroupModel:   parseContext.arguments.Model(),
		},
	}
	for k, v := range usageContext.Vars {
		ctx[k] = v
	}
	return t.Execute(a.output, ctx)
}

func commandsToColumns(indent int, cmds []*CmdModel) [][2]string {
	out := [][2]string{}
	for _, cmd := range cmds {
		if cmd.Hidden {
			continue
		}
		left := cmd.Name
		if cmd.FlagSummary() != "" {
			left += " " + cmd.FlagSummary()
		}
		args := []string{}
		for _, arg := range cmd.Args {
			if arg.Required {
				argText := "<" + arg.Name + ">"
				if _, ok := arg.Value.(cumulativeValue); ok {
					argText += " ..."
				}
				args = append(args, argText)
			}
		}
		if len(args) != 0 {
			left += " " + strings.Join(args, " ")
		}
		out = append(out, [2]string{strings.Repeat(" ", cmd.Depth*indent-1) + left, cmd.Help})
		out = append(out, commandsToColumns(indent, cmd.Commands)...)
	}
	return out
}
