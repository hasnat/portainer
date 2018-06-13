package bolt

import (
	"github.com/portainer/portainer"
	"github.com/portainer/portainer/bolt/internal"

	"github.com/boltdb/bolt"
)

// CommandService represents a service for managing commands.
type CommandService struct {
	store *Store
}

// Command returns an command by ID.
func (service *CommandService) Command(ID portainer.CommandID) (*portainer.Command, error) {
	var data []byte
	err := service.store.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(commandBucketName))
		value := bucket.Get(internal.Itob(int(ID)))
		if value == nil {
			return portainer.ErrCommandNotFound
		}

		data = make([]byte, len(value))
		copy(data, value)
		return nil
	})
	if err != nil {
		return nil, err
	}

	var command portainer.Command
	err = internal.UnmarshalCommand(data, &command)
	if err != nil {
		return nil, err
	}
	return &command, nil
}

// Commands return an array containing all the commands.
func (service *CommandService) Commands() ([]portainer.Command, error) {
	var commands = make([]portainer.Command, 0)
	err := service.store.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(commandBucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var command portainer.Command
			err := internal.UnmarshalCommand(v, &command)
			if err != nil {
				return err
			}
			commands = append(commands, command)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return commands, nil
}

// Synchronize creates, updates and deletes commands inside a single transaction.
func (service *CommandService) Synchronize(toCreate, toUpdate, toDelete []*portainer.Command) error {
	return service.store.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(commandBucketName))

		for _, command := range toCreate {
			err := storeNewCommand(command, bucket)
			if err != nil {
				return err
			}
		}

		for _, command := range toUpdate {
			err := marshalAndStoreCommand(command, bucket)
			if err != nil {
				return err
			}
		}

		for _, command := range toDelete {
			err := bucket.Delete(internal.Itob(int(command.ID)))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// CreateCommand assign an ID to a new command and saves it.
func (service *CommandService) CreateCommand(command *portainer.Command) error {
	return service.store.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(commandBucketName))
		err := storeNewCommand(command, bucket)
		if err != nil {
			return err
		}
		return nil
	})
}

// UpdateCommand updates an command.
func (service *CommandService) UpdateCommand(ID portainer.CommandID, command *portainer.Command) error {
	data, err := internal.MarshalCommand(command)
	if err != nil {
		return err
	}

	return service.store.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(commandBucketName))
		err = bucket.Put(internal.Itob(int(ID)), data)
		if err != nil {
			return err
		}
		return nil
	})
}

// DeleteCommand deletes an command.
func (service *CommandService) DeleteCommand(ID portainer.CommandID) error {
	return service.store.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(commandBucketName))
		err := bucket.Delete(internal.Itob(int(ID)))
		if err != nil {
			return err
		}
		return nil
	})
}

func marshalAndStoreCommand(command *portainer.Command, bucket *bolt.Bucket) error {
	data, err := internal.MarshalCommand(command)
	if err != nil {
		return err
	}

	err = bucket.Put(internal.Itob(int(command.ID)), data)
	if err != nil {
		return err
	}
	return nil
}

func storeNewCommand(command *portainer.Command, bucket *bolt.Bucket) error {
	id, _ := bucket.NextSequence()
	command.ID = portainer.CommandID(id)
	return marshalAndStoreCommand(command, bucket)
}
