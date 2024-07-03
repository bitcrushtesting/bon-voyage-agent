package connection

import (
	"bon-voyage-agent/models"
	"encoding/json"
	"fmt"
)

const pluginName = "serial"

// Implements interface Plugin
type PluginAgent struct {
	Name         string
	Information  models.AgentInformation
	Capabilities []string
	Plugins      []string
}

func (p PluginAgent) Init(params any) string {

	switch v := params.(type) {
	case models.AgentInformation:
		PluginInstance.Information = v
	default:
		fmt.Printf("Invalid plugin params")
	}
	p.Name = pluginName
	return p.Name
}

func AgentCall(request models.RPCRequest, response *models.RPCResponse) {

	switch request.Method {
	case "agent_get_information":
		i, err := json.Marshal(PluginInstance.Information)
		if err != nil {
			response.Error = err.Error()
			return
		}
		response.Result = string(i)
	case "agent_set_name":
		var name models.SetNameMethod
		err := json.Unmarshal(request.Params, &name)
		if err != nil {
			response.Error = err.Error()
		}
		PluginInstance.Information.Name = name.Name

	case "agent_get_capabilities":
		i, err := json.Marshal(PluginInstance.Capabilities)
		if err != nil {
			response.Error = err.Error()
			return
		}
		response.Result = string(i)
	case "agent_get_plugins":
		i, err := json.Marshal(PluginInstance.Plugins)
		if err != nil {
			response.Error = err.Error()
			return
		}
		response.Result = string(i)
	default:
		response.Error = "Unknown method"
	}
}

// Exported symbol
var PluginInstance PluginAgent
