package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/csaldiasdev/distws/internal/repository/db"
	"github.com/csaldiasdev/distws/internal/repository/model"
	stateRaft "github.com/csaldiasdev/distws/internal/repository/raft"
	"github.com/csaldiasdev/distws/internal/repository/raft/fsm"
	"github.com/csaldiasdev/distws/internal/repository/rpc/client"
	"github.com/csaldiasdev/distws/internal/repository/rpc/server"

	"github.com/google/uuid"
	"github.com/hashicorp/raft"
)

type Repository interface {
	GetByUserId(id uuid.UUID) ([]model.Connection, error)
	Insert(connectionId uuid.UUID, userId uuid.UUID, nodeId uuid.UUID) error
	DeleteConnection(connectionId uuid.UUID) error
	DeleteAllInNode(nodeId uuid.UUID) error
	AddNode(nodeId string, ip string, raftPort uint, rpcPort uint) error
	RemoveNode(nodeId string) error
}

type inMemoryRepository struct {
	nodeId   string
	raftNode *raft.Raft
	memoryDb *db.MemoryDb
	mu       sync.Mutex
	rpcNodes map[string]*client.RepositoryRpcClient
}

func (r *inMemoryRepository) GetByUserId(id uuid.UUID) ([]model.Connection, error) {
	return r.memoryDb.GetByUserId(id)
}

func (r *inMemoryRepository) Insert(connectionId uuid.UUID, userId uuid.UUID, nodeId uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	ev := fsm.ElementValue{
		ConnectionId: connectionId,
		UserId:       userId,
		NodeId:       nodeId,
	}

	evBytes, _ := json.Marshal(ev)

	command := fsm.CommandPayload{
		Operation: fsm.InsertElement,
		Value:     evBytes,
	}

	data, _ := json.Marshal(command)

	_, leaderId := r.raftNode.LeaderWithID()

	if r.nodeId == string(leaderId) {
		r.raftNode.Apply(data, time.Millisecond*500)
		return nil
	}

	rpcNode := r.rpcNodes[string(leaderId)]

	return rpcNode.RaftApplyCommand(data)
}

func (r *inMemoryRepository) DeleteConnection(connectionId uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	ev := fsm.DeleteConnectionValue{
		ConnectionId: connectionId,
	}

	evBytes, _ := json.Marshal(ev)

	command := fsm.CommandPayload{
		Operation: fsm.DeleteElement,
		Value:     evBytes,
	}

	data, _ := json.Marshal(command)

	_, leaderId := r.raftNode.LeaderWithID()

	if r.nodeId == string(leaderId) {
		r.raftNode.Apply(data, time.Millisecond*500)
		return nil
	}

	rpcNode := r.rpcNodes[string(leaderId)]

	return rpcNode.RaftApplyCommand(data)
}

func (r *inMemoryRepository) DeleteAllInNode(nodeId uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	ev := fsm.DeleteAllValue{NodeId: nodeId}

	evBytes, _ := json.Marshal(ev)

	command := fsm.CommandPayload{
		Operation: fsm.DeleteAll,
		Value:     evBytes,
	}

	data, _ := json.Marshal(command)

	_, leaderId := r.raftNode.LeaderWithID()

	if r.nodeId == string(leaderId) {
		r.raftNode.Apply(data, time.Millisecond*500)
		return nil
	}

	rpcNode := r.rpcNodes[string(leaderId)]

	return rpcNode.RaftApplyCommand(data)
}

func (r *inMemoryRepository) AddNode(nodeId string, ip string, raftPort uint, rpcPort uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.raftNode.State() == raft.Leader {
		f := r.raftNode.AddVoter(raft.ServerID(nodeId), raft.ServerAddress(fmt.Sprintf("%s:%d", ip, raftPort)), 0, 0)

		if err := f.Error(); err != nil {
			return err
		}
	}

	cli, err := client.NewRepositoryRpcClient(fmt.Sprintf("%s:%d", ip, rpcPort))

	if err != nil {
		return err
	}

	r.rpcNodes[nodeId] = cli

	return nil
}

func (r *inMemoryRepository) RemoveNode(nodeId string) error {
	if r.raftNode.State() == raft.Leader {
		r.raftNode.DemoteVoter(raft.ServerID(nodeId), 0, 0)
		r.raftNode.RemoveServer(raft.ServerID(nodeId), 0, 0)
	}

	return nil
}

func NewRepository(nodeId string, nodeIp string, raftPort uint, rpcPort uint) (Repository, error) {
	baseFolder := fmt.Sprintf("raft-data/node-%s", nodeId)

	os.MkdirAll(baseFolder, 0777)

	memoryDb := db.NewDb()
	raftNode, err := stateRaft.NewRaftServer(nodeId, baseFolder, nodeIp, raftPort, memoryDb)

	if err != nil {
		return nil, err
	}

	go server.ListenAndServeRepositoryRpc(nodeIp, rpcPort, raftNode)

	return &inMemoryRepository{
		nodeId:   nodeId,
		raftNode: raftNode,
		memoryDb: memoryDb,
		rpcNodes: make(map[string]*client.RepositoryRpcClient),
	}, nil
}
