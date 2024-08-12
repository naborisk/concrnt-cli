package main

import (
	"log"
	"os"
	// "strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

func InitModel() model {
	ti := textinput.New()
	ti.Placeholder = "What's on your mind!?"

	items := []list.Item{
		item{title: "Loading...", desc: "🦊"},
	}

	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.SetShowHelp(false)
	list.Title = "#arrival-lounge"
	list.SetShowStatusBar(false)
	list.SetShowFilter(false)
	list.SetShowTitle(true)

	client := resty.New()
	client.SetBaseURL(os.Getenv("BASE_URL"))

	return model{
		client:    client,
		help:      help.New(),
		keys:      keys,
		textinput: ti,
		list:      list,
		spinner:   spinner.New(),
	}
}

func main() {
	godotenv.Load()
	p := tea.NewProgram(InitModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		fetchPost(m.client),
		textinput.Blink,
	)
}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := message.(type) {

	case tea.WindowSizeMsg:
		m.windowSize = WindowSize{msg.Width, msg.Height}

		var width, height int

		width = m.windowSize.width

		if m.help.ShowAll {
			// 3 is the height of draft 4 is the height of help 1 is bottom padding
			height = m.windowSize.height - 3 - 4 - 1
		} else {
			height = m.windowSize.height - 3 - 1 - 1
		}

		m.list.SetSize(width, height)

	case tea.KeyMsg:
		if m.textinput.Focused() {
			switch {
			case key.Matches(msg, m.keys.Unfocus):
				m.textinput.Blur()
			case key.Matches(msg, m.keys.Enter):
				m.textinput.Blur()
				return m, newMessage(m.client, m.textinput.Value())
			}
		} else {
			switch {
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.keys.Refresh):
				m.text = "refreshing..."
				m.list.SetItems([]list.Item{item{title: "Loading...", desc: "🦊"}})
				return m, fetchPost(m.client)
			case key.Matches(msg, m.keys.NewMessage):
				m.help.ShortHelpView([]key.Binding{m.keys.Enter, m.keys.Quit})
				return m, m.textinput.Focus()
			case key.Matches(msg, m.keys.Up):
				m.list.CursorUp()
			case key.Matches(msg, m.keys.Down):
				m.list.CursorDown()
			case key.Matches(msg, m.keys.Help):
				m.help.ShowAll = !m.help.ShowAll

				if m.help.ShowAll {
					// 3 is the height of draft 4 is the height of help 1 is bottom padding
					m.list.SetSize(m.windowSize.width, m.windowSize.height-3-4-1)
				} else {
					m.list.SetSize(m.windowSize.height, m.windowSize.height-3-1-1)
				}
			}
		}

	case fetchMsg:
		return m, tea.Batch(
			m.list.SetItems(msg.list),
			fetchPost(m.client),
		)

	case newMsg:
		m.textinput.SetValue("")
		m.text = msg.text
	}

	m.textinput, cmd = m.textinput.Update(message)
	return m, cmd
}

func (m model) View() string {
	helpView := m.help.View(m.keys)

	return " " +
		// strconv.Itoa(m.windowSize.width) + "x" + strconv.Itoa(m.windowSize.height) +
		"\n" +
		" " + m.textinput.View() +
		"\n\n" + m.list.View() +
		// fmt.Sprintf("\n\n%v", m.text) +
		// strings.Repeat("\n", height) +
		// "\n" +
		// m.spinner.View() +
		"\n" + helpView
}
