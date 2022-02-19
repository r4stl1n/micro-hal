package structs

// ServoCalibrationItem stores servo calibration information
type ServoCalibrationItem struct {
	Alias           string
	PinId           int
	ActuationRange  int
	MinPulse        float32
	MaxPulse        float32
	DefaultPosition int
}

// Map of servo calibration information
type ServoCalibrationMap struct {
	Servos []ServoCalibrationItem
}
