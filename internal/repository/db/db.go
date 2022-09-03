package db

import (
	"github.com/csaldiasdev/distws/internal/repository/model"

	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
)

const (
	connectionTableName = "connection"
	idIndexName         = "id"
	idIndexFieldName    = "Id"
	userIdIndexName     = "userId"
	userIdFieldName     = "UserId"
	nodeIdIndexName     = "nodeId"
	nodeIdFieldName     = "NodeId"
)

type MemoryDb struct {
	database *memdb.MemDB
}

func (m *MemoryDb) GetByUserId(id uuid.UUID) ([]model.Connection, error) {
	txn := m.database.Txn(false)

	it, err := txn.Get(connectionTableName, userIdIndexName, id.String())

	if err != nil {
		return nil, err
	}

	result := make([]model.Connection, 0)

	for obj := it.Next(); obj != nil; obj = it.Next() {
		nodeRep := obj.(model.Connection)
		result = append(result, nodeRep)
	}

	return result, nil
}

func (m *MemoryDb) Insert(connectionId uuid.UUID, userId uuid.UUID, nodeId uuid.UUID) error {
	txn := m.database.Txn(true)

	err := txn.Insert(connectionTableName, model.Connection{
		Id:     connectionId.String(),
		UserId: userId.String(),
		NodeId: nodeId.String(),
	})

	if err != nil {
		txn.Abort()
		return err
	}

	txn.Commit()
	return nil
}

func (m *MemoryDb) DeleteConnection(connectionId uuid.UUID) error {
	txn := m.database.Txn(true)

	_, err := txn.DeleteAll(connectionTableName, idIndexName, connectionId.String())

	if err != nil {
		txn.Abort()
		return err
	}

	txn.Commit()
	return nil
}

func (m *MemoryDb) DeleteAllInNode(nodeId uuid.UUID) error {
	txn := m.database.Txn(true)

	_, err := txn.DeleteAll(connectionTableName, nodeIdIndexName, nodeId.String())

	if err != nil {
		txn.Abort()
		return err
	}

	txn.Commit()
	return nil
}

func NewDb() *MemoryDb {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			connectionTableName: {
				Name: connectionTableName,
				Indexes: map[string]*memdb.IndexSchema{
					idIndexName: {
						Name:    idIndexName,
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: idIndexFieldName},
					},
					userIdIndexName: {
						Name:    userIdIndexName,
						Indexer: &memdb.StringFieldIndex{Field: userIdFieldName},
					},
					nodeIdIndexName: {
						Name:    nodeIdIndexName,
						Indexer: &memdb.StringFieldIndex{Field: nodeIdFieldName},
					},
				},
			},
		},
	}

	memdb, _ := memdb.NewMemDB(schema)

	return &MemoryDb{database: memdb}
}
