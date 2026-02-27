package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Exchange struct {
	Rates map[string]float32 `json:"rates"`
}

type ratesMsg map[string]float32
type errMsg error

type model struct {
	baseChoices          []string
	targetChoices        []string
	selectedBaseChoice   int
	selectedTargetChoice int
	amount               textinput.Model
	cursor               int
	result               float32
	rates                map[string]float32
	loading              bool
	err                  error
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "100"
	ti.CharLimit = 15
	ti.Focus()

	return model{
		amount:               ti,
		baseChoices:          []string{"INR", "USD", "EUR", "GBP", "JPY"},
		targetChoices:        []string{"INR", "USD", "EUR", "GBP", "JPY"},
		selectedBaseChoice:   1,
		selectedTargetChoice: 0,
	}
}

func (m model) Init() tea.Cmd {
	return fetchRates
}

func fetchRates() tea.Msg {
	rates, err := GetRates()
	if err != nil {
		return errMsg(err)
	}
	return ratesMsg(rates)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor == 0 && m.selectedBaseChoice > 0 {
				m.selectedBaseChoice--
			} else if m.cursor == 1 && m.selectedTargetChoice > 0 {
				m.selectedTargetChoice--
			}

		case "down", "j":
			if m.cursor == 0 && m.selectedBaseChoice < len(m.baseChoices)-1 {
				m.selectedBaseChoice++
			} else if m.cursor == 1 && m.selectedTargetChoice < len(m.targetChoices)-1 {
				m.selectedTargetChoice++
			}

		case "tab":
			m.cursor = (m.cursor + 1) % 3
			if m.cursor == 2 {
				m.amount.Focus()
			} else {
				m.amount.Blur()
			}

		case "enter":
			if m.rates != nil {
				m.result = m.calculateConversion()
			}

		case "r":
			m.loading = true
			m.err = nil
			return m, fetchRates
		}

	case ratesMsg:
		m.rates = msg
		m.loading = false
		m.err = nil
		m.result = m.calculateConversion()
		return m, nil

	case errMsg:
		m.err = msg
		m.loading = false
		return m, nil
	}

	if m.cursor == 2 {
		m.amount, cmd = m.amount.Update(msg)
	}

	return m, cmd
}

func (m model) calculateConversion() float32 {
	if m.rates == nil {
		return 0
	}

	amt, err := strconv.ParseFloat(m.amount.Value(), 32)
	if err != nil || amt <= 0 {
		return 0
	}

	base := m.baseChoices[m.selectedBaseChoice]
	target := m.targetChoices[m.selectedTargetChoice]

	return float32(amt) * (m.rates[target] / m.rates[base])
}

func (m model) View() string {
	var s strings.Builder

	s.WriteString("╔════════════════════════════════════╗\n")
	s.WriteString("║    CURRENCY CONVERTER              ║\n")
	s.WriteString("╚════════════════════════════════════╝\n\n")

	if m.loading {
		s.WriteString("Loading exchange rates...\n\n")
	} else if m.err != nil {
		s.WriteString("Error: " + m.err.Error() + "\n\n")
	}

	s.WriteString("1. Base Currency:\n")
	for i, choice := range m.baseChoices {
		cursor := " "
		if m.cursor == 0 && m.selectedBaseChoice == i {
			cursor = "→"
		}
		check := "○"
		if m.selectedBaseChoice == i {
			check = "●"
		}
		s.WriteString("   " + cursor + " " + check + " " + choice + "\n")
	}

	s.WriteString("\n2. Target Currency:\n")
	for i, choice := range m.targetChoices {
		cursor := " "
		if m.cursor == 1 && m.selectedTargetChoice == i {
			cursor = "→"
		}
		check := "○"
		if m.selectedTargetChoice == i {
			check = "●"
		}
		s.WriteString("   " + cursor + " " + check + " " + choice + "\n")
	}

	s.WriteString("\n3. Amount:\n")
	prefix := "  "
	if m.cursor == 2 {
		prefix = "→ "
	}
	s.WriteString(prefix + m.amount.View() + "\n")

	if m.result > 0 && m.rates != nil {
		base := m.baseChoices[m.selectedBaseChoice]
		target := m.targetChoices[m.selectedTargetChoice]
		s.WriteString("\n┌────────────────────────────────────┐\n")
		fmt.Fprintf(&s, "│ %.2f %s = %.2f %s\n",
			mustParseFloat(m.amount.Value()), base, m.result, target)
		fmt.Fprintf(&s, "│ Rate: 1 %s = %.4f %s\n",
			base, m.rates[target]/m.rates[base], target)
		s.WriteString("└────────────────────────────────────┘\n")
	}

	s.WriteString("\n")
	s.WriteString("  ⌨  Navigation: ↑/↓ or j/k\n")
	s.WriteString("  ⇥  Tab: Next field\n")
	s.WriteString("  ↵  Enter: Convert\n")
	s.WriteString("  r: Refresh rates\n")
	s.WriteString("  q: Quit\n")

	return s.String()
}

func mustParseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func GetRates() (map[string]float32, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	requestURL := "https://open.er-api.com/v6/latest/USD"

	res, err := client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", res.StatusCode)
	}

	var exchange Exchange
	if err := json.NewDecoder(res.Body).Decode(&exchange); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return exchange.Rates, nil
}
