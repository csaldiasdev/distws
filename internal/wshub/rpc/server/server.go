package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"

	"github.com/csaldiasdev/distws/internal/wshub/rpc/model"
)

type HubRpc struct {
	handler func(userId string, message []byte) error
}

func (h *HubRpc) MessageToUserHandler(data []byte, ack *bool) error {

	var msgUser = model.MessageToUser{}

	if err := json.Unmarshal(data, &msgUser); err != nil {
		fmt.Fprint(os.Stderr, "[RPC - MessageToUserHandler] error marshalling data struct")
		return err
	}

	h.handler(msgUser.UserId, msgUser.Message)

	return nil
}

func ListenAndServeHubRpc(ip string, port uint, handler func(userId string, message []byte) error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))

	if err != nil {
		log.Fatal(err)
	}

	inbound, err := net.ListenTCP("tcp", addr)

	if err != nil {
		log.Fatal(err)
	}

	hRpc := HubRpc{
		handler: handler,
	}

	nRpc := rpc.NewServer()

	nRpc.Register(&hRpc)
	nRpc.Accept(inbound)
}
