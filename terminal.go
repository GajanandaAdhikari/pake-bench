// Terminal UI to control OPAQUE client/server app
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
	"net"

	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	startMenu state = iota
	runningServer
	clientRegister
	clientAuth
	exitApp
)

type model struct {
	state         state
	menuIndex     int
	menuOptions   []string
	serverRunning bool
	keyPreview    string
	serverCmd     *exec.Cmd
}

func initialModel() model {
	return model{
		state:     startMenu,
		menuIndex: 0,
		menuOptions: []string{
			"Start OPAQUE Server",
			"Register New User",
			"Authenticate User",
			"Exit",
		},
		serverRunning: false,
		keyPreview:    "",
	}
}

func (m model) Init() tea.Cmd {
	return checkServerStatus()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.serverCmd != nil {
				_ = m.serverCmd.Process.Kill()
			}
			return m, tea.Quit
		case "up", "k":
			if m.menuIndex > 0 {
				m.menuIndex--
			}
		case "down", "j":
			if m.menuIndex < len(m.menuOptions)-1 {
				m.menuIndex++
			}
		case "enter":
			switch m.menuIndex {
			case 0:
				cmd := runServer()
				m.serverCmd = cmd
				m.serverRunning = true
			case 1:
				m.keyPreview = runClient("-pwreg")
			case 2:
				m.keyPreview = runClient("-auth")
			case 3:
				if m.serverCmd != nil {
					_ = m.serverCmd.Process.Kill()
				}
				return m, tea.Quit
			}
		}
	case serverStatusMsg:
		m.serverRunning = bool(msg)
	}
	return m, checkServerStatus()
}

func (m model) View() string {
	s := "\n🔐 OPAQUE Protocol Terminal UI\n"
	if m.serverRunning {
		s += "[✔] Server running on :9999\n"
	} else {
		s += "[✗] Server not running\n"
	}
	s += "Use ↑/↓ and press Enter to select:\n\n"
	for i, choice := range m.menuOptions {
		cursor := " "
		if m.menuIndex == i {
			cursor = ">"
		}
		s += fmt.Sprintf(" %s %s\n", cursor, choice)
	}
	s += "\nPress q to quit.\n"
	if m.keyPreview != "" {
		s += "\n🔑 Session Key Preview: " + m.keyPreview + "\n"
	}
	return s
}

type serverStatusMsg bool

func checkServerStatus() tea.Cmd {
	return func() tea.Msg {
		conn, err := net.DialTimeout("tcp", "localhost:9999", time.Millisecond*300)
		if err != nil {
			return serverStatusMsg(false)
		}
		conn.Close()
		return serverStatusMsg(true)
	}
}

func runServer() *exec.Cmd {
	cmd := exec.Command("go", "run", "cmd/server/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting server:", err)
		return nil
	}
	time.Sleep(1 * time.Second)
	fmt.Println("[✔] Server started in background.")
	return cmd
}

func runClient(mode string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')

	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	args := []string{"run", "cmd/client/main.go", mode, "-username", username, "-password", password}
	cmd := exec.Command("go", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("[✗] Client error:", err)
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Received") || strings.Contains(line, "Shared") {
			return line
		}
	}
	return ""
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

