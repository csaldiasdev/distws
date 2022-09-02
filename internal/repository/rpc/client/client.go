package client

import (
	"net/rpc"
)

type RepositoryRpcClient struct {
	rpcClient *rpc.Client
}

func (r *RepositoryRpcClient) RaftApplyCommand(data []byte) error {
	var reply bool
	return r.rpcClient.Call("RepositoryRpc.RaftApplyCommandHandler", data, &reply)
}

func NewRepositoryRpcClient(addr string) (*RepositoryRpcClient, error) {
	client, err := rpc.Dial("tcp", addr)

	if err != nil {
		return nil, err
	}

	return &RepositoryRpcClient{
		rpcClient: client,
	}, nil
}
