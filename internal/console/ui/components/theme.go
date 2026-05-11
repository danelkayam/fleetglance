package components

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	ColorDivider       = "#1A2426"
	ColorTextPrimary   = "#E6E6E6"
	ColorTextSecondary = "#B8B8B8"
	ColorTextMuted     = "#8A8F91"
	ColorTextDim       = "#5B6366"
	ColorOnline        = "#29E64A"
	ColorFailed        = "#FF426D"
	ColorContainers    = "#A855F7"
	ColorCPU           = "#29E64A"
	ColorRAM           = "#0084FF"
	ColorDisk          = "#F2C300"
	ColorNeutralIcon   = "#8A8F91"
)

var shipAccentColors = []string{
	"#29E64A",
	"#00D6D6",
	"#0084FF",
	"#A855F7",
	"#F2C300",
	"#FF8A00",
	"#FF426D",
	"#B7F000",
}

type IconSet struct {
	Ships      string
	Containers string
	Ship       string
	Status     string
	CPU        string
	RAM        string
	Disk       string
	Uptime     string
}

var Icons = IconSet{
	Ships:      "",
	Containers: "",
	Ship:       "",
	Status:     "",
	CPU:        "",
	RAM:        "󰘚",
	Disk:       "",
	Uptime:     "",
}

var (
	BackgroundStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorTextPrimary))

	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorTextPrimary)).
			Bold(true)

	SummaryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorTextSecondary))

	ContentStyle = lipgloss.NewStyle()

	TopBarStyle = lipgloss.NewStyle()

	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder())

	HeaderStyle = lipgloss.NewStyle().
			Bold(true)

	LabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorTextSecondary))

	ValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorTextPrimary))

	MutedValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorTextMuted))

	DimValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorTextDim))

	DividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorDivider))

	OnlineStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorOnline))

	FailedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorFailed))

	NeutralIconStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorNeutralIcon))

	RowStyle = lipgloss.NewStyle()
)

func ShipAccentByIndex(index int) string {
	return shipAccentColors[index%len(shipAccentColors)]
}

func StatusLabel(status Status) string {
	switch status {
	case StatusOnline:
		return "ONLINE"
	case StatusFailed:
		return "FAILED"
	default:
		return "PENDING"
	}
}

func StatusValue(status Status) string {
	switch status {
	case StatusOnline:
		return "OK"
	case StatusFailed:
		return "FAILED"
	default:
		return "--"
	}
}

func StatusStyle(status Status) lipgloss.Style {
	switch status {
	case StatusOnline:
		return OnlineStyle
	case StatusFailed:
		return FailedStyle
	default:
		return MutedValueStyle
	}
}
