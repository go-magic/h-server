package work

import (
	"errors"
	"github.com/go-magic/h-server/task"
)

type Worker struct {
	WorkerPool     chan chan task.CheckRequest
	requestChannel chan task.CheckRequest
	quit           chan bool
}

func NewWorker(workerPool chan chan task.CheckRequest) Worker {
	return Worker{
		WorkerPool:     workerPool,
		requestChannel: make(chan task.CheckRequest),
		quit:           make(chan bool)}
}

func (w Worker) Start() {
	go w.start()
}

func (w Worker) start() {
	for {
		//注册当前worker到worker队列
		w.WorkerPool <- w.requestChannel
		select {
		case request := <-w.requestChannel:
			// 接收到一个工作请求
			w.check(request)
		case <-w.quit:
			// 接收到停止工作信号
			return
		}
	}
}

func (w Worker) check(request task.CheckRequest) {
	if request.Tasker == nil {
		request.CheckResultChan <- task.CheckResult{
			SubResult: task.ErrorResult(request.SubTask,"taskType not exist"),
			Error: errors.New("tasker not exist"),
		}
		return
	}
	result, err := request.Tasker.Check(request.SubTask.SubTask)
	if err != nil {
		request.CheckResultChan <- task.CheckResult{
			SubResult: task.ErrorResult(request.SubTask,err.Error()),
			Error: err,
		}
		return
	}
	request.CheckResultChan <- task.CheckResult{
		SubResult: task.SuccessResult(request.SubTask,result),
		Error: err,
	}
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
