# pciex

PCI Explorer is a PCIe topology explorer and visualizer. The generated tree can be visualized and each element inspected using a simple terminal UI designed with the bubbletea framework.

<img width="1070" alt="Screenshot 2025-01-29 at 7 03 50 PM" src="https://github.com/user-attachments/assets/92bcf5d1-6f77-4374-bee1-edc491c80024" />

## Installation

This can be installed using:

```
GOPROXY=direct go install github.com/LambdaLabs/pciex@latest
```

Or you may run it directly via:

```
GOPROXY=direct go run github.com/LambdaLabs/pciex@latest
```

## Dependencies

This module requires the `lshw` utility to be installed.

### Ubuntu

```
sudo apt-get install lshw
```

### Redhat

```
sudo yum install lshw
```
