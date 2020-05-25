package consul

import (
	"context"
	"github.com/hashicorp/vault/sdk/physical"
	"time"
)

// Verify ConsulBackend satisfies the correct interfaces
var _ physical.Backend = (*LocalReplica)(nil)


// LocalReplica just lies to you and does nothing.
type LocalReplica struct {
	backend       physical.Backend
	local         physical.Backend
	cacheLifetime time.Duration
}

func (n *LocalReplica) Put(ctx context.Context, entry *physical.Entry) error {
	return nil
}

func (n *LocalReplica) Get(ctx context.Context, key string) (*physical.Entry, error) {
	return nil, nil
}

func (n *LocalReplica) Delete(ctx context.Context, key string) error {
	return nil
}

func (n *LocalReplica) List(ctx context.Context, prefix string) ([]string, error) {
	return nil, nil
}
