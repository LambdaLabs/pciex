package models

import (
	"fmt"
	"time"

	"github.com/LambdaLabs/pciex/pcie"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/tree"
	"gopkg.in/yaml.v2"
)

type Node struct {
	Name   string
	Detail pcie.Details
	Parent *Node
	// Idx is the index of Node in Parent's Children field.
	Idx      int
	children Children
	model    *TreeModel
}

func NewNode(name string, detail pcie.Details, parent *Node, model *TreeModel) *Node {
	return &Node{
		Name:     name,
		Detail:   detail,
		Parent:   parent,
		children: Children{},
		model:    model,
	}
}

func (n *Node) GetDetail() string {
	out, err := yaml.Marshal(n.Detail)
	if err != nil {
		panic(fmt.Errorf("unmarshalling yaml: %w", err))
	}
	return string(out)
}

func (n *Node) RefreshDetail() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(5 * time.Second)
		return nil
	}
}

func (n *Node) AddChild(name string, detail pcie.Details) *Node {
	child := &Node{
		Name:   name,
		Detail: detail,
		Parent: n,
		model:  n.model,
	}
	n.children = append(n.children, child)
	child.Idx = len(n.children) - 1
	return child
}

func (n Node) String() string {
	return n.Name
}

func (n Node) Value() string {
	return n.Name
}

func (n Node) Children() tree.Children {
	return n.children
}

func (n Node) Hidden() bool {
	return false
}

type Children []*Node

func (c Children) At(index int) tree.Node {
	return c[index]
}

func (c Children) Length() int {
	return len(c)
}
