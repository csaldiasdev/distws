package rpc

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/csaldiasdev/distws/internal/repository/db"
	stateRaft "github.com/csaldiasdev/distws/internal/repository/raft"
	"github.com/csaldiasdev/distws/internal/repository/raft/fsm"
	"github.com/csaldiasdev/distws/internal/repository/rpc/client"
	"github.com/csaldiasdev/distws/internal/repository/rpc/server"

	"github.com/google/uuid"
	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/assert"
)

func TestSendDataViaRpc(t *testing.T) {
	baseFolder := "raft-data"

	memoryDb := db.NewDb()

	nodeId := uuid.New()

	os.MkdirAll(baseFolder, 0777)

	r, _ := stateRaft.NewRaftServer(nodeId.String(), "localhost", baseFolder, 5000, memoryDb)

	time.Sleep(time.Second * 2)

	state := r.State()

	assert.Equal(t, raft.Leader, state)

	defer os.RemoveAll(baseFolder)

	go server.ListenAndServeRepositoryRpc("localhost", 5001, r)

	time.Sleep(time.Millisecond * 500)

	cli, err := client.NewRepositoryRpcClient("localhost:5001")

	assert.NoError(t, err)

	userId := uuid.New()

	ev := fsm.ElementValue{
		ConnectionId: uuid.New(),
		UserId:       userId,
		NodeId:       nodeId,
	}

	evBytes, _ := json.Marshal(ev)

	command := fsm.CommandPayload{
		Operation: fsm.InsertElement,
		Value:     evBytes,
	}

	data, _ := json.Marshal(command)

	cli.RaftApplyCommand(data)

	found, _ := memoryDb.GetByUserId(userId)

	assert.NotEmpty(t, len(found))
}
