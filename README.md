# Vault Local Read Replica 
Usually Hashicorp Vault is a single-instance application; even high-availability mode simply enables warm standby instances. 
This instance is a single point of failure. Single points of failure are bad. 
This storage backend removes that single point of failure by giving each host a local read-only Vault that caches secrets. 

## How does that work? 
If you go through the github.com/hashicorp/vault/physical folder you'll see a number of different implementations 
of Vault storage backends. Each of them verifies their own integrity by calling `ExerciseBackend`, which then
does just that and makes sure that the backend works as expected. 

Well, guess what. A read-only backend does not pass the test. And come to think of it, that makes sense. 
A local replica might need to make temporary changes to the data, and if we can do that then we have a 
Vault backend that actually passes the `ExerciseBackend` test suite. 

### The `physical.Backend` interface 
```go
type Backend interface {
	Put(ctx context.Context, entry *Entry) error
	Get(ctx context.Context, key string) (*Entry, error)
	Delete(ctx context.Context, key string) error
	List(ctx context.Context, prefix string) ([]string, error)
}
```

Pretty straightforward. Put and Delete are mutating, and Get and List are read-only. 
Let's look at the struct for `LocalReplica`. 

### The `LocalReplica` struct

```go
type LocalReplica struct {
	backend       physical.Backend
	local         physical.Backend
	cacheLifetime time.Duration
}
```

`LocalReplica` is simply two backends and a cache lifetime. `backend` is the "real" backend.
It will be used as the authoritative source for data. `local` is just an `inmem` backend that 
plays the role of a local cache. `cacheLifetime` is exactly what it sounds like.
