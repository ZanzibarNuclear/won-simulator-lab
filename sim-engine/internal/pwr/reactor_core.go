package pwr

import (
	"fmt"

	"worldofnuclear.com/internal/simworks"
)

type ReactorCore struct {
	simworks.BaseComponent
}

func NewReactorCore(name string, description string) *ReactorCore {
	return &ReactorCore{
		BaseComponent: *simworks.NewBaseComponent(name, description),
	}
}

func (r *ReactorCore) Status() map[string]interface{} {
	return map[string]interface{}{
		"about": r.BaseComponent.Status(),
	}
}

func (r *ReactorCore) Print() {
	fmt.Printf("=> Reactor Core\n")
	r.BaseComponent.Print()
}

func (r *ReactorCore) Update(s *simworks.Simulator) (map[string]interface{}, error) {
	return r.Status(), nil
}
