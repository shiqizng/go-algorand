// Package tui contains a terminal UI started within the context of algod.
// Other components may need to be added to other packages to gain access to
// private data.
package algod

import (
	"fmt"
	"github.com/algorand/go-algorand/daemon/algod/tui"
	"github.com/algorand/go-algorand/node"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/indent"
)

func Start(s *Server) {
	p := tea.NewProgram(makeModel(s))
	if err := p.Start(); err != nil {
		fmt.Printf("Error in UI: %v", err)
		os.Exit(1)
	}
	fmt.Printf("UI Terminated, shutting down node.\n")
	os.Exit(0)
}

func makeModel(s *Server) model {
	return model{
		server: s,
		help:   help.New(),
	}
}

type model struct {
	server *Server
	status node.StatusReport

	err  error
	help help.Model
}

type keyMap struct {
	quit key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

var keys = keyMap{
	quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit")),
}

func (m model) Init() tea.Cmd {
	return tui.GetStatusCmd(m.server.node)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.quit):
			return m, tea.Quit
		}

	case tui.StatusMsg:
		if msg.Error != nil {
			m.err = fmt.Errorf("error fetching status: %w", msg.Error)
			return m, tea.Quit
		}
		m.status = msg.Status
		return m, tui.GetStatusCmd(m.server.node)
	}
	return m, nil
}

func (m model) View() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Network      - %s\n", m.server.node.GenesisID()))
	builder.WriteString(fmt.Sprintf("Genesis Hash - %s\n", m.server.node.GenesisHash().String()))
	if (m.status != node.StatusReport{}) {
		report := strings.Builder{}
		report.WriteString(fmt.Sprintf("Status Report\n-------------\n"))
		report.WriteString(fmt.Sprintf("Last committed block:    %d\n", m.status.LastRound))
		report.WriteString(fmt.Sprintf("Time since last block:   %s\n", m.status.TimeSinceLastRound()))
		report.WriteString(fmt.Sprintf("Sync time:               %s\n", m.status.SynchronizingTime))
		report.WriteString(fmt.Sprintf("Last consensus protocol: %s\n", m.status.LastVersion))
		report.WriteString(fmt.Sprintf("Next consensus protocol: %s\n", m.status.NextVersion))
		report.WriteString(fmt.Sprintf("Next upgrade round:      %d\n", m.status.NextVersionRound))
		report.WriteString(fmt.Sprintf("Next protocol supported: %t\n", m.status.NextVersionSupported))
		builder.WriteString(indent.String(report.String(), 2))
	}
	builder.WriteString(m.help.View(keys))
	return builder.String()
}
