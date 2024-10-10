package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	currentLines   []string
	targetText     string
	userInput      string
	cursorPos      int
	showPrompt     bool
	prompt         string
	countdownTimer timer.Model
	timerStarted   bool
}

func NewModel() Model {
	randomSentences := generateSentences(generateWords())
	return Model{
		targetText:     randomSentences[0],
		userInput:      "",
		cursorPos:      0,
		currentLines:   randomSentences,
		showPrompt:     false,
		prompt:         "",
		countdownTimer: timer.NewWithInterval(5*time.Second, time.Second),
		timerStarted:   false,
	}
}

func generateWords() []string {
	wordCorpus := []string{
		"the", "of", "and", "to", "a", "in", "is", "it", "you", "that", "he", "was", "for", "on", "are", "with",
		"as", "I", "his", "they", "be", "at", "one", "have", "this", "from", "or", "had", "by", "not", "word",
		"but", "what", "some", "we", "can", "out", "other", "were", "all", "there", "when", "up", "use", "your",
		"how", "said", "an", "each", "she", "which", "do", "their", "time", "if", "will", "way", "about", "many",
		"then", "them", "write", "would", "like", "so", "these", "her", "long", "make", "thing", "see", "him",
		"two", "has", "look", "more", "day", "could", "go", "come", "did", "number", "sound", "no", "most",
		"people", "my", "over", "know", "water", "than", "call", "first", "who", "may", "down", "side", "been",
		"now", "find", "any", "new", "work", "part", "take", "get", "place", "made", "live", "where", "after",
		"back", "little", "only", "round", "man",
	}
	return wordCorpus
}

func generateSentences(wordCorpus []string) []string {
	rand.Seed(time.Now().UnixNano())

	sentences := []string{}
	lastUsed := map[string]int{}
	totalWords := len(wordCorpus)
	numSentences := 3
	sentenceLength := 12

	for i := 0; i < numSentences; i++ {
		sentence := ""
		wordsInSentence := 0
		for wordsInSentence < sentenceLength {
			wordIndex := rand.Intn(totalWords)
			word := wordCorpus[wordIndex]

			if wordsInSentence-lastUsed[word] > 5 || lastUsed[word] == 0 {
				if wordsInSentence == 0 {
					sentence += word
				} else {
					sentence += " " + word
				}

				lastUsed[word] = wordsInSentence + 1
				wordsInSentence++
			}
		}
		sentences = append(sentences, sentence+" ")
	}
	return sentences
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.timerStarted {
			m.timerStarted = true
			cmds = append(cmds, m.countdownTimer.Init())
		}

		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "backspace":
			if len(m.userInput) > 0 {
				m.userInput = m.userInput[:len(m.userInput)-1]
				m.cursorPos--
			}
		case "tab":
			m.showPrompt = true
			m.prompt = "🔄 Restart test?"

		case "enter":
			if m.showPrompt == true {
				m.countdownTimer.Stop()
				m = NewModel()
				return m, nil
			}
		default:
			if len(m.userInput) < len(m.targetText) {
				m.userInput += msg.String()
				m.cursorPos++
				if len(m.userInput) == len(m.targetText) && m.userInput[len(m.userInput)-1] == m.targetText[len(m.targetText)-1] {
					m.targetText = m.currentLines[1]
					m.cursorPos = 0
					m.userInput = ""
					m.currentLines[1], m.currentLines[2] = m.currentLines[2], generateSentences(generateWords())[0]
				}
			}
		}

	case timer.TickMsg:
		var cmd tea.Cmd
		m.countdownTimer, cmd = m.countdownTimer.Update(msg)
		cmds = append(cmds, cmd)

	case timer.TimeoutMsg:
		m = NewModel()
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var output strings.Builder

	for i, r := range m.targetText {
		if i < len(m.userInput) {
			if rune(m.userInput[i]) == r {
				output.WriteString(
					lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Render(string(r)),
				)
			} else {
				output.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(string(r)))
			}
		} else {
			output.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Render(string(r)))
		}
	}
	m.targetText = generateSentences(generateWords())[2]

	return fmt.Sprintf("%s\n\n%s\n%s\n%s\n\n%s",
		m.countdownTimer.View(),
		output.String(),
		m.currentLines[1],
		m.currentLines[2],
		m.prompt,
	)
}

func main() {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
