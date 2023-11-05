package task_repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/huynhtruongson/simple-todo/lib"
	task_entity "github.com/huynhtruongson/simple-todo/services/task/entity"
	"github.com/huynhtruongson/simple-todo/utils"
)

type TaskRepo struct{}

func NewTaskRepo() *TaskRepo {
	return &TaskRepo{}
}

func (t TaskRepo) CreateTask(ctx context.Context, db lib.QueryExecer, task task_entity.Task) (int, error) {
	fields, values := task.FieldMap()
	query := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES (%s) RETURNING task_id`,
		task.TableName(),
		//remove the task_id when creating
		strings.Join(fields[1:], ","),
		utils.GeneratePlaceHolders(len(fields[1:])),
	)
	var taskID int
	if err := db.QueryRow(ctx, query, values[1:]...).Scan(&taskID); err != nil {
		return taskID, err
	}
	return taskID, nil
}

func (t TaskRepo) UpdateTask(ctx context.Context, db lib.QueryExecer, task task_entity.Task) error {
	fields, values := task.FieldMap()
	query := fmt.Sprintf(
		`UPDATE %s SET %s,updated_at = now() WHERE task_id = $%d`,
		task.TableName(),
		//remove the task_id when updating
		utils.GenerateUpdatePlaceHolders(fields[1:]),
		len(fields[1:])+1,
	)
	var args []interface{} = values[1:]
	args = append(args, values[0])
	if _, err := db.Exec(ctx, query, args...); err != nil {
		return err
	}
	return nil
}

func (t TaskRepo) GetTasksByIds(ctx context.Context, db lib.QueryExecer, ids []int) ([]task_entity.Task, error) {
	task := &task_entity.Task{}
	fields, _ := task.FieldMap()
	query := fmt.Sprintf(
		`SELECT %s FROM %s WHERE task_id = ANY($1) AND deleted_at IS NULL`,
		strings.Join(fields, ","),
		task_entity.Task{}.TableName(),
	)

	tasks := []task_entity.Task{}
	rows, err := db.Query(ctx, query, ids)
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		var task task_entity.Task
		_, values := task.FieldMap()
		if err := rows.Scan(values...); err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (t TaskRepo) DeleteTask(ctx context.Context, db lib.QueryExecer, id int) error {
	query := fmt.Sprintf(
		`UPDATE %s SET deleted_at = now() WHERE task_id = $1`,
		task_entity.Task{}.TableName(),
	)

	if _, err := db.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}

func (t TaskRepo) CountTask(ctx context.Context, db lib.QueryExecer) (int, error) {
	query := fmt.Sprintf(
		`SELECT count(task_id) FROM %s WHERE deleted_at IS NULL`,
		task_entity.Task{}.TableName(),
	)
	var count int
	if err := db.QueryRow(ctx, query).Scan(&count); err != nil {
		return count, err
	}

	return count, nil
}

func (t TaskRepo) GetTasksWithFilter(ctx context.Context, db lib.QueryExecer, limit, offset int) ([]task_entity.Task, error) {
	task := &task_entity.Task{}
	fields, _ := task.FieldMap()
	query := fmt.Sprintf(
		`SELECT %s FROM %s WHERE deleted_at IS NULL ORDER BY created_at LIMIT $1 OFFSET $2`,
		strings.Join(fields, ","),
		task_entity.Task{}.TableName(),
	)

	tasks := []task_entity.Task{}
	rows, err := db.Query(ctx, query, limit, offset)
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		var task task_entity.Task
		_, values := task.FieldMap()
		if err := rows.Scan(values...); err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
