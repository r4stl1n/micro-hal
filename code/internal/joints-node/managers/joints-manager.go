package managers

import (
	"encoding/json"
	"github.com/r4stl1n/micro-hal/code/pkg/components"
	"github.com/r4stl1n/micro-hal/code/pkg/consts"
	"github.com/r4stl1n/micro-hal/code/pkg/drivers"
	base "github.com/r4stl1n/micro-hal/code/pkg/drivers/base"
	"github.com/r4stl1n/micro-hal/code/pkg/messages"
	"github.com/r4stl1n/micro-hal/code/pkg/mq"
	"github.com/r4stl1n/micro-hal/code/pkg/structs"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

type JointsManager struct {
	nats *mq.Nats

	baseI2CConn *base.I2C
	pcaDriver   *drivers.PCA9685

	servoMap                   map[string]*components.Servo
	currentJointsPosition      messages.Joints
	defaultServoCalibrationMap structs.ServoCalibrationMap
}

func (jointsManager *JointsManager) Init() (*JointsManager, error) {

	*jointsManager = JointsManager{
		nats: new(mq.Nats).Init(*new(structs.NatsConfig).Defaults()),
	}

	err := jointsManager.connectI2C()

	if err != nil {
		return nil, err
	}

	err = jointsManager.loadServoMap()

	return jointsManager, err
}

func (jointsManager *JointsManager) connectI2C() error {
	// We create a connection to the i2c interface on the raspberry pi
	logrus.Infof("Attempting to connect to the i2c address: /dev/i2c-1")
	i2c, err := new(base.I2C).Init(drivers.DefaultPCA9685Address, "/dev/i2c-1", base.DEFAULT_I2C_ADDRESS)
	if err != nil {
		return err
	}

	jointsManager.baseI2CConn = i2c

	// Next we create the needed driver to connect to the pca9685
	logrus.Info("Creating new connection to pca9685")
	pca, err := new(drivers.PCA9685).Init(i2c, nil)

	if err != nil {
		return err
	}

	jointsManager.pcaDriver = pca

	return nil

}

func (jointsManager *JointsManager) loadServoMap() error {

	// Attempt to load the servo map
	servoMapData, err := ioutil.ReadFile("./ServoMap.json")

	if err != nil {
		return err
	}

	err = json.Unmarshal(servoMapData, &jointsManager.defaultServoCalibrationMap)

	if err != nil {
		return err
	}

	for _, element := range jointsManager.defaultServoCalibrationMap.Servos {
		jointsManager.servoMap[element.Alias] = new(components.Servo).Init(jointsManager.pcaDriver, element.PinId, &components.ServoOptions{
			ActuationRange: element.ActuationRange,
			MinPulse:       element.MinPulse,
			MaxPulse:       element.MaxPulse,
		})

		jointsManager.servoMap[element.Alias].Angle(element.DefaultPosition)
	}

	return err
}

func (jointsManager *JointsManager) connectToNats() error {
	return jointsManager.nats.Connect()
}

func (jointsManager *JointsManager) HandleJointsMessage(joints *messages.Joints) {

}

func (jointsManager *JointsManager) Process() error {

	connectToNatsError := jointsManager.connectToNats()

	if connectToNatsError != nil {
		return connectToNatsError
	}

	receiveChannel := make(chan *[]byte, 100)
	_, bindError := jointsManager.nats.EncodedConn.BindRecvChan(consts.MQPoseSetChannel, receiveChannel)
	if bindError != nil {
		return bindError
	}

	_, bindError = jointsManager.nats.EncodedConn.BindRecvChan(consts.MQJointSetChannel, receiveChannel)
	if bindError != nil {
		return bindError
	}

	logrus.Info("service started waiting for messages")

	for {
		select {
		case receiveData := <-receiveChannel:
			requestMessage := new(messages.Message)
			requestError := requestMessage.Unpack(*receiveData)
			if requestError != nil {
				logrus.Error(requestError)
				continue
			}

			logrus.Debugf("Received Message: %+v", requestMessage)

			switch requestMessage.Type {

			case messages.JointsMessage:
				sMessage := new(messages.Joints)

				unpackError := sMessage.Unpack(requestMessage.Data)
				if unpackError != nil {
					logrus.Error(unpackError)
					continue
				}

				jointsManager.HandleJointsMessage(sMessage)
			default:
				logrus.Errorf("unknown message received %v+", requestMessage)
			}

		}
	}
}
