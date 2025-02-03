package models

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mistakenelf/teacup/statusbar"
)

type view string
type progressTick struct{}

var (
	progressPercentPerTick = 0.10
	progressTickDuration   = 100 * time.Millisecond
)

const (
	treeView    view = "tree"
	detailsView      = "details"
)

type rootKeymap struct {
	tab, quit, refresh, debug key.Binding
}

type viewportKeymaps struct {
	HalfPageUp   key.Binding
	HalfPageDown key.Binding
}

func newViewportKeymaps() viewportKeymaps {
	return viewportKeymaps{
		HalfPageUp: key.NewBinding(
			key.WithKeys("w"),
			key.WithHelp("w", "page up"),
		),
		HalfPageDown: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "page down"),
		),
	}
}

func (v viewportKeymaps) reassignViewportKeymap(m *viewport.KeyMap) {
	m.HalfPageUp = v.HalfPageUp
	m.HalfPageDown = v.HalfPageDown
}

type RootModel struct {
	Tree            *TreeModel
	treeViewport    viewport.Model
	treeStyle       lipgloss.Style
	detailsViewport viewport.Model
	detailsStyle    lipgloss.Style
	view            view
	height          int
	width           int
	help            help.Model
	keymap          rootKeymap
	viewportKeymap  viewportKeymaps
	progress        progress.Model
	showProgress    bool
	status          statusbar.Model
	hostname        string
}

func NewRootModel() (*RootModel, error) {
	details := viewport.New(0, 0)
	tree := viewport.New(0, 0)

	viewportKeymap := newViewportKeymaps()
	viewportKeymap.reassignViewportKeymap(&tree.KeyMap)
	viewportKeymap.reassignViewportKeymap(&details.KeyMap)

	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("getting hostname from os: %w", err)
	}

	m := &RootModel{
		Tree:            NewTreeModel(),
		treeViewport:    tree,
		detailsViewport: details,
		view:            treeView,
		help:            help.New(),
		keymap: rootKeymap{
			tab: key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "next window"),
			),
			quit: key.NewBinding(
				key.WithKeys("esc", "ctrl+c"),
				key.WithHelp("esc", "quit"),
			),
			refresh: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "refresh"),
			),
			debug: key.NewBinding(
				key.WithKeys("d"),
				key.WithHelp("d", "debug"),
			),
		},
		viewportKeymap: viewportKeymap,
		progress:       progress.New(progress.WithDefaultScaledGradient()),
		showProgress:   true,
		status: statusbar.New(
			statusbar.ColorConfig{
				Foreground: lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
				Background: lipgloss.AdaptiveColor{Light: string(foregroundColor), Dark: string(foregroundColor)},
			},
			statusbar.ColorConfig{
				Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
				Background: lipgloss.AdaptiveColor{Light: "#3c3836", Dark: "#3c3836"},
			},
			statusbar.ColorConfig{
				Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
				Background: lipgloss.AdaptiveColor{Light: "#A550DF", Dark: "#A550DF"},
			},
			statusbar.ColorConfig{
				Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
				Background: lipgloss.AdaptiveColor{Light: string(lambdaPurple), Dark: string(lambdaPurple)},
			},
		),
		hostname: hostname,
	}
	m.status.SetContent("one", "two", "three", "four")
	return m, nil
}

func (m *RootModel) Init() tea.Cmd {
	return tea.Batch(
		m.Tree.Init(),
		m.detailsViewport.Init(),
		m.progress.Init(),
		func() tea.Msg {
			return progressTick{}
		},
	)
}

func (m *RootModel) activeViewport() *viewport.Model {
	if m.view == treeView {
		return &m.treeViewport
	}
	return &m.detailsViewport
}

func (m *RootModel) updateViewports(msg tea.Msg) tea.Cmd {
	newViewport, cmd := m.activeViewport().Update(msg)
	*m.activeViewport() = newViewport
	return cmd
}

