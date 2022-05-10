# Raspberry Pi 4 Model B setup

Here is how we will setup our raspberry pi to work with our current features and setup.

TODO: Clean this up and make it more detailed


## Installation

Using the Raspberry PI Imager flash a copy of ubuntu server.

Run the following commands to install dependencies

```bash
sudo apt-get update
sudo apt-get upgrade

sudo apt-get install build-essential
sudo apt-get install git
sudo apt-get install i2c-tools
sudo apt-get install python-smbus
sudo apt-get install screen

sudo snap install go --classic
```

### Change i2c from root only to work from our user

In order to communicate with our various i2c devices we need to change the i2c group permissions

```bash
sudo groupadd i2c
sudo chown :i2c /dev/i2c-1
sudo chmod g+rw /dev/i2c-1
sudo usermod -aG i2c <YourUserHere>
su root
echo 'KERNEL=="i2c-[0-9]*", GROUP="i2c"' >> /etc/udev/rules.d/10-local_i2c_group.rules
sudo reboot now
```

Once logged in we can verify it all works by doing

```bash
i2cdetect -y 1
```

## Advanced Configuration

These changes are designed to be ran on startup rather than modifying the boot config patterns

### Enable additional i2c ports

In order to connect all our devices we need to expand the i2c ports. We could have all our ports inline with one another and reference them by address that way.
But i found it easier to just enable the extra ports we need.

To do this we need to modify the /boot/firmware/usercfg.txt to include new dtoverlays to turn specific gpio pins on the pi to I2C pins.

Open the /boot/firmware/usercfg.txt using nano

```bash
sudo nano /boot/firmware/usercfg.txt
```

Now we are just going to enable the additional i2c pins we need for the MPU6050 and the SSD1306 we use the original i2c pins to power the PCA9685
Add the following lines to the bottom of the usercfg.txt

```bash
dtoverlay=i2c5
dtoverlay=i2c6
```

Then hit Ctrl+O and then Ctrl+X

What this will do is now enable i2c on the following raspberry pi ports

SDA5 = 32
SCL5 = 33

SDA6 = 15
SCL6 = 16


Now just reboot and the new i2c ports should be avalible
