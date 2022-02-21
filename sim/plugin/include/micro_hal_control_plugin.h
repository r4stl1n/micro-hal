#ifndef PLUGIN_MICRO_HAL_CONTROL_PLUGIN_H
#define PLUGIN_MICRO_HAL_CONTROL_PLUGIN_H


#include <memory>
#include <string>
#include <vector>

#include "gazebo/common/common.hh"
#include "gazebo/physics/Model.hh"

namespace micro_hal_control {

    class MicroHalControlPlugin : public gazebo::ModelPlugin {
    public:
        MicroHalControlPlugin();

        ~MicroHalControlPlugin() override;

        // Overloaded Gazebo entry point
        void Load(gazebo::physics::ModelPtr parent, sdf::ElementPtr sdf) override;

    };
}  // namespace micro_hal_control

#endif  // PLUGIN_MICRO_HAL_CONTROL_PLUGIN_H