package main

import (
	"fmt"
	"strings"
	// "log"
	// "time"

	tea "github.com/charmbracelet/bubbletea"
  	mack "github.com/andybrewer/mack"
)

type model struct {
	albums []string
	playlists []string
	songs []string
	player string
	currentPage	string
	cursor int
	selected map[int]struct{}
}

func initialModel(playlistsStr string, albumsStr string, songsStr string) model {
	return model{
		playlists: strings.Split(playlistsStr, ", "),
		albums: removeDuplicates(strings.Split(albumsStr, ", ")),
		songs: strings.Split(songsStr, ", "),
		player: "Music",
		currentPage: "home",
		cursor: 0,
	}
}

func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate element.
		} else {
			// Append to result slice if it is a unique element.
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}

func (m model) View() string {
	var choices []string
	
	var s string = fmt.Sprintf("\n\nApple Music Player\n\n")
	switch m.currentPage {
	case "home":
		choices = []string{"playlists", "albums", "songs"}
	case "playlists":
		choices = m.playlists
	case "albums":
		choices = m.albums
	case "songs":
		choices = m.songs
	}
    for i, choice := range choices {
        cursor := " " // no cursor
        if m.cursor == i {
            cursor = ">" // cursor!
        }

        s += fmt.Sprintf("%s %s\n", cursor, choice)
    }

    s += "\nPress q to quit.\n"

    // Send the UI for rendering
    return s
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var pageLen int
	switch m.currentPage {
	case "home":
		pageLen = 3
	case "playlists":
		pageLen = len(m.playlists)
	case "albums":
		pageLen = len(m.albums)
	case "songs":
		pageLen = len(m.songs)
	}
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }
        case "down", "j":
            if m.cursor < pageLen-1 {
                m.cursor++
            }
        // The "enter" key and the spacebar (a literal space) toggle
        // the selected state for the item that the cursor is pointing at.
        case "enter", " ":
            _, ok := m.selected[m.cursor]
            if ok {
                delete(m.selected, m.cursor)
            } else {
                m.selected[m.cursor] = struct{}{}
            }
        }
    }

    // Return the updated model to the Bubble Tea runtime for processing.
    // Note that we're not returning a command.
    return m, nil
}

func main() {

	// mack.Tell("Music", "playpause")
	playlistsStr, err := mack.Tell("Music", "set allPlaylists to name of every playlist")
	if err != nil {
		fmt.Println(err)
		return
	}
	albumsStr, err := mack.Tell("Music", "set albumsWithDups to album of every track")
	if err != nil {
		fmt.Println(err)
		return
	}
	songsStr, err := mack.Tell("Music", "set allSongs to name of every track")
	if err != nil {
		fmt.Println(err)
		return
	}

	p := tea.NewProgram(initialModel(playlistsStr, albumsStr, songsStr))
	if err := p.Start(); err != nil {
		fmt.Println("Oh noes!", err)
	}
}

