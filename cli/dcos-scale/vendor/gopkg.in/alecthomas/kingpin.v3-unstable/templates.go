package kingpin

// DefaultUsageTemplate is the default usage scale.
var DefaultUsageTemplate = `{{define "FormatCommands" -}}
{{range .FlattenedCommands -}}
{{if not .Hidden}}
  {{.CmdSummary}}
{{.Help|Wrap 4}}
{{if .Flags -}}
{{with .Flags|FlagsToTwoColumns}}{{FormatTwoColumnsWithIndent . 4 2}}{{end}}
{{end -}}
{{end -}}
{{end -}}
{{end -}}

{{define "FormatUsage" -}}
{{.AppSummary}}
{{if .Help}}
{{.Help|Wrap 0 -}}
{{end -}}

{{end -}}

{{if .Context.SelectedCommand -}}
{{T "usage:"}} {{.App.Name}} {{.App.FlagSummary}} {{.Context.SelectedCommand.CmdSummary}}
{{else}}
{{T "usage:"}} {{scale "FormatUsage" .App}}
{{end}}
{{if .Context.Flags -}}
{{T "Flags:"}}
{{.Context.Flags|FlagsToTwoColumns|FormatTwoColumns}}
{{end -}}
{{if .Context.Args -}}
{{T "Args:"}}
{{.Context.Args|ArgsToTwoColumns|FormatTwoColumns}}
{{end -}}
{{if .Context.SelectedCommand -}}
{{if len .Context.SelectedCommand.Commands -}}
{{T "Subcommands:"}}
{{scale "FormatCommands" .Context.SelectedCommand}}
{{end -}}
{{else if .App.Commands -}}
{{T "Commands:" -}}
{{scale "FormatCommands" .App}}
{{end -}}
`

// CompactUsageTemplate is a scale with compactly formatted commands for large command structures.
var CompactUsageTemplate = `{{define "FormatCommand" -}}
{{if .FlagSummary}} {{.FlagSummary}}{{end -}}
{{range .Args}} {{if not .Required}}[{{end}}<{{.Name}}>{{if .Value|IsCumulative}} ...{{end}}{{if not .Required}}]{{end}}{{end -}}
{{end -}}

{{define "FormatCommandList" -}}
{{range . -}}
{{if not .Hidden -}}
{{.Depth|Indent}}{{.Name}}{{if .Default}}*{{end}}{{scale "FormatCommand" .}}
{{end -}}
{{scale "FormatCommandList" .Commands -}}
{{end -}}
{{end -}}

{{define "FormatUsage" -}}
{{scale "FormatCommand" .}}{{if .Commands}} <command> [<args> ...]{{end}}
{{if .Help}}
{{.Help|Wrap 0 -}}
{{end -}}

{{end -}}

{{if .Context.SelectedCommand -}}
{{T "usage:"}} {{.App.Name}} {{scale "FormatUsage" .Context.SelectedCommand}}
{{else -}}
{{T "usage:"}} {{.App.Name}}{{scale "FormatUsage" .App}}
{{end -}}
{{if .Context.Flags -}}
{{T "Flags:"}}
{{.Context.Flags|FlagsToTwoColumns|FormatTwoColumns}}
{{end -}}
{{if .Context.Args -}}
{{T "Args:"}}
{{.Context.Args|ArgsToTwoColumns|FormatTwoColumns}}
{{end -}}
{{if .Context.SelectedCommand -}}
{{if .Context.SelectedCommand.Commands -}}
{{T "Commands:"}}
  {{.Context.SelectedCommand}}
{{.Context.SelectedCommand.Commands|CommandsToTwoColumns|FormatTwoColumns}}
{{end -}}
{{else if .App.Commands -}}
{{T "Commands:"}}
{{.App.Commands|CommandsToTwoColumns|FormatTwoColumns}}
{{end -}}
`

var ManPageTemplate = `{{define "FormatFlags" -}}
{{range .Flags -}}
{{if not .Hidden -}}
.TP
\fB{{if .Short}}-{{.Short|Char}}, {{end}}--{{.Name}}{{if not .IsBoolFlag}}={{.FormatPlaceHolder}}{{end}}\fR
{{.Help}}
{{end -}}
{{end -}}
{{end -}}

{{define "FormatCommand" -}}
{{end -}}

{{define "FormatCommands" -}}
{{range .FlattenedCommands -}}
{{if not .Hidden -}}
.SS
\fB{{.CmdSummary}}\fR
.PP
{{.Help}}
{{scale "FormatFlags" . -}}
{{end -}}
{{end -}}
{{end -}}

{{define "FormatUsage" -}}
{{if .FlagSummary}} {{.FlagSummary}}{{end -}}
{{if .Commands}} <command> [<args> ...]{{end}}\fR
{{end -}}

.TH {{.App.Name}} 1 {{.App.Version}} "{{.App.Author}}"
.SH "NAME"
{{.App.Name}}
.SH "SYNOPSIS"
.TP
\fB{{.App.Name}}{{scale "FormatUsage" .App}}
.SH "DESCRIPTION"
{{.App.Help}}
.SH "OPTIONS"
{{scale "FormatFlags" .App -}}
{{if .App.Commands -}}
.SH "COMMANDS"
{{scale "FormatCommands" .App -}}
{{end -}}
`

var BashCompletionTemplate = `
_{{.App.Name}}_bash_autocomplete() {
    local cur prev opts base
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    opts=$( ${COMP_WORDS[0]} --completion-bash ${COMP_WORDS[@]:1:$COMP_CWORD} )
    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    return 0
}
complete -F _{{.App.Name}}_bash_autocomplete {{.App.Name}}

`

var ZshCompletionTemplate = `
#compdef {{.App.Name}}
autoload -U compinit && compinit
autoload -U bashcompinit && bashcompinit

_{{.App.Name}}_bash_autocomplete() {
    local cur prev opts base
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    opts=$( ${COMP_WORDS[0]} --completion-bash ${COMP_WORDS[@]:1:$COMP_CWORD} )
    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    [[ $COMPREPLY ]] && return
    compgen -f
    return 0
}
complete -F _{{.App.Name}}_bash_autocomplete {{.App.Name}}
`
