package virtualmachine

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/jonaprince/envy/cloudhypervisor"
)

type Status int

const (
	StatusStopped Status = iota
	StatusUninitialized
	StatusInitialized
	StatusStarting
	StatusRunning
	StatusStopping
	StatusError
)

type Virtualmachine struct {
	// gorm.Model
	ID            string `gorm:"primaryKey"`
	Name          string
	CPU           int
	Memory        int
	MachineSocket string
	Status        Status
	Client        cloudhypervisor.Client
	PID           int
	Disk          string
}

var statusNames = map[Status]string{
	StatusStopped:       "Stopped",
	StatusUninitialized: "Uninitialized",
	StatusInitialized:   "Initialized",
	StatusStarting:      "Starting",
	StatusRunning:       "Running",
	StatusStopping:      "Stopping",
	StatusError:         "Error",
}

func NewVirtualmachine(name string, cpu, memory int) *Virtualmachine {
	vm := &Virtualmachine{
		ID:     uuid.New().String(),
		Name:   name,
		CPU:    cpu,
		Memory: memory,
		Status: StatusStopped,
	}
	// pid, err := createCloudHypervisorVM("/var/run/envy/")
	// vm.PID = pid
	// if err != nil {
	// 	return nil, err
	// }
	return vm
}

// Flash an image to the VM disk
func (vm *Virtualmachine) FlashDisk(image string) error {
	src, err := os.Open(image)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.OpenFile(vm.Disk, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer dst.Close()
	// 4MB buffer for good performances
	buf := make([]byte, 4*1024*1024)
	_, err = io.CopyBuffer(dst, src, buf)
	return err
}

// Create the cloud hypervisor thread and return the PID of the detached process
func (vm *Virtualmachine) Init() (int, error) {
	cmd := exec.Command("cloud-hypervisor", "--api-socket", vm.MachineSocket)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil
	err := cmd.Start()
	if err != nil {
		return -1, err
	}
	// Ensure the socket file is created
	maxRetry := 10
	for i := 0; i < maxRetry; i++ {
		if _, err := os.Stat(vm.MachineSocket); err == nil {
			return cmd.Process.Pid, nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return -1, fmt.Errorf("socket %s was not created", vm.MachineSocket)
}

// Create a new virtual machine using the cloud hypervisor API
func (vm *Virtualmachine) Create() error {
	vmConfig := cloudhypervisor.VMConfig{
		Cpus: &cloudhypervisor.CpusConfig{
			BootVcpus: vm.CPU,
			MaxVcpus:  vm.CPU,
		},
		Disks: []cloudhypervisor.DiskConfig{
			{Path: vm.Disk},
		},
	}
	vm.Client.CreateVM(vmConfig)
	return nil
}

// Start the virtual machine
func (vm *Virtualmachine) Start() error {
	return nil
}

// Stop the virtual machine
func (vm *Virtualmachine) Stop() error {
	return nil
}

// Destroy the virtual machine
func (vm *Virtualmachine) Destroy() error {
	return nil
}

// Contact the cloud hypervisor API to check the vm status
func (vm *Virtualmachine) UpdateStatus(status Status) {
	vm.Status = status
}

// Reconcile the VM state between desired and actual state
func (vm *Virtualmachine) Reconcile() error {
	return nil
}