func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.tab):
			if m.view == treeView {
				m.view = detailsView
			} else {
				m.view = treeView
			}
			m.resizeElements()
		case key.Matches(msg, m.viewportKeymap.HalfPageUp):
			cmds = append(cmds, m.updateViewports(msg))
		case key.Matches(msg, m.viewportKeymap.HalfPageDown):
			cmds = append(cmds, m.updateViewports(msg))
		case key.Matches(msg, m.keymap.refresh):
			cmds = append(cmds, m.Tree.CurNode.RefreshDetail())
		}

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		if m.showProgress {
			p := m.progress
			p.Width = m.width
			m.progress = p
			//cmds = append(cmds, m.image.SetSize(m.width-10, m.height-10))
		}
		m.resizeElements()
	case progressTick:
		p := m.progress
		if p.Percent() >= 1.0 {
			m.showProgress = false
		} else {
			cmds = append(cmds, p.IncrPercent(progressPercentPerTick))
			m.progress = p
			cmds = append(cmds, func() tea.Msg {
				time.Sleep(progressTickDuration)
				return progressTick{}
			})
		}
		return m, tea.Batch(cmds...)
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		cmds = append(cmds, cmd)
		m.progress = progressModel.(progress.Model)
		return m, tea.Batch(cmds...)
	}

	if !m.showProgress && m.view == treeView {
		newTree, cmd := m.Tree.Update(msg)
		m.Tree = newTree.(*TreeModel)
		cmds = append(cmds, cmd)
	}

	scrollPercent := m.activeViewport().ScrollPercent()
	m.status.SetContent(
		m.Tree.CurNode.Detail.Businfo,
		m.hostname,
		fmt.Sprintf("%d", int(scrollPercent*100))+"%",
		string(m.view),
	)
	return m, tea.Batch(cmds...)
}

func (m *RootModel) resizeElements() {
	var (
		height              int = m.height - 4
		width               int = m.width - 4
		detailsHeightOffset int = 1
		treeHeightOffset    int = 1
		treeStyle           lipgloss.Style
		detailsStyle        lipgloss.Style
	)

	if m.view == treeView {
		treeStyle = focusedModelStyle
		detailsStyle = modelStyle
	} else {
		treeStyle = modelStyle
		detailsStyle = focusedModelStyle
	}

	m.treeStyle = treeStyle.
		Height(height - treeHeightOffset).
		Width(width / 2)
	m.treeViewport.Height = m.treeStyle.GetHeight()
	m.treeViewport.Width = m.treeStyle.GetWidth()

	m.detailsStyle = detailsStyle.
		Height(height - detailsHeightOffset).
		Width((width - m.treeStyle.GetWidth()))
	m.detailsViewport.Height = m.detailsStyle.GetHeight()
	m.detailsViewport.Width = m.detailsStyle.GetWidth()
	m.status.SetSize(m.width)
}

func (m *RootModel) View() string {
	var s string
	if m.showProgress {
		return m.progress.View()
	}
	m.detailsViewport.SetContent(m.Tree.CurNode.GetDetail())
	m.treeViewport.SetContent(m.Tree.View())

	help := m.help.ShortHelpView([]key.Binding{
		m.keymap.tab,
		m.keymap.quit,
		m.keymap.refresh,
		m.keymap.debug,
	})
	treeHelp := m.help.ShortHelpView([]key.Binding{
		m.treeViewport.KeyMap.HalfPageUp,
		m.treeViewport.KeyMap.HalfPageDown,
		m.Tree.Keymap.Left,
		m.Tree.Keymap.Right,
		m.Tree.Keymap.Up,
		m.Tree.Keymap.Down,
	})
	detailsHelp := m.help.ShortHelpView([]key.Binding{
		m.detailsViewport.KeyMap.HalfPageUp,
		m.detailsViewport.KeyMap.HalfPageDown,
	})
	treeViewport := lipgloss.JoinVertical(lipgloss.Top, m.treeViewport.View(), treeHelp)
	detailsViewport := lipgloss.JoinVertical(lipgloss.Top, m.detailsViewport.View(), detailsHelp)
	s += lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.treeStyle.Render(treeViewport),
		m.detailsStyle.Render(detailsViewport),
	)
	return lipgloss.JoinVertical(
		lipgloss.Top,
		s,
		help,
		m.status.View(),
	)
}
