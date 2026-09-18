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
	case TabGamemodes:
		sb.WriteString(m.renderGamemodesTab(width))
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

	tabNames := []string{"1. Dashboard", "2. Gamemodes", "3. Rutas & DLC", "4. Red"}
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

	modeName := getSpecialModeName(m.Config.ServerCfg.SpecialMode)
	srvContent := fmt.Sprintf(
		"%s %-18s %s\n%s %-18s %s\n%s %-18s %s\n\n%-15s %s",
		LabelStyle.Render("●"), "Servidor Dedicado", srvStatus,
		LabelStyle.Render("●"), "Milk Bar Launcher", milkStatus,
		LabelStyle.Render("●"), "Cemu 1.26.2", cemuStatus,
		LabelStyle.Render("Modo Servidor:"), ValueStyle.Render(modeName),
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
		"%-15s %s\n%-15s %s\n%-15s %s\n\n%-15s %s",
		LabelStyle.Render("IP Tailscale:"), ValueStyle.Render(tsIP),
		LabelStyle.Render("IP Local LAN:"), ValueStyle.Render(m.Network.LocalIP),
		LabelStyle.Render("Puerto:"), ValueStyle.Render(m.Config.ServerCfg.Port+" (UDP/TCP)"),
		LabelStyle.Render("Clave Servidor:"), ValueStyle.Render(nonEmpty(m.Config.ServerCfg.Password, "Sin clave")),
	)
	networkBox := BoxStyle.Width(halfWidth).Render(
		TitleStyle.Render("CONECTIVIDAD") + "\n" + netContent,
	)

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, servicesBox, " ", networkBox)

	// Console / Logs Viewport
	consoleTitle := TitleStyle.Render("CONSOLA & REGISTROS (SERVIDOR DEDICADO)")
	vpContent := m.Viewport.View()

	consoleBox := ActiveBoxStyle.Width(width - 4).Render(
		consoleTitle + "\n" + vpContent,
	)

	return topRow + "\n" + consoleBox
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
			return lipgloss.NewStyle().Foreground(ColorMuted).Render("[-] Bloqueado")
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
		statusBanner = lipgloss.NewStyle().Foreground(ColorWarning).Bold(true).Render(">> SERVIDOR ACTIVO: Los cambios se aplicaran reiniciando el servidor al guardar.")
	} else {
		statusBanner = lipgloss.NewStyle().Foreground(ColorSuccess).Render(">> SERVIDOR DETENIDO: Configure sus reglas y presione [s] para iniciar.")
	}

	var syncHeader string
	isLockedProgress := mode != 0
	isLockedEnemies := mode == 2

	if mode == 0 {
		syncHeader = LabelStyle.Render("OPCIONES DE SINCRONIZACION (MODO LIBRE / PERSONALIZABLE):")
	} else if mode == 1 {
		syncHeader = LabelStyle.Render("OPCIONES DE SINCRONIZACION (PROGRESO INDIVIDUAL - HUNTER VS SPEEDRUNNER):")
	} else {
		syncHeader = LabelStyle.Render("OPCIONES DE SINCRONIZACION (DESACTIVADAS POR SUPERVIVENCIA - DEATHSWAP):")
	}

	items := []string{
		statusBanner,
		"",
		fmt.Sprintf("%s: %s", cur(0, "Modo de Juego Especial"), KeyStyle.Render(modeStr)),
		"",
		syncHeader,
		fmt.Sprintf("%s %s Sincronizar Misiones (QuestSync)", cur(1, chk(m.Config.ServerCfg.QuestSync, isLockedProgress)), ""),
		fmt.Sprintf("%s %s Sincronizar Santuarios (ShrineSync)", cur(2, chk(m.Config.ServerCfg.ShrineSync, isLockedProgress)), ""),
		fmt.Sprintf("%s %s Sincronizar Torres (TowerSync)", cur(3, chk(m.Config.ServerCfg.TowerSync, isLockedProgress)), ""),
		fmt.Sprintf("%s %s Sincronizar Kologs (KorokSync)", cur(4, chk(m.Config.ServerCfg.KorokSync, isLockedProgress)), ""),
		fmt.Sprintf("%s %s Sincronizar Enemigos (EnemySync)", cur(5, chk(m.Config.ServerCfg.EnemySync, isLockedEnemies)), ""),
		fmt.Sprintf("%s %s Sincronizar Mazmorras (DungeonSync)", cur(6, chk(m.Config.ServerCfg.DungeonSync, isLockedProgress)), ""),
		fmt.Sprintf("%s %s Sincronizar Ubicaciones (LocationSync)", cur(7, chk(m.Config.ServerCfg.LocationSync, isLockedProgress)), ""),
		"",
		LabelStyle.Render("PARAMETROS DEL SERVIDOR:"),
		fmt.Sprintf("%s: %s", cur(8, "Contrasena"), m.ServerPassIn.View()),
		fmt.Sprintf("%s: %s", cur(9, "Descripcion"), m.ServerDescIn.View()),
		"",
		DescStyle.Render("[Espacio/Enter: Alternar opcion]  [s / Enter en inputs: Guardar y Aplicar]  [Esc: Volver]"),
	}

	return BoxStyle.Width(boxWidth).Render(
		TitleStyle.Render("CONFIGURACION DE GAMEMODES Y SERVIDOR DEDICADO") + "\n\n" + strings.Join(items, "\n"),
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
		return StatusWarn.String()
	}

	fieldHeader := func(idx int, title string, ok bool) string {
		prefix := "  "
		if m.PathFocusIdx == idx {
			prefix = lipgloss.NewStyle().Bold(true).Foreground(ColorHighlight).Render("▶ ")
		}
		return fmt.Sprintf("%s%s %s  %s", prefix, LabelStyle.Render(title), badge(ok), KeyStyle.Render("[Ctrl+O / F2: Explorar]"))
	}

	content := fmt.Sprintf(
		"%s\n  %s\n  %s %s\n\n"+
			"%s\n  %s\n  %s %s\n\n"+
			"%s\n  %s\n  %s %s\n\n"+
			"%s",
		fieldHeader(0, "1. Carpeta del Juego Base (Zelda BotW):", baseValid),
		m.BaseInput.View(),
		DescStyle.Render("Estado:"), ValueStyle.Render(baseMsg),

		fieldHeader(1, "2. Carpeta de Actualizacion (Update v208):", upValid),
		m.UpdateInput.View(),
		DescStyle.Render("Estado:"), ValueStyle.Render(upMsg),

		fieldHeader(2, "3. Carpeta de DLC (v80):", dlcValid),
		m.DLCInput.View(),
		DescStyle.Render("Estado:"), ValueStyle.Render(dlcMsg),

		DescStyle.Render("[Tab/Flechas: Cambiar campo]  [Ctrl+O / F2: Abrir Explorador GUI]  [Enter: Guardar]  [Esc: Volver]"),
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
			"  - Puerto por defecto: %s (TCP/UDP)\n\n"+
			"%s",
		LabelStyle.Render("INTERFACES DETECTADAS:"),
		"Tailscale VPN:", ValueStyle.Render(nonEmpty(m.Network.TailscaleIP, "Inactivo")), tsBadge,
		"Red Local (LAN):", ValueStyle.Render(m.Network.LocalIP),
		"ZeroTier VPN:", ValueStyle.Render(nonEmpty(m.Network.ZeroTierIP, "Inactivo")), ztBadge,

		LabelStyle.Render("INSTRUCCIONES DE CONEXION MULTIJUGADOR:"),
		KeyStyle.Render("127.0.0.1"), KeyStyle.Render(m.Config.ServerCfg.Port),
		KeyStyle.Render(nonEmpty(m.Network.TailscaleIP, m.Network.LocalIP)),
		KeyStyle.Render(m.Config.ServerCfg.Port),

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
		KeyStyle.Render("[g]") + " " + DescStyle.Render("Gamemodes"),
		KeyStyle.Render("[r]") + " " + DescStyle.Render("Rutas"),
		KeyStyle.Render("[t]") + " " + DescStyle.Render("Copiar IP"),
		KeyStyle.Render("[Tab]") + " " + DescStyle.Render("Pestana"),
		KeyStyle.Render("[q]") + " " + DescStyle.Render("Salir"),
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, strings.Join(keys, "  "))
}

func getSpecialModeName(code int) string {
	switch code {
	case 1:
		return "Hunter vs Speedrunner"
	case 2:
		return "DeathSwap"
	default:
		return "Cooperativo Estandar (Libre)"
	}
}

func nonEmpty(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
