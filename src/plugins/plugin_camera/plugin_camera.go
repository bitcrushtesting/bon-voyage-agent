package main

import (
	"bon-voyage-agent/models"
	"bon-voyage-agent/shared"
)

const pluginName = "camera"

// Implements interface Plugin
type PluginCamera struct {
	Name string
}

func (p PluginCamera) Init(params any) string {
	shared.Debug("Camera Init Call")
	p.Name = pluginName
	return p.Name
}

func (p PluginCamera) Call(rpcReq models.RPCRequest, response *models.RPCResponse) {
	shared.Debug("Camera Plugin Call")
	switch rpcReq.Method {
	case "camera_information":

	case "camera_start_stream":

	case "camera_stop_stream":

	case "camera_status_stream":
	default:
		response.Error = "Unknown method"
	}
}

func (p PluginCamera) Deinit() {
	shared.Debug("Camera Deinit Call")
}

// Exported symbol
var PluginInstance PluginCamera
