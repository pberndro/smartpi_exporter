package ade7878

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

type Config struct {
	// [base]
	Serial               string
	Name                 string
	LogLevel             log.Level
	MetricsListenAddress string

	// [location]
	Lat float64
	Lng float64

	// [device]
	I2CDevice            string
	PowerFrequency       float64
	Samplerate           int
	Integrator           bool
	CTType               map[Phase]string
	CTTypePrimaryCurrent map[Phase]int
	CurrentDirection     map[Phase]bool
	MeasureCurrent       map[Phase]bool
	MeasureVoltage       map[Phase]bool
	Voltage              map[Phase]float64

	// [calibration]
	CalibrationfactorI map[Phase]float64
	CalibrationfactorU map[Phase]float64
}

var cfg *ini.File
var err error

func (p *Config) ReadParameterFromFile() {

	cfg, err = ini.Load("smartpi.ini")
	if err != nil {
		panic(err)
	}

	// [base]
	p.Serial = cfg.Section("base").Key("serial").String()
	p.Name = cfg.Section("base").Key("name").MustString("smartpi")

	// Handle logging levels
	p.LogLevel, err = log.ParseLevel(cfg.Section("base").Key("loglevel").MustString("info"))
	if err != nil {
		panic(err)
	}

	p.MetricsListenAddress = cfg.Section("base").Key("metrics_listen_address").MustString(":9246")

	// [location]
	p.Lat = cfg.Section("location").Key("lat").MustFloat64(52.3667)
	p.Lng = cfg.Section("location").Key("lng").MustFloat64(9.7167)

	// [device]
	p.I2CDevice = cfg.Section("device").Key("i2c_device").MustString("/dev/i2c-1")
	p.PowerFrequency = cfg.Section("device").Key("power_frequency").MustFloat64(50)
	p.Samplerate = cfg.Section("device").Key("samplerate").MustInt(1)
	//p.Integrator = cfg.Section("device").Key("integrator").MustBool(false)

	p.CTType = make(map[Phase]string)
	p.CTType[PhaseA] = cfg.Section("device").Key("ct_type_1").MustString("YHDC_SCT013")
	p.CTType[PhaseB] = cfg.Section("device").Key("ct_type_2").MustString("YHDC_SCT013")
	p.CTType[PhaseC] = cfg.Section("device").Key("ct_type_3").MustString("YHDC_SCT013")
	p.CTType[PhaseN] = cfg.Section("device").Key("ct_type_4").MustString("YHDC_SCT013")

	p.CTTypePrimaryCurrent = make(map[Phase]int)
	//p.CTTypePrimaryCurrent[PhaseA] = cfg.Section("device").Key("ct_type_1_primary_current").MustInt(100)
	//p.CTTypePrimaryCurrent[PhaseB] = cfg.Section("device").Key("ct_type_2_primary_current").MustInt(100)
	//p.CTTypePrimaryCurrent[PhaseC] = cfg.Section("device").Key("ct_type_3_primary_current").MustInt(100)
	//p.CTTypePrimaryCurrent[PhaseN] = cfg.Section("device").Key("ct_type_4_primary_current").MustInt(100)

	p.CTTypePrimaryCurrent[PhaseA] = cfg.Section("device").Key("ct_type_1_primary_current").MustInt(200)
	p.CTTypePrimaryCurrent[PhaseB] = cfg.Section("device").Key("ct_type_2_primary_current").MustInt(200)
	p.CTTypePrimaryCurrent[PhaseC] = cfg.Section("device").Key("ct_type_3_primary_current").MustInt(200)
	p.CTTypePrimaryCurrent[PhaseN] = cfg.Section("device").Key("ct_type_4_primary_current").MustInt(200)

	p.CurrentDirection = make(map[Phase]bool)
	p.CurrentDirection[PhaseA] = cfg.Section("device").Key("change_current_direction_1").MustBool(false)
	p.CurrentDirection[PhaseB] = cfg.Section("device").Key("change_current_direction_2").MustBool(false)
	p.CurrentDirection[PhaseC] = cfg.Section("device").Key("change_current_direction_3").MustBool(false)
	p.CurrentDirection[PhaseN] = cfg.Section("device").Key("change_current_direction_4").MustBool(false)

	p.MeasureCurrent = make(map[Phase]bool)
	p.MeasureCurrent[PhaseA] = cfg.Section("device").Key("measure_current_1").MustBool(true)
	p.MeasureCurrent[PhaseB] = cfg.Section("device").Key("measure_current_2").MustBool(true)
	p.MeasureCurrent[PhaseC] = cfg.Section("device").Key("measure_current_3").MustBool(true)
	p.MeasureCurrent[PhaseN] = cfg.Section("device").Key("measure_current_4").MustBool(true)

	p.MeasureVoltage = make(map[Phase]bool)
	p.MeasureVoltage[PhaseA] = cfg.Section("device").Key("measure_voltage_1").MustBool(false)
	p.MeasureVoltage[PhaseB] = cfg.Section("device").Key("measure_voltage_2").MustBool(false)
	p.MeasureVoltage[PhaseC] = cfg.Section("device").Key("measure_voltage_3").MustBool(false)
	p.MeasureCurrent[PhaseN] = cfg.Section("device").Key("measure_current_4").MustBool(false)


	//p.MeasureVoltage[PhaseA] = cfg.Section("device").Key("measure_voltage_1").MustBool(true)
	//p.MeasureVoltage[PhaseB] = cfg.Section("device").Key("measure_voltage_2").MustBool(true)
	//p.MeasureVoltage[PhaseC] = cfg.Section("device").Key("measure_voltage_3").MustBool(true)

	p.Voltage = make(map[Phase]float64)
	p.Voltage[PhaseA] = cfg.Section("device").Key("voltage_1").MustFloat64(230)
	p.Voltage[PhaseB] = cfg.Section("device").Key("voltage_2").MustFloat64(230)
	p.Voltage[PhaseC] = cfg.Section("device").Key("voltage_3").MustFloat64(230)

	// [calibration]
	p.CalibrationfactorI = make(map[Phase]float64)
	p.CalibrationfactorI[PhaseA] = cfg.Section("calibration").Key("calibrationfactorI_1").MustFloat64(1)
	p.CalibrationfactorI[PhaseB] = cfg.Section("calibration").Key("calibrationfactorI_2").MustFloat64(1)
	p.CalibrationfactorI[PhaseC] = cfg.Section("calibration").Key("calibrationfactorI_3").MustFloat64(1)
	p.CalibrationfactorI[PhaseN] = cfg.Section("calibration").Key("calibrationfactorI_4").MustFloat64(1)

	p.CalibrationfactorU = make(map[Phase]float64)
	p.CalibrationfactorU[PhaseA] = cfg.Section("calibration").Key("calibrationfactorU_1").MustFloat64(1)
	p.CalibrationfactorU[PhaseB] = cfg.Section("calibration").Key("calibrationfactorU_2").MustFloat64(1)
	p.CalibrationfactorU[PhaseC] = cfg.Section("calibration").Key("calibrationfactorU_3").MustFloat64(1)

}

func NewConfig() *Config {
	t := new(Config)
	t.ReadParameterFromFile()
	return t
}
