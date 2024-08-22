import simpy
import random
import matplotlib.pyplot as plt
from matplotlib.animation import FuncAnimation


class RandomWalker:
    def __init__(self, env):
        self.env = env
        self.x = 0
        self.y = 0
        self.history = [(self.x, self.y)]
        self.directions = ["east", "west", "north", "south"]

    def walk(self):
        while True:
            direction = random.choice(self.directions)
            if direction == "east":
                self.x += 1
            elif direction == "west":
                self.x -= 1
            elif direction == "north":
                self.y += 1
            elif direction == "south":
                self.y -= 1

            self.history.append((self.x, self.y))
            yield self.env.timeout(1)


def update_plot(frame):
    x, y = zip(*walker.history)
    ax.clear()
    ax.plot(x, y, "b-")
    ax.plot(x[-1], y[-1], "ro")  # Red dot for current position
    ax.set_title(f"Random Walk (Step {frame})")
    ax.set_xlabel("X")
    ax.set_ylabel("Y")
    ax.grid(True)

    # Adjust limits to keep the walker centered
    max_range = max(max(abs(max(x)), abs(min(x))), max(abs(max(y)), abs(min(y))))
    ax.set_xlim(-max_range - 1, max_range + 1)
    ax.set_ylim(-max_range - 1, max_range + 1)


# Set up the simulation
env = simpy.Environment()
walker = RandomWalker(env)

# Set up the plot
fig, ax = plt.subplots(figsize=(8, 8))

# Run the simulation and animation
duration = 100
env.process(walker.walk())

anim = FuncAnimation(fig, update_plot, frames=duration, interval=100, repeat=False)
env.run(until=duration)

plt.show()

# Print final results
print("\nSimulation complete.")
print(f"Final position: {walker.history[-1]}")
print(
    f"Distance from origin: {abs(walker.history[-1][0]) + abs(walker.history[-1][1])}"
)
print(f"Total steps taken: {len(walker.history) - 1}")
