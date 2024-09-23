package simworks

import (
	"fmt"
	"time"
)

// SimComponent represents a component in the simulation.
// Takes a Simulator to access the clock, environment, and other components.
type SimComponent interface {
	ID() string
	Name() string
	Description() string
	Status() map[string]interface{}
	Update(s *Simulator) (map[string]interface{}, error)
}

type BaseComponent struct {
	id           string
	name         string
	description  string
	latestMoment time.Time
}

func NewBaseComponent(name, description string) *BaseComponent {
	return &BaseComponent{
		id:           GenerateRandomID(12),
		name:         name,
		description:  description,
		latestMoment: time.Time{},
	}
}

func (bc *BaseComponent) ID() string {
	return bc.id
}

func (bc *BaseComponent) Name() string {
	return bc.name
}

func (bc *BaseComponent) Description() string {
	return bc.description
}

func (bc *BaseComponent) LatestMoment() time.Time {
	return bc.latestMoment
}

func (bc *BaseComponent) SetLatestMoment(t time.Time) {
	bc.latestMoment = t
}

func (bc *BaseComponent) Status() map[string]interface{} {
	// Implement the status retrieval logic for the component
	return map[string]interface{}{
		"ID":          bc.ID(),
		"Name":        bc.Name(),
		"Description": bc.Description(),
	}
}

func (bc *BaseComponent) Update(s *Simulator) (map[string]interface{}, error) {
	// Implement the update logic for the component
	fmt.Println("Update logic goes here for component:", bc.Name())
	fmt.Println("Simulator time:", s.Clock.SimNow())
	bc.SetLatestMoment(s.Clock.SimNow())
	return bc.Status(), nil
}
