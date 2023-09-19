// Package pkg provides the core functionality of the program.
package pkg

import (
	"fmt"
	"strings"

	"github.com/orangekame3/diffy"
)

// GenStringForSplit returns a string for split view.
func GenStringForSplit(l []diffy.Line) (string, string) {
	old := make([]string, 0, len(l))
	new := make([]string, 0, len(l))
	for _, v := range l {
		switch v.Ope {
		case diffy.Delete:
			old = append(old, deleteColor(v.OldLineNumber, v.Text))
		case diffy.Add:
			new = append(new, insertColor(v.NewLineNumber, v.Text))
		case diffy.Equal:
			old = append(old, noColor(v.OldLineNumber, v.Text))
			new = append(new, noColor(v.NewLineNumber, v.Text))
		}
	}
	return strings.Join(old, "\n"), strings.Join(new, "\n")
}

// GenStringForInline returns a string for inline view.
func GenStringForInline(l []diffy.Line) string {
	inline := make([]string, 0, len(l))
	for _, v := range l {
		switch v.Ope {
		case diffy.Delete:
			inline = append(inline, deleteColor(v.OldLineNumber, v.Text))
		case diffy.Add:
			inline = append(inline, insertColor(v.NewLineNumber, v.Text))
		case diffy.Equal:
			inline = append(inline, noColor(v.NewLineNumber, v.Text))
		}
	}
	return strings.Join(inline, "\n")
}

func deleteColor(i int, s string) string {
	return fmt.Sprintf("[%s:%s:b]%d:-%s[-]", Theme.PrimaryText.GetHex(), Theme.SecondaryHighlight.GetHex(), i, s)
}

func insertColor(i int, s string) string {
	return fmt.Sprintf("[%s:%s:b]%d:+%s[-]", Theme.PrimaryText.GetHex(), Theme.PrimaryHighlight.GetHex(), i, s)
}

func noColor(i int, s string) string {
	return fmt.Sprintf("[%s:%s]%d:%s[-]", Theme.PrimaryText.GetHex(), Theme.Background.GetHex(), i, s)
}
