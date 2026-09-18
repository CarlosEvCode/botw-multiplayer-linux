package ui

import (
	"fmt"
	"strings"
	"time"

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
		// When typing in text inputs
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
			m.CurrentTab = (m.CurrentTab + 1) % 3
			if m.CurrentTab == TabPaths {
				m.PathFocusIdx = 0
				m.updatePathFocus()
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

		case "i", "/":
			if m.CurrentTab == TabDashboard {
				m.InputFocused = true
				m.CmdInput.Focus()
				return m, textinput.Blink
			}

		case "r":
			m.CurrentTab = TabPaths
			m.PathFocusIdx = 0
			m.updatePathFocus()
			return m, nil
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
