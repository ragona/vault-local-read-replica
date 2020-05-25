package local_replica

import (
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
	"testing"
)

func TestLocalReplicaBackend(t *testing.T) {
	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewLocalReplicaBackend(map[string]string{
		"storage_type": "s3",
		"region": "us-west-2",
    	"bucket": "ragona-vault-test",
	}, logger)

	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseBackend(t, b)
}