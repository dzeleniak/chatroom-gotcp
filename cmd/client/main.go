package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dzeleniak/chatroom-gotcp/pkg/user"
	"github.com/gofrs/uuid"
)

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	userId, _ := uuid.NewV7();

	u := &user.User{
		Conn: conn,
		Username: os.Args[1],
		Id: userId,
	}

	u.Conn.Write([]byte(u.Username+string('\n')))

	p := tea.NewProgram(initialModel(u))

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type errMsg error;

type model struct {
	viewport viewport.Model
	messages []string
	textarea textarea.Model
	senderStyle lipgloss.Style
	err error
	user user.User
}

func initialModel(user *user.User) model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus();

	ta.Prompt = "| "
	ta.CharLimit = 200

	ta.SetWidth(30)
	ta.SetHeight(3);
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false;

	vp := viewport.New(100, 5)
	vp.SetContent(`Welcome! Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea: ta,
		viewport: vp,
		messages: []string{},
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err: nil,
		user: *user,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg) 

	switch t := msg.(type) {
	case tea.KeyMsg: 
		switch t.Type {
			case tea.KeyCtrlC, tea.KeyEsc: 
				fmt.Println(m.textarea.Value())
				return m, tea.Quit
			case tea.KeyEnter:
				m.user.Conn.Write([]byte(m.textarea.Value() + string('\n')))
				m.messages = append(m.messages, m.senderStyle.Render("You: ") + m.textarea.Value())
				m.viewport.SetContent(strings.Join(m.messages, "\n"))
				m.textarea.Reset()
				m.viewport.GotoBottom()
		}
	case errMsg:
		m.err = t
		return m, nil
	}	

	return m, tea.Batch(tiCmd, vpCmd);
}