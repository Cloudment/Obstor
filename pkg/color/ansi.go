package color

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// NoColor disables color output globally when set to true.
var NoColor bool

// Output is the writer where colored output goes.
var Output io.Writer = os.Stdout

// Attribute defines a single ANSI color attribute.
type Attribute int

// ANSI color attributes.
// Attributes that would conflict with exported color helper variables
// in color.go are prefixed with "Attr" (e.g. AttrBold, AttrFgRed).
const (
	Reset        Attribute = 0
	AttrBold     Attribute = 1
	Faint        Attribute = 2
	Italic       Attribute = 3
	Underline    Attribute = 4
	CrossedOut   Attribute = 9
	FgBlack      Attribute = 30
	AttrFgRed    Attribute = 31
	FgGreen      Attribute = 32
	FgYellow     Attribute = 33
	FgBlue       Attribute = 34
	FgMagenta    Attribute = 35
	FgCyan       Attribute = 36
	AttrFgWhite  Attribute = 37
	AttrBgRed    Attribute = 41
	BgGreen      Attribute = 42
	AttrBgYellow Attribute = 43
	FgHiRed      Attribute = 91
	FgHiGreen    Attribute = 92
	FgHiYellow   Attribute = 93
	FgHiBlue     Attribute = 94
	FgHiCyan     Attribute = 96
	FgHiWhite    Attribute = 97
	BgHiYellow   Attribute = 103
)

// Color wraps ANSI escape code attributes.
type Color struct {
	attrs []Attribute
}

// New creates a new Color with the given attributes.
func New(attrs ...Attribute) *Color {
	return &Color{attrs: attrs}
}

func (c *Color) escape() string {
	if NoColor || len(c.attrs) == 0 {
		return ""
	}
	parts := make([]string, len(c.attrs))
	for i, a := range c.attrs {
		parts[i] = fmt.Sprintf("%d", int(a))
	}
	return fmt.Sprintf("\033[%sm", strings.Join(parts, ";"))
}

func (c *Color) reset() string {
	if NoColor || len(c.attrs) == 0 {
		return ""
	}
	return "\033[0m"
}

// Sprint formats using the default formats for its operands.
func (c *Color) Sprint(a ...interface{}) string {
	return c.escape() + fmt.Sprint(a...) + c.reset()
}

// Sprintf formats according to a format specifier.
func (c *Color) Sprintf(format string, a ...interface{}) string {
	return c.escape() + fmt.Sprintf(format, a...) + c.reset()
}

// SprintFunc returns a function that formats using the default formats.
func (c *Color) SprintFunc() func(a ...interface{}) string {
	return c.Sprint
}

// SprintfFunc returns a function that formats according to a format specifier.
func (c *Color) SprintfFunc() func(format string, a ...interface{}) string {
	return c.Sprintf
}

// Print formats using the default formats and writes to Output.
func (c *Color) Print(a ...interface{}) (n int, err error) {
	return fmt.Fprint(Output, c.Sprint(a...))
}

// Printf formats according to a format specifier and writes to Output.
func (c *Color) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprint(Output, c.Sprintf(format, a...))
}

// Println formats using the default formats and writes to Output followed by a newline.
func (c *Color) Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(Output, c.Sprint(a...))
}
