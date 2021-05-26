package storage

import (
	"github.com/nuetoban/go-pomodoro/dto"
)

type Null struct{}

func NewNull() *Null                                  { return &Null{} }
func (Null) Create(dto.Task) (id int, err error)      { return 0, nil }
func (Null) Update(int, map[string]interface{}) error { return nil }
func (Null) Get(id int) (dto.Task, error)             { return dto.Task{}, nil }
func (Null) Delete(id int) error                      { return nil }
