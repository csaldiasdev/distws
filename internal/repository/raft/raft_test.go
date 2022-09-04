package raft

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/csaldiasdev/distws/internal/repository/db"
	"github.com/csaldiasdev/distws/internal/repository/raft/fsm"

	"github.com/google/uuid"
	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/assert"
)

func TestCreateRaftServer(t *testing.T) {
	serverOneId := uuid.NewString()
	serverOneDataFolder := fmt.Sprintf("node-%s", serverOneId)

	os.MkdirAll(serverOneDataFolder, 0777)

	defer os.Remove(serverOneDataFolder)

	db := db.NewDb()
	svr, err := NewRaftServer(serverOneId, serverOneDataFolder, "localhost", 9001, db)

	time.Sleep(time.Second * 2)

	state := svr.State()

	assert.Equal(t, raft.Leader, state)

	assert.NoError(t, err)
	assert.NotNil(t, svr)

	os.RemoveAll(serverOneDataFolder)
}

func TestCreateRaftCluster(t *testing.T) {
	baseFolder := "raft-data"

	serverOneId := uuid.NewString()
	serverOneDataFolder := fmt.Sprintf("%s/node-%s", baseFolder, serverOneId)

	serverTwoId := uuid.NewString()
	serverTwoDataFolder := fmt.Sprintf("%s/node-%s", baseFolder, serverTwoId)

	serverThreeId := uuid.NewString()
	serverThreeDataFolder := fmt.Sprintf("%s/node-%s", baseFolder, serverThreeId)

	serverFourId := uuid.NewString()
	serverFourDataFolder := fmt.Sprintf("%s/node-%s", baseFolder, serverFourId)

	serverFiveId := uuid.NewString()
	serverFiveDataFolder := fmt.Sprintf("%s/node-%s", baseFolder, serverFiveId)

	os.MkdirAll(serverOneDataFolder, 0777)
	os.MkdirAll(serverTwoDataFolder, 0777)
	os.MkdirAll(serverThreeDataFolder, 0777)
	os.MkdirAll(serverFourDataFolder, 0777)
	os.MkdirAll(serverFiveDataFolder, 0777)

	defer os.RemoveAll(baseFolder)

	serverOneRepo := db.NewDb()
	serverTwoRepo := db.NewDb()
	serverThreeRepo := db.NewDb()
	serverFourRepo := db.NewDb()
	serverFiveRepo := db.NewDb()

	serverOne, _ := NewRaftServer(serverOneId, serverOneDataFolder, "localhost", 9001, serverOneRepo)

	time.Sleep(time.Second * 2)

	NewRaftServer(serverTwoId, serverTwoDataFolder, "localhost", 9002, serverTwoRepo)
	NewRaftServer(serverThreeId, serverThreeDataFolder, "localhost", 9003, serverThreeRepo)
	NewRaftServer(serverFourId, serverFourDataFolder, "localhost", 9004, serverFourRepo)
	NewRaftServer(serverFiveId, serverFiveDataFolder, "localhost", 9005, serverFiveRepo)

	f := serverOne.AddVoter(raft.ServerID(serverTwoId), raft.ServerAddress("localhost:9002"), 0, 0)

	assert.NoError(t, f.Error())

	f = serverOne.AddVoter(raft.ServerID(serverThreeId), raft.ServerAddress("localhost:9003"), 0, 0)

	assert.NoError(t, f.Error())

	f = serverOne.AddVoter(raft.ServerID(serverFourId), raft.ServerAddress("localhost:9004"), 0, 0)

	assert.NoError(t, f.Error())

	f = serverOne.AddVoter(raft.ServerID(serverFiveId), raft.ServerAddress("localhost:9005"), 0, 0)

	assert.NoError(t, f.Error())
}

func TestReplicateData(t *testing.T) {
	baseFolder := "raft-data"

	serverOneId := uuid.NewString()
	serverOneDataFolder := fmt.Sprintf("%s/node-%s", baseFolder, serverOneId)

	serverTwoId := uuid.NewString()
	serverTwoDataFolder := fmt.Sprintf("%s/node-%s", baseFolder, serverTwoId)

	serverThreeId := uuid.NewString()
	serverThreeDataFolder := fmt.Sprintf("%s/node-%s", baseFolder, serverThreeId)

	serverFourId := uuid.NewString()
	serverFourDataFolder := fmt.Sprintf("%s/node-%s", baseFolder, serverFourId)

	serverFiveId := uuid.NewString()
	serverFiveDataFolder := fmt.Sprintf("%s/node-%s", baseFolder, serverFiveId)

	os.MkdirAll(serverOneDataFolder, 0777)
	os.MkdirAll(serverTwoDataFolder, 0777)
	os.MkdirAll(serverThreeDataFolder, 0777)
	os.MkdirAll(serverFourDataFolder, 0777)
	os.MkdirAll(serverFiveDataFolder, 0777)

	defer os.RemoveAll(baseFolder)

	serverOneRepo := db.NewDb()
	serverTwoRepo := db.NewDb()
	serverThreeRepo := db.NewDb()
	serverFourRepo := db.NewDb()
	serverFiveRepo := db.NewDb()

	serverOne, _ := NewRaftServer(serverOneId, serverOneDataFolder, "localhost", 9001, serverOneRepo)

	time.Sleep(time.Second * 2)

	NewRaftServer(serverTwoId, serverTwoDataFolder, "localhost", 9002, serverTwoRepo)
	NewRaftServer(serverThreeId, serverThreeDataFolder, "localhost", 9003, serverThreeRepo)
	NewRaftServer(serverFourId, serverFourDataFolder, "localhost", 9004, serverFourRepo)
	NewRaftServer(serverFiveId, serverFiveDataFolder, "localhost", 9005, serverFiveRepo)

	f := serverOne.AddVoter(raft.ServerID(serverTwoId), raft.ServerAddress("localhost:9002"), 0, 0)

	assert.NoError(t, f.Error())

	f = serverOne.AddVoter(raft.ServerID(serverThreeId), raft.ServerAddress("localhost:9003"), 0, 0)

	assert.NoError(t, f.Error())

	f = serverOne.AddVoter(raft.ServerID(serverFourId), raft.ServerAddress("localhost:9004"), 0, 0)

	assert.NoError(t, f.Error())

	f = serverOne.AddVoter(raft.ServerID(serverFiveId), raft.ServerAddress("localhost:9005"), 0, 0)

	assert.NoError(t, f.Error())

	ev := fsm.ElementValue{
		UserId: uuid.New(),
		NodeId: uuid.MustParse(serverOneId),
	}

	evBytes, _ := json.Marshal(ev)

	command := fsm.CommandPayload{
		Operation: fsm.InsertElement,
		Value:     evBytes,
	}

	data, _ := json.Marshal(command)

	af := serverOne.Apply(data, time.Millisecond*500)

	assert.NoError(t, af.Error())

	assert.NotNil(t, af.Response())

	time.Sleep(time.Second * 1)

	results, err := serverFiveRepo.GetByUserId(ev.UserId)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(results))

	results, err = serverFourRepo.GetByUserId(ev.UserId)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(results))

	results, err = serverThreeRepo.GetByUserId(ev.UserId)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(results))

	results, err = serverTwoRepo.GetByUserId(ev.UserId)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(results))

	command = fsm.CommandPayload{
		Operation: fsm.DeleteElement,
		Value:     evBytes,
	}

	data, _ = json.Marshal(command)

	af = serverOne.Apply(data, time.Millisecond*500)

	assert.NoError(t, af.Error())

	assert.NotNil(t, af.Response())

	time.Sleep(time.Second * 1)

	results, err = serverFiveRepo.GetByUserId(ev.UserId)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(results))

	results, err = serverFourRepo.GetByUserId(ev.UserId)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(results))

	results, err = serverThreeRepo.GetByUserId(ev.UserId)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(results))

	results, err = serverTwoRepo.GetByUserId(ev.UserId)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(results))
}
