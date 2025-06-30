// Terminal UI to control OPAQUE client/server app
package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

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
	state       state
	menuIndex   int
	menuOptions []string
	serverCmd   *exec.Cmd
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
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
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
				go runServer()
				m.state = runningServer
			case 1:
				runClient("-pwreg")
				m.state = clientRegister
			case 2:
				runClient("-auth")
				m.state = clientAuth
			case 3:
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "\n🔐 OPAQUE Protocol Terminal UI\n"
	s += "Use ↑/↓ and press Enter to select:\n\n"
	for i, choice := range m.menuOptions {
		cursor := " "
		if m.menuIndex == i {
			cursor = ">"
		}
		s += fmt.Sprintf(" %s %s\n", cursor, choice)
	}
	s += "\nPress q to quit.\n"
	return s
}

func runServer() {
	go func() {
		cmd := exec.Command("go", "run", "cmd/server/main.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			fmt.Println("Error starting server:", err)
			return
		}
		time.Sleep(1 * time.Second)
		fmt.Println("[✔] Server started.")
	}()
}

func runClient(mode string) {
	fmt.Print("Username: ")
	var username string
	fmt.Scanln(&username)
	fmt.Print("Password: ")
	var password string
	fmt.Scanln(&password)

	args := []string{"run", "cmd/client/main.go", mode, "-username", username, "-password", password}
	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Println("[✗] Client error:", err)
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
