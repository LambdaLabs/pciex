package pcie

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/chigopher/pathlib"
)

type AdditionalDetails struct {
	NumaNode     *int
	LocalCPUList *string
}

func NewAdditionalDetailsFromSysfs(sysfsPath *pathlib.Path) (AdditionalDetails, error) {
	d := AdditionalDetails{}
	numa := sysfsPath.Join("numa_node")
	b, err := numa.ReadFile()
	if err != nil {
		if !os.IsNotExist(err) {
			return d, fmt.Errorf("reading numa_node: %w", err)
		}
	} else {
		numaInt, err := strconv.Atoi(strings.TrimSuffix(string(b), "\n"))
		if err != nil {
			return d, fmt.Errorf("parsing numa_node into int: %w", err)
		}
		d.NumaNode = &numaInt
	}

	localCPUList := sysfsPath.Join("local_cpulist")
	cpuListBytes, err := localCPUList.ReadFile()
	if err != nil {
		if !os.IsNotExist(err) {
			return d, fmt.Errorf("reading local_cpulist: %w", err)
		}
	} else {
		asStr := string(cpuListBytes)
		d.LocalCPUList = &asStr
	}

	return d, nil
}

type Details struct {
	AdditionalDetails // Details not provided by lshw that we need to scrape ourselves
	Id                string
	Class             string
	Claimed           bool
	Handle            string
	Description       string
	Product           string
	Vendor            string
	Physid            string
	Businfo           string
	Version           string
	Width             int
	Clock             int
	Serial            string
	Slot              string
	Units             string
	Size              int
	Configuration     map[string]any
	Capabilities      map[string]any
}

func (d *Details) GetAdditionalDetails() error {
	addressSplit := strings.Split(d.Businfo, "@")
	if len(addressSplit) != 2 {
		return nil
	}
	address := addressSplit[1]
	sysfsPath := pathlib.NewPath("/sys/bus/pci/devices/" + address)
	details, err := NewAdditionalDetailsFromSysfs(sysfsPath)
	if err != nil {
		return err
	}
	d.AdditionalDetails = details
	return nil
}

func (d Details) String() string {
	var s string
	s += d.Class + " | "
	switch d.Class {
	case "bridge":
		s += d.Handle
	case "bus":
		s += d.Description
	case "display", "memory", "communication", "generic", "network":
		s += d.Product
	default:
		s += d.Description
	}
	return s
}
