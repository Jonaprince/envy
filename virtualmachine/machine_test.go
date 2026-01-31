package virtualmachine_test

import (
	"os"
	"testing"

	"github.com/jonaprince/envy/virtualmachine"
)

func createVirtualMachine() *virtualmachine.Virtualmachine {
	vm := virtualmachine.NewVirtualmachine("TestVM", 2, 2048, "/tmp/testdisk.img", "/tmp/hypervisor-fw")
	return vm
}

// TODO: Check the disk integrity after flashing
func TestFlashDisk(t *testing.T) {
	machine := &virtualmachine.Virtualmachine{
		Disk: "/tmp/testdisk.img",
	}

	err := machine.FlashDisk("/path/to/image.img")

	if err != nil {
		t.Fatalf("FlashDisk failed: %v", err)
	}
}

func TestVirtualMachineLifecycle(t *testing.T) {
	machine := virtualmachine.NewVirtualmachine("toto", 1, 2048, "/tmp/testdisk.img", "/tmp/hypervisor-fw")
	// Flash a disk image only if testdisk does not exist
	if _, err := os.Stat(machine.Disk); os.IsNotExist(err) {
		err := machine.FlashDisk("../jammy-server-cloudimg-amd64.raw")
		if err != nil {
			t.Fatalf("Failed to flash disk: %v", err)
		}
	}
	_, err := machine.Init()
	defer machine.Destroy()
	if err != nil {
		t.Fatalf("Failed to init VM: %v", err)
	}

	err = machine.Create()
	if err != nil {
		t.Fatalf("Failed to create VM: %v", err)
	}

	err = machine.Start()
	if err != nil {
		t.Fatalf("Failed to start VM: %v", err)
	}
	if machine.Status != virtualmachine.StatusRunning {
		t.Fatalf("Expected machine status to be Running, got %v", machine.Status)
	}

	err = machine.Shutdown()
	if machine.Status != virtualmachine.StatusStopped {
		t.Fatalf("Expected machine status to be Stopped, got %v", machine.Status)
	}
	if err != nil {
		t.Fatalf("Failed to stop VM: %v", err)
	}
	// machine.UpdateStatus(virtualmachine.StatusRunning)

	// if machine.Status != virtualmachine.StatusStopped {
	// 	t.Fatalf("Expected machine status to be Running, got %v", machine.Status)
	// }
}
