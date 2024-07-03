package main

import (
	"bon-voyage-agent/models"
	"bon-voyage-agent/shared"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

const pluginName = "serial"

// Implements interface Plugin
type PluginSerial struct {
	Name string
}

func (p PluginSerial) Init(params any) string {

	shared.Debug("Serial Plugin Init")
	p.Name = pluginName
	return p.Name
}

func (p PluginSerial) Call(rpcReq models.RPCRequest, response *models.RPCResponse) {

	shared.Debug("Serial Plugin Call")

	switch rpcReq.Method {
	case "serial_get_ports":
		data, _ := json.Marshal(serialPortList())
		response.Result = string(data)

	case "serial_open":
		response.Result = "Serial port opened"
		serialPort.open()

	case "serial_close":
		response.Result = "Serial port closed"
		serialPort.close()

	case "serial_set_baudrate":
		var params map[string]int
		if err := json.Unmarshal(rpcReq.Params, &params); err != nil {
			response.Error = "Invalid params"
			return
		}
		serialPort.Mode.BaudRate = params["baudrate"]
		serialPort.port.SetMode(&serialPort.Mode)
		response.Result = "Baudrate set to " + strconv.Itoa(params["baudrate"])

	case "serial_set_databits":
		var params map[string]int
		if err := json.Unmarshal(rpcReq.Params, &params); err != nil {
			response.Error = "Invalid params"
			return
		}
		serialPort.Mode.DataBits = params["databits"]
		serialPort.port.SetMode(&serialPort.Mode)
		response.Result = "Data bits set to " + strconv.Itoa(params["databits"])

	case "serial_set_parity":
		var params map[string]string
		if err := json.Unmarshal(rpcReq.Params, &params); err != nil {
			response.Error = "Invalid params"
			return
		}
		if err := serialPort.setParity(params["parity"]); err != nil {
			response.Error = "Invalid params"
			return
		}
		response.Result = "Parity set to " + params["parity"]

	case "serial_set_stopbits":
		var params map[string]int
		if err := json.Unmarshal(rpcReq.Params, &params); err != nil {
			response.Error = "Invalid params"
			return
		}
		serialPort.setStopBit(params["stopbits"])
		response.Result = "Stopbits set to " + strconv.Itoa(params["stopbits"])

	default:
		response.Error = "Unknown method"
	}
}

func (p PluginSerial) Deinit() {
	shared.Debug("Serial Plugin Deinit")
}

// Exported symbol
var PluginInstance PluginSerial

func serialPortList() (list []string) {

	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
		return
	}
	if len(ports) == 0 {
		fmt.Println("No serial ports found!")
		return
	}
	fmt.Println("Found serial ports:")
	for _, port := range ports {
		fmt.Printf("    - %s\n", port.Name)
		list = append(list, port.Name)
		if port.IsUSB {
			fmt.Printf("      * USB ID     %s:%s\n", port.VID, port.PID)
			fmt.Printf("      * USB serial %s\n", port.SerialNumber)
		}
	}
	return
}

var serialPort SerialPort

type SerialPort struct {
	PortName string
	Mode     serial.Mode
	port     serial.Port
}

func (s *SerialPort) setParity(p string) error {

	switch strings.ToLower(p) {
	case "odd":
		s.Mode.Parity = serial.OddParity
	case "even":
		s.Mode.Parity = serial.EvenParity
	case "no":
		s.Mode.Parity = serial.NoParity
	default:
		return fmt.Errorf("unknown parity %s", p)
	}
	return nil
}

func (s *SerialPort) setStopBit(b int) error {

	switch b {
	case 1:
		s.Mode.StopBits = serial.OneStopBit
	case 2:
		s.Mode.StopBits = serial.TwoStopBits
	default:
		return fmt.Errorf("stop bit count %d not supported", b)
	}
	return nil
}

func (s *SerialPort) open() (err error) {
	s.port, err = serial.Open(s.PortName, &s.Mode)
	return
}

func (s *SerialPort) close() error {
	return s.port.Close()
}
