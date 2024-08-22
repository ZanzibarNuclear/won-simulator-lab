import simpy
import random

def flip_coin():
    return random.choice(['Heads', 'Tails'])

class FuelConsumption:
    def __init__(self, env):
        self.env = env
        self.action = env.process(self.run())
        self.rateOfIncrease = 1
        self.rateOfConsumption = 1
        self.consumed = 0
        self.generator_speed = None # will be set later

    def run(self):
        while True:
            # Add your fuel consumption logic here
            if not self.generator_speed.isRunningAtMax():
                self.rateOfConsumption += self.rateOfIncrease
            self.consumed += self.rateOfConsumption
            print(f"Consuming {self.rateOfConsumption} units of fuel. Total consumption at {self.consumed}.")
            yield self.env.timeout(1)  # Wait for 1 time unit

class GeneratorSpeed:
    def __init__(self, env):
        self.env = env
        self.action = env.process(self.run())
        self.speed = 0
        self.max_speed = 3600 # max of 3,600 RPM for 2 pole generator OR 1,800 RPM for 4 pole generator
        self.fuel_consumption = None # will be set later

    def run(self):
        while True:
            # Add your generator speed logic here
            self.speed = min(self.max_speed, self.fuel_consumption.rateOfConsumption * 500)
            print(f"Adjusting generator to run at {self.speed} RPMs")
            yield self.env.timeout(1)  # Wait for 1 time unit
    
    def isRunningAtMax(self):
        return self.speed >= self.max_speed

class ElectricityDemand:
    def __init__(self, env):
        self.env = env
        self.action = env.process(self.run())
        self.demand = 100

    def run(self):
        while True:
            # Add your electricity demand logic here
            # assumes random small fluxuation
            increase = random.randint(-5, 5)
            self.demand += increase
            if self.demand < 0:
                self.demand = 0
            print(f"Updating electricity demand changed by {increase} and is at {self.demand} kW")
            yield self.env.timeout(1)  # Wait for 1 time unit

def setup(env):
    fuel_consumption = FuelConsumption(env)
    generator_speed = GeneratorSpeed(env)

    fuel_consumption.generator_speed = generator_speed
    generator_speed.fuel_consumption = fuel_consumption

    electricity_demand = ElectricityDemand(env)

# Create SimPy environment
env = simpy.Environment()

# Setup and start the simulation
setup(env)

# Run the simulation
print("Starting simulation...")
env.run(until=20)
print("Simulation complete.")
