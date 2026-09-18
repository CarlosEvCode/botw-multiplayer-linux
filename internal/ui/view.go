package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	width := m.Width
	if width < 80 {
		width = 80
	}

	var sb strings.Builder

	// 1. Header & Tabs
	title := TitleStyle.Render("ZELDA: BREATH OF THE WILD MULTIPLAYER MANAGER")
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
	case TabPaths:
		sb.WriteString(m.renderPathsTab(width))
	case TabNetwork:
		sb.WriteString(m.renderNetworkTab(width))
	}

	// 4. Footer Keybindings
	sb.WriteString("\n" + m.renderFooter(width))

	return sb.String()
}

func (m Model) renderTabs() string {
	var tabs []string

	tabNames := []string{"1. Dashboard", "2. Rutas & DLC", "3. Red & Conectividad"}
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

	var milkStatus string
	if m.Process.IsMilkBarRunning() {
		milkStatus = StatusRunning.String()
	} else {
		milkStatus = StatusStopped.String()
	}

	var cemuStatus string
	if m.Process.IsCemuRunning() {
		cemuStatus = StatusRunning.String()
	} else {
		cemuStatus = StatusStopped.String()
	}

	srvContent := fmt.Sprintf(
		"%s %-20s %s\n%s %-20s %s\n%s %-20s %s",
		LabelStyle.Render("●"), "Servidor Dedicado", srvStatus,
		LabelStyle.Render("●"), "Milk Bar Launcher", milkStatus,
		LabelStyle.Render("●"), "Cemu 1.26.2", cemuStatus,
	)
	servicesBox := BoxStyle.Width(halfWidth).Render(
		TitleStyle.Render("ESTADO DE SERVICIOS") + "\n" + srvContent,
	)

	// Network Quick Info
	tsIP := m.Network.TailscaleIP
	if tsIP == "" {
		tsIP = "No detectado"
	}
	netContent := fmt.Sprintf(
		"%-15s %s\n%-15s %s\n%-15s %s",
		LabelStyle.Render("IP Tailscale:"), ValueStyle.Render(tsIP),
		LabelStyle.Render("IP Local LAN:"), ValueStyle.Render(m.Network.LocalIP),
		LabelStyle.Render("Puerto:"), ValueStyle.Render("5050 (UDP/TCP)"),
	)
	networkBox := BoxStyle.Width(halfWidth).Render(
		TitleStyle.Render("CONECTIVIDAD") + "\n" + netContent,
	)

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, servicesBox, " ", networkBox)

	// Console / Logs Viewport
	consoleTitle := TitleStyle.Render("CONSOLA & REGISTROS (SERVIDOR DEDICADO)")
	vpContent := m.Viewport.View()

	inputLine := ""
	if m.InputFocused {
		inputLine = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true).Render("> ") + m.CmdInput.View()
	} else {
		inputLine = DescStyle.Render("[Presione 'i' para enviar comando al servidor]")
	}

	consoleBox := ActiveBoxStyle.Width(width - 4).Render(
		consoleTitle + "\n" + vpContent + "\n" + inputLine,
	)

	return topRow + "\n" + consoleBox
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
		return StatusWarn.String()
	}

	content := fmt.Sprintf(
		"%s %s\n%s\n  %s %s\n\n"+
			"%s %s\n%s\n  %s %s\n\n"+
			"%s %s\n%s\n  %s %s\n\n"+
			"%s",
		LabelStyle.Render("1. Carpeta del Juego Base:"), badge(baseValid),
		m.BaseInput.View(),
		DescStyle.Render("Estado:"), ValueStyle.Render(baseMsg),

		LabelStyle.Render("2. Carpeta de Actualizacion (Update v208):"), badge(upValid),
		m.UpdateInput.View(),
		DescStyle.Render("Estado:"), ValueStyle.Render(upMsg),

		LabelStyle.Render("3. Carpeta de DLC (v80):"), badge(dlcValid),
		m.DLCInput.View(),
		DescStyle.Render("Estado:"), ValueStyle.Render(dlcMsg),

		DescStyle.Render("[Tab/Flechas: Cambiar campo]  [Enter/Ctrl+S: Guardar y Aplicar]  [Esc: Volver]"),
	)

	return BoxStyle.Width(boxWidth).Render(
		TitleStyle.Render("GESTOR DE RUTAS DEL JUEGO, UPDATE Y DLC") + "\n\n" + content,
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
			"  - Host (Servidor): Conecta localmente a %s en el puerto %s\n"+
			"  - Clientes remotos: Ingresan tu IP de Tailscale (%s) en Milk Bar Launcher\n"+
			"  - Puerto por defecto: 5050 (TCP/UDP)\n\n"+
			"%s",
		LabelStyle.Render("INTERFACES DETECTADAS:"),
		"Tailscale VPN:", ValueStyle.Render(nonEmpty(m.Network.TailscaleIP, "Inactivo")), tsBadge,
		"Red Local (LAN):", ValueStyle.Render(m.Network.LocalIP),
		"ZeroTier VPN:", ValueStyle.Render(nonEmpty(m.Network.ZeroTierIP, "Inactivo")), ztBadge,

		LabelStyle.Render("INSTRUCCIONES DE CONEXION MULTIJUGADOR:"),
		KeyStyle.Render("127.0.0.1"), KeyStyle.Render("5050"),
		KeyStyle.Render(nonEmpty(m.Network.TailscaleIP, m.Network.LocalIP)),

		DescStyle.Render("[t: Copiar IP al portapapeles]  [Tab: Siguiente pestana]"),
	)

	return BoxStyle.Width(boxWidth).Render(
		TitleStyle.Render("DETALLES DE RED & MULTIJUGADOR") + "\n\n" + content,
	)
}

func (m Model) renderFooter(width int) string {
	keys := []string{
		KeyStyle.Render("[s]") + " " + DescStyle.Render("Server"),
		KeyStyle.Render("[m]") + " " + DescStyle.Render("MilkBar"),
		KeyStyle.Render("[c]") + " " + DescStyle.Render("Cemu"),
		KeyStyle.Render("[t]") + " " + DescStyle.Render("Copiar IP"),
		KeyStyle.Render("[r]") + " " + DescStyle.Render("Rutas"),
		KeyStyle.Render("[i]") + " " + DescStyle.Render("Comando"),
		KeyStyle.Render("[Tab]") + " " + DescStyle.Render("Pestana"),
		KeyStyle.Render("[q]") + " " + DescStyle.Render("Salir"),
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, strings.Join(keys, "  "))
}

func nonEmpty(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
