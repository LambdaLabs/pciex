package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/tree"
)

type TreeKeymap struct {
	Left  key.Binding
	Right key.Binding
	Up    key.Binding
	Down  key.Binding
}

type TreeModel struct {
	Root    *Node
	CurNode *Node
	Keymap  TreeKeymap
}

func NewTreeModel() *TreeModel {
	m := &TreeModel{
		Keymap: TreeKeymap{
			Left: key.NewBinding(
				key.WithKeys("left"),
				key.WithHelp("←", "left"),
			),
			Right: key.NewBinding(
				key.WithKeys("right"),
				key.WithHelp("→", "right"),
			),
			Up: key.NewBinding(
				key.WithKeys("up"),
				key.WithHelp("↑", "up"),
			),
			Down: key.NewBinding(
				key.WithKeys("down"),
				key.WithHelp("↓", "down"),
			),
		},
	}
	return m
}

func (m *TreeModel) Init() tea.Cmd {
	m.CurNode = m.Root.children[0]
	return nil
}

func findClosestRelative(curNode *Node) *Node {
	if curNode.Parent == nil {
		return nil
	}
	if curNode.Idx < len(curNode.Parent.children)-1 {
		return curNode.Parent.children[curNode.Idx+1]
	}
	return findClosestRelative(curNode.Parent)
}

func (m *TreeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)
	debug(fmt.Sprintf("CurNode: %s \tCurNodeIdx: %d \tlen(Children): %d \n", m.CurNode.Name, m.CurNode.Idx, len(m.CurNode.children)))
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		debug(fmt.Sprintf("Pressed Key: %s\n", msg.String()))
		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up":
			if m.CurNode.Idx == 0 {
				if m.CurNode.Parent != nil && m.CurNode.Parent.Parent != nil {
					m.CurNode = m.CurNode.Parent
				}
			} else {
				m.CurNode = m.CurNode.Parent.children[m.CurNode.Idx-1]
			}
		case "down":
			if m.CurNode.Parent != nil && m.CurNode.Idx >= len(m.CurNode.Parent.children)-1 {
				// See if our parent has a sibling that we can traverse to.
				closestRelative := findClosestRelative(m.CurNode)
				if closestRelative == nil {
					break
				}
				m.CurNode = closestRelative
				break
			}
			m.CurNode = m.CurNode.Parent.children[m.CurNode.Idx+1]
		case "right":
			if len(m.CurNode.children) == 0 {
				break
			}
			m.CurNode = m.CurNode.children[0]
		case "left":
			if m.CurNode.Parent == nil || m.CurNode.Parent.Parent == nil {
				break
			}
			m.CurNode = m.CurNode.Parent
		}

	}

	debug(fmt.Sprintf("CurNode: %s \tCurNodeIdx: %d \tlen(Children): %d \n", m.CurNode.Name, m.CurNode.Idx, len(m.CurNode.children)))
	return m, tea.Batch(cmds...)
}

func (m *TreeModel) View() string {
	t := tree.Root(m.Root.Name).
		Enumerator(tree.DefaultEnumerator).
		EnumeratorStyle(enumeratorStyle).
		RootStyle(rootStyle).
		ItemStyleFunc(itemStyleFunc).
		Child(m.Root.children)
	return t.String()
}
