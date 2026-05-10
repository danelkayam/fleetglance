package ui

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	colorDivider       = "#1A2426"
	colorTextPrimary   = "#E6E6E6"
	colorTextSecondary = "#B8B8B8"
	colorTextMuted     = "#8A8F91"
	colorTextDim       = "#5B6366"
	colorOnline        = "#29E64A"
	colorFailed        = "#FF426D"
	colorContainers    = "#A855F7"
	colorCPU           = "#29E64A"
	colorRAM           = "#0084FF"
	colorDisk          = "#F2C300"
	colorNeutralIcon   = "#8A8F91"
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

type iconSet struct {
	ships      string
	containers string
	ship       string
	status     string
	cpu        string
	ram        string
	disk       string
	uptime     string
}

var icons = iconSet{
	ships:      "",
	containers: "",
	ship:       "",
	status:     "",
	cpu:        "",
	ram:        "󰘚",
	disk:       "",
	uptime:     "",
}

var (
	backgroundStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextPrimary))

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextPrimary)).
			Bold(true)

	summaryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextSecondary))

	contentStyle = lipgloss.NewStyle()

	topBarStyle = lipgloss.NewStyle()

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder())

	headerStyle = lipgloss.NewStyle().
			Bold(true)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextSecondary))

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextPrimary))

	mutedValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextMuted))

	dimValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextDim))

	dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorDivider))

	onlineStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorOnline))

	failedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorFailed))

	neutralIconStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorNeutralIcon))

	rowStyle = lipgloss.NewStyle()
)

func shipAccentByIndex(index int) string {
	return shipAccentColors[index%len(shipAccentColors)]
}
