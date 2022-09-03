package fsm

import "github.com/google/uuid"

type CommandOperation int

const (
	InsertElement CommandOperation = iota + 1
	DeleteElement
	DeleteAll
)

type CommandPayload struct {
	Operation CommandOperation
	Value     []byte
}

type ElementValue struct {
	ConnectionId uuid.UUID
	UserId       uuid.UUID
	NodeId       uuid.UUID
}

type DeleteConnectionValue struct {
	ConnectionId uuid.UUID
}

type DeleteAllValue struct {
	NodeId uuid.UUID
}
