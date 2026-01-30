package virtualmachine

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	errors        chan (error)
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

func NewVirtualmachine(name string, cpu, memory int, db *gorm.DB) (*Virtualmachine, error) {
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
	err := saveVirtualmachine(vm, db)
	return vm, err
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
