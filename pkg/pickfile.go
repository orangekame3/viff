// Package pkg provides the core functionality of the program.
package pkg

import (
	"os"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

// Picker is a model for the file picker.
type Picker struct {
	Filepicker   filepicker.Model
	SelectedFile string
	Quitting     bool
	Err          error
}

// Init initializes the file picker.
func (p *Picker) Init() tea.Cmd {
	return p.Filepicker.Init()
}

// NewPicker returns a new file picker.
func NewPicker() Picker {
	fp := filepicker.New()
	fp.CurrentDirectory, _ = os.Getwd()
	return Picker{
		Filepicker: fp,
	}
}

// Update updates the file picker.
func (p *Picker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			p.Quitting = true
			return p, tea.Quit
		}
	}

	var cmd tea.Cmd
	p.Filepicker, cmd = p.Filepicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := p.Filepicker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		p.SelectedFile = path
		p.Quitting = true
		return p, tea.Quit
	}

	return p, cmd
}

// View returns the file picker view.
func (p *Picker) View() string {
	if p.Quitting {
		return ""
	}
	var s string
	if p.Err != nil {
		s = p.Filepicker.Styles.DisabledFile.Render(p.Err.Error())
	} else if p.SelectedFile == "" {
		s = "Pick a file:"
	} else {
		s = "Selected file: " + p.Filepicker.Styles.Selected.Render(p.SelectedFile)
	}
	return s + "\n\n" + p.Filepicker.View() + "\n"
}
