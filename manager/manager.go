package manager

import (
	"context"
	"errors"

	virtualmachine "github.com/jonaprince/envy/virtualmachine"
	"gorm.io/gorm"
)

type VMManager struct {
	db  *gorm.DB
	vms map[string]*virtualmachine.Virtualmachine
}

func (vmm *VMManager) DeleteVirtualMachine(vm *virtualmachine.Virtualmachine) error {
	// Ensure the VM is destroyed before deleting from the database
	if vm.State != virtualmachine.Destroyed {
		return errors.New("You can't delete a virtual machine from the db which is not destroyed")
	}
	tx := vmm.db.Begin()
	tx.Delete(vm)
	_, err := gorm.G[virtualmachine.Virtualmachine](tx).Where("ID = ?", vm.ID).Delete(context.Background())
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (vmm *VMManager) SaveVirtualmachine(vm *virtualmachine.Virtualmachine) error {
	ctx := context.Background()
	err := gorm.G[virtualmachine.Virtualmachine](vmm.db).Create(ctx, vm)
	return err
}

func (vmm *VMManager) ReconcileVirtualMachine(vm *virtualmachine.Virtualmachine) error {
	for _, machine := range vmm.vms {
		machine.Reconcile()
		if machine.State == virtualmachine.Destroyed {
			vmm.DeleteVirtualMachine(machine)
			delete(vmm.vms, machine.ID)
		} else {
			vmm.SaveVirtualmachine(machine)
		}
	}
	return nil
}

func (vmm *VMManager) RetrieveAllVirtualmachines() ([]virtualmachine.Virtualmachine, error) {
	var vms []virtualmachine.Virtualmachine
	ctx := context.Background()
	vms, err := gorm.G[virtualmachine.Virtualmachine](vmm.db).Find(ctx)
	if err != nil {
		return nil, err
	}
	return vms, nil
}
