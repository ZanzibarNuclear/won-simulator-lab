import time

class SimulationComponent:
    def update(self, delta_time):
        pass

class FuelConsumption(SimulationComponent):
    def update(self, delta_time):
        print(f"Updating fuel consumption. Delta time: {delta_time}")
        # Add your fuel consumption logic here

class GeneratorSpeed(SimulationComponent):
    def update(self, delta_time):
        print(f"Updating generator speed. Delta time: {delta_time}")
        # Add your generator speed logic here

class ElectricityDemand(SimulationComponent):
    def update(self, delta_time):
        print(f"Updating electricity demand. Delta time: {delta_time}")
        # Add your electricity demand logic here

class Simulation:
    def __init__(self, time_step):
        self.time_step = time_step
        self.components = []
        self.current_time = 0

    def add_component(self, component):
        self.components.append(component)

    def run(self, duration):
        end_time = self.current_time + duration
        while self.current_time < end_time:
            self.step()
            time.sleep(self.time_step)  # This line simulates the passage of real time

    def step(self):
        print(f"\nTime step: {self.current_time}")
        for component in self.components:
            component.update(self.time_step)
        self.current_time += self.time_step

# Usage
if __name__ == "__main__":
    # Create the simulation with a time step of 1 second
    sim = Simulation(time_step=1)

    # Add components
    sim.add_component(FuelConsumption())
    sim.add_component(GeneratorSpeed())
    sim.add_component(ElectricityDemand())

    # Run the simulation for 10 seconds
    sim.run(duration=10)
