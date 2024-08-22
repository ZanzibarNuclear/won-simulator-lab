# World of Nuclear Simulator Lab

Development area for experimentation and assembly of nuclear energy simulations.

The first project is to simulate a PWR fission power plant.

For some warm-up exercises, let's build some prototypes, mini-simulations, to figure out a solid yet flexible core architecture for the simulations.

- Random walk: featuring a core process, use of a random number generator, and graphics that show what is happening.
- Power plant simulator base (pp-sim-base): meant to be a high-level prototype. That means we start with some of the bigger loops, using simplifying assumptions. As we like, we can iterate on the assumptions to make them more realistic. We can also add more of the interdependent parts.

Ultimately, we want to model things like:

- how much fuel is consumed to heat how much water to what temperature per unit time
- how hot the water get => how much steam is generated per unit time
- how fast the turbine spins => how much electricity the generator creates
- how much water is flowing in the primary loop, how much in the steam loop
- demand curves for power: one that assumes typical metro usage patterns, one with "renewables" that work intermittently

Eventually, we want to add models that mimic, say, reactivity of the core, that change with control rod position and boron levels. This could be tied to real calculations that operators make to keep things stable and producing an optimal level of power.

## Code

The first prototype uses Python. Config parameters could come from a file. We can run the Python code as a serverless function. We will also need a UI, probably using NuxtJS like the rest of the website. Context for a given sim should be stored in a database and loaded to continue when the user returns. A JSON column in Postgres should be fine, or maybe store as files in S3.

## As a game

A player can provision their power plant (an instance of the simulator), and let it run for some number of days or weeks, sim time. That would generate interesting data to view. Or maybe we just track a cumulative score based on how the sim ran.
