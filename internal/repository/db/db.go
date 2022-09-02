package db

import (
	"github.com/csaldiasdev/distws/internal/repository/model"

	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
	"github.com/samber/lo"
)

const (
	userNodeTableName = "userNode"
	userIdIndexName   = "userId"
	userIdFieldName   = "UserId"
	nodeIdIndexName   = "nodeId"
	nodeIdFieldName   = "NodeId"
)

type MemoryDb struct {
	database *memdb.MemDB
}

func (m *MemoryDb) GetByUserId(id uuid.UUID) ([]model.UserNode, error) {
	txn := m.database.Txn(false)

	it, err := txn.Get(userNodeTableName, userIdIndexName, id.String())

	if err != nil {
		return nil, err
	}

	result := make([]model.UserNode, 0)

	for obj := it.Next(); obj != nil; obj = it.Next() {
		nodeRep := obj.(model.UserNode)
		result = append(result, nodeRep)
	}

	return result, nil
}

func (m *MemoryDb) GetByNodeId(id uuid.UUID) ([]model.UserNode, error) {
	txn := m.database.Txn(false)

	it, err := txn.Get(userNodeTableName, nodeIdIndexName, id.String())

	if err != nil {
		return nil, err
	}

	result := make([]model.UserNode, 0)

	for obj := it.Next(); obj != nil; obj = it.Next() {
		nodeRep := obj.(model.UserNode)
		result = append(result, nodeRep)
	}

	return result, nil
}

func (m *MemoryDb) Insert(userId uuid.UUID, nodeId uuid.UUID) error {
	txn := m.database.Txn(true)

	err := txn.Insert(userNodeTableName, model.UserNode{
		Id:     uuid.NewString(),
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

func (m *MemoryDb) DeleteUserWithNode(userId uuid.UUID, nodeId uuid.UUID) error {
	users, err := m.GetByUserId(userId)

	if err != nil {
		return err
	}

	txn := m.database.Txn(true)

	listToDelete := lo.Filter(users, func(u model.UserNode, _ int) bool {
		return u.NodeId == nodeId.String()
	})

	for _, v := range listToDelete {
		err := txn.Delete(userNodeTableName, v)
		if err != nil {
			txn.Abort()
			return err
		}
	}

	txn.Commit()
	return nil
}

func (m *MemoryDb) DeleteAllInNode(nodeId uuid.UUID) error {
	txn := m.database.Txn(true)

	_, err := txn.DeleteAll(userNodeTableName, nodeIdIndexName, nodeId.String())

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
			userNodeTableName: {
				Name: userNodeTableName,
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Id"},
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
