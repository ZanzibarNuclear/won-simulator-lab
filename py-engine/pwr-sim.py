import simpy # type: ignore

def process(env):
    while True:
        print(f"Time: {env.now}")
        yield env.timeout(1)

def run_simulation(duration):
    env = simpy.Environment()
    env.process(process(env))
    env.run(until=duration)

if __name__ == "__main__":
    simulation_duration = 10
    run_simulation(simulation_duration)
