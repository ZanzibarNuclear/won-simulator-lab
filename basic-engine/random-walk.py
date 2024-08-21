import random


def random_walk(steps):
    x, y = 0, 0  # Starting position
    history = [(x, y)]
    directions = ["east", "west", "north", "south"]

    for _ in range(steps):
        direction = random.choice(directions)
        if direction == "east":
            x += 1
        elif direction == "west":
            x -= 1
        elif direction == "north":
            y += 1
        elif direction == "south":
            y -= 1

        history.append((x, y))

    return history


# Run the simulation
steps = 100
result = random_walk(steps)

# Print the results
print(f"Random Walk of {steps} steps:")
for i, position in enumerate(result):
    print(f"Step {i}: {position}")

print(f"\nFinal position: {result[-1]}")
print(f"Distance from origin: {abs(result[-1][0]) + abs(result[-1][1])}")
