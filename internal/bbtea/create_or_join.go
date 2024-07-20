package bbtea

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	createNewGame = "Create a new game"
	joinGame      = "Join an existing game"
)

type createOrJoinModel struct {
	choices       []string
	cursor        int
	Selected      string
	joinTextInput textinput.Model
	err           error
}

func initialCreateOrJoinModel() createOrJoinModel {
	return createOrJoinModel{
		choices: []string{createNewGame, joinGame},
		err:     nil,
	}
}

func (coj createOrJoinModel) Init() tea.Cmd {
	return nil
}

func (coj createOrJoinModel) Update(msg tea.KeyMsg) (createOrJoinModel, tea.Cmd) {
	var cmd tea.Cmd

	switch coj.Selected {
	case "":
		coj, cmd = coj.handleChoices(msg)
	case "Join an existing game":
		coj, cmd = coj.handleJoinInput(msg)
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
		coj.Selected = coj.choices[coj.cursor]
		if coj.Selected == createNewGame {
			cmd = setPhaseCmd(AddPlayerName)
		} else {
			ti := textinput.New()
			ti.Focus()
			ti.CharLimit = 156
			ti.Width = 20
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
		cmd = setPhaseCmd(AddPlayerName)
	} else {
		coj.joinTextInput, cmd = coj.joinTextInput.Update(msg)
	}
	return coj, cmd
}

func (coj createOrJoinModel) View() string {
	var s string
	switch coj.Selected {
	case "":
		s = "Do you want to create or join a game?\n\n"

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
		s += "\nPress esc to quit.\n"
	case joinGame:
		s = fmt.Sprintf(
			"Enter game uuid to join:\n\n %s\n\n%s",
			coj.joinTextInput.View(),
			"(esc to quit)",
		) + "\n"
	}
	return s
}
