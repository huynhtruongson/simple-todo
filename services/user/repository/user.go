package user_repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/huynhtruongson/simple-todo/lib"
	user_entity "github.com/huynhtruongson/simple-todo/services/user/entity"
	"github.com/huynhtruongson/simple-todo/utils"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (repo *UserRepo) CreateUser(ctx context.Context, db lib.QueryExecer, user user_entity.User) (int, error) {
	fields, values := user.FieldMap()
	query := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES (%s) RETURNING user_id`,
		user.TableName(),
		//remove the user_id when creating
		strings.Join(fields[1:], ","),
		utils.GeneratePlaceHolders(len(fields[1:])),
	)
	var userID int
	if err := db.QueryRow(ctx, query, values[1:]...).Scan(&userID); err != nil {
		return userID, err
	}
	return userID, nil
}

func (repo *UserRepo) GetUsersByUsername(ctx context.Context, db lib.QueryExecer, username string) ([]user_entity.User, error) {
	var user user_entity.User
	fields, _ := user.FieldMap()
	query := fmt.Sprintf(
		`SELECT %s FROM %s WHERE username = $1 AND deleted_at IS NULL`,
		strings.Join(fields, ","),
		user_entity.User{}.TableName(),
	)

	var users []user_entity.User
	rows, err := db.Query(ctx, query, username)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user user_entity.User
		_, values := user.FieldMap()
		if err := rows.Scan(values...); err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (repo *UserRepo) GetUsersByEmail(ctx context.Context, db lib.QueryExecer, email string) ([]user_entity.User, error) {
	var user user_entity.User
	fields, _ := user.FieldMap()
	query := fmt.Sprintf(
		`SELECT %s FROM %s WHERE email = $1 AND deleted_at IS NULL`,
		strings.Join(fields, ","),
		user_entity.User{}.TableName(),
	)

	var users []user_entity.User
	rows, err := db.Query(ctx, query, email)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user user_entity.User
		_, values := user.FieldMap()
		if err := rows.Scan(values...); err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (repo *UserRepo) GetUsersByUserIds(ctx context.Context, db lib.QueryExecer, ids []int) ([]user_entity.User, error) {
	var user user_entity.User
	fields, _ := user.FieldMap()
	query := fmt.Sprintf(
		`SELECT %s FROM %s WHERE user_id = ANY($1) AND deleted_at IS NULL`,
		strings.Join(fields, ","),
		user_entity.User{}.TableName(),
	)

	var users []user_entity.User
	rows, err := db.Query(ctx, query, ids)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user user_entity.User
		_, values := user.FieldMap()
		if err := rows.Scan(values...); err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, nil
}
