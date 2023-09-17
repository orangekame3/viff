// Package pkg provides the core functionality of the program.
package pkg

import (
	"fmt"
	"strings"

	"github.com/shibukawa/cdiff"
)

// GenStringForSplit returns a string for split view.
func GenStringForSplit(r cdiff.Result) (string, string) {
	old := make([]string, 0, len(r.Lines))
	new := make([]string, 0, len(r.Lines))
	for _, v := range r.Lines {
		switch v.Ope {
		case cdiff.Delete:
			old = append(old, deleteColor(v.OldLineNumber, v.String()))
		case cdiff.Insert:
			new = append(new, insertColor(v.NewLineNumber, v.String()))
		case cdiff.Keep:
			old = append(old, noColor(v.OldLineNumber, v.String()))
			new = append(new, noColor(v.NewLineNumber, v.String()))
		}
	}
	return strings.Join(old, "\n"), strings.Join(new, "\n")
}

// GenStringForInline returns a string for inline view.
func GenStringForInline(r cdiff.Result) string {
	inline := make([]string, 0, len(r.Lines))
	for _, v := range r.Lines {
		switch v.Ope {
		case cdiff.Delete:
			inline = append(inline, deleteColor(v.OldLineNumber, v.String()))
		case cdiff.Insert:
			inline = append(inline, insertColor(v.NewLineNumber, v.String()))
		case cdiff.Keep:
			inline = append(inline, noColor(v.NewLineNumber, v.String()))
		}
	}
	return strings.Join(inline, "\n")
}

func deleteColor(i int, s string) string {
	return fmt.Sprintf("[%s:%s:b]%d:-%s[-]", Flavour.Base().Hex, Flavour.Red().Hex, i, s)
}

func insertColor(i int, s string) string {
	return fmt.Sprintf("[%s:%s:b]%d:+%s[-]", Flavour.Base().Hex, Flavour.Teal().Hex, i, s)
}

func noColor(i int, s string) string {
	return fmt.Sprintf("[%s:%s]%d:%s[-]", Flavour.Text().Hex, Flavour.Surface1().Hex, i, s)
}
