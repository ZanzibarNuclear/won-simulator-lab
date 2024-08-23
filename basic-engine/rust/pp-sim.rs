use std::cmp;
use rand::Rng;
use std::thread;
use std::time::Duration;

struct FuelConsumption {
    rate_of_increase: i32,
    rate_of_consumption: i32,
    consumed: i32,
}

struct GeneratorSpeed {
    speed: i32,
    max_speed: i32,
}

struct ElectricityDemand {
    demand: i32,
}

struct Simulation {
    fuel_consumption: FuelConsumption,
    generator_speed: GeneratorSpeed,
    electricity_demand: ElectricityDemand,
    iterations: i32,
}

impl FuelConsumption {
    fn update(&mut self, generator_running_at_max: bool) {
        if !generator_running_at_max {
            self.rate_of_consumption += self.rate_of_increase;
        }
        self.consumed += self.rate_of_consumption;
        println!("Consuming {} units of fuel. Total consumption at {}.", self.rate_of_consumption, self.consumed);
    }
}

impl GeneratorSpeed {
    fn update(&mut self, fuel_consumption_rate: i32) {
        self.speed = cmp::min(self.max_speed, fuel_consumption_rate * 500);
        println!("Adjusting generator to run at {} RPMs", self.speed);
    }

    fn is_running_at_max(&self) -> bool {
        self.speed >= self.max_speed
    }
}

impl ElectricityDemand {
    fn update(&mut self) {
        let mut rng = rand::thread_rng();
        let increase: i32 = rng.gen_range(-5..=5);
        self.demand += increase;
        self.demand = cmp::max(0, self.demand);
        println!("Updating electricity demand changed by {} and is at {} kW", increase, self.demand);
    }
}

impl Simulation {
    fn new(iterations: i32) -> Self {
        Simulation {
            fuel_consumption: FuelConsumption {
                rate_of_increase: 1,
                rate_of_consumption: 1,
                consumed: 0,
            },
            generator_speed: GeneratorSpeed {
                speed: 0,
                max_speed: 3600,
            },
            electricity_demand: ElectricityDemand {
                demand: 100,
            },
            iterations,
        }
    }

    fn run_iteration(&mut self, iteration: i32) {
        println!("\n--- Iteration {} ---", iteration);
        self.fuel_consumption.update(self.generator_speed.is_running_at_max());
        self.generator_speed.update(self.fuel_consumption.rate_of_consumption);
        self.electricity_demand.update();
    }

    fn run(&mut self) {
        for i in 0..self.iterations {
            self.run_iteration(i + 1);
            thread::sleep(Duration::from_secs(1)); // Simulate passage of time
        }
    }
}

fn main() {
    println!("Starting simulation...");
    let mut simulation = Simulation::new(20);
    simulation.run();
    println!("Simulation complete.");
}