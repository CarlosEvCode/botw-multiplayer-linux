package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/CarlosEvCode/botw-multiplayer-linux/internal/config"
	"github.com/CarlosEvCode/botw-multiplayer-linux/internal/network"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Viewport.Width = msg.Width - 6
		m.Viewport.Height = msg.Height - 17
		if m.Viewport.Height < 6 {
			m.Viewport.Height = 6
		}
		inputW := msg.Width - 10
		if inputW > 120 {
			inputW = 120
		} else if inputW < 40 {
			inputW = 40
		}
		m.BaseInput.Width = inputW
		m.UpdateInput.Width = inputW
		m.DLCInput.Width = inputW
		m.ServerPassIn.Width = 35
		m.ServerDescIn.Width = 50

	case PickedDirMsg:
		if msg.Err == nil && msg.Path != "" {
			cleanPath := config.NormalizeGameDir(msg.Path)
			switch msg.Field {
			case 0:
				m.BaseInput.SetValue(cleanPath)
			case 1:
				m.UpdateInput.SetValue(cleanPath)
			case 2:
				m.DLCInput.SetValue(cleanPath)
			}
			err := m.Config.SavePaths(m.BaseInput.Value(), m.UpdateInput.Value(), m.DLCInput.Value())
			if err != nil {
				m.SetNotification(fmt.Sprintf("Error guardando ruta: %v", err), 3*time.Second)
			} else {
				m.SetNotification("Ruta seleccionada y aplicada correctamente!", 3*time.Second)
			}
		}

	case TickMsg:
		m.Network = network.GetNetworkInfo()
		logs := m.Process.GetLogs()
		if len(logs) > 0 {
			m.Viewport.SetContent(strings.Join(logs, "\n"))
			m.Viewport.GotoBottom()
		}
		if time.Now().After(m.NotifyTimeout) {
			m.Notification = ""
		}
		cmds = append(cmds, tickEvery(1*time.Second))

	case tea.KeyMsg:
		// When typing in command input on dashboard
		if m.InputFocused {
			switch msg.String() {
			case "esc":
				m.InputFocused = false
				m.CmdInput.Blur()
				return m, nil
			case "enter":
				val := strings.TrimSpace(m.CmdInput.Value())
				if val != "" {
					_ = m.Process.SendServerInput(val)
					m.CmdInput.SetValue("")
				}
				return m, nil
			}
			m.CmdInput, cmd = m.CmdInput.Update(msg)
			return m, cmd
		}

		// When in Gamemodes Tab
		if m.CurrentTab == TabGamemodes {
			if m.GamemodeCursor == 8 { // Password input
				switch msg.String() {
				case "esc", "down", "tab":
					m.ServerPassIn.Blur()
					m.GamemodeCursor = (m.GamemodeCursor + 1) % 10
					m.updateGamemodeFocus()
					return m, nil
				case "up", "shift+tab":
					m.ServerPassIn.Blur()
					m.GamemodeCursor = (m.GamemodeCursor + 9) % 10
					m.updateGamemodeFocus()
					return m, nil
				case "enter", "ctrl+s":
					m.saveGamemodeSettings()
					return m, nil
				}
				m.ServerPassIn, cmd = m.ServerPassIn.Update(msg)
				return m, cmd
			} else if m.GamemodeCursor == 9 { // Description input
				switch msg.String() {
				case "esc", "down", "tab":
					m.ServerDescIn.Blur()
					m.GamemodeCursor = (m.GamemodeCursor + 1) % 10
					m.updateGamemodeFocus()
					return m, nil
				case "up", "shift+tab":
					m.ServerDescIn.Blur()
					m.GamemodeCursor = (m.GamemodeCursor + 9) % 10
					m.updateGamemodeFocus()
					return m, nil
				case "enter", "ctrl+s":
					m.saveGamemodeSettings()
					return m, nil
				}
				m.ServerDescIn, cmd = m.ServerDescIn.Update(msg)
				return m, cmd
			}

			switch msg.String() {
			case "esc":
				m.CurrentTab = TabDashboard
				return m, nil
			case "down", "j":
				m.GamemodeCursor = (m.GamemodeCursor + 1) % 10
				m.updateGamemodeFocus()
				return m, nil
			case "up", "k":
				m.GamemodeCursor = (m.GamemodeCursor + 9) % 10
				m.updateGamemodeFocus()
				return m, nil
			case " ", "enter":
				m.toggleGamemodeOption()
				return m, nil
			case "ctrl+s", "s":
				m.saveGamemodeSettings()
				return m, nil
			case "tab":
				m.CurrentTab = (m.CurrentTab + 1) % 4
				return m, nil
			}
		}

		// When in Paths Tab
		if m.CurrentTab == TabPaths {
			switch msg.String() {
			case "esc":
				m.CurrentTab = TabDashboard
				m.blurPathInputs()
				return m, nil
			case "tab", "down":
				m.PathFocusIdx = (m.PathFocusIdx + 1) % 3
				m.updatePathFocus()
				return m, nil
			case "shift+tab", "up":
				m.PathFocusIdx = (m.PathFocusIdx + 2) % 3
				m.updatePathFocus()
				return m, nil
			case "ctrl+o", "f2", "alt+o":
				var title, initDir string
				switch m.PathFocusIdx {
				case 0:
					title = "Selecciona Carpeta del Juego Base (Zelda BotW)"
					initDir = m.BaseInput.Value()
				case 1:
					title = "Selecciona Carpeta de Update (v208)"
					initDir = m.UpdateInput.Value()
				case 2:
					title = "Selecciona Carpeta de DLC (v80)"
					initDir = m.DLCInput.Value()
				}
				m.SetNotification("Abriendo explorador de carpetas...", 2*time.Second)
				return m, pickDirCmd(m.PathFocusIdx, title, initDir)
			case "enter", "ctrl+s":
				err := m.Config.SavePaths(m.BaseInput.Value(), m.UpdateInput.Value(), m.DLCInput.Value())
				if err != nil {
					m.SetNotification(fmt.Sprintf("Error guardando rutas: %v", err), 3*time.Second)
				} else {
					m.SetNotification("Rutas actualizadas y enlazadas correctamente!", 3*time.Second)
				}
				m.CurrentTab = TabDashboard
				m.blurPathInputs()
				return m, nil
			}

			switch m.PathFocusIdx {
			case 0:
				m.BaseInput, cmd = m.BaseInput.Update(msg)
			case 1:
				m.UpdateInput, cmd = m.UpdateInput.Update(msg)
			case 2:
				m.DLCInput, cmd = m.DLCInput.Update(msg)
			}
			return m, cmd
		}

		// Global Keybindings
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "tab":
			m.CurrentTab = (m.CurrentTab + 1) % 4
			if m.CurrentTab == TabPaths {
				m.PathFocusIdx = 0
				m.updatePathFocus()
			} else if m.CurrentTab == TabGamemodes {
				m.updateGamemodeFocus()
			}
			return m, nil

		case "s":
			if m.Process.IsServerRunning() {
				_ = m.Process.StopServer()
				m.SetNotification("Deteniendo servidor dedicado...", 2*time.Second)
			} else {
				err := m.Process.StartServer(m.Config.PrefixDir)
				if err != nil {
					m.SetNotification(fmt.Sprintf("Error iniciando servidor: %v", err), 3*time.Second)
				} else {
					m.SetNotification("Servidor iniciado en puerto 5050", 2*time.Second)
				}
			}
			return m, nil

		case "m":
			err := m.Process.StartMilkBar(m.Config.PrefixDir)
			if err != nil {
				m.SetNotification(fmt.Sprintf("Error iniciando MilkBar: %v", err), 3*time.Second)
			} else {
				m.SetNotification("Lanzando Milk Bar Launcher...", 2*time.Second)
			}
			return m, nil

		case "c":
			err := m.Process.StartCemu(m.Config.PrefixDir)
			if err != nil {
				m.SetNotification(fmt.Sprintf("Error iniciando Cemu: %v", err), 3*time.Second)
			} else {
				m.SetNotification("Lanzando Cemu 1.26.2...", 2*time.Second)
			}
			return m, nil

		case "t":
			if m.Network.TailscaleIP != "" {
				_ = network.CopyToClipboard(m.Network.TailscaleIP)
				m.SetNotification(fmt.Sprintf("IP de Tailscale copiada: %s", m.Network.TailscaleIP), 2*time.Second)
			} else if m.Network.LocalIP != "" {
				_ = network.CopyToClipboard(m.Network.LocalIP)
				m.SetNotification(fmt.Sprintf("IP Local copiada: %s", m.Network.LocalIP), 2*time.Second)
			}
			return m, nil

		case "g":
			m.CurrentTab = TabGamemodes
			m.GamemodeCursor = 0
			m.updateGamemodeFocus()
			return m, nil

		case "r":
			m.CurrentTab = TabPaths
			m.PathFocusIdx = 0
			m.updatePathFocus()
			return m, nil

		case "i", "/":
			if m.CurrentTab == TabDashboard {
				m.InputFocused = true
				m.CmdInput.Focus()
				return m, textinput.Blink
			}
		}
	}

	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) SetNotification(text string, dur time.Duration) {
	m.Notification = text
	m.NotifyTimeout = time.Now().Add(dur)
}

