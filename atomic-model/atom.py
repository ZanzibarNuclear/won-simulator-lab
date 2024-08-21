import matplotlib.pyplot as plt
import numpy as np


class Atom:
    def __init__(self, name, protons, neutrons, electrons):
        self.name = name
        self.protons = protons
        self.neutrons = neutrons
        self.electrons = electrons

    def draw(self):
        fig, ax = plt.subplots()

        # Draw nucleus
        nucleus = plt.Circle((0, 0), 0.5, color="r", fill=True)
        ax.add_artist(nucleus)

        # Draw electron shells
        shells = [1, 2, 3]  # Simplified model with 3 shells
        for i, shell in enumerate(shells):
            circle = plt.Circle((0, 0), (i + 1) * 1.5, color="b", fill=False)
            ax.add_artist(circle)

        # Place electrons
        electrons_placed = 0
        for i, shell in enumerate(shells):
            if electrons_placed >= self.electrons:
                break
            electrons_in_shell = min(
                2 * (i + 1) ** 2, self.electrons - electrons_placed
            )
            angles = np.linspace(0, 2 * np.pi, electrons_in_shell, endpoint=False)
            x = (i + 1) * 1.5 * np.cos(angles)
            y = (i + 1) * 1.5 * np.sin(angles)
            ax.plot(x, y, "bo", markersize=5)
            electrons_placed += electrons_in_shell

        ax.set_xlim(-5, 5)
        ax.set_ylim(-5, 5)
        ax.set_aspect("equal")
        ax.axis("off")
        plt.title(f"{self.name}: {self.protons}p, {self.neutrons}n, {self.electrons}e")
        plt.show()


elements = [
    ["Hydrogen", 1, 1, 1],
    ["Helium", 2, 2, 2],
    ["Lithium", 3, 4, 3],
    ["Beryllium", 4, 5, 4],
    ["Boron", 5, 5, 5],
    ["Carbon", 6, 6, 6],
    ["Nitrogen", 7, 7, 7],
    ["Oxygen", 8, 8, 8],
    ["Fluorine", 9, 10, 9],
    ["Neon", 10, 10, 10],
    ["Sodium", 11, 12, 11],
    ["Magnesium", 12, 12, 12],
    ["Aluminum", 13, 14, 13],
    ["Silicon", 14, 14, 14],
    ["Phosphorus", 15, 16, 15],
    ["Sulfur", 16, 16, 16],
    ["Chlorine", 17, 18, 17],
    ["Argon", 18, 22, 18],
]


def drawElements():
    for element in elements:
        name, protons, neutrons, electrons = element
        element = Atom(
            name=name, protons=protons, neutrons=neutrons, electrons=electrons
        )
        element.draw()
