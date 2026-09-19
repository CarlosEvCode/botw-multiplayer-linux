package ui

import (
	"fmt"
	"strings"

	"github.com/CarlosEvCode/botw-multiplayer-linux/internal/i18n"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	width := m.Width
	if width < 80 {
		width = 80
	}

	var sb strings.Builder

	// 1. Header & Tabs
	title := TitleStyle.Render(i18n.T("app.title"))
	tabs := m.renderTabs()
	header := lipgloss.JoinHorizontal(lipgloss.Center, title, "  ", tabs)
	sb.WriteString(header + "\n\n")

	// 2. Notification bar
	if m.Notification != "" {
		notif := lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorHighlight).
			Padding(0, 1).
			Render(">> " + m.Notification)
		sb.WriteString(notif + "\n")
	}

	// 3. Tab Contents
	switch m.CurrentTab {
	case TabDashboard:
		sb.WriteString(m.renderDashboard(width))
	case TabClient:
		sb.WriteString(m.renderClientTab(width))
	case TabGamemodes:
		sb.WriteString(m.renderGamemodesTab(width))
	case TabPaths:
		sb.WriteString(m.renderPathsTab(width))
	case TabNetwork:
		sb.WriteString(m.renderNetworkTab(width))
	case TabSettings:
		sb.WriteString(m.renderSettingsTab(width))
	}

	// 4. Footer Keybindings
	sb.WriteString("\n" + m.renderFooter(width))

	return sb.String()
}

