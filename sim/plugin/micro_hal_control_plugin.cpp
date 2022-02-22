#include <string>
#include <memory>
#include <utility>
#include <vector>

#include "include/micro_hal_control_plugin.h"

using namespace std::chrono_literals;

namespace micro_hal_control {

    MicroHalControlPlugin::MicroHalControlPlugin() {
        printf("Hello World!\n");
    };

    MicroHalControlPlugin::~MicroHalControlPlugin() = default;

    // Overloaded Gazebo entry point
    void MicroHalControlPlugin::Load(gazebo::physics::ModelPtr parent, sdf::ElementPtr sdf) {
    }

    // Register this plugin with the simulator
    GZ_REGISTER_MODEL_PLUGIN(MicroHalControlPlugin)
}  // namespace micro_hal_control
