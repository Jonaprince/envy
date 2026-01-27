package virtualmachine

import (
	"context"
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
	StatusStarting
	StatusRunning
	StatusStopping
	StatusError
)

type Virtualmachine struct {
	gorm.Model
	ID            string
	Name          string
	CPU           int
	Memory        int
	MachineSocket string
	Status        Status
}

var statusNames = map[Status]string{
	StatusStopped:  "Stopped",
	StatusStarting: "Starting",
	StatusRunning:  "Running",
	StatusStopping: "Stopping",
	StatusError:    "Error",
}

func NewVirtualmachine(name string, cpu, memory int, db *gorm.DB) (*Virtualmachine, error) {
	vm := &Virtualmachine{
		ID:     uuid.New().String(),
		Name:   name,
		CPU:    cpu,
		Memory: memory,
		Status: StatusStopped,
	}
	err := createCloudHypervisorVM("/tmp/ch-id")
	if err != nil {
		return nil, err
	}
	err = saveVirtualmachine(vm, db)
	return vm, err
}

func createCloudHypervisorVM(chSocket string) error {
	cmd := exec.Command("cloud-hypervisor", "--api-socket", chSocket)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}

	// cmd.Stdout = nil
	// cmd.Stderr = nil
	// cmd.Stdin = nil

	err := cmd.Start()
	if err != nil {
		return err
	}

	// Ensure the socket file is created
	maxRetry := 10
	for i := 0; i < maxRetry; i++ {
		if _, err := os.Stat(chSocket); err == nil {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("socket %s was not created", chSocket)
}

func saveVirtualmachine(vm *Virtualmachine, db *gorm.DB) error {
	ctx := context.Background()
	err := gorm.G[Virtualmachine](db).Create(ctx, vm)
	return err
}

func RetrieveVirtualmachineByName(name string, db *gorm.DB) (*Virtualmachine, error) {
	var vm Virtualmachine
	ctx := context.Background()
	vm, err := gorm.G[Virtualmachine](db).Where("name = ?", name).First(ctx)
	if err != nil {
		return nil, err
	}
	return &vm, nil
}
