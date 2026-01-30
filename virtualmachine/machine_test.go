package virtualmachine_test

import (
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
	machine := &virtualmachine.Virtualmachine{}
	machine.FlashDisk("")
	machine.Start()

	if machine.Status != virtualmachine.StatusRunning {
		t.Fatalf("Expected machine status to be Running, got %v", machine.Status)
	}
}
