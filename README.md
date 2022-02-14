# Spot Micro - Hal Edition

This take on the spot-micro attempts to pull pieces together from various different spot projects into one complete package.

The goal of this project is to not only pull the best pieces from the various spot projects but also to have a single repo that contains everything you need to get up and running smoothly.

Each project that pieces are taken from will be documented in this README so anyone who is interested can visit the parent projects.

All current code samples, models, and rewrites used for this project will be posted within this repo.

## Project Structure

Due to the amount of different data involved the project will be stored in various directories. This goes as follows:


```bash
spot-micro-hal
├── parts                     - Parts Necessary
│ └── 3d-prints               - Directory for the 3d prints
│     ├── freecad             - Original Free Cad Files
│     └── stl                 - 3d Printable STL Files
│         └── templates       - Template STLS to modify for different hardware
└── scripts                   - Utility scripts for various things
    └── servo-test            - Servo test script for calibration
```

## Reference && Community

Deok-Yeon Kim (KDY0523) - Original design of the sport micro: https://www.thingiverse.com/thing:3445283

Michael Kubina (michaelkubina) - Easy print redesign of the spot micro: https://github.com/michaelkubina/SpotMicroESP32/

CHVMP - Original quad kinematics and node code - https://github.com/chvmp/champ