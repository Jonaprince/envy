package virtualmachine

import (
	"context"

	"gorm.io/gorm"
)

type VMManager struct {
	db *gorm.DB
}

func DestroyVirtualMachine(vm *Virtualmachine, db *gorm.DB) error {
	tx := db.Begin()
	tx.Delete(vm)
	_, err := gorm.G[Virtualmachine](tx).Where("ID = ?", vm.ID).Delete(context.Background())
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func CreateDisk(image string, vm *Virtualmachine) string {
	// Copy the base image to a new disk file
	diskPath := "/var/lib/envy/vms/" + vm.ID + "/disk.img"
	// os.MkdirAll("/var/lib/envy/vms/"+vm.ID, 0755)
	// cmd := exec.Command("cp", image, diskPath)
	// err := cmd.Run()
	// if err != nil {
	// 	return "", err
	// }
	return diskPath
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
