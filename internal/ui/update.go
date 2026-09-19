package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/CarlosEvCode/botw-multiplayer-linux/internal/config"
	"github.com/CarlosEvCode/botw-multiplayer-linux/internal/i18n"
	"github.com/CarlosEvCode/botw-multiplayer-linux/internal/network"
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
		m.ClientIPIn.Width = 35
		m.ClientPortIn.Width = 15
		m.ClientPassIn.Width = 35
		m.ClientNameIn.Width = 35
		m.ClientModelIn.Width = 35
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
				m.SetNotification(fmt.Sprintf("%s: %v", i18n.T("paths.val_not_set"), err), 3*time.Second)
			} else {
				m.SetNotification(i18n.T("notif.paths_saved"), 3*time.Second)
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
		// Global Tab navigation across all 6 tabs
		switch msg.String() {
		case "tab":
			m.CurrentTab = (m.CurrentTab + 1) % 6
			m.handleTabSwitch()
			return m, nil

		case "shift+tab":
			m.CurrentTab = (m.CurrentTab + 5) % 6
			m.handleTabSwitch()
			return m, nil

		case "1":
			m.CurrentTab = TabDashboard
			m.handleTabSwitch()
			return m, nil
		case "2":
			m.CurrentTab = TabClient
			m.handleTabSwitch()
			return m, nil
		case "3":
			m.CurrentTab = TabGamemodes
			m.handleTabSwitch()
			return m, nil
		case "4":
			m.CurrentTab = TabPaths
			m.handleTabSwitch()
			return m, nil
		case "5":
			m.CurrentTab = TabNetwork
			m.handleTabSwitch()
			return m, nil
		case "6":
			m.CurrentTab = TabSettings
			m.handleTabSwitch()
			return m, nil
		}

		// When in Client & Connect Tab
		if m.CurrentTab == TabClient {
			switch msg.String() {
			case "esc":
				m.CurrentTab = TabDashboard
				m.blurClientInputs()
				return m, nil
			case "down":
				m.ClientCursor = (m.ClientCursor + 1) % 5
				m.updateClientFocus()
				return m, nil
			case "up":
				m.ClientCursor = (m.ClientCursor + 4) % 5
				m.updateClientFocus()
				return m, nil
			case "l":
				if m.ClientCursor == 0 {
					m.ClientIPIn.SetValue("127.0.0.1")
					m.SetNotification(i18n.T("notif.client_set_local"), 2*time.Second)
					return m, nil
				}
			case "ctrl+s", "s":
				m.saveClientSettings()
				return m, nil
			case "enter", "c":
				m.launchClient()
				return m, nil
			case "x":
				if m.Process.IsClientRunning() {
					_ = m.Process.StopClient()
					m.SetNotification(i18n.T("notif.client_stopping"), 2*time.Second)
				}
				return m, nil
			}

			switch m.ClientCursor {
			case 0:
				m.ClientIPIn, cmd = m.ClientIPIn.Update(msg)
			case 1:
				m.ClientPortIn, cmd = m.ClientPortIn.Update(msg)
			case 2:
				m.ClientPassIn, cmd = m.ClientPassIn.Update(msg)
			case 3:
				m.ClientNameIn, cmd = m.ClientNameIn.Update(msg)
			case 4:
				m.ClientModelIn, cmd = m.ClientModelIn.Update(msg)
			}
			return m, cmd
		}

		// When in Gamemodes Tab
		if m.CurrentTab == TabGamemodes {
			if m.GamemodeCursor == 8 { // Password input
				switch msg.String() {
				case "esc":
					m.ServerPassIn.Blur()
					m.CurrentTab = TabDashboard
					return m, nil
				case "down", "enter":
					m.ServerPassIn.Blur()
					m.GamemodeCursor = 9
					m.updateGamemodeFocus()
					return m, nil
				case "up":
					m.ServerPassIn.Blur()
					m.GamemodeCursor = 7
					m.updateGamemodeFocus()
					return m, nil
				case "ctrl+s":
					m.saveGamemodeSettings()
					return m, nil
				}
				m.ServerPassIn, cmd = m.ServerPassIn.Update(msg)
				return m, cmd
			} else if m.GamemodeCursor == 9 { // Description input
				switch msg.String() {
				case "esc":
					m.ServerDescIn.Blur()
					m.CurrentTab = TabDashboard
					return m, nil
				case "down", "enter":
					m.ServerDescIn.Blur()
					m.GamemodeCursor = 0
					m.updateGamemodeFocus()
					return m, nil
				case "up":
					m.ServerDescIn.Blur()
					m.GamemodeCursor = 8
					m.updateGamemodeFocus()
					return m, nil
				case "ctrl+s":
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
			case "right", "l":
				if m.GamemodeCursor == 0 {
					m.changeSpecialMode(1)
				} else if m.GamemodeCursor >= 1 && m.GamemodeCursor <= 7 {
					m.toggleSyncOption(m.GamemodeCursor)
				}
				return m, nil
			case "left", "h":
				if m.GamemodeCursor == 0 {
					m.changeSpecialMode(-1)
				} else if m.GamemodeCursor >= 1 && m.GamemodeCursor <= 7 {
					m.toggleSyncOption(m.GamemodeCursor)
				}
				return m, nil
			case " ", "enter":
				if m.GamemodeCursor == 0 {
					m.changeSpecialMode(1)
				} else if m.GamemodeCursor >= 1 && m.GamemodeCursor <= 7 {
					m.toggleSyncOption(m.GamemodeCursor)
				}
				return m, nil
			case "ctrl+s", "s":
				m.saveGamemodeSettings()
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
			case "down":
				m.PathFocusIdx = (m.PathFocusIdx + 1) % 3
				m.updatePathFocus()
				return m, nil
			case "up":
				m.PathFocusIdx = (m.PathFocusIdx + 2) % 3
				m.updatePathFocus()
				return m, nil
			case "f", "o", "ctrl+o", "f2", "alt+o":
				var title, initDir string
				switch m.PathFocusIdx {
				case 0:
					title = i18n.T("paths.base")
					initDir = m.BaseInput.Value()
				case 1:
					title = i18n.T("paths.update")
					initDir = m.UpdateInput.Value()
				case 2:
					title = i18n.T("paths.dlc")
					initDir = m.DLCInput.Value()
				}
				m.SetNotification(i18n.T("notif.picker_opening"), 2*time.Second)
				return m, pickDirCmd(m.PathFocusIdx, title, initDir)
			case "enter", "ctrl+s":
				err := m.Config.SavePaths(m.BaseInput.Value(), m.UpdateInput.Value(), m.DLCInput.Value())
				if err != nil {
					m.SetNotification(fmt.Sprintf("Error: %v", err), 3*time.Second)
				} else {
					m.SetNotification(i18n.T("notif.paths_saved"), 3*time.Second)
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

		// When in Settings Tab
		if m.CurrentTab == TabSettings {
			switch msg.String() {
			case "esc":
				m.CurrentTab = TabDashboard
				return m, nil
			case "left", "right", "h", "l", " ", "enter":
				if i18n.CurrentLang == i18n.LangEN {
					i18n.SetLanguage(i18n.LangES)
				} else {
					i18n.SetLanguage(i18n.LangEN)
				}
				m.SetNotification(i18n.T("notif.lang_changed"), 2*time.Second)
				return m, nil
			}
		}

		// Global Keybindings (Dashboard & Top-level)
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "s":
			if m.Process.IsServerRunning() {
				_ = m.Process.StopServer()
				m.SetNotification(i18n.T("notif.srv_stopping"), 2*time.Second)
			} else {
				err := m.Process.StartServer(m.Config.PrefixDir)
				if err != nil {
					m.SetNotification(fmt.Sprintf(i18n.T("notif.srv_error"), err), 3*time.Second)
				} else {
					m.SetNotification(fmt.Sprintf(i18n.T("notif.srv_started"), m.Config.ServerCfg.Port), 2*time.Second)
				}
			}
			return m, nil

		case "c":
			m.launchClient()
			return m, nil

		case "x":
			if m.Process.IsClientRunning() {
				_ = m.Process.StopClient()
				m.SetNotification(i18n.T("notif.client_stopping"), 2*time.Second)
			}
			return m, nil

		case "e":
			err := m.Process.StartCemu(m.Config.PrefixDir)
			if err != nil {
				m.SetNotification(fmt.Sprintf("%s: %v", i18n.T("dash.srv_cemu"), err), 3*time.Second)
			} else {
				m.SetNotification(i18n.T("notif.cemu_starting"), 2*time.Second)
			}
			return m, nil

		case "t":
			if m.Network.TailscaleIP != "" {
				_ = network.CopyToClipboard(m.Network.TailscaleIP)
				m.SetNotification(fmt.Sprintf(i18n.T("notif.ip_copied"), m.Network.TailscaleIP), 2*time.Second)
			} else if m.Network.ZeroTierIP != "" {
				_ = network.CopyToClipboard(m.Network.ZeroTierIP)
				m.SetNotification(fmt.Sprintf(i18n.T("notif.zt_ip_copied"), m.Network.ZeroTierIP), 2*time.Second)
			} else if m.Network.LocalIP != "" {
				_ = network.CopyToClipboard(m.Network.LocalIP)
				m.SetNotification(fmt.Sprintf(i18n.T("notif.local_ip_copied"), m.Network.LocalIP), 2*time.Second)
			}
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

func (m *Model) changeSpecialMode(delta int) {
	newMode := (m.Config.ServerCfg.SpecialMode + delta) % 3
	if newMode < 0 {
		newMode += 3
	}
	m.Config.ServerCfg.SpecialMode = newMode
	switch m.Config.ServerCfg.SpecialMode {
	case 0: // Cooperativo Estándar
		if !m.Config.ServerCfg.QuestSync && !m.Config.ServerCfg.ShrineSync {
			m.Config.ServerCfg.QuestSync = true
			m.Config.ServerCfg.ShrineSync = true
			m.Config.ServerCfg.TowerSync = true
			m.Config.ServerCfg.KorokSync = true
			m.Config.ServerCfg.EnemySync = true
			m.Config.ServerCfg.DungeonSync = true
			m.Config.ServerCfg.LocationSync = true
		}
		m.SetNotification(i18n.T("notif.mode_coop"), 2*time.Second)
	case 1: // Hunter vs Speedrunner
		m.Config.ServerCfg.QuestSync = false
		m.Config.ServerCfg.ShrineSync = false
		m.Config.ServerCfg.TowerSync = false
		m.Config.ServerCfg.KorokSync = false
		m.Config.ServerCfg.DungeonSync = false
		m.Config.ServerCfg.LocationSync = false
		m.Config.ServerCfg.EnemySync = true
		m.SetNotification(i18n.T("notif.mode_hunter"), 2*time.Second)
	case 2: // DeathSwap
		m.Config.ServerCfg.QuestSync = false
		m.Config.ServerCfg.ShrineSync = false
		m.Config.ServerCfg.TowerSync = false
		m.Config.ServerCfg.KorokSync = false
		m.Config.ServerCfg.DungeonSync = false
		m.Config.ServerCfg.LocationSync = false
		m.Config.ServerCfg.EnemySync = false
		m.SetNotification(i18n.T("notif.mode_deathswap"), 2*time.Second)
	}
	_ = m.Config.SaveServerConfig(m.Config.ServerCfg)
}

func (m *Model) toggleSyncOption(cursor int) {
	switch cursor {
	case 1:
		if m.Config.ServerCfg.SpecialMode != 0 {
			m.SetNotification(i18n.T("notif.sync_locked"), 2*time.Second)
			return
		}
		m.Config.ServerCfg.QuestSync = !m.Config.ServerCfg.QuestSync
	case 2:
		if m.Config.ServerCfg.SpecialMode != 0 {
			m.SetNotification(i18n.T("notif.sync_locked"), 2*time.Second)
			return
		}
		m.Config.ServerCfg.ShrineSync = !m.Config.ServerCfg.ShrineSync
	case 3:
		if m.Config.ServerCfg.SpecialMode != 0 {
			m.SetNotification(i18n.T("notif.sync_locked"), 2*time.Second)
			return
		}
		m.Config.ServerCfg.TowerSync = !m.Config.ServerCfg.TowerSync
	case 4:
		if m.Config.ServerCfg.SpecialMode != 0 {
			m.SetNotification(i18n.T("notif.sync_locked"), 2*time.Second)
			return
		}
		m.Config.ServerCfg.KorokSync = !m.Config.ServerCfg.KorokSync
	case 5:
		if m.Config.ServerCfg.SpecialMode == 2 {
			m.SetNotification(i18n.T("notif.sync_locked"), 2*time.Second)
			return
		}
		m.Config.ServerCfg.EnemySync = !m.Config.ServerCfg.EnemySync
	case 6:
		if m.Config.ServerCfg.SpecialMode != 0 {
			m.SetNotification(i18n.T("notif.sync_locked"), 2*time.Second)
			return
		}
		m.Config.ServerCfg.DungeonSync = !m.Config.ServerCfg.DungeonSync
	case 7:
		if m.Config.ServerCfg.SpecialMode != 0 {
			m.SetNotification(i18n.T("notif.sync_locked"), 2*time.Second)
			return
		}
		m.Config.ServerCfg.LocationSync = !m.Config.ServerCfg.LocationSync
	}
	_ = m.Config.SaveServerConfig(m.Config.ServerCfg)
}

func (m *Model) saveGamemodeSettings() {
	m.Config.ServerCfg.Password = strings.TrimSpace(m.ServerPassIn.Value())
	m.Config.ServerCfg.Description = strings.TrimSpace(m.ServerDescIn.Value())
	err := m.Config.SaveServerConfig(m.Config.ServerCfg)
	if err != nil {
		m.SetNotification(fmt.Sprintf("Error: %v", err), 3*time.Second)
	} else {
		if m.Process.IsServerRunning() {
			_ = m.Process.StopServer()
			prefix := m.Config.PrefixDir
			go func() {
				time.Sleep(600 * time.Millisecond)
				_ = m.Process.StartServer(prefix)
			}()
			m.SetNotification(i18n.T("notif.cfg_saved_reload"), 3*time.Second)
		} else {
			m.SetNotification(i18n.T("notif.cfg_saved"), 3*time.Second)
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

func (m *Model) handleTabSwitch() {
	m.blurClientInputs()
	m.ServerPassIn.Blur()
	m.ServerDescIn.Blur()
	m.blurPathInputs()

	switch m.CurrentTab {
	case TabClient:
		m.updateClientFocus()
	case TabGamemodes:
		m.updateGamemodeFocus()
	case TabPaths:
		m.updatePathFocus()
	}
}

func (m *Model) saveClientSettings() {
	m.Config.ClientCfg.TargetIP = strings.TrimSpace(m.ClientIPIn.Value())
	m.Config.ClientCfg.TargetPort = strings.TrimSpace(m.ClientPortIn.Value())
	m.Config.ClientCfg.Password = strings.TrimSpace(m.ClientPassIn.Value())
	m.Config.ClientCfg.PlayerName = strings.TrimSpace(m.ClientNameIn.Value())
	m.Config.ClientCfg.CharacterModel = strings.TrimSpace(m.ClientModelIn.Value())

	err := m.Config.SaveClientConfig(m.Config.ClientCfg)
	if err != nil {
		m.SetNotification(fmt.Sprintf("Error: %v", err), 3*time.Second)
	} else {
		m.SetNotification(i18n.T("notif.client_saved"), 2*time.Second)
	}
}

func (m *Model) launchClient() {
	m.saveClientSettings()
	ip := m.Config.ClientCfg.TargetIP
	port := m.Config.ClientCfg.TargetPort
	pass := m.Config.ClientCfg.Password
	name := m.Config.ClientCfg.PlayerName
	model := m.Config.ClientCfg.CharacterModel

	err := m.Process.StartClient(m.Config.PrefixDir, ip, port, pass, name, model)
	if err != nil {
		m.SetNotification(fmt.Sprintf("Error: %v", err), 3*time.Second)
	} else {
		m.SetNotification(fmt.Sprintf(i18n.T("notif.client_starting"), ip, port), 3*time.Second)
	}
	m.CurrentTab = TabDashboard
	m.handleTabSwitch()
}

func (m *Model) updateClientFocus() {
	m.blurClientInputs()
	switch m.ClientCursor {
	case 0:
		m.ClientIPIn.Focus()
	case 1:
		m.ClientPortIn.Focus()
	case 2:
		m.ClientPassIn.Focus()
	case 3:
		m.ClientNameIn.Focus()
	case 4:
		m.ClientModelIn.Focus()
	}
}

func (m *Model) blurClientInputs() {
	m.ClientIPIn.Blur()
	m.ClientPortIn.Blur()
	m.ClientPassIn.Blur()
	m.ClientNameIn.Blur()
	m.ClientModelIn.Blur()
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
