package pwr

import (
	"testing"

	"worldofnuclear.com/internal/simworks"
)

func TestNewSteamGenerator(t *testing.T) {
	sg := NewSteamGenerator("SG1", "Test Steam Generator", nil, nil)
	s := simworks.NewSimulator("Test Simulator", "Testing steam generator")
	s.AddComponent(sg)

	// FIXME: make sure this complains.
	if sg.Name() != "SG1" {
		t.Errorf("expected name to be SG1, got %s", sg.Name())
	}

	if sg.Description() != "Test Steam Generator" {
		t.Errorf("expected description to be 'Test Steam Generator', got %s", sg.Description())
	}

	_, err := sg.Update(s)
	if err == nil {
		t.Errorf("Should complain about not having a primary or secondary loop\n")
	}
}

func TestSteamGeneratorStatus(t *testing.T) {
	primaryLoop := &PrimaryLoop{}
	secondaryLoop := &SecondaryLoop{}
	sg := NewSteamGenerator("SG1", "Test Steam Generator", primaryLoop, secondaryLoop)
	s := simworks.NewSimulator("Test Simulator", "Testing steam generator")
	s.AddComponent(primaryLoop)
	s.AddComponent(secondaryLoop)
	s.AddComponent(sg)

	// FIXME: try to set up normal conditions.
	status := sg.Status()

	if status["primaryInletTemp"] != sg.PrimaryInletTemp() {
		t.Errorf("expected primaryInletTemp to be %f, got %f", sg.PrimaryInletTemp(), status["primaryInletTemp"])
	}

	if status["primaryOutletTemp"] != sg.PrimaryOutletTemp() {
		t.Errorf("expected primaryOutletTemp to be %f, got %f", sg.PrimaryOutletTemp(), status["primaryOutletTemp"])
	}

	if status["secondaryInletTemp"] != sg.SecondaryInletTemp() {
		t.Errorf("expected secondaryInletTemp to be %f, got %f", sg.SecondaryInletTemp(), status["secondaryInletTemp"])
	}

	if status["secondaryOutletTemp"] != sg.SecondaryOutletTemp() {
		t.Errorf("expected secondaryOutletTemp to be %f, got %f", sg.SecondaryOutletTemp(), status["secondaryOutletTemp"])
	}

	if status["heatTransferRate"] != sg.HeatTransferRate() {
		t.Errorf("expected heatTransferRate to be %f, got %f", sg.HeatTransferRate(), status["heatTransferRate"])
	}

	if status["steamFlowRate"] != sg.SteamFlowRate() {
		t.Errorf("expected steamFlowRate to be %f, got %f", sg.SteamFlowRate(), status["steamFlowRate"])
	}
}

func TestSteamGeneratorUpdate(t *testing.T) {
	primaryLoop := &PrimaryLoop{}
	secondaryLoop := &SecondaryLoop{}
	sg := NewSteamGenerator("SG1", "Test Steam Generator", primaryLoop, secondaryLoop)
	simulator := &simworks.Simulator{}

	status, err := sg.Update(simulator)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if status["primaryInletTemp"] != sg.PrimaryInletTemp() {
		t.Errorf("expected primaryInletTemp to be %f, got %f", sg.PrimaryInletTemp(), status["primaryInletTemp"])
	}

	if status["primaryOutletTemp"] != sg.PrimaryOutletTemp() {
		t.Errorf("expected primaryOutletTemp to be %f, got %f", sg.PrimaryOutletTemp(), status["primaryOutletTemp"])
	}

	if status["secondaryInletTemp"] != sg.SecondaryInletTemp() {
		t.Errorf("expected secondaryInletTemp to be %f, got %f", sg.SecondaryInletTemp(), status["secondaryInletTemp"])
	}

	if status["secondaryOutletTemp"] != sg.SecondaryOutletTemp() {
		t.Errorf("expected secondaryOutletTemp to be %f, got %f", sg.SecondaryOutletTemp(), status["secondaryOutletTemp"])
	}

	if status["heatTransferRate"] != sg.HeatTransferRate() {
		t.Errorf("expected heatTransferRate to be %f, got %f", sg.HeatTransferRate(), status["heatTransferRate"])
	}

	if status["steamFlowRate"] != sg.SteamFlowRate() {
		t.Errorf("expected steamFlowRate to be %f, got %f", sg.SteamFlowRate(), status["steamFlowRate"])
	}
}
