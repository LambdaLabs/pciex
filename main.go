package main

// These imports will be used later on the tutorial. If you save the file
// now, Go might complain they are unused, but that's fine.
// You may also need to run `go mod tidy` to download bubbletea and its
// dependencies.
import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/LambdaLabs/pciex/models"
	"github.com/LambdaLabs/pciex/pcie"
	tea "github.com/charmbracelet/bubbletea"
)

type Slot struct {
	Domain string
	Bus    string
	Device string
	Func   string
}

func NewSlotFromString(s string) Slot {
	split := strings.Split(s, ":")
	deviceFunc := split[2]
	deviceFuncSplit := strings.Split(deviceFunc, ".")

	return Slot{
		Domain: split[0],
		Bus:    split[1],
		Device: deviceFuncSplit[0],
		Func:   deviceFuncSplit[1],
	}
}

type LshwElem struct {
	pcie.Details
	Children []LshwElem
}

func buildPCIETreeHelper(parent *models.Node, children []LshwElem) error {
	if parent == nil {
		panic("parent is nil")
	}
	for _, child := range children {
		if err := child.Details.GetAdditionalDetails(); err != nil {
			return fmt.Errorf("getting additional details after json unmarshal: %w", err)
		}
		childNode := parent.AddChild(child.Details.String(), child.Details)
		if err := buildPCIETreeHelper(childNode, child.Children); err != nil {
			return err
		}
	}
	return nil
}

func buildPCIETree(tree *models.TreeModel) error {
	cmd := exec.Command("/usr/bin/lshw", "-json")
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("reading command output: %w", err)
	}
	lshw := LshwElem{}
	if err := json.Unmarshal(out, &lshw); err != nil {
		// Some versions of lshw wrap the JSON in a single-element array.
		// Why? I don't know. See if we can successfully unmarshal into an
		// array.
		lshwArray := []LshwElem{}
		if newErr := json.Unmarshal(out, &lshwArray); newErr != nil {
			return fmt.Errorf("unmarshalling json: %w", err)
		}
		lshw = lshwArray[0]
	}
	tree.Root = models.NewNode("root", pcie.Details{}, nil, tree)
	// Find PCI element
	for _, child := range lshw.Children[0].Children {
		if !strings.HasPrefix(child.Id, "pci") {
			continue
		}
		if err := child.Details.GetAdditionalDetails(); err != nil {
			return fmt.Errorf("getting additional details after json unmarshal: %w", err)
		}
		childNode := tree.Root.AddChild(child.String(), child.Details)
		if err := buildPCIETreeHelper(childNode, child.Children); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	rootModel, err := models.NewRootModel()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	if err := buildPCIETree(rootModel.Tree); err != nil {
		fmt.Printf("Error occurred: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(rootModel)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error occurred: %v\n", err)
		os.Exit(1)
	}
}