func (m Model) renderTabs() string {
	var tabs []string

	tabNames := []string{
		i18n.T("tab.dashboard"),
		i18n.T("tab.client"),
		i18n.T("tab.gamemodes"),
		i18n.T("tab.paths"),
		i18n.T("tab.network"),
		i18n.T("tab.settings"),
	}
	for i, name := range tabNames {
		if Tab(i) == m.CurrentTab {
			tabs = append(tabs, TabActiveStyle.Render(name))
		} else {
			tabs = append(tabs, TabInactiveStyle.Render(name))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

func (m Model) renderDashboard(width int) string {
	halfWidth := (width - 6) / 2
	if halfWidth < 35 {
		halfWidth = 35
	}

	// Services Panel
	var srvStatus string
	if m.Process.IsServerRunning() {
		srvStatus = StatusRunning.String()
	} else {
		srvStatus = StatusStopped.String()
	}

	var clientStatus string
	if m.Process.IsClientRunning() {
		clientStatus = StatusRunning.String()
	} else {
		clientStatus = StatusStopped.String()
	}

	var cemuStatus string
	if m.Process.IsCemuRunning() {
		cemuStatus = StatusRunning.String()
	} else {
		cemuStatus = StatusStopped.String()
	}

	modeName := getSpecialModeName(m.Config.ServerCfg.SpecialMode)
	srvContent := fmt.Sprintf(
		"%s %-22s %s\n%s %-22s %s\n%s %-22s %s\n\n%-15s %s",
		LabelStyle.Render("●"), i18n.T("dash.srv_dedicated"), srvStatus,
		LabelStyle.Render("●"), i18n.T("dash.srv_client"), clientStatus,
		LabelStyle.Render("●"), i18n.T("dash.srv_cemu"), cemuStatus,
		LabelStyle.Render(i18n.T("dash.srv_mode")), ValueStyle.Render(modeName),
	)
	servicesBox := BoxStyle.Width(halfWidth).Render(
		TitleStyle.Render(i18n.T("dash.services_title")) + "\n" + srvContent,
	)

	// Network Quick Info
	var vpnLabel, vpnIP string
	if m.Network.TailscaleIP != "" {
		vpnLabel = "IP Tailscale:"
		vpnIP = m.Network.TailscaleIP
	} else if m.Network.ZeroTierIP != "" {
		vpnLabel = "IP ZeroTier:"
		vpnIP = m.Network.ZeroTierIP
	} else {
		vpnLabel = "IP VPN:"
		vpnIP = i18n.T("dash.net_notfound")
	}

	netContent := fmt.Sprintf(
		"%-15s %s\n%-15s %s\n%-15s %s\n\n%-15s %s",
		LabelStyle.Render(vpnLabel), ValueStyle.Render(vpnIP),
		LabelStyle.Render(i18n.T("dash.net_local")), ValueStyle.Render(m.Network.LocalIP),
		LabelStyle.Render(i18n.T("dash.net_port")), ValueStyle.Render(m.Config.ServerCfg.Port+" (UDP/TCP)"),
		LabelStyle.Render(i18n.T("dash.net_pass")), ValueStyle.Render(nonEmpty(m.Config.ServerCfg.Password, i18n.T("dash.net_nopass"))),
	)
	networkBox := BoxStyle.Width(halfWidth).Render(
		TitleStyle.Render(i18n.T("dash.net_title")) + "\n" + netContent,
	)

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, servicesBox, " ", networkBox)

	// Console / Logs Viewport
	consoleTitle := TitleStyle.Render(i18n.T("dash.console_title"))
	vpContent := m.Viewport.View()

	consoleBox := ActiveBoxStyle.Width(width - 4).Render(
		consoleTitle + "\n" + vpContent,
	)

	return topRow + "\n" + consoleBox
}

func (m Model) renderClientTab(width int) string {
	boxWidth := width - 4

	cur := func(idx int, text string) string {
		if m.ClientCursor == idx {
			return lipgloss.NewStyle().Bold(true).Foreground(ColorHighlight).Render("▶ " + text)
		}
		return "  " + text
	}

	var clientStatusStr string
	if m.Process.IsClientRunning() {
		clientStatusStr = StatusRunning.String() + " " + lipgloss.NewStyle().Foreground(ColorSuccess).Render(i18n.T("client.status_running"))
	} else {
		clientStatusStr = StatusStopped.String() + " " + lipgloss.NewStyle().Foreground(ColorMuted).Render(i18n.T("client.status_ready"))
	}

	items := []string{
		fmt.Sprintf("%s %s", LabelStyle.Render(i18n.T("client.status")), clientStatusStr),
		"",
		LabelStyle.Render(i18n.T("client.header_params")),
		fmt.Sprintf("%s: %s", cur(0, i18n.T("client.target_ip")), m.ClientIPIn.View()),
		fmt.Sprintf("%s: %s", cur(1, i18n.T("client.target_port")), m.ClientPortIn.View()),
		fmt.Sprintf("%s: %s", cur(2, i18n.T("client.password")), m.ClientPassIn.View()),
		fmt.Sprintf("%s: %s", cur(3, i18n.T("client.player_name")), m.ClientNameIn.View()),
		fmt.Sprintf("%s: %s", cur(4, i18n.T("client.model")), m.ClientModelIn.View()),
		"",
		LabelStyle.Render(i18n.T("client.quick_actions")),
		fmt.Sprintf("  %s   %s   %s   %s",
			KeyStyle.Render(i18n.T("client.btn_connect")),
			KeyStyle.Render(i18n.T("client.btn_local")),
			KeyStyle.Render(i18n.T("client.btn_vpn")),
			KeyStyle.Render(i18n.T("client.btn_save")),
		),
		"",
		DescStyle.Render(i18n.T("client.footer_hint")),
	}

	return BoxStyle.Width(boxWidth).Render(
		TitleStyle.Render(i18n.T("client.title")) + "\n\n" + strings.Join(items, "\n"),
	)
}

func (m Model) renderGamemodesTab(width int) string {
	boxWidth := width - 4

	cur := func(idx int, text string) string {
		if m.GamemodeCursor == idx {
			return lipgloss.NewStyle().Bold(true).Foreground(ColorHighlight).Render("▶ " + text)
		}
		return "  " + text
	}

	chk := func(b bool, locked bool) string {
		if locked {
			return lipgloss.NewStyle().Foreground(ColorMuted).Render(i18n.T("gm.locked"))
		}
		if b {
			return lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true).Render("[X]")
		}
		return lipgloss.NewStyle().Foreground(ColorMuted).Render("[ ]")
	}

	mode := m.Config.ServerCfg.SpecialMode
	modeStr := fmt.Sprintf("%s (%d)", getSpecialModeName(mode), mode)

	var statusBanner string
	if m.Process.IsServerRunning() {
		statusBanner = lipgloss.NewStyle().Foreground(ColorWarning).Bold(true).Render(i18n.T("gm.srv_running"))
	} else {
		statusBanner = lipgloss.NewStyle().Foreground(ColorSuccess).Render(i18n.T("gm.srv_stopped"))
	}

	var syncHeader string
	isLockedProgress := mode != 0
	isLockedEnemies := mode == 2

	if mode == 0 {
		syncHeader = LabelStyle.Render(i18n.T("gm.header_coop"))
	} else if mode == 1 {
		syncHeader = LabelStyle.Render(i18n.T("gm.header_hunter"))
	} else {
		syncHeader = LabelStyle.Render(i18n.T("gm.header_deathswap"))
	}

	items := []string{
		statusBanner,
		"",
		fmt.Sprintf("%s: %s  %s", cur(0, i18n.T("gm.special_mode")), KeyStyle.Render("◄  "+modeStr+"  ►"), DescStyle.Render("[← / →]")),
		"",
		syncHeader,
		fmt.Sprintf("%s %s %s", cur(1, chk(m.Config.ServerCfg.QuestSync, isLockedProgress)), i18n.T("gm.sync_quest"), ""),
		fmt.Sprintf("%s %s %s", cur(2, chk(m.Config.ServerCfg.ShrineSync, isLockedProgress)), i18n.T("gm.sync_shrine"), ""),
		fmt.Sprintf("%s %s %s", cur(3, chk(m.Config.ServerCfg.TowerSync, isLockedProgress)), i18n.T("gm.sync_tower"), ""),
		fmt.Sprintf("%s %s %s", cur(4, chk(m.Config.ServerCfg.KorokSync, isLockedProgress)), i18n.T("gm.sync_korok"), ""),
		fmt.Sprintf("%s %s %s", cur(5, chk(m.Config.ServerCfg.EnemySync, isLockedEnemies)), i18n.T("gm.sync_enemy"), ""),
		fmt.Sprintf("%s %s %s", cur(6, chk(m.Config.ServerCfg.DungeonSync, isLockedProgress)), i18n.T("gm.sync_dungeon"), ""),
		fmt.Sprintf("%s %s %s", cur(7, chk(m.Config.ServerCfg.LocationSync, isLockedProgress)), i18n.T("gm.sync_location"), ""),
		"",
		LabelStyle.Render(i18n.T("gm.server_params")),
		fmt.Sprintf("%s: %s", cur(8, i18n.T("gm.password")), m.ServerPassIn.View()),
		fmt.Sprintf("%s: %s", cur(9, i18n.T("gm.description")), m.ServerDescIn.View()),
		"",
		DescStyle.Render(i18n.T("gm.footer_hint")),
	}

	return BoxStyle.Width(boxWidth).Render(
		TitleStyle.Render(i18n.T("gm.title")) + "\n\n" + strings.Join(items, "\n"),
	)
}

func (m Model) renderPathsTab(width int) string {
	boxWidth := width - 4

	baseValid, baseMsg := m.Config.ValidateBaseGame()
	upValid, upMsg := m.Config.ValidateUpdate()
	dlcValid, dlcMsg := m.Config.ValidateDLC()

	badge := func(ok bool) string {
		if ok {
			return StatusOk.String()
		}
		return StatusError.String()
	}

	formatMsg := func(ok bool, msg string) string {
		if ok {
			return lipgloss.NewStyle().Foreground(ColorSuccess).Render(msg)
		}
		return lipgloss.NewStyle().Foreground(ColorDanger).Render(msg)
	}

	fieldHeader := func(idx int, title string, ok bool) string {
		prefix := "  "
		if m.PathFocusIdx == idx {
			prefix = lipgloss.NewStyle().Bold(true).Foreground(ColorHighlight).Render("▶ ")
		}
		return fmt.Sprintf("%s%s %s  %s", prefix, LabelStyle.Render(title), badge(ok), KeyStyle.Render(i18n.T("paths.browse")))
	}

	content := fmt.Sprintf(
		"%s\n  %s\n  %s %s\n\n"+
			"%s\n  %s\n  %s %s\n\n"+
			"%s\n  %s\n  %s %s\n\n"+
			"%s",
		fieldHeader(0, i18n.T("paths.base"), baseValid),
		m.BaseInput.View(),
		DescStyle.Render(i18n.T("paths.status")), formatMsg(baseValid, baseMsg),

		fieldHeader(1, i18n.T("paths.update"), upValid),
		m.UpdateInput.View(),
		DescStyle.Render(i18n.T("paths.status")), formatMsg(upValid, upMsg),

		fieldHeader(2, i18n.T("paths.dlc"), dlcValid),
		m.DLCInput.View(),
		DescStyle.Render(i18n.T("paths.status")), formatMsg(dlcValid, dlcMsg),

		DescStyle.Render(i18n.T("paths.footer_hint")),
	)

	return BoxStyle.Width(boxWidth).Render(
		TitleStyle.Render(i18n.T("paths.title")) + "\n\n" + content,
	)
}

func (m Model) renderNetworkTab(width int) string {
	boxWidth := width - 4

	tsBadge := StatusStopped.String()
	if m.Network.TailscaleIP != "" {
		tsBadge = StatusActive.String()
	}

	ztBadge := StatusStopped.String()
	if m.Network.ZeroTierIP != "" {
		ztBadge = StatusActive.String()
	}

	content := fmt.Sprintf(
		"%s\n"+
			"  %-18s %-25s %s\n"+
			"  %-18s %-25s\n"+
			"  %-18s %-25s %s\n\n"+
			"%s\n"+
			"  "+fmt.Sprintf(i18n.T("net.host_info"), KeyStyle.Render("127.0.0.1"), KeyStyle.Render(m.Config.ServerCfg.Port))+"\n"+
			"  "+fmt.Sprintf(i18n.T("net.client_info"), KeyStyle.Render(nonEmpty(m.Network.TailscaleIP, m.Network.LocalIP)))+"\n"+
			"  "+fmt.Sprintf(i18n.T("net.port_info"), KeyStyle.Render(m.Config.ServerCfg.Port))+"\n\n"+
			"%s",
		LabelStyle.Render(i18n.T("net.interfaces")),
		i18n.T("net.tailscale"), ValueStyle.Render(nonEmpty(m.Network.TailscaleIP, i18n.T("net.inactive"))), tsBadge,
		i18n.T("net.local"), ValueStyle.Render(m.Network.LocalIP),
		i18n.T("net.zerotier"), ValueStyle.Render(nonEmpty(m.Network.ZeroTierIP, i18n.T("net.inactive"))), ztBadge,

		LabelStyle.Render(i18n.T("net.instructions")),
		DescStyle.Render(i18n.T("net.footer_hint")),
	)

	return BoxStyle.Width(boxWidth).Render(
		TitleStyle.Render(i18n.T("net.title")) + "\n\n" + content,
	)
}

func (m Model) renderSettingsTab(width int) string {
	boxWidth := width - 4

	var langDisplay string
	if i18n.CurrentLang == i18n.LangEN {
		langDisplay = "English [Default]"
	} else {
		langDisplay = "Español [Spanish]"
	}

	content := fmt.Sprintf(
		"%s\n"+
			"  %s  %s  %s\n\n"+
			"%s\n"+
			"  ● %s\n"+
			"  ● %s\n"+
			"  ● %s\n"+
			"  ● %s\n"+
			"  ● %s\n\n"+
			"%s",
		LabelStyle.Render(i18n.T("set.pref_title")),
		i18n.T("set.lang_label"), KeyStyle.Render("◄  "+langDisplay+"  ►"), DescStyle.Render("[← / → / Space]"),

		LabelStyle.Render(i18n.T("set.about_title")),
		ValueStyle.Render(i18n.T("set.app_name")),
		ValueStyle.Render(i18n.T("set.version")+"  •  "+i18n.T("set.author")),
		ValueStyle.Render(i18n.T("set.github")),
		ValueStyle.Render(i18n.T("set.stack")),
		DescStyle.Render(i18n.T("set.desc")),

		DescStyle.Render(i18n.T("set.footer_hint")),
	)

	return BoxStyle.Width(boxWidth).Render(
		TitleStyle.Render(i18n.T("set.title")) + "\n\n" + content,
	)
}

func (m Model) renderFooter(width int) string {
	keys := []string{
		KeyStyle.Render("[s]") + " " + DescStyle.Render(i18n.T("foot.server")),
		KeyStyle.Render("[c]") + " " + DescStyle.Render(i18n.T("foot.client")),
		KeyStyle.Render("[e]") + " " + DescStyle.Render(i18n.T("foot.cemu")),
		KeyStyle.Render("[t]") + " " + DescStyle.Render(i18n.T("foot.copy_ip")),
		KeyStyle.Render("[Tab]") + " " + DescStyle.Render(i18n.T("foot.tab")),
		KeyStyle.Render("[q]") + " " + DescStyle.Render(i18n.T("foot.quit")),
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, strings.Join(keys, "  "))
}

func getSpecialModeName(code int) string {
	switch code {
	case 1:
		return i18n.T("gm.mode_hunter")
	case 2:
		return i18n.T("gm.mode_deathswap")
	default:
		return i18n.T("gm.mode_coop")
	}
}

func nonEmpty(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

