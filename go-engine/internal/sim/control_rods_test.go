package sim

import (
	"fmt"
	"testing"
)

func TestInitiateShutdownBankWithdrawal(t *testing.T) {
	cr := NewControlRods()
	fmt.Printf("Initial Status: %v\n", cr.Status())

	cr.InitiateShutdownBankWithdrawal()

	for i := 0; i < 3; i++ {
		cr.Update()
		status := cr.Status()
		fmt.Printf("Update %d:\n", i+1)
		fmt.Printf("Control Banks: %v\n", status["controlBanks"])
		fmt.Printf("Shutdown Banks: %v\n", status["shutdownBanks"])
		fmt.Printf("Withdraw Shutdown Banks: %v\n", status["withdrawShutdownBanks"])
		fmt.Printf("Insert Shutdown Banks: %v\n", status["insertShutdownBanks"])
		fmt.Println()
	}
}
