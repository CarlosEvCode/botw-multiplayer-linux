package process

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

type ProcessManager struct {
	mu           sync.RWMutex
	serverCmd    *exec.Cmd
	serverStdin  io.WriteCloser
	serverLogs   []string
	maxLogs      int
	milkBarCmd   *exec.Cmd
	cemuCmd      *exec.Cmd
	onLogMessage func(line string)
}

func NewProcessManager(maxLogs int) *ProcessManager {
	if maxLogs <= 0 {
		maxLogs = 500
	}
	return &ProcessManager{
		serverLogs: make([]string, 0, maxLogs),
		maxLogs:    maxLogs,
	}
}

func (p *ProcessManager) SetLogCallback(cb func(line string)) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.onLogMessage = cb
}

func (p *ProcessManager) AddLog(format string, args ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()

	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("15:04:05")
	line := fmt.Sprintf("[%s] %s", timestamp, msg)

	if len(p.serverLogs) >= p.maxLogs {
		p.serverLogs = p.serverLogs[1:]
	}
	p.serverLogs = append(p.serverLogs, line)

	if p.onLogMessage != nil {
		go p.onLogMessage(line)
	}
}

func (p *ProcessManager) GetLogs() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	copied := make([]string, len(p.serverLogs))
	copy(copied, p.serverLogs)
	return copied
}

func (p *ProcessManager) IsServerRunning() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.serverCmd != nil && p.serverCmd.Process != nil && p.serverCmd.ProcessState == nil
}

func (p *ProcessManager) StartServer(prefixDir string) error {
	p.mu.Lock()
	if p.serverCmd != nil && p.serverCmd.Process != nil && p.serverCmd.ProcessState == nil {
		p.mu.Unlock()
		return fmt.Errorf("el servidor ya se encuentra en ejecucion")
	}
	p.mu.Unlock()

	serverExe := filepath.Join(prefixDir, "drive_c/MilkBarLauncher/DedicatedServer/MBL.DedicatedServer.exe")
	workDir := filepath.Dir(serverExe)

	if _, err := os.Stat(serverExe); err != nil {
		return fmt.Errorf("no se encontro el binario del servidor en %s", serverExe)
	}

	cmd := exec.Command("wine", serverExe)
	cmd.Dir = workDir
	cmd.Env = append(os.Environ(),
		"WINEARCH=win64",
		fmt.Sprintf("WINEPREFIX=%s", prefixDir),
		"WINEDEBUG=-all",
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("error creando stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creando stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error creando stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error iniciando servidor: %w", err)
	}

	p.mu.Lock()
	p.serverCmd = cmd
	p.serverStdin = stdin
	p.mu.Unlock()

	p.AddLog("Servidor dedicado iniciado (PID: %d)", cmd.Process.Pid)

	// Stream stdout & stderr
	go p.streamPipe(stdout)
	go p.streamPipe(stderr)

	go func() {
		_ = cmd.Wait()
		p.mu.Lock()
		p.serverCmd = nil
		p.serverStdin = nil
		p.mu.Unlock()
		p.AddLog("Servidor dedicado detenido.")
	}()

	return nil
}

func (p *ProcessManager) streamPipe(r io.Reader) {
	buf := make([]byte, 1024)
	var lineAcc string

	for {
		n, err := r.Read(buf)
		if n > 0 {
			chunk := string(buf[:n])
			lineAcc += chunk
			for strings.Contains(lineAcc, "\n") {
				parts := strings.SplitN(lineAcc, "\n", 2)
				trimmed := strings.TrimRight(parts[0], "\r")
				if trimmed != "" {
					p.AddLog("%s", trimmed)
				}
				lineAcc = parts[1]
			}
			if len(lineAcc) > 0 && (strings.HasSuffix(lineAcc, ": ") || strings.HasSuffix(lineAcc, "? ")) {
				p.AddLog("%s", strings.TrimRight(lineAcc, "\r"))
				lineAcc = ""
			}
		}
		if err != nil {
			if lineAcc != "" {
				p.AddLog("%s", strings.TrimRight(lineAcc, "\r\n"))
			}
			break
		}
	}
}

func (p *ProcessManager) SendServerInput(input string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.serverStdin == nil {
		return fmt.Errorf("el servidor no esta aceptando entradas")
	}

	_, err := io.WriteString(p.serverStdin, input+"\n")
	if err == nil {
		timestamp := time.Now().Format("15:04:05")
		line := fmt.Sprintf("[%s] > %s", timestamp, input)
		if len(p.serverLogs) >= p.maxLogs {
			p.serverLogs = p.serverLogs[1:]
		}
		p.serverLogs = append(p.serverLogs, line)
	}
	return err
}

func (p *ProcessManager) StopServer() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.serverCmd == nil || p.serverCmd.Process == nil {
		return nil
	}

	_ = p.serverCmd.Process.Signal(syscall.SIGTERM)
	go func(proc *os.Process) {
		time.Sleep(2 * time.Second)
		_ = proc.Kill()
	}(p.serverCmd.Process)

	return nil
}

func (p *ProcessManager) IsMilkBarRunning() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.milkBarCmd != nil && p.milkBarCmd.Process != nil && p.milkBarCmd.ProcessState == nil
}

func (p *ProcessManager) StartMilkBar(prefixDir string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.milkBarCmd != nil && p.milkBarCmd.Process != nil && p.milkBarCmd.ProcessState == nil {
		return fmt.Errorf("Milk Bar Launcher ya esta en ejecucion")
	}

	milkExe := filepath.Join(prefixDir, "drive_c/MilkBarLauncher/Milk Bar Launcher.exe")
	if _, err := os.Stat(milkExe); err != nil {
		return fmt.Errorf("no se encontro Milk Bar Launcher en %s", milkExe)
	}

	cmd := exec.Command("wine", milkExe)
	cmd.Dir = filepath.Dir(milkExe)
	cmd.Env = append(os.Environ(),
		"WINEARCH=win64",
		fmt.Sprintf("WINEPREFIX=%s", prefixDir),
		"WINEDEBUG=-all",
		"WINE_FULLSCREEN_FSR=1",
	)

	if err := cmd.Start(); err != nil {
		return err
	}

	p.milkBarCmd = cmd
	go func() {
		_ = cmd.Wait()
		p.mu.Lock()
		p.milkBarCmd = nil
		p.mu.Unlock()
	}()

	return nil
}

func (p *ProcessManager) IsCemuRunning() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.cemuCmd != nil && p.cemuCmd.Process != nil && p.cemuCmd.ProcessState == nil
}

func (p *ProcessManager) StartCemu(prefixDir string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cemuCmd != nil && p.cemuCmd.Process != nil && p.cemuCmd.ProcessState == nil {
		return fmt.Errorf("Cemu ya esta en ejecucion")
	}

	cemuExe := filepath.Join(prefixDir, "drive_c/cemu_1.26.2/Cemu.exe")
	if _, err := os.Stat(cemuExe); err != nil {
		return fmt.Errorf("no se encontro Cemu.exe en %s", cemuExe)
	}

	cmd := exec.Command("wine", cemuExe)
	cmd.Dir = filepath.Dir(cemuExe)
	cmd.Env = append(os.Environ(),
		"WINEARCH=win64",
		fmt.Sprintf("WINEPREFIX=%s", prefixDir),
		"WINEDEBUG=-all",
	)

	if err := cmd.Start(); err != nil {
		return err
	}

	p.cemuCmd = cmd
	go func() {
		_ = cmd.Wait()
		p.mu.Lock()
		p.cemuCmd = nil
		p.mu.Unlock()
	}()

	return nil
}
