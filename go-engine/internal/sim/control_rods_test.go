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

	fmt.Printf("Status: %+v\n", cr.Status())

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
}
