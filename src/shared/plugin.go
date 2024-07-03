package shared

import "bon-voyage-agent/models"

type Plugin interface {
	Init(params any) string //returns name
	Call(rpcReq models.RPCRequest, response *models.RPCResponse)
	Deinit()
}
