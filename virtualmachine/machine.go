package virtualmachine

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/jonaprince/envy/cloudhypervisor"
)

type Status int

const (
	Stopped Status = iota
	Uninitialized
	Created
	Starting
	Running
	Stopping
	Error
	Destroyed
)

type Virtualmachine struct {
	// gorm.Model
	ID            string `gorm:"primaryKey"`
	Name          string
	CPU           int
	Memory        int
	MachineSocket string
	SerialSocket  string
	State         Status
	DesiredState  Status
	Client        *cloudhypervisor.Client
	PID           int
	Disk          string
	Firmware      string
}

var statusNames = map[Status]string{
	Stopped:       "stopped",
	Uninitialized: "uninitialized",
	Created:       "created",
	Starting:      "starting",
	Running:       "running",
	Stopping:      "stopping",
	Error:         "error",
	Destroyed:     "destroyed",
}

var statusValues = map[string]Status{
	"stopped":       Stopped,
	"uninitialized": Uninitialized,
	"created":       Created,
	"starting":      Starting,
	"running":       Running,
	"stopping":      Stopping,
	"error":         Error,
	"destroyed":     Destroyed,
}

func NewVirtualmachine(name string, cpu, memory int, disk, firmware string) *Virtualmachine {
	id := uuid.New().String()
	os.MkdirAll("/var/run/envy", 0755)
	chSocket := fmt.Sprintf("/var/run/envy/ch-%s.sock", id)
	serialSocket := fmt.Sprintf("/var/run/envy/ch-%s.console", id)
	vm := &Virtualmachine{
		ID:            id,
		Name:          name,
		CPU:           cpu,
		Memory:        memory,
		State:         Stopped,
		DesiredState:  Stopped,
		Disk:          disk,
		MachineSocket: chSocket,
		SerialSocket:  serialSocket,
		Client:        cloudhypervisor.NewClient(chSocket),
		Firmware:      firmware,
	}
	return vm
}

// Flash an image to the VM disk
func (vm *Virtualmachine) FlashDisk(image string) error {
	src, err := os.Open(image)
	if err != nil {
		return err
	}
	defer src.Close()
	// Open the destination disk file and create it if it does not exist
	dst, err := os.OpenFile(vm.Disk, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
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
			vm.PID = cmd.Process.Pid
			vm.State = Created
			return cmd.Process.Pid, nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return -1, fmt.Errorf("socket %s was not created", vm.MachineSocket)
}

// Create a new virtual machine using the cloud hypervisor API
func (vm *Virtualmachine) Create() error {
	vmConfig := cloudhypervisor.VMConfig{
		Cpus: cloudhypervisor.CpusConfig{
			BootVcpus: vm.CPU,
			MaxVcpus:  vm.CPU,
		},
		Disks: []cloudhypervisor.DiskConfig{
			{Path: vm.Disk},
		},
		Payload: cloudhypervisor.PayloadConfig{
			Firmware: vm.Firmware,
		},
		Serial: cloudhypervisor.ConsoleConfig{
			Mode:   "Socket",
			Socket: vm.SerialSocket,
		},
		MemoryConfig: cloudhypervisor.MemoryConfig{
			Size: int64(vm.Memory * 1024 * 1024),
		},
	}
	err := vm.Client.CreateVM(vmConfig)
	return err
}

// Start the v
// irtual machine
// TODO: Check the VM status after starting
func (vm *Virtualmachine) Start() error {
	if vm.State != Created {
		return fmt.Errorf("VM is not initialized")
	}
	err := vm.Client.BootVM()
	if err != nil {
		return err
	}
	vm.State = Running
	return nil
}

// Stop the virtual machine
// TODO: Check the VM status after shutting down
// TODO: Implement a soft and hard shutdown
func (vm *Virtualmachine) Shutdown() error {
	if vm.State != Running {
		return fmt.Errorf("VM is not running")
	}
	err := vm.Client.ShutdownVM()
	if err != nil {
		return err
	}
	vm.State = Stopped
	return nil
}

// Destroy the virtual machine
func (vm *Virtualmachine) Destroy() error {
	process, err := os.FindProcess(vm.PID)
	if err != nil {
		return err
	}
	err = process.Kill()
	if err != nil {
		return err
	}
	// Remove the disk file
	err = os.Remove(vm.Disk)
	if err != nil {
		return err
	}

	// Clean up the socket file
	err = os.Remove(vm.MachineSocket)
	if err != nil {
		return err
	}
	// Mark the VM as destroyed
	vm.State = Destroyed
	return nil
}

// Contact the cloud hypervisor API to check the vm status
func (vm *Virtualmachine) UpdateStatus() {
	info, err := vm.Client.GetVMInfo()
	if err != nil {
		vm.State = Error
		return
	}
	vm.State = statusValues[(string)(info.State)]
}

// Reconcile the VM state between desired and actual state
func (vm *Virtualmachine) Reconcile() {
	vm.UpdateStatus()
	if vm.State == vm.DesiredState {
		return
	}
	if vm.State == Error {
		// TODO; Handle error state, maybe try to restart the VM or mark it for deletion
		// Maybe should add a function mitigiateError() to the VM struct to handle this case
		slog.Error("VM is in error state", "vm", vm.Name)
		return
	}
	slog.Info("VM state not matching desired state", "vm", vm.Name, "from", statusNames[vm.State], "to", statusNames[vm.DesiredState])
	switch vm.DesiredState {
	case Running:
		vm.Start()
	case Stopped:
		vm.Shutdown()
	case Destroyed:
		vm.Destroy()
	}
}
