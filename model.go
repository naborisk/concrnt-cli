package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/go-resty/resty/v2"
	"github.com/totegamma/concurrent/core"
)

type model struct {
	client   *resty.Client
	text     string
	timeline *core.Timeline
	messages []core.Message

	keys       keyMap
	windowSize WindowSize

	// components
	help      help.Model
	textinput textinput.Model
	spinner   spinner.Model
	list      list.Model
}

type User struct {
	Username    string
	Description string
}

type Body struct {
	Username string `json:"username"`
}

type WindowSize struct {
	width  int
	height int
}

type fetchMsg struct {
	text string
	list []list.Item
}

type newMsg struct {
	text string
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// tmp
type ProfileOverride struct {
	Username    string `json:"username,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	Description string `json:"description,omitempty"`
	Link        string `json:"link,omitempty"`
	CharacterID string `json:"characterID,omitempty"`
}

type MarkdownMessage struct {
	Body            string           `json:"body"`
	Emojis          map[string]Emoji `json:"emojis,omitempty"`
	ProfileOverride ProfileOverride  `json:"profileOverride,omitempty"`
}

type Emoji struct {
	ImageURL string `json:"imageURL"`
}
