package task

import (
	"context"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type Task struct {
	sync.Mutex
	Jobs map[string]cron.EntryID
	cron *cron.Cron
}

type Job struct {
	Name string
	Spec string
	Func func()
}

var (
	t    *Task
	once sync.Once
)

func GetManager() *Task {
	once.Do(func() {
		t = &Task{
			Jobs: map[string]cron.EntryID{},
			cron: cron.New(),
		}
	})
	return t
}

func (t *Task) AddJob(job Job) error {
	id, err := t.cron.AddFunc(job.Spec, job.Func)
	if err != nil {
		logrus.Errorf("Task AddJob failed! name:%s, spec:%s, err:%s", job.Name, job.Spec, err.Error())
		return err
	}
	t.Lock()
	t.Jobs[job.Name] = id
	t.Unlock()
	logrus.Infof("Task AddJob success! name:%s", job.Name)
	return nil
}

func (t *Task) RemoveJob(job string) {
	t.Lock()
	defer t.Unlock()
	id, ok := t.Jobs[job]
	if ok {
		t.cron.Remove(id)
		delete(t.Jobs, job)
	}
}

func (t *Task) Start() {
	t.cron.Start()
}

func (t *Task) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	select {
	case <-ctx.Done():
		logrus.Info("Task Stop Timeout!")
	case <-t.cron.Stop().Done():
		logrus.Info("Task Stopped!")
	}
}
