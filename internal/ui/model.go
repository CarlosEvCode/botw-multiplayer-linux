package ui

import (
	"time"

	"github.com/CarlosEvCode/botw-multiplayer-linux/internal/config"
	"github.com/CarlosEvCode/botw-multiplayer-linux/internal/i18n"
	"github.com/CarlosEvCode/botw-multiplayer-linux/internal/network"
	"github.com/CarlosEvCode/botw-multiplayer-linux/internal/process"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Tab int

const (
	TabDashboard Tab = iota
	TabClient
	TabGamemodes
	TabPaths
	TabNetwork
	TabSettings
)

type LogMsg string
type TickMsg time.Time

type PickedDirMsg struct {
	Field int
	Path  string
	Err   error
}

type Model struct {
	Config       *config.ManagerConfig
	Process      *process.ProcessManager
	Network      network.NetworkInfo
	CurrentTab   Tab
	Width        int
	Height       int
	Viewport     viewport.Model

	// Client & Connect inputs
	ClientCursor  int
	ClientIPIn    textinput.Model
	ClientPortIn  textinput.Model
	ClientPassIn  textinput.Model
	ClientNameIn  textinput.Model
	ClientModelIn textinput.Model

	// Gamemodes cursor
	GamemodeCursor int
	ServerPassIn   textinput.Model
	ServerDescIn   textinput.Model

	// Paths edit inputs
	BaseInput     textinput.Model
	UpdateInput   textinput.Model
	DLCInput      textinput.Model
	PathFocusIdx  int
	Notification  string
	NotifyTimeout time.Time
}

func NewModel(cfg *config.ManagerConfig, proc *process.ProcessManager) Model {
	i18n.LoadLanguage()

	// Client inputs
	cIP := textinput.New()
	cIP.Placeholder = "127.0.0.1 or friend's ZeroTier/VPN IP"
	cIP.SetValue(cfg.ClientCfg.TargetIP)

	cPort := textinput.New()
	cPort.Placeholder = "5050"
	cPort.SetValue(cfg.ClientCfg.TargetPort)

	cPass := textinput.New()
	cPass.Placeholder = "Empty if server has no password"
	cPass.SetValue(cfg.ClientCfg.Password)

	cName := textinput.New()
	cName.Placeholder = "Link"
	cName.SetValue(cfg.ClientCfg.PlayerName)

	cModel := textinput.New()
	cModel.Placeholder = "Link:Link"
	cModel.SetValue(cfg.ClientCfg.CharacterModel)

	baseIn := textinput.New()
	baseIn.Placeholder = "/path/to/The Legend of Zelda Breath of the Wild"
	baseIn.SetValue(cfg.BaseGame)

	upIn := textinput.New()
	upIn.Placeholder = "/path/to/Update v208 (optional if already in Cemu)"
	upIn.SetValue(cfg.UpdatePath)

	dlcIn := textinput.New()
	dlcIn.Placeholder = "/path/to/DLC v80 (optional if already in Cemu)"
	dlcIn.SetValue(cfg.DLCPath)

	passIn := textinput.New()
	passIn.Placeholder = "Empty if no password required"
	passIn.SetValue(cfg.ServerCfg.Password)

	descIn := textinput.New()
	descIn.Placeholder = "Server description"
	descIn.SetValue(cfg.ServerCfg.Description)

	vp := viewport.New(80, 14)
	vp.SetContent(i18n.T("dash.console_title") + "...")

	m := Model{
		Config:         cfg,
		Process:        proc,
		Network:        network.GetNetworkInfo(),
		CurrentTab:     TabDashboard,
		Viewport:       vp,
		ClientCursor:   0,
		ClientIPIn:     cIP,
		ClientPortIn:   cPort,
		ClientPassIn:   cPass,
		ClientNameIn:   cName,
		ClientModelIn:  cModel,
		BaseInput:      baseIn,
		UpdateInput:    upIn,
		DLCInput:       dlcIn,
		ServerPassIn:   passIn,
		ServerDescIn:   descIn,
		GamemodeCursor: 0,
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		tickEvery(1*time.Second),
	)
}

func tickEvery(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
