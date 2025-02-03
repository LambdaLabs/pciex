# pciex

PCI Explorer is a PCIe topology explorer and visualizer. The generated tree can be visualized and each element inspected using a simple terminal UI designed with the bubbletea framework.

![Screenshot 2025-02-03 at 5 24 40 PM](https://github.com/user-attachments/assets/4a0fcaf6-50de-4a78-a7bb-430c89f41f46)

## Installation

This can be installed using:

```
go install github.com/LambdaLabs/pciex@latest
```

Or you may run it directly via:

```
go run github.com/LambdaLabs/pciex@latest
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
