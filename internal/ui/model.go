package ui

import (
	"time"

	"github.com/CarlosEvCode/botw-multiplayer-linux/internal/config"
	"github.com/CarlosEvCode/botw-multiplayer-linux/internal/network"
	"github.com/CarlosEvCode/botw-multiplayer-linux/internal/process"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Tab int

const (
	TabDashboard Tab = iota
	TabPaths
	TabNetwork
)

type LogMsg string
type TickMsg time.Time

type Model struct {
	Config       *config.ManagerConfig
	Process      *process.ProcessManager
	Network      network.NetworkInfo
	CurrentTab   Tab
	Width        int
	Height       int
	Viewport     viewport.Model
	CmdInput     textinput.Model
	InputFocused bool

	// Paths edit inputs
	BaseInput     textinput.Model
	UpdateInput   textinput.Model
	DLCInput      textinput.Model
	PathFocusIdx  int
	Notification  string
	NotifyTimeout time.Time
}

func NewModel(cfg *config.ManagerConfig, proc *process.ProcessManager) Model {
	cmdIn := textinput.New()
	cmdIn.Placeholder = "Enviar comando al servidor (ej: 0, 1, Prop Hunt start)..."
	cmdIn.CharLimit = 120

	baseIn := textinput.New()
	baseIn.Placeholder = "/ruta/a/The Legend of Zelda Breath of the Wild"
	baseIn.SetValue(cfg.BaseGame)

	upIn := textinput.New()
	upIn.Placeholder = "/ruta/a/Update v208 (opcional si ya esta en Cemu)"
	upIn.SetValue(cfg.UpdatePath)

	dlcIn := textinput.New()
	dlcIn.Placeholder = "/ruta/a/DLC v80 (opcional si ya esta en Cemu)"
	dlcIn.SetValue(cfg.DLCPath)

	vp := viewport.New(80, 14)
	vp.SetContent("Iniciando consola de logs...")

	m := Model{
		Config:      cfg,
		Process:     proc,
		Network:     network.GetNetworkInfo(),
		CurrentTab:  TabDashboard,
		Viewport:    vp,
		CmdInput:    cmdIn,
		BaseInput:   baseIn,
		UpdateInput: upIn,
		DLCInput:    dlcIn,
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
