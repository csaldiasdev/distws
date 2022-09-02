package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"

	"github.com/hashicorp/raft"
)

type RepositoryRpc struct {
	raft *raft.Raft
}

func (r *RepositoryRpc) RaftApplyCommandHandler(data []byte, ack *bool) error {
	if r.raft.State() != raft.Leader {
		return errors.New("node isn't leader")
	}

	af := r.raft.Apply(data, time.Millisecond*500)

	err := af.Error()

	if err != nil {
		return err
	}

	return nil
}

func ListenAndServeRepositoryRpc(port uint, raft *raft.Raft) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatal(err)
	}

	inbound, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	rRpc := RepositoryRpc{
		raft: raft,
	}

	rpc.Register(&rRpc)
	rpc.Accept(inbound)
}
