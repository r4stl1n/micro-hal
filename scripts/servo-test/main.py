import time
import argparse

from adafruit_servokit import ServoKit 

def main():
    
    parser = argparse.ArgumentParser(description='Spot-Hal Servo Test Utility')

    parser.add_argument("--servo", type=int, required=True,            help="The servo you want to test starting at 0")
    parser.add_argument("--min",   type=int, default=500,  nargs="?",  help="minimum impulse value")
    parser.add_argument("--max",   type=int, default=2500, nargs="?",  help="maximum impulse value")
    parser.add_argument("--angle", type=int, default=0,    nargs="?",  help="angle to move to")

    args = parser.parse_args()

    print("Using Values:")
    print("Min Impulse: " + str(args.min))
    print("Max Impulse: " + str(args.max))
    print("Angle: " + str(args.angle) + "\n")

    # Attempt to connect to the PCA9685
    pca = ServoKit(channels=16)

    # Set the pulse width range for the servo
    pca.servo[args.servo].set_pulse_width_range(args.min, args.max)
    
    # Move the servo to the angle
    pca.servo[args.servo].angle = args.angle

    print("Movement completed")

if __name__ == '__main__':
    main()