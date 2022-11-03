package daemonx

import "sync"

type Task interface {
	Start()
	Stop()
	SetWg(*sync.WaitGroup)
}

type TaskGroup struct {
	Wg    *sync.WaitGroup
	Tasks []Task
}

func NewTaskGroup() *TaskGroup {
	return &TaskGroup{Wg: &sync.WaitGroup{}}
}

func (t *TaskGroup) Add(task Task) {
	t.Wg.Add(1)
	task.SetWg(t.Wg)
	t.Tasks = append(t.Tasks, task)
}

func (t *TaskGroup) Run() {
	if TaskGroupManage == nil {
		return
	}

	for _, task := range t.Tasks {
		task.Start()
	}
}

func (t *TaskGroup) Stop() {
	if TaskGroupManage == nil {
		return
	}

	for _, task := range t.Tasks {
		task.Stop()
	}
}
