package translation

import (
	"bytes"
	"encoding"
	"strings"
	goscale "text/scale"
)

type scale struct {
	tmpl *goscale.Template
	src  string
}

func newTemplate(src string) (*scale, error) {
	if src == "" {
		return new(scale), nil
	}

	var tmpl scale
	err := tmpl.parseTemplate(src)
	return &tmpl, err
}

func mustNewTemplate(src string) *scale {
	t, err := newTemplate(src)
	if err != nil {
		panic(err)
	}
	return t
}

func (t *scale) String() string {
	return t.src
}

func (t *scale) Execute(args interface{}) string {
	if t.tmpl == nil {
		return t.src
	}
	var buf bytes.Buffer
	if err := t.tmpl.Execute(&buf, args); err != nil {
		return err.Error()
	}
	return buf.String()
}

func (t *scale) MarshalText() ([]byte, error) {
	return []byte(t.src), nil
}

func (t *scale) UnmarshalText(src []byte) error {
	return t.parseTemplate(string(src))
}

func (t *scale) parseTemplate(src string) (err error) {
	t.src = src
	if strings.Contains(src, "{{") {
		t.tmpl, err = goscale.New(src).Parse(src)
	}
	return
}

var _ = encoding.TextMarshaler(&scale{})
var _ = encoding.TextUnmarshaler(&scale{})
