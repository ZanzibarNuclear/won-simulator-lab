package sim

import (
	"fmt"
	"testing"
)

func TestInitiateShutdownBankWithdrawal(t *testing.T) {
	cr := NewControlRods()

	cr.Update() // do nothing

	// Ensure all shutdown banks are initially fully inserted
	if !cr.ShutdownBanksFullyInserted() {
		t.Errorf("Expected shutdown banks to be fully inserted initially")
	}

	cr.InitiateShutdownBankWithdrawal()

	// Update until shutdown banks are fully withdrawn
	for i := 0; i < 1000 && !cr.ShutdownBanksFullyWithdrawn(); i++ {
		cr.Update()
	}

	fmt.Printf("Status after withdrawal: %+v\n", cr.Status())

	// Check if shutdown banks are fully withdrawn
	if !cr.ShutdownBanksFullyWithdrawn() {
		t.Errorf("Expected shutdown banks to be fully withdrawn after updates")
	}

	// Verify that control banks were not affected
	for i, bank := range cr.controlBanks {
		if bank.Position() != 0 {
			t.Errorf("Control bank %d position changed unexpectedly: got %d, want 0", i+1, bank.Position())
		}
	}

	// Now test shutdown bank insertion
	cr.InitiateShutdownBankInsertion()

	// Update until shutdown banks are fully inserted
	for i := 0; i < 1000 && !cr.ShutdownBanksFullyInserted(); i++ {
		cr.Update()
	}

	fmt.Printf("Status after insertion: %+v\n", cr.Status())

	// Check if shutdown banks are fully inserted
	if !cr.ShutdownBanksFullyInserted() {
		t.Errorf("Expected shutdown banks to be fully inserted after updates")
	}

	// Verify that all shutdown banks are at position 0
	for i, bank := range cr.shutdownBanks {
		if bank.Position() != 0 {
			t.Errorf("Shutdown bank %d not fully inserted: got position %d, want 0", i+1, bank.Position())
		}
	}

	// Verify that control banks were still not affected
	for i, bank := range cr.controlBanks {
		if bank.Position() != 0 {
			t.Errorf("Control bank %d position changed unexpectedly: got %d, want 0", i+1, bank.Position())
		}
	}
}

func TestShutdownBankPartialWithdrawalAndInsertion(t *testing.T) {
	cr := NewControlRods()

	// Initiate shutdown bank withdrawal
	cr.InitiateShutdownBankWithdrawal()

	// Update a few times to partially withdraw
	for i := 0; i < 12; i++ {
		cr.Update()
	}

	// Check partial withdrawal
	if cr.ShutdownBanksFullyWithdrawn() {
		t.Errorf("Wanted partial withdrawal, but went all the way")
	}

	fmt.Printf("Status after partial withdrawal: %+v\n", cr.Status())

	// Initiate shutdown bank insertion
	cr.InitiateShutdownBankInsertion()

	// Update until fully inserted
	for i := 0; i < 100 && !cr.ShutdownBanksFullyInserted(); i++ {
		cr.Update()
	}

	fmt.Printf("Status after insertion: %+v\n", cr.Status())

	// Check if shutdown banks are fully inserted
	if !cr.ShutdownBanksFullyInserted() {
		t.Errorf("Expected shutdown banks to be fully inserted after updates")
	}

	// Verify that all shutdown banks are at position 0
	for i, bank := range cr.shutdownBanks {
		if bank.Position() != 0 {
			t.Errorf("Shutdown bank %d not fully inserted: got position %d, want 0", i+1, bank.Position())
		}
	}

	// Verify that control banks were not affected
	for i, bank := range cr.controlBanks {
		if bank.Position() != 0 {
			t.Errorf("Control bank %d position changed unexpectedly: got %d, want 0", i+1, bank.Position())
		}
	}
}

func TestAdjustControlBankPosition(t *testing.T) {
	cr := NewControlRods()

	// Test raising control rods in group 3 to position 50
	cr.AdjustControlBankPosition(3, 50)
	cr.AdjustControlBankPosition(4, 75)

	// Update until the target position is reached or a timeout occurs
	for i := 0; i < 100 && cr.controlBanks[2].Position() != 50; i++ {
		cr.Update()
	}

	fmt.Printf("Status after adjustment to control bank 3: %+v\n", cr.Status())

	// Check if the control bank 3 reached the target position
	if cr.controlBanks[2].Position() != 50 {
		t.Errorf("Control bank 3 did not reach target position: got %d, want 50", cr.controlBanks[2].Position())
	}

	if cr.controlBanks[3].Position() != 75 {
		t.Errorf("Control bank 4 did not reach target position: got %d, want 75", cr.controlBanks[3].Position())
	}

	// Verify that other control banks were not affected
	for i, bank := range cr.controlBanks {
		if i != 2 && i != 3 && bank.Position() != 0 {
			t.Errorf("Control bank %d position changed unexpectedly: got %d, want 0", i+1, bank.Position())
		}
	}

	// Verify that shutdown banks were not affected
	for i, bank := range cr.shutdownBanks {
		if bank.Position() != 0 {
			t.Errorf("Shutdown bank %d position changed unexpectedly: got %d, want 0", i+1, bank.Position())
		}
	}
}
