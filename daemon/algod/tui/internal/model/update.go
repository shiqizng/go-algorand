package model

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/algorand/go-algorand/daemon/algod"
	"github.com/algorand/go-algorand/daemon/algod/tui/internal/constants"
)

func networkFromID(genesisID string) string {
	return strings.Split(genesisID, "-")[0]
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case algod.NetworkMsg:
		m.network = msg

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, constants.Keys.Catchup):
			return m, algod.StartFastCatchup(m.Server, networkFromID(m.Status.Network.GenesisID))
		case key.Matches(msg, constants.Keys.AbortCatchup):
			return m, algod.StopFastCatchup(m.Server, networkFromID(m.Status.Network.GenesisID))
		case key.Matches(msg, constants.Keys.Section):
			m.active += 1
			m.active %= 4
			m.Tabs.SetActiveIndex(int(m.active))
			return m, nil
		}
		switch m.active {
		case explorerTab:
			var explorerCommand tea.Cmd
			m.BlockExplorer, explorerCommand = m.BlockExplorer.Update(msg)
			return m, explorerCommand
		case accountTab:
		case configTab:
		case helpTab:
		}

	case tea.WindowSizeMsg:
		m.lastResize = msg
	}

	m.Status, cmd = m.Status.Update(msg)
	cmds = append(cmds, cmd)

	m.Accounts, cmd = m.Accounts.Update(msg)
	cmds = append(cmds, cmd)

	m.BlockExplorer, cmd = m.BlockExplorer.Update(msg)
	cmds = append(cmds, cmd)

	m.Configs, cmd = m.Configs.Update(msg)
	cmds = append(cmds, cmd)

	m.Footer, cmd = m.Footer.Update(msg)
	cmds = append(cmds, cmd)

	m.Tabs, cmd = m.Tabs.Update(msg)
	cmds = append(cmds, cmd)

	m.About, cmd = m.About.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
