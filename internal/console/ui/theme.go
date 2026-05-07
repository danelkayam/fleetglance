package ui

import (
	"hash/fnv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	colorBackground       = "#000808"
	colorPanelBackground  = colorBackground
	colorHeaderBackground = colorBackground
	colorDivider          = "#1A2426"
	colorTextPrimary      = "#E6E6E6"
	colorTextSecondary    = "#B8B8B8"
	colorTextMuted        = "#8A8F91"
	colorTextDim          = "#5B6366"
	colorOnline           = "#29E64A"
	colorFailed           = "#FF426D"
	colorCPU              = "#29E64A"
	colorRAM              = "#0084FF"
	colorDisk             = "#F2C300"
	colorNeutralIcon      = "#C8CDD0"
)

var shipAccentColors = []string{
	"#29E64A",
	"#0084FF",
	"#A855F7",
	"#00D6D6",
	"#F2C300",
	"#FF426D",
	"#FF8A00",
}

var namedShipAccents = map[string]string{
	"donnager":  "#29E64A",
	"rocinante": "#0084FF",
	"romulus":   "#A855F7",
	"nostromo":  "#00D6D6",
	"tycho":     "#F2C300",
	"betty":     "#FF426D",
	"serenity":  "#FF8A00",
}

var namedShipIcons = map[string]string{
	"donnager":  "✣",
	"rocinante": "✦",
	"romulus":   "✚",
	"nostromo":  "☼",
	"tycho":     "⬡",
	"betty":     "◎",
	"serenity":  "✤",
}

type iconSet struct {
	ships      string
	containers string
	status     string
	cpu        string
	ram        string
	disk       string
	uptime     string
}

var icons = iconSet{
	ships:      "✦",
	containers: "◇",
	status:     "◎",
	cpu:        "▣",
	ram:        "▤",
	disk:       "▬",
	uptime:     "◷",
}

var (
	backgroundStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(colorBackground)).
			Foreground(lipgloss.Color(colorTextPrimary))

	titleStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(colorBackground)).
			Foreground(lipgloss.Color(colorTextPrimary)).
			Bold(true)

	subtitleStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(colorBackground)).
			Foreground(lipgloss.Color(colorTextMuted))

	summaryStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(colorBackground)).
			Foreground(lipgloss.Color(colorTextSecondary))

	contentStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(colorBackground))

	topBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(colorBackground))

	panelStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(colorPanelBackground)).
			BorderBackground(lipgloss.Color(colorBackground)).
			Border(lipgloss.RoundedBorder())

	headerStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(colorHeaderBackground)).
			Bold(true)

	labelStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(colorBackground)).
			Foreground(lipgloss.Color(colorTextSecondary))

	valueStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(colorBackground)).
			Foreground(lipgloss.Color(colorTextPrimary))

	mutedValueStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(colorBackground)).
			Foreground(lipgloss.Color(colorTextMuted))

	dimValueStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(colorBackground)).
			Foreground(lipgloss.Color(colorTextDim))

	dividerStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(colorBackground)).
			Foreground(lipgloss.Color(colorDivider))

	onlineStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(colorBackground)).
			Foreground(lipgloss.Color(colorOnline))

	failedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(colorBackground)).
			Foreground(lipgloss.Color(colorFailed))

	neutralIconStyle = lipgloss.NewStyle().
				Background(lipgloss.Color(colorBackground)).
				Foreground(lipgloss.Color(colorNeutralIcon))
)

func shipAccent(name string) string {
	normalized := strings.ToLower(strings.TrimSpace(name))
	if color, ok := namedShipAccents[normalized]; ok {
		return color
	}

	hash := fnv.New32a()
	_, _ = hash.Write([]byte(normalized))

	return shipAccentColors[int(hash.Sum32())%len(shipAccentColors)]
}

func shipIcon(name string) string {
	normalized := strings.ToLower(strings.TrimSpace(name))
	if icon, ok := namedShipIcons[normalized]; ok {
		return icon
	}

	return "✦"
}
