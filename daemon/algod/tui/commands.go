package tui

import (
	"github.com/algorand/go-algorand/node"
	tea "github.com/charmbracelet/bubbletea"
)

// commands are used to asynchronously fetch data for the UI.

type StatusMsg struct {
	Status node.StatusReport
	Error  error
}

func GetStatusCmd(fullNode *node.AlgorandFullNode) tea.Cmd {
	return func() tea.Msg {
		s, err := fullNode.Status()
		return StatusMsg{
			Status: s,
			Error:  err,
		}
	}
}
