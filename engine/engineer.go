package engine

import "github.com/go-magic/h-server/task"

type Producer interface {
	CreateServerTask() (*task.ServerTask, error)
}

type Consumer interface {
	DestroyServerResult(result *task.ServerResult) error
}
