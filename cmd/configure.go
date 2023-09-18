/*
Copyright © 2023 Takafumi Miyanaga <miya.org.0309@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	catppuccin "github.com/catppuccin/go"

	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mritd/bubbles/common"
	"github.com/mritd/bubbles/selector"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure the system settings",
	Long:  `This command allows users to configure system settings by choosing from a list of options.`,
	Run: func(cmd *cobra.Command, args []string) {
		runSelector()
	},
}

func runSelector() {
	for _, flavour := range []catppuccin.Flavour{
		catppuccin.Macchiato,
		catppuccin.Latte,
	} {

		fmt.Println(lipgloss.NewStyle().Bold(true).Render(flavour.Name() + ":"))
		format("background", flavour.Surface1(), flavour.Text())
		format("delete", flavour.Red(), flavour.Base())
		format("insert", flavour.Teal(), flavour.Base())
		format("text1", flavour.Text(), flavour.Base())
		format("text2", flavour.Base(), flavour.Text())
		fmt.Println()
	}
	m := &model{
		sl: selector.Model{
			Data: []interface{}{
				"Macchiato",
				"Latte",
			},
			PerPage:    4,
			HeaderFunc: selector.DefaultHeaderFuncWithAppend("Select Theme:"),
			SelectedFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
				t := obj.(string)
				return common.FontColor(fmt.Sprintf("[%d] %s", gdIndex+1, t), selector.ColorSelected)
			},
			UnSelectedFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
				t := obj.(string)
				return common.FontColor(fmt.Sprintf(" %d. %s", gdIndex+1, t), selector.ColorUnSelected)
			},
			FooterFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
				t := m.Selected().(string)
				return common.FontColor(fmt.Sprintf("Theme: %s", t), selector.ColorFooter)
			},
			FinishedFunc: func(s interface{}) string {
				return common.FontColor("Current selected: ", selector.ColorFinished) + fmt.Sprintf("%v", s) + "\n"
			},
		},
	}

	p := tea.NewProgram(m)
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
	if !m.sl.Canceled() {
		selectedTheme := m.sl.Selected().(string)
		config := Config{
			Theme: selectedTheme,
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("could not find user home directory: %v", err)
		}

		viffDir := filepath.Join(homeDir, ".viff")
		if err := os.MkdirAll(viffDir, 0755); err != nil {
			log.Fatalf("could not create .viff directory: %v", err)
		}

		configFile := filepath.Join(viffDir, "config.toml")
		file, err := os.OpenFile(configFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatalf("could not open configure.toml: %v", err)
		}
		defer file.Close()

		data, err := toml.Marshal(config)
		if err != nil {
			log.Fatalf("could not marshal theme data: %v", err)
		}

		if _, err := file.Write(data); err != nil {
			log.Fatalf("could not write to config.toml: %v", err)
		}
	} else {
		log.Println("user canceled...")
	}

}

// Config is a struct for configure.toml
type Config struct {
	Theme string `toml:"theme"`
}

// model is the Bubble Tea model
type model struct {
	sl selector.Model
}

// init initializes the model
func (m model) Init() tea.Cmd {
	return nil
}

// update updates the model
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	_, cmd := m.sl.Update(msg)
	switch msg {
	case common.DONE:
		return m, tea.Quit
	}
	return m, cmd
}

func (m model) View() string {
	return m.sl.View()
}

func init() {
	rootCmd.AddCommand(configureCmd)
}

func format(s string, c, txt catppuccin.Color) {
	fmt.Print(lipgloss.NewStyle().
		Background(lipgloss.Color(c.Hex)).
		Foreground(lipgloss.Color(txt.Hex)).
		Align(lipgloss.Center).
		Width(22).
		Render(s))
}
