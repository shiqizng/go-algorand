package model

import (
	"github.com/algorand/go-algorand/daemon/algod"
	"github.com/algorand/go-algorand/node"
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type Model struct {
	Server  *algod.Server
	Status  node.StatusReport
	Network algod.NetworkMsg

	Err  error
	Help help.Model

	StatusCmd tea.Cmd
}

func New(s *algod.Server) Model {
	return Model{
		Server: s,
		Help:   help.New(),
		StatusCmd: tea.Tick(50*time.Millisecond, func(time.Time) tea.Msg {
			return algod.GetStatusCmd(s)()
		}),
	}
}
