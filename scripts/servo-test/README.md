# Servo Test Script

This script provides a easy way to test servo PWM signal values as well has adjust the angle for building.
We take advantage of the Raspberry Pi connected and wired up properly to the PCA9685 

The default values for the DS3218MG servos are as follows:

```bash
Minimum Impulse: 500
Maximum Impulse: 2500

Minimum Angle: 0
Maximum Angle: 180
```

Note: Please ensure you have gone through the raspberry pi setup instructions

## Dependencies

```bash
virtualenv .env
source .env/bin/activate
pip install -r requirements.txt
```

## Usage

You can modify the min and max impulse values to see at what point the servo stops responding. This can be done by simply setting the angle to 0 then modifying the --min value and trying to increase the angle by 1 if it works then you know that min value works. Do the inverse to calculate max impulse


```bash
usage: main.py [-h] --servo SERVO [--min [MIN]] [--max [MAX]] [--angle [ANGLE]]

Spot-Hal Servo Test Utility

optional arguments:
  -h, --help       show this help message and exit
  --servo SERVO    The servo you want to test starting at 0
  --min   [MIN]    minimum impulse value
  --max   [MAX]    maximum impulse value
  --angle [ANGLE]  angle to move to
```

