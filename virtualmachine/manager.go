package virtualmachine

import (
	"context"

	"gorm.io/gorm"
)

type VMManager struct {
	db *gorm.DB
}

func (vmm *VMManager) DestroyVirtualMachine(vm *Virtualmachine) error {
	tx := vmm.db.Begin()
	tx.Delete(vm)
	_, err := gorm.G[Virtualmachine](tx).Where("ID = ?", vm.ID).Delete(context.Background())
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (vmm *VMManager) SaveVirtualmachine(vm *Virtualmachine) error {
	ctx := context.Background()
	err := gorm.G[Virtualmachine](vmm.db).Create(ctx, vm)
	return err
}

func (vmm *VMManager) RetrieveVirtualmachineByName(name string) (*Virtualmachine, error) {
	var vm Virtualmachine
	ctx := context.Background()
	vm, err := gorm.G[Virtualmachine](vmm.db).Where("name = ?", name).First(ctx)
	if err != nil {
		return nil, err
	}
	return &vm, nil
}
