package connection

import (
	"bon-voyage-agent/models"
	"encoding/json"
	"fmt"
	"strings"
)

type Handler func(models.RPCRequest, *models.RPCResponse)

type Route struct {
	handler Handler
}

func NewRoute() *Route {
	return &Route{}
}

type Router struct {
	NotFoundHandler Handler
	namedRoutes     map[string]*Route
}

func NewRouter() *Router {
	r := Router{namedRoutes: make(map[string]*Route)}
	r.HandleFunc("agent", AgentCall)
	return &r
}

func (r *Router) HandleFunc(name string, f func(request models.RPCRequest, response *models.RPCResponse)) *Route {
	route := &Route{handler: f}
	r.namedRoutes[name] = route
	return route
}

func (r *Router) ParseMessage(data []byte) ([]byte, error) {

	var request models.RPCRequest
	err := json.Unmarshal(data, &request)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}
	if request.Jsonrpc != "2.0" {
		return nil, fmt.Errorf("jsonrpc field not '2.0'")
	}

	response := models.RPCResponse{
		Jsonrpc: "2.0",
		ID:      request.ID,
	}
	if len(strings.Split(request.Method, "_")) == 0 {
		return nil, fmt.Errorf("method malformed")
	}
	routeName := strings.Split(request.Method, "_")[0]

	route := r.namedRoutes[routeName]
	route.handler(request, &response)

	resBytes, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %v", err)
	}
	return resBytes, nil
}
