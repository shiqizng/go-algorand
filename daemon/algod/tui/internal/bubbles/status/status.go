package status

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/indent"

	"github.com/algorand/go-algorand/daemon/algod"
	"github.com/algorand/go-algorand/node"
)

type Model struct {
	Status  node.StatusReport
	Network algod.NetworkMsg
	Err     error

	server *algod.Server
}

func NewModel(server *algod.Server) Model {
	return Model{
		server: server,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		algod.GetNetworkCmd(m.server),
		algod.GetStatusCmd(m.server))
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case algod.StatusMsg:
		if msg.Error != nil {
			m.Err = fmt.Errorf("error fetching status: %w", msg.Error)
			return m, tea.Quit
		}
		m.Status = msg.Status
		return m, tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg {
			return algod.GetStatusCmd(m.server)()
		})

	case algod.NetworkMsg:
		m.Network = msg
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

func (m Model) View() string {
	builder := strings.Builder{}

	// general information
	builder.WriteString(fmt.Sprintf("Network      - %s\n", m.Network.GenesisID))
	builder.WriteString(fmt.Sprintf("Genesis Hash - %s\n", m.Network.GenesisHash.String()))

	// status
	if (m.Status != node.StatusReport{}) {
		nextVersion := formatNextVersion(
			string(m.Status.LastVersion),
			string(m.Status.NextVersion),
			uint64(m.Status.NextVersionRound))
		report := strings.Builder{}
		report.WriteString(fmt.Sprintf("Status Report\n-------------                                     \n"))
		report.WriteString(fmt.Sprintf("Last committed block:    %d\n", m.Status.LastRound))
		report.WriteString(fmt.Sprintf("Time since last block:   %s\n", m.Status.TimeSinceLastRound()))
		report.WriteString(fmt.Sprintf("Sync time:               %s\n", m.Status.SynchronizingTime))
		report.WriteString(fmt.Sprintf("Last consensus protocol: %s\n", formatVersion(string(m.Status.LastVersion))))
		report.WriteString(fmt.Sprintf("Next consensus protocol: %s\n", formatVersion(string(m.Status.NextVersion))))
		report.WriteString(fmt.Sprintf("Next upgrade round:      %s\n", nextVersion))
		report.WriteString(fmt.Sprintf("Next protocol supported: %t\n", m.Status.NextVersionSupported))
		report.WriteString(fmt.Sprintf("catch point total blocks: %d\n", m.Status.CatchpointCatchupTotalBlocks))
		report.WriteString(fmt.Sprintf("catchpoint: %s\n", m.Status.Catchpoint))
		report.WriteString(fmt.Sprintf("catchpoint processed accounts: %d\n", m.Status.CatchpointCatchupProcessedAccounts))
		report.WriteString(fmt.Sprintf("catchpoint acquired block: %d\n", m.Status.CatchpointCatchupAcquiredBlocks))
		report.WriteString(fmt.Sprintf("catchup verified accounts: %d\n", m.Status.CatchpointCatchupVerifiedAccounts))
		report.WriteString(fmt.Sprintf("catchup total accounts: %d\n", m.Status.CatchpointCatchupTotalAccounts))
		builder.WriteString(indent.String(report.String(), 4))
	}

	return builder.String()
}
