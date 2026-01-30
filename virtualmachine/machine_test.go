package virtualmachine_test

import (
	"os"
	"testing"

	"github.com/jonaprince/envy/virtualmachine"
)

func TestFlashDisk(t *testing.T) {
	machine := &virtualmachine.Virtualmachine{
		Disk: "/tmp/testdisk.img",
	}
	os.Open()
	err := machine.FlashDisk("/path/to/image.img")

	if err != nil {
		t.Fatalf("FlashDisk failed: %v", err)
	}
}
