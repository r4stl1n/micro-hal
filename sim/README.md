# Sim

This folder contains everything needed to run micro-hal in simulation using gazebo.

TODO: Fill this out alot better


## Setup

Below are the rough steps for getting the simulation running

### Building the micro-hal plugin
In order for our simulation to work and be able to communicate with the micro-hal stack we first have to build our plugin so it can work.
Aside from needing the gazebo-dev libraries we will also needs a nats and msgpack library

TODO: Add the lib names here for debian

To build the project simply do the following.

```bash
cd plugin
mkdir build
cd build
cmake ../
make
```

### Exporting the plugin path

In order for gazebo to be able to find our plugin we need to export the plugin location. In a terminal do the following

```bash
export GAZEBO_PLUGIN_PATH=${GAZEBO_PLUGIN_PATH}:<Full Path to the cmake-build-debug>
```

### GZServer
To start the gzserver context with our selected world do the following.

```bash
export GAZEBO_PLUGIN_PATH=${GAZEBO_PLUGIN_PATH}:<Full Path to the cmake-build-debug>
gzserver world/hal.world --verbose
```

### GZClient
To start the client run the following

```bash
export GAZEBO_PLUGIN_PATH=${GAZEBO_PLUGIN_PATH}:<Full Path to the cmake-build-debug>
gzclient --verbose


```

### Spawn robot
Finally all that is left is to spawn our robot. Gazebo can load a urdf file directly so no need to convert it over to a sdf

```bash
export GAZEBO_PLUGIN_PATH=${GAZEBO_PLUGIN_PATH}:<Full Path to the cmake-build-debug>
gz model --spawn-file=./urdf/micro-hal.urdf --model-name=micro-hal
```