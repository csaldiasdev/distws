package raft

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/csaldiasdev/distws/internal/repository/db"
	"github.com/csaldiasdev/distws/internal/repository/raft/fsm"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

const (
	// The maxPool controls how many connections we will pool.
	maxPool = 5

	// The timeout is used to apply I/O deadlines. For InstallSnapshot, we multiply
	// the timeout by (SnapshotSize / TimeoutScale).
	// https://github.com/hashicorp/raft/blob/v1.1.2/net_transport.go#L177-L181
	tcpTimeout = 10 * time.Second

	// The `retain` parameter controls how many
	// snapshots are retained. Must be at least 1.
	raftSnapShotRetain = 2

	// raftLogCacheSize is the maximum number of logs to cache in-memory.
	// This is used to reduce disk I/O for the recently committed entries.
	raftLogCacheSize = 512
)

func NewRaftServer(svrId string, dataDir string, port uint, d *db.MemoryDb) (*raft.Raft, error) {
	raftConf := raft.DefaultConfig()
	raftConf.LocalID = raft.ServerID(svrId)
	raftConf.SnapshotThreshold = 1024

	fsmStore := fsm.NewFsm(d)

	fullPathBolt := filepath.Join(dataDir, "bolt")

	store, err := raftboltdb.NewBoltStore(fullPathBolt)

	if err != nil {
		return nil, err
	}

	cacheStore, err := raft.NewLogCache(raftLogCacheSize, store)
	if err != nil {
		return nil, err
	}

	snapshotStore, err := raft.NewFileSnapshotStore(dataDir, raftSnapShotRetain, os.Stdout)
	if err != nil {
		return nil, err
	}

	addr := fmt.Sprintf("localhost:%d", port)

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	transport, err := raft.NewTCPTransport(addr, tcpAddr, maxPool, tcpTimeout, os.Stdout)
	if err != nil {
		return nil, err
	}

	raftServer, err := raft.NewRaft(raftConf, fsmStore, cacheStore, store, snapshotStore, transport)
	if err != nil {
		return nil, err
	}

	// always start single server as a leader
	configuration := raft.Configuration{
		Servers: []raft.Server{
			{
				ID:      raft.ServerID(svrId),
				Address: transport.LocalAddr(),
			},
		},
	}

	raftServer.BootstrapCluster(configuration)

	return raftServer, nil
}
