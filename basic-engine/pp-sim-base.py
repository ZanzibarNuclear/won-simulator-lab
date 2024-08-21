import simpy
import random

def flip_coin():
    return random.choice(['Heads', 'Tails'])

class FuelConsumption:
    def __init__(self, env, generator_speed):
        self.env = env
        self.action = env.process(self.run())
        self.generator_speed = generator_speed
        self.rateOfConsumption = 1
        self.consumed = 0

    def run(self):
        while True:
            self.rateOfConsumption += 1
            self.consumed += self.rateOfConsumption
            print(f"Updating fuel consumption at {self.consumed} ({self.rateOfConsumption} units per tick)")
            # Add your fuel consumption logic here
            yield self.env.timeout(1)  # Wait for 1 time unit

class GeneratorSpeed:
    def __init__(self, env, fuel_consumption):
        self.env = env
        self.fuel_consumption = fuel_consumption
        self.action = env.process(self.run())
        self.speed = 0
        self.max = 3600

    def run(self):
        while True:
            # max of 3,600 RPM for 2 pole generator OR 1,800 RPM for 4 pole generator
            self.speed = min(self.max, self.fuel_consumption.rateOfConsumption * 500)
            print(f"Updating generator speed at {self.speed}")
            # Add your generator speed logic here
            yield self.env.timeout(1)  # Wait for 1 time unit

class ElectricityDemand:
    def __init__(self, env):
        self.env = env
        self.action = env.process(self.run())
        self.demand = 100

    def run(self):
        while True:
            # assumes random small fluxuation
            increase = random.randint(-5, 5)
            self.demand += increase
            if self.demand < 0:
                self.demand = 0
            print(f"Updating electricity demand changed by {increase} and is at {self.demand}")
            # Add your electricity demand logic here
            yield self.env.timeout(1)  # Wait for 1 time unit

def setup(env):
    fuel_consumption = FuelConsumption(env, generator_speed)
    generator_speed = GeneratorSpeed(env, fuel_consumption)
    electricity_demand = ElectricityDemand(env)

# Create SimPy environment
env = simpy.Environment()

# Setup and start the simulation
setup(env)

# Run the simulation
print("Starting simulation...")
env.run(until=10)  # Run for 10 time units
print("Simulation complete.")
