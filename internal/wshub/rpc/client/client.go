package client

import (
	"encoding/json"
	"net/rpc"

	"github.com/csaldiasdev/distws/internal/wshub/rpc/model"
)

type HubRpcClient struct {
	rpcClient *rpc.Client
}

func (h *HubRpcClient) MessageToUser(userId string, message []byte) error {
	var reply bool
	data, _ := json.Marshal(model.MessageToUser{UserId: userId, Message: message})
	return h.rpcClient.Call("HubRpc.MessageToUserHandler", data, &reply)
}

func NewHubRpcClient(addr string) (*HubRpcClient, error) {
	client, err := rpc.Dial("tcp", addr)

	if err != nil {
		return nil, err
	}

	return &HubRpcClient{rpcClient: client}, nil
}
