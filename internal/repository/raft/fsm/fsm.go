package fsm

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/csaldiasdev/distws/internal/repository/db"

	"github.com/hashicorp/raft"
)

type fsm struct {
	db *db.MemoryDb
}

func (i fsm) Apply(log *raft.Log) interface{} {
	switch log.Type {
	case raft.LogCommand:
		var command = CommandPayload{}
		if err := json.Unmarshal(log.Data, &command); err != nil {
			fmt.Fprintf(os.Stderr, "error marshalling store payload %s\n", err.Error())
			return nil
		}

		switch command.Operation {
		case InsertElement:

			var ieValue = ElementValue{}

			if err := json.Unmarshal(command.Value, &ieValue); err != nil {
				fmt.Fprint(os.Stderr, "[InsertElement] error marshalling command value struct")
				return nil
			}

			return &ApplyResponse{
				Error: i.db.Insert(ieValue.UserId, ieValue.NodeId),
				Data:  ieValue,
			}

		case DeleteElement:

			var deValue = ElementValue{}

			if err := json.Unmarshal(command.Value, &deValue); err != nil {
				fmt.Fprint(os.Stderr, "[DeleteElement] error marshalling command value struct")
				return nil
			}

			return &ApplyResponse{
				Error: i.db.DeleteUserWithNode(deValue.UserId, deValue.NodeId),
				Data:  deValue,
			}

		case DeleteAll:

			var daValue = DeleteAllValue{}

			if err := json.Unmarshal(command.Value, &daValue); err != nil {
				fmt.Fprint(os.Stderr, "[DeleteAll] error marshalling command value struct")
				return nil
			}

			return &ApplyResponse{
				Error: i.db.DeleteAllInNode(daValue.NodeId),
				Data:  daValue,
			}
		}
	}

	_, _ = fmt.Fprintf(os.Stderr, "not raft log command type\n")

	return nil
}

func (i fsm) Snapshot() (raft.FSMSnapshot, error) {
	return newSnapshotNoop()
}

func (i fsm) Restore(rClose io.ReadCloser) error {
	defer func() {
		if err := rClose.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "[FINALLY RESTORE] close error %s\n", err.Error())
		}
	}()

	_, _ = fmt.Fprintf(os.Stdout, "[START RESTORE] read all message from snapshot\n")
	var totalRestored int

	decoder := json.NewDecoder(rClose)
	for decoder.More() {

		var command = &CommandPayload{}
		err := decoder.Decode(command)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "[END RESTORE] error decode data %s\n", err.Error())
			return err
		}

		var ieValue ElementValue

		if err := json.Unmarshal(command.Value, &ieValue); err != nil {
			fmt.Fprint(os.Stderr, "[WHILE RESTORE] error marshalling command value struct")
			return nil
		}

		if err := i.db.Insert(ieValue.UserId, ieValue.NodeId); err != nil {
			return err
		}

		totalRestored++
	}

	// read closing bracket
	_, err := decoder.Token()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "[END RESTORE] error %s\n", err.Error())
		return err
	}

	_, _ = fmt.Fprintf(os.Stdout, "[END RESTORE] success restore %d messages in snapshot\n", totalRestored)
	return nil
}

func NewFsm(d *db.MemoryDb) raft.FSM {
	return fsm{
		db: d,
	}
}
