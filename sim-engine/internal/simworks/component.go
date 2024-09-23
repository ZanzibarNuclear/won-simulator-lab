package simworks

import (
	"fmt"
)

// SimComponent represents a component in the simulation.
// Takes a Simulator to access the clock, environment, and other components.
type SimComponent interface {
	Update(s *Simulator) (map[string]interface{}, error)
	GetStatus() map[string]interface{}
}

type BaseComponent struct {
	ID   string
	Name string
}

func NewBaseComponent(name string) *BaseComponent {
	return &BaseComponent{
		ID:   GenerateRandomID(12),
		Name: name,
	}
}

func (bc *BaseComponent) Update(s *Simulator) (map[string]interface{}, error) {
	// Implement the update logic for the component
	fmt.Println("Update logic goes here for component:", bc.Name)
	return nil, nil
}

func (bc *BaseComponent) GetStatus() map[string]interface{} {
	// Implement the status retrieval logic for the component
	return map[string]interface{}{
		"ID":   bc.ID,
		"Name": bc.Name,
	}
}
