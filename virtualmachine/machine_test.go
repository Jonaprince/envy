package virtualmachine_test

import (
	"os"
	"testing"

	"github.com/jonaprince/envy/virtualmachine"
)

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

func TestVirtualCreation(t *testing.T) {
	machine := virtualmachine.NewVirtualmachine("toto", 1, 512, "/tmp/testdisk.img")
	// Flash a disk image only if testdisk does not exist
	if _, err := os.Stat(machine.Disk); os.IsNotExist(err) {
		err := machine.FlashDisk("jammy-server-cloudimg-amd64.raw")
		if err != nil {
			t.Fatalf("Failed to flash disk: %v", err)
		}
	}
	_, err := machine.Init()
	// defer machine.Destroy()
	if err != nil {
		t.Fatalf("Failed to init VM: %v", err)
	}

	err = machine.Create()
	if err != nil {
		t.Fatalf("Failed to create VM: %v", err)
	}

	// machine.UpdateStatus(virtualmachine.StatusRunning)

	// if machine.Status != virtualmachine.StatusStopped {
	// 	t.Fatalf("Expected machine status to be Running, got %v", machine.Status)
	// }
}
