package main

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/concrnt/ccworld-ap-bridge/world"
	"github.com/go-resty/resty/v2"
	"github.com/totegamma/concurrent/core"
)

func newMessage(client *resty.Client, message string) tea.Cmd {
	return func() tea.Msg {

		timelines := strings.Split(os.Getenv("TIMELINES"), ",")

		doc := core.MessageDocument[MarkdownMessage]{
			DocumentBase: core.DocumentBase[MarkdownMessage]{
				Signer: os.Getenv("CCID"),
				Type:   "message",
				Schema: world.MarkdownMessageSchema,
				Body: MarkdownMessage{
					Body: message,
				},
				SignedAt: time.Now(),
				KeyID:    os.Getenv("CKID"),
			},
			Timelines: timelines,
		}

		j, err := json.MarshalIndent(doc, "", "  ")

		if err != nil {
			return newMsg{
				text: err.Error(),
			}
		}

		signed, err := core.SignBytes(j, os.Getenv("PRIVATE_KEY"))

		hex := hex.EncodeToString(signed)

		request := core.Commit{
			Document:  string(j),
			Signature: string(hex),
		}

		_, err = client.R().
			SetBody(request).
			Post("/commit")

		if err != nil {
			return newMsg{
				text: err.Error(),
			}
		}

		return newMsg{
			text: message,
		}
	}
}

func fetchPost(client *resty.Client) tea.Cmd {
	return func() tea.Msg {
		timeline := strings.Split(os.Getenv("TIMELINES"), ",")[0]

		resp, err := client.R().Get("/timelines/recent?timelines=" + timeline)

		if err != nil {
			return fetchMsg{
				text: err.Error(),
			}
		}

		var (
			j core.ResponseBase[[]core.TimelineItem]
		)

		err = json.Unmarshal(resp.Body(), &j)

		if err != nil {
			return fetchMsg{
				text: err.Error(),
			}
		}

		var users map[string]string

		var listToReturn []list.Item

		for _, v := range j.Content {
			// fetch each post
			resp, err := client.R().Get("/message/" + v.ResourceID)

			if err != nil {
				return fetchMsg{
					text: err.Error(),
				}
			}

			var message core.ResponseBase[core.Message]

			err = json.Unmarshal(resp.Body(), &message)

			if err != nil {
				return fetchMsg{
					text: err.Error(),
				}
			}

			if message.Content.Schema != world.MarkdownMessageSchema {
				continue
			}

			var document core.MessageDocument[world.MarkdownMessage]

			err = json.Unmarshal([]byte(message.Content.Document), &document)

			if err != nil {
				return fetchMsg{
					text: err.Error(),
				}
			}

			if users[document.Signer] == "" && document.Signer != "" {
				resp, err := client.R().Get("/profile/" + document.Signer + "/world.concrnt.p")

				if err != nil {
					log.Fatal(err)
				}

				var user core.ResponseBase[core.Profile]

				err = json.Unmarshal(resp.Body(), &user)

				if err != nil {
					return fetchMsg{
						text: err.Error(),
					}
				}

				if users == nil {
					users = make(map[string]string)
				}

				var u core.ProfileDocument[Body]

				err = json.Unmarshal([]byte(user.Content.Document), &u)

				if err != nil {
					return fetchMsg{
						text: err.Error(),
					}
				}

				users[document.Signer] = u.Body.Username
			}

			listToReturn = append(listToReturn, item{
				title: users[document.Signer],
				desc:  document.Body.Body,
			})
		}

		return fetchMsg{
			list: listToReturn,
		}
	}
}
