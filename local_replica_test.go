package local_replica

import (
	"github.com/hashicorp/vault/sdk/physical"
	"testing"
)

func TestNilBackend(t *testing.T) {
	physical.ExerciseBackend(t, &LocalReplicaBackend{})
}
