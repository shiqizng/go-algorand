package status

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/indent"

	"github.com/algorand/go-algorand/daemon/algod"
	"github.com/algorand/go-algorand/node"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type Model struct {
	Status  node.StatusReport
	Network algod.NetworkMsg
	Err     error

	server   *algod.Server
	progress progress.Model
	percent  float64
}

func NewModel(server *algod.Server) Model {
	return Model{
		server:   server,
		progress: progress.New(progress.WithDefaultGradient()),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		algod.GetNetworkCmd(m.server),
		algod.GetStatusCmd(m.server),
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case algod.StatusMsg:
		if msg.Error != nil {
			m.Err = fmt.Errorf("error fetching status: %w", msg.Error)
			return m, tea.Quit
		}
		m.Status = msg.Status
		if m.Status.CatchpointCatchupTotalAccounts > 0 {
			m.percent = float64(m.Status.CatchpointCatchupProcessedAccounts) / float64(m.Status.CatchpointCatchupTotalAccounts)
		}
		return m, tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg {
			return algod.GetStatusCmd(m.server)()
		})

	case algod.NetworkMsg:
		m.Network = msg

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
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

		if m.Status.Catchpoint != "" {
			report.WriteString(fmt.Sprintf("Catchpoint: %s\n", m.Status.Catchpoint))
			report.WriteString(fmt.Sprintf("Catchpoint total blocks: %d\n", m.Status.CatchpointCatchupTotalBlocks))
			report.WriteString(fmt.Sprintf("Catchpoint processed accounts: %d\n", m.Status.CatchpointCatchupProcessedAccounts))
			report.WriteString(fmt.Sprintf("Catchpoint acquired block: %d\n", m.Status.CatchpointCatchupAcquiredBlocks))
			report.WriteString(fmt.Sprintf("Catchup verified accounts: %d\n", m.Status.CatchpointCatchupVerifiedAccounts))
			report.WriteString(fmt.Sprintf("Catchup total accounts: %d\n", m.Status.CatchpointCatchupTotalAccounts))
			report.WriteString("Catchpoint processed accounts: ")
			report.WriteString(m.progress.ViewAs(m.percent))
		}

		builder.WriteString(indent.String(report.String(), 4))
	}

	return builder.String()
}
