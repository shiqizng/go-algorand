// Copyright (C) 2019-2022 Algorand, Inc.
// This file is part of go-algorand
//
// go-algorand is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// go-algorand is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with go-algorand.  If not, see <https://www.gnu.org/licenses/>.

// Package tui contains a terminal UI started within the context of algod.
// Other components may need to be added to other packages to gain access to
// private data.
package tui

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/algorand/go-algorand/daemon/algod"
	"github.com/algorand/go-algorand/node"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/indent"
)

type model struct {
	server  *algod.Server
	status  node.StatusReport
	network algod.NetworkMsg

	err  error
	help help.Model

	statusCmd tea.Cmd
}

func makeModel(s *algod.Server) model {
	return model{
		server: s,
		help:   help.New(),
		statusCmd: tea.Tick(50*time.Millisecond, func(time.Time) tea.Msg {
			return algod.GetStatusCmd(s)()
		}),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		algod.GetNetworkCmd(m.server),
		algod.GetStatusCmd(m.server))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.quit):
			return m, tea.Quit
		}

	case algod.NetworkMsg:
		m.network = msg

	case algod.StatusMsg:
		if msg.Error != nil {
			m.err = fmt.Errorf("error fetching status: %w", msg.Error)
			return m, tea.Quit
		}
		m.status = msg.Status
		return m, m.statusCmd
	}
	return m, nil
}

func formatVersion(v string) string {
	i := strings.LastIndex(v, "/")
	if i != 0 {
		i++
	}
	return v[i:]
}

func formatNextVersion(last, next string, round uint64) string {
	if last == next {
		return "N/A"
	}
	return strconv.FormatUint(round, 10)
}

func (m model) View() string {
	builder := strings.Builder{}

	// general information
	builder.WriteString(fmt.Sprintf("Network      - %s\n", m.network.GenesisID))
	builder.WriteString(fmt.Sprintf("Genesis Hash - %s\n", m.network.GenesisHash.String()))

	// status
	if (m.status != node.StatusReport{}) {
		nextVersion := formatNextVersion(
			string(m.status.LastVersion),
			string(m.status.NextVersion),
			uint64(m.status.NextVersionRound))
		report := strings.Builder{}
		report.WriteString(fmt.Sprintf("Status Report\n-------------                                     \n"))
		report.WriteString(fmt.Sprintf("Last committed block:    %d\n", m.status.LastRound))
		report.WriteString(fmt.Sprintf("Time since last block:   %s\n", m.status.TimeSinceLastRound()))
		report.WriteString(fmt.Sprintf("Sync time:               %s\n", m.status.SynchronizingTime))
		report.WriteString(fmt.Sprintf("Last consensus protocol: %s\n", formatVersion(string(m.status.LastVersion))))
		report.WriteString(fmt.Sprintf("Next consensus protocol: %s\n", formatVersion(string(m.status.NextVersion))))
		report.WriteString(fmt.Sprintf("Next upgrade round:      %s\n", nextVersion))
		report.WriteString(fmt.Sprintf("Next protocol supported: %t\n", m.status.NextVersionSupported))
		builder.WriteString(indent.String(report.String(), 4))
	}

	// help
	builder.WriteString(m.help.View(keys))
	return builder.String()
}

func Start(s *algod.Server) {
	p := tea.NewProgram(makeModel(s), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Error in UI: %v", err)
		os.Exit(1)
	}
	fmt.Printf("\nUI Terminated, shutting down node.\n")
	os.Exit(0)
}
