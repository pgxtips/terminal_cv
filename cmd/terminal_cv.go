package terminal_cv

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
)

type model struct {
	Tabs       []string
	TabContent []string
	activeTab  int
	w  int
	h  int
	Session ssh.Session 
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w = msg.Width
		m.h = msg.Height
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l", "n", "tab":
			m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
	}

	return m, nil
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

func (m model) View() string {
	
	if m.w == 0 {
		return ""
	}

	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.activeTab
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(
		windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).
		Render(m.TabContent[m.activeTab]))

	// center the whole document
	docStyle := lipgloss.NewStyle().
		Width(m.w).
		Height(m.h).
		Align(lipgloss.Center, lipgloss.Center).
		Padding(1, 2, 1, 2)

	return docStyle.Render(doc.String())
}

func terminalcv(width int, height int, session ssh.Session) model {
	tabs := []string{"Personal", "Education", "Work Experience", "Projects", "About Me"}

	personalContent := `name: Brandon Lee Gill
	location: Manchester
	email: brandongill123@gmail.com
	linkedin: linkedin.com/in/brandon-lee-gill/
	`
	
	educationContent := `
	MSc in Artificial Intelligence -> Manchester Metropolitan University -> 2025 - 2026

	BSC (Hons) in Computer Science -> Manchester Metropolitan University -> 2021 - 2025

	A-Levels -> Computer Science -> Statistics -> Psychology

	GCSEs -> 9 GCSEs A*-C -> Including English, Maths & Science
	`

	workExperienceContent := `
	Software Developer - Chippy Digital
	 -> 2023 - Present
	 -> Developed a multitude of web apps and tools for
	 clients in the education domain.
	 Working within a team of 4 engineers, I've developed
	 full stack applications using various technologies 
	 and have overseen many projects from start to finish.
	 -> Technologies used: HTML, CSS, JavaScript, PHP,
	 Rust, SQL, Postgres, Docker, Github Actions
	 and Digital Ocean.
	`

	projectsContent := `The Binding of Isaac Achievement Tracker -> An online platform that uses your Steam data to display your achievement progress per character and overall in the hit game The Binding of Isaac by Ed McMillen. 

	beam -> An easy to adopt, flexible, open source recommendation engine

	beam-demo -> A real-world application to demonstrate the use of beam

	py-tools -> An online tool for executing python scripts with a long life time

	gnews-rs -> a rust library to handle google news fetching and processing 
	`

	aboutmeContent := `
	From a young age, I’ve been passionate about computers and technology. I built my first desktop PC at age 11, and programming quickly became my main hobby, which is something I’ve pursued seriously for over a decade. This early interest developed into a deep understanding of software design, system architecture, and a strong proficiency in multiple programming languages including C++, Rust, JavaScript, and Python.
	`

	tabContent := []string{ personalContent, educationContent, workExperienceContent, projectsContent, aboutmeContent}
	m := model{Tabs: tabs, TabContent: tabContent}

	m.Session = session
	m.w = width
	m.h = height

	return m

	// if _, err := tea.NewProgram(m).Run(); err != nil {
	// 	fmt.Println("Error running program:", err)
	// 	os.Exit(1)
	// }
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
