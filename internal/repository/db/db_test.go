package db

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInsertData(t *testing.T) {
	database := NewDb()

	dummyUserId := uuid.New()
	dummyNodeId := uuid.New()

	assert.NoError(t, database.Insert(uuid.New(), dummyUserId, dummyNodeId))

	found, err := database.GetByUserId(dummyUserId)

	assert.NoError(t, err)
	assert.NotEmpty(t, len(found))
}

func TestDeleteData(t *testing.T) {
	database := NewDb()

	dummyUserId := uuid.New()
	dummyNodeId := uuid.New()
	dummyConnectionId := uuid.New()

	assert.NoError(t, database.Insert(dummyConnectionId, dummyUserId, dummyNodeId))

	found, err := database.GetByUserId(dummyUserId)

	assert.NoError(t, err)
	assert.NotEmpty(t, len(found))

	assert.NoError(t, database.DeleteConnection(dummyConnectionId))

	found, err = database.GetByUserId(dummyUserId)

	assert.NoError(t, err)
	assert.Empty(t, len(found))
}
