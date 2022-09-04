package repository

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDataReplication(t *testing.T) {
	defer os.RemoveAll("raft-data")

	repoOneId := uuid.New()
	repoOne, err := NewRepository(repoOneId.String(), "localhost", 5001, 5002)

	assert.NoError(t, err)

	repoTwoId := uuid.New()
	repoTwo, err := NewRepository(repoTwoId.String(), "localhost", 6001, 6002)

	assert.NoError(t, err)

	repoThreeId := uuid.New()
	repoThree, err := NewRepository(repoThreeId.String(), "localhost", 7001, 7002)

	assert.NoError(t, err)

	time.Sleep(time.Second * 2)

	assert.NoError(t, repoOne.AddNode(repoTwoId.String(), "localhost", 6001, 6002))
	assert.NoError(t, repoOne.AddNode(repoThreeId.String(), "localhost", 7001, 7002))

	assert.NoError(t, repoTwo.AddNode(repoOneId.String(), "localhost", 5001, 5002))
	assert.NoError(t, repoTwo.AddNode(repoThreeId.String(), "localhost", 7001, 7002))

	assert.NoError(t, repoThree.AddNode(repoTwoId.String(), "localhost", 6001, 6002))
	assert.NoError(t, repoThree.AddNode(repoOneId.String(), "localhost", 5001, 5002))

	time.Sleep(time.Second * 2)

	userId := uuid.New()

	connId := uuid.New()

	assert.NoError(t, repoOne.Insert(uuid.New(), userId, repoOneId))
	assert.NoError(t, repoTwo.Insert(connId, userId, repoTwoId))
	assert.NoError(t, repoThree.Insert(uuid.New(), userId, repoThreeId))

	time.Sleep(time.Second * 2)

	foundOne, _ := repoOne.GetByUserId(userId)
	foundTwo, _ := repoTwo.GetByUserId(userId)
	founoThree, _ := repoThree.GetByUserId(userId)

	assert.Equal(t, 3, len(foundOne))
	assert.Equal(t, 3, len(foundTwo))
	assert.Equal(t, 3, len(founoThree))

	assert.NoError(t, repoThree.DeleteConnection(connId))

	time.Sleep(time.Millisecond * 500)

	foundOne, _ = repoOne.GetByUserId(userId)
	foundTwo, _ = repoTwo.GetByUserId(userId)
	founoThree, _ = repoThree.GetByUserId(userId)

	assert.Equal(t, 2, len(foundOne))
	assert.Equal(t, 2, len(foundTwo))
	assert.Equal(t, 2, len(founoThree))
}