func (m *Model) toggleGamemodeOption() {
	switch m.GamemodeCursor {
	case 0:
		m.Config.ServerCfg.SpecialMode = (m.Config.ServerCfg.SpecialMode + 1) % 3
	case 1:
		m.Config.ServerCfg.QuestSync = !m.Config.ServerCfg.QuestSync
	case 2:
		m.Config.ServerCfg.ShrineSync = !m.Config.ServerCfg.ShrineSync
	case 3:
		m.Config.ServerCfg.TowerSync = !m.Config.ServerCfg.TowerSync
	case 4:
		m.Config.ServerCfg.KorokSync = !m.Config.ServerCfg.KorokSync
	case 5:
		m.Config.ServerCfg.EnemySync = !m.Config.ServerCfg.EnemySync
	case 6:
		m.Config.ServerCfg.DungeonSync = !m.Config.ServerCfg.DungeonSync
	case 7:
		m.Config.ServerCfg.LocationSync = !m.Config.ServerCfg.LocationSync
	}
	_ = m.Config.SaveServerConfig(m.Config.ServerCfg)
}

func (m *Model) saveGamemodeSettings() {
	m.Config.ServerCfg.Password = strings.TrimSpace(m.ServerPassIn.Value())
	m.Config.ServerCfg.Description = strings.TrimSpace(m.ServerDescIn.Value())
	err := m.Config.SaveServerConfig(m.Config.ServerCfg)
	if err != nil {
		m.SetNotification(fmt.Sprintf("Error guardando ServerConfig: %v", err), 3*time.Second)
	} else {
		if m.Process.IsServerRunning() {
			_ = m.Process.StopServer()
			prefix := m.Config.PrefixDir
			go func() {
				time.Sleep(600 * time.Millisecond)
				_ = m.Process.StartServer(prefix)
			}()
			m.SetNotification("Configuracion guardada! Reiniciando servidor con las nuevas reglas...", 3*time.Second)
		} else {
			m.SetNotification("Configuracion guardada! Lista para el inicio del servidor.", 3*time.Second)
		}
	}
	m.CurrentTab = TabDashboard
}

func (m *Model) updateGamemodeFocus() {
	m.ServerPassIn.Blur()
	m.ServerDescIn.Blur()
	if m.GamemodeCursor == 8 {
		m.ServerPassIn.Focus()
	} else if m.GamemodeCursor == 9 {
		m.ServerDescIn.Focus()
	}
}

func (m *Model) updatePathFocus() {
	m.BaseInput.Blur()
	m.UpdateInput.Blur()
	m.DLCInput.Blur()
	switch m.PathFocusIdx {
	case 0:
		m.BaseInput.Focus()
	case 1:
		m.UpdateInput.Focus()
	case 2:
		m.DLCInput.Focus()
	}
}

func (m *Model) blurPathInputs() {
	m.BaseInput.Blur()
	m.UpdateInput.Blur()
	m.DLCInput.Blur()
}

func pickDirCmd(field int, title, initial string) tea.Cmd {
	return func() tea.Msg {
		p, err := config.PickDirectory(title, initial)
		return PickedDirMsg{Field: field, Path: p, Err: err}
	}
}
