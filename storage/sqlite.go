package storage

import (
	"fmt"
	"os"
	"os/user"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/nuetoban/go-pomodoro/dto"
)

type Sqlite struct {
	db *gorm.DB
}

func NewSqlite() (*Sqlite, error) {
	usr, _ := user.Current()
	dir := usr.HomeDir

	err := os.MkdirAll(dir+"/.local/share/go-pomodoro", os.ModePerm)
	if err != nil {
		return nil, err
	}

	var s Sqlite

	s.db, err = gorm.Open(sqlite.Open(dir+"/.local/share/go-pomodoro/pom.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = s.db.AutoMigrate(&dto.Task{})
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *Sqlite) Create(task dto.Task) (id int, err error) {
	err = s.db.Create(&task).Error
	return task.ID, err
}

func (s *Sqlite) Update(id int, updates map[string]interface{}) error {
	return s.db.Model(&dto.Task{}).Where("id = ?", id).Updates(updates).Error
}

func (s *Sqlite) Get(id int) (dto.Task, error) { return dto.Task{}, fmt.Errorf("not implemented yet") }
func (s *Sqlite) Delete(id int) error          { return fmt.Errorf("not implemented yet") }
