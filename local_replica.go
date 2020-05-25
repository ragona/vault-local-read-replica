package local_replica

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	"time"
)

// Verify LocalReplicaBackend satisfies the correct interfaces
var _ physical.Backend = (*LocalReplicaBackend)(nil)


type LocalReplicaBackend struct {
	backend       physical.Backend
	local         physical.Backend
	cacheLifetime time.Duration
}

func NewLocalReplicaBackend(conf map[string]string, logger hclog.Logger) (physical.Backend, error) {
	// make sure the user defined the target storage type
	storageType, ok := conf["storage_type"]
	if !ok {
		return nil, errors.New("no 'storage_type' specified in config")
	}

	// remove this key so the target backend doesn't get confused by an extra value
	delete(conf, "storage_type")

	// grab the factory method out of the command list
	factory, ok := command.PhysicalBackends[storageType]
	if !ok {
		return nil, fmt.Errorf("unknown background type: %s", conf["storage_type"])
	}

	// the real backend that our nodes will call
	backend, err := factory(conf, logger)
	if err != nil {
		return nil, err
	}

	// the in memory cache we'll use for temporary mutations
	local, err := inmem.NewInmem(nil, logger)
	if err != nil {
		return nil, err
	}

	r := &LocalReplicaBackend{
		backend: backend,
		local: local,
		cacheLifetime: time.Minute * 5, // todo: Move to conf
	}

	return r, nil
}

func (n *LocalReplicaBackend) Put(ctx context.Context, entry *physical.Entry) error {
	return nil
}

func (n *LocalReplicaBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	return nil, nil
}

func (n *LocalReplicaBackend) Delete(ctx context.Context, key string) error {
	return nil
}

func (n *LocalReplicaBackend) List(ctx context.Context, prefix string) ([]string, error) {
	return nil, nil
}
