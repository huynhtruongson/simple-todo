package auth_repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/huynhtruongson/simple-todo/lib"
	auth_entity "github.com/huynhtruongson/simple-todo/services/auth/entity"
	"github.com/huynhtruongson/simple-todo/utils"
)

type SessionRepo struct{}

func NewSessionRepo() *SessionRepo {
	return &SessionRepo{}
}

func (s SessionRepo) CreateSession(ctx context.Context, db lib.QueryExecer, session auth_entity.Session) error {
	fields, values := session.FieldMap()
	query := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES (%s)`,
		session.TableName(),
		strings.Join(fields, ","),
		utils.GeneratePlaceHolders(len(fields)),
	)
	_, err := db.Exec(ctx, query, values...)

	return err
}

func (s SessionRepo) GetSessionByIds(ctx context.Context, db lib.QueryExecer, ids uuid.UUIDs) ([]auth_entity.Session, error) {
	session := &auth_entity.Session{}
	fields, _ := session.FieldMap()
	query := fmt.Sprintf(
		`SELECT %s FROM %s WHERE session_id = ANY($1)`,
		strings.Join(fields, ","),
		session.TableName(),
	)

	sessions := []auth_entity.Session{}
	rows, err := db.Query(ctx, query, ids)
	if err != nil {
		return sessions, err
	}
	defer rows.Close()

	for rows.Next() {
		var session auth_entity.Session
		_, values := session.FieldMap()
		if err := rows.Scan(values...); err != nil {
			return sessions, err
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}
