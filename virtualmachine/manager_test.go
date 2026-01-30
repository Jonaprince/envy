// package virtualmachine_test

// import (
// 	"testing"

// 	"github.com/jonaprince/envy/virtualmachine"
// 	"gorm.io/driver/sqlite"
// 	"gorm.io/gorm"
// )

// func createSampleVirtualMachine(db *gorm.DB) *virtualmachine.Virtualmachine {
// 	return virtualmachine.NewVirtualmachine("TestVM", 4, 8192)
// }

// func createTestDB(t *testing.T) *gorm.DB {
// 	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
// 	if err != nil {
// 		t.Fatalf("Failed to open database: %v", err)
// 	}
// 	err = db.AutoMigrate(&virtualmachine.Virtualmachine{})
// 	if err != nil {
// 		t.Fatalf("Failed to migrate database: %v", err)
// 	}
// 	return db
// }

// func TestCreateVirtualmachine(t *testing.T) {
// 	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
// 	if err != nil {
// 		t.Fatalf("Failed to open database: %v", err)
// 	}
// 	db.AutoMigrate(&virtualmachine.Virtualmachine{})

// 	vm, err := createSampleVirtualMachine(db)
// 	if err != nil {
// 		t.Fatalf("Failed to create virtual machine: %v", err)
// 	}

// 	if vm.ID == "" {
// 		t.Errorf("Expected non-empty ID, got '%s'", vm.ID)
// 	}
// 	if vm.Name != "TestVM" {
// 		t.Errorf("Expected Name 'TestVM', got '%s'", vm.Name)
// 	}
// 	if vm.CPU != 4 {
// 		t.Errorf("Expected CPU 4, got %d", vm.CPU)
// 	}
// 	if vm.Memory != 8192 {
// 		t.Errorf("Expected Memory 8192, got %d", vm.Memory)
// 	}
// 	if vm.Status != virtualmachine.StatusStopped {
// 		t.Errorf("Expected Status 'Stopped', got %d", vm.Status)
// 	}
// }

// func TestRetrieveVirtualmachineByName(t *testing.T) {
// 	db := createTestDB(t)
// 	vmCreated, err := createSampleVirtualMachine(db)
// 	if err != nil {
// 		t.Fatalf("Failed to create virtual machine: %v", err)
// 	}

// 	vmRetrieved, err := virtualmachine.RetrieveVirtualmachineByName("TestVM", db)
// 	if err != nil {
// 		t.Fatalf("Failed to retrieve virtual machine: %v", err)
// 	}

// 	if vmRetrieved.ID != vmCreated.ID {
// 		t.Errorf("Expected ID '%s', got '%s'", vmCreated.ID, vmRetrieved.ID)
// 	}
// 	if vmRetrieved.Name != vmCreated.Name {
// 		t.Errorf("Expected Name '%s', got '%s'", vmCreated.Name, vmRetrieved.Name)
// 	}
// 	if vmRetrieved.CPU != vmCreated.CPU {
// 		t.Errorf("Expected CPU %d, got %d", vmCreated.CPU, vmRetrieved.CPU)
// 	}
// 	if vmRetrieved.Memory != vmCreated.Memory {
// 		t.Errorf("Expected Memory %d, got %d", vmCreated.Memory, vmRetrieved.Memory)
// 	}
// 	if vmRetrieved.Status != vmCreated.Status {
// 		t.Errorf("Expected Status %d, got %d", vmCreated.Status, vmRetrieved.Status)
// 	}
// }

// func TestDestroyVirtualMachine(t *testing.T) {
// 	db := createTestDB(t)
// 	vm, err := createSampleVirtualMachine(db)
// 	if err != nil {
// 		t.Fatalf("Failed to create virtual machine: %v", err)
// 	}

// 	err = virtualmachine.DestroyVirtualMachine(vm, db)
// 	if err != nil {
// 		t.Fatalf("Failed to destroy virtual machine: %v", err)
// 	}

// 	_, err = virtualmachine.RetrieveVirtualmachineByName("TestVM", db)
// 	if err == nil {
// 		t.Errorf("Expected error when retrieving destroyed virtual machine, got none")
// 	}
// }
