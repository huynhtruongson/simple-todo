package task_entity

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/field"
)

type Task struct {
	TaskID      field.Int    `json:"task_id"`
	UserID      field.Int    `json:"user_id"`
	Title       field.String `json:"title"`
	Status      TaskStatus   `json:"status"`
	Description field.String `json:"description"`
	common.SQLModel
}

//go:generate stringer -type TaskStatus -trimprefix TaskStatus
type TaskStatus int

func (t Task) TableName() string {
	return "tasks"
}

func (t *Task) FieldMap() ([]string, []interface{}) {
	return []string{
			"task_id",
			"user_id",
			"title",
			"status",
			"description",
		}, []interface{}{
			&t.TaskID,
			&t.UserID,
			&t.Title,
			&t.Status,
			&t.Description,
		}
}

const (
	TaskStatusNotStart TaskStatus = iota
	TaskStatusDoing
	TaskStatusFinish
)

var (
	TaskStatusString = []string{"Not start", "Doing", "Finish"}
)

var (
	ErrorTitleIsEmpty  = errors.New("Title is empty")
	ErrorUserIsEmpty   = errors.New("User is empty")
	ErrorInvalidStatus = errors.New("Status is invalid")
	ErrorUserNotFound  = errors.New("User not found")
	ErrorTaskNotFound  = errors.New("Task not found")
)

func parseTaskStatusToEnum(s string) (TaskStatus, error) {
	for i, status := range TaskStatusString {
		if s == status {
			return TaskStatus(i), nil
		}
	}
	return TaskStatus(0), errors.New("invalid task status string")
}

// implement the database/sql Scan and Value interfaces
func (status *TaskStatus) Scan(src interface{}) error {
	bytes, ok := src.(string)
	if !ok {
		return errors.New(fmt.Sprintf("fail to scan data from sql: %s", src))
	}
	value, err := parseTaskStatusToEnum(string(bytes))
	if err != nil {
		return err
	}
	*status = value

	return nil
}

func (status *TaskStatus) Value() (driver.Value, error) {
	if status == nil {
		return nil, nil
	}

	statusStr := status.String()

	return []byte(statusStr), nil
}

type TasksWithPaging struct {
	Tasks  []Task        `json:"tasks"`
	Paging common.Paging `json:"pagination"`
}
