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
	accessHistory map[string]time.Time
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
		backend:       backend,
		local:         local,
		cacheLifetime: time.Minute * 5, // todo: Move to conf
		accessHistory: map[string]time.Time{},
	}

	return r, nil
}

func (n *LocalReplicaBackend) Put(ctx context.Context, entry *physical.Entry) error {
	err := n.local.Put(ctx, entry)
	if err != nil {
		return err
	}

	n.accessHistory[entry.Key] = time.Now()

	return nil
}

func (n *LocalReplicaBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	if n.warm(key) {
		return n.local.Get(ctx, key)
	}

	entry, err := n.backend.Get(ctx, key)
	if err == nil {
		err = n.local.Put(ctx, entry)
		n.accessHistory[key] = time.Now()
		if err != nil {
			// our inmem backend is broken; this should never happen
			return nil, err
		}

		return entry, nil
	}

	// in this case we're falling back to a cold entry because we failed to update from backend
	// todo consider logging a warning?
	entry, err = n.local.Get(ctx, key)
	if err != nil {
		// something truly funky is happening; suspect the accessHistory
		return nil, err
	}

	return entry, nil
}

// Keep in mind that this only deletes for the TTL of the entry.
// Subsequent Get operations will repopulate the cache.
// To be honest this one is weird for a read only replica, probably don't do this.
// todo: Should this warn?
func (n *LocalReplicaBackend) Delete(ctx context.Context, key string) error {
	err := n.local.Delete(ctx, key)
	if err != nil {
		return err
	}

	return nil
}

func (n *LocalReplicaBackend) List(ctx context.Context, prefix string) ([]string, error) {
	return n.backend.List(ctx, prefix)
}

func (n *LocalReplicaBackend) cached(key string) bool {
	_, ok := n.accessHistory[key]

	return ok
}

func (n *LocalReplicaBackend) warm(key string) bool {
	if !n.cached(key) {
		return false
	}

	lastAccessed, _ := n.accessHistory[key]

	if time.Now().After(lastAccessed.Add(n.cacheLifetime)) {
		return false
	}

	return true
}
