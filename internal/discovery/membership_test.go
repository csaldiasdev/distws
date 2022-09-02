package discovery

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMembershipJoins(t *testing.T) {

	joinFunc := func(np *NodeProps) {}
	leaveFunc := func(np *NodeProps) {}

	nodeOne, err := NewMembership(10000, NodeProps{
		GrpcPort: 10001,
		RaftPort: 10003,
	}, joinFunc, leaveFunc)

	require.NoError(t, err)

	_, err = NewMembership(20000, NodeProps{
		GrpcPort: 20001,
		RaftPort: 20003,
	}, joinFunc, leaveFunc, "127.0.0.1:10000")

	require.NoError(t, err)

	_, err = NewMembership(30000, NodeProps{
		GrpcPort: 30001,
		RaftPort: 30003,
	}, joinFunc, leaveFunc, "127.0.0.1:10000")

	require.NoError(t, err)

	_, err = NewMembership(40000, NodeProps{
		GrpcPort: 40001,
		RaftPort: 40003,
	}, joinFunc, leaveFunc, "127.0.0.1:10000")

	require.NoError(t, err)

	require.Equal(t, 4, len(nodeOne.serf.Members()))
}

func TestMembershipLeaveGracefully(t *testing.T) {

	joinFunc := func(np *NodeProps) {}
	leaveFunc := func(np *NodeProps) {}

	nodeOne, err := NewMembership(10000, NodeProps{
		GrpcPort: 10001,
		RaftPort: 10003,
	}, joinFunc, leaveFunc)

	require.NoError(t, err)

	nodeTwo, err := NewMembership(20000, NodeProps{
		GrpcPort: 20001,
		RaftPort: 20003,
	}, joinFunc, leaveFunc, "127.0.0.1:10000")

	require.NoError(t, err)

	_, err = NewMembership(30000, NodeProps{
		GrpcPort: 30001,
		RaftPort: 30003,
	}, joinFunc, leaveFunc, "127.0.0.1:10000")

	require.NoError(t, err)

	nodeFour, err := NewMembership(40000, NodeProps{
		GrpcPort: 40001,
		RaftPort: 40003,
	}, joinFunc, leaveFunc, "127.0.0.1:10000")

	require.NoError(t, err)

	require.Equal(t, 4, len(nodeOne.serf.Members()))

	nodeFour.Leave()

	time.Sleep(time.Second * 1)

	require.Equal(t, 3, len(nodeTwo.AliveMembers()))
}
