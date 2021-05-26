package pomodoro

import (
	"context"
	"time"

	"github.com/nuetoban/go-pomodoro/dto"
)

type TaskRepo interface {
	Create(dto.Task) (id int, err error)
	Update(int, map[string]interface{}) error
	Get(id int) (dto.Task, error)
	Delete(id int) error
}

type Service struct {
	tasks TaskRepo
}

func NewService(taskRepo TaskRepo) *Service {
	return &Service{
		tasks: taskRepo,
	}
}

func (s *Service) NewTask(name string) *dto.Task {
	return &dto.Task{
		Name:     name,
		Interval: time.Minute * 25,
	}
}

func (s *Service) RunTask(ctx context.Context, task dto.Task, progress, done chan struct{}) error {
	task.State = dto.STARTED
	id, err := s.tasks.Create(task)
	if err != nil {
		return err
	}

	now := time.Now()

	// If canceled
	go func() {
		<-ctx.Done()
		s.tasks.Update(id, map[string]interface{}{"state": dto.CANCELED, "finished_at": time.Now()})
		done <- struct{}{}
	}()

	for {
		time.Sleep(time.Second)
		progress <- struct{}{}

		if now.Add(task.Interval).Before(time.Now()) {
			s.tasks.Update(
				id,
				map[string]interface{}{
					"state":       dto.DONE,
					"finished_at": time.Now(),
				},
			)
			if err != nil {
				return err
			}
			break
		}
	}

	done <- struct{}{}

	return nil
}
