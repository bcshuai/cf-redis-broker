package recovery

import "github.com/bcshuai/cf-redis-broker/recovery/task"

type Snapshotter interface {
	Snapshot() (task.Artifact, error)
}
