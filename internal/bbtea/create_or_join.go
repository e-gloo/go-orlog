package bbtea

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	c "github.com/e-gloo/orlog/internal/client"
)

const (
	createNewGame    = "Create a new game"
	joinExistingGame = "Join an existing game"
)

type createOrJoinModel struct {
	client        c.Client
	choices       []string
	cursor        int
	selected      string
	joinTextInput textinput.Model
	validated     bool
	err           error
}

func initialCreateOrJoinModel(client c.Client) createOrJoinModel {
	return createOrJoinModel{
		client:    client,
		choices:   []string{createNewGame, joinExistingGame},
		validated: false,
	}
}

func (coj createOrJoinModel) Init() tea.Cmd {
	return nil
}

func (coj createOrJoinModel) Update(msg tea.Msg) (createOrJoinModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch coj.selected {
		case "":
			coj, cmd = coj.handleChoices(msg)
		case joinExistingGame:
			coj, cmd = coj.handleJoinInput(msg)
		}
	}

	return coj, cmd
}

func (coj createOrJoinModel) handleChoices(msg tea.KeyMsg) (createOrJoinModel, tea.Cmd) {
	var cmd tea.Cmd
	// Cool, what was the actual key pressed?
	switch msg.Type {

	// The "up" and "k" keys move the cursor up
	case tea.KeyUp:
		if coj.cursor > 0 {
			coj.cursor--
		}

	// The "down" and "j" keys move the cursor down
	case tea.KeyDown:
		if coj.cursor < len(coj.choices)-1 {
			coj.cursor++
		}

	// The "enter" key and the spacebar (a literal space) toggle
	// the selected state for the item that the cursor is pointing at.
	case tea.KeyEnter, tea.KeySpace:
		coj.selected = coj.choices[coj.cursor]
		if coj.selected == createNewGame {
			coj.err = coj.client.CreateGame()
		} else {
			ti := textinput.New()
			ti.Focus()
			ti.CharLimit = 156
			ti.Width = 256
			coj.joinTextInput = ti
			cmd = textinput.Blink
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return coj, cmd
}

func (coj createOrJoinModel) handleJoinInput(msg tea.KeyMsg) (createOrJoinModel, tea.Cmd) {
	var cmd tea.Cmd

	if msg.Type == tea.KeyEnter {
		coj.err = coj.client.JoinGame(coj.joinTextInput.Value())
		if coj.err == nil {
			coj.validated = true
		}
	} else {
		coj.joinTextInput, cmd = coj.joinTextInput.Update(msg)
	}
	return coj, cmd
}

func (coj createOrJoinModel) View() string {
	var s string

	if coj.err != nil {
		s = coj.err.Error() + "\n"
	}

	if ph.Phase() == c.CreateOrJoin && coj.client.Error() != "" {
		s += coj.client.Error() + "\n"
	}

	switch coj.selected {
	case "":
		s += "Do you want to create or join a game?\n\n"

		// Iterate over our choices
		for i, choice := range coj.choices {

			// Is the cursor pointing at this choice?
			cursor := " " // no cursor
			if coj.cursor == i {
				cursor = ">" // cursor!
			}

			// Render the row
			s += fmt.Sprintf("\t%s %s\n", cursor, choice)
		}

		// The footer
	case createNewGame:
		if coj.client.GameUuid() != "" {
			s += fmt.Sprintf("Game created with uuid %s\n", coj.client.GameUuid())
		} else {
			s += "Creating game...\n"
		}
	case joinExistingGame:
		if coj.validated && coj.client.Error() == "" {
			if coj.client.GameUuid() != "" {
				s += fmt.Sprintf("Game with uuid %s joined\n", coj.client.GameUuid())
			} else {
				s += fmt.Sprintf("Joining game with uuid %s...\n", coj.joinTextInput.Value())
			}
		} else {
			s = fmt.Sprintf(
				"Enter game uuid to join:\n\n %s\n",
				coj.joinTextInput.View(),
			)
		}
	}
	return s
}
