# Micro Hal - Spot Micro bot Variation

This take on the spot-micro attempts to pull pieces together from various different spot projects into one complete package.

The goal of this project is to not only pull the best pieces from the various spot projects but also to have a single repo that contains everything you need to get up and running smoothly.

Each project that pieces are taken from will be documented in this README so anyone who is interested can visit the parent projects.

All current code samples, models, and rewrites used for this project will be posted within this repo.

Note: This is a in-progress working spot micro bot written in GO

## Project Structure

Due to the amount of different data involved the project will be stored in various directories. This goes as follows:

```bash
micro-hal
├── code                       - Code repository for working with micro hal
│   ├── cmd                    - Micro hal Applications
│   │    └── hal-utilities     - Utility Application
│   │
├── parts                      - Parts Necessary
│    └── 3d-prints             
│       ├── freecad            - Original Free Cad Files
│       └── stl                - 3d Printable STL Files
│           └── templates      - Template STLS to modify for different hardware
│
└── sim                        - Sim folder containing everything required to run in gazebo
    ├── meshes                 
    │   └── stl                - STL files used to render the bot in gazebo
    ├── plugin                 - Gazebo plugin source directory
    │   └── include
    ├── urdf                   - The urdf and xacro file for the micro-hal bot
    └── world                  - Folder containing the world defenition file for gazebo


```

## Reference && Community

Deok-Yeon Kim (KDY0523) - Original design of the sport micro: https://www.thingiverse.com/thing:3445283

Michael Kubina (michaelkubina) - Easy print redesign of the spot micro: https://github.com/michaelkubina/SpotMicroESP32/

CHVMP - Original quad kinematics and node code - https://github.com/chvmp/champ