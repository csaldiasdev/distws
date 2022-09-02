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

	assert.NoError(t, database.Insert(dummyUserId, dummyNodeId))

	found, err := database.GetByUserId(dummyUserId)

	assert.NoError(t, err)
	assert.NotEmpty(t, len(found))
}

func TestDeleteData(t *testing.T) {
	database := NewDb()

	dummyUserId := uuid.New()
	dummyNodeId := uuid.New()

	assert.NoError(t, database.Insert(dummyUserId, dummyNodeId))

	found, err := database.GetByUserId(dummyUserId)

	assert.NoError(t, err)
	assert.NotEmpty(t, len(found))

	assert.NoError(t, database.DeleteUserWithNode(dummyUserId, dummyNodeId))

	found, err = database.GetByUserId(dummyUserId)

	assert.NoError(t, err)
	assert.Empty(t, len(found))
}
