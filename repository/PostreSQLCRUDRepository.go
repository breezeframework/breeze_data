package repository

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/breezeframework/breeze_data/client/db"
	"github.com/pkg/errors"
)

const (
	WebPageTableName               = "WebPage"
	NoteTableName                  = "Note"
	UsersTableName                 = "Users"
	JoinWebPage                    = "WebPage ON webpage.id = webpage_id"
	JoinUser                       = "Users ON users.id = user_id"
	idColumn                       = "id"
	notesCountColumn               = "notes_count"
	urlColumn                      = "url"
	userIdColumn                   = "user_id"
	userKeyColumn                  = "user_key"
	userSNColumn                   = "user_sn"
	webpageIdColumn                = "webpage_Id"
	webpageURLColumn               = "url"
	contentColumn                  = "content"
	trustColumn                    = "trust"
	createdAtColumn                = "created_at"
	updatedAtColumn                = "updated_at"
	QUERY_NAME_CREATE_NOTE         = "NoteRepository.CreateNote"
	QUERY_NAME_GET_NOTES           = "NoteRepository.GetAllNotes"
	QUERY_NAME_GET_NOTE_BY_URL     = "NoteRepository.GetNoteByURL"
	QUERY_NAME_GET_NOTE_BY_ID      = "NoteRepository.GetNoteById"
	QUERY_NAME_CREATE_WEBPAGE      = "NoteRepository.CreateWebPage"
	QUERY_NAME_GET_WEB_PAGE_BY_URL = "NoteRepository.GetNoteByURL"
	QUERY_NAME_GET_WEB_PAGE_BY_ID  = "NoteRepository.GetNoteById"
	RETURNING_ID                   = "RETURNING id"
)

type PostgreSQLCRUDRepository[T any] struct {
	dbConnection  db.DBConnection
	insertBuilder sq.InsertBuilder
	selectBuilder sq.SelectBuilder
	updateBuilder sq.UpdateBuilder
	deleteBuilder sq.DeleteBuilder
	scanner       func(row pgx.Row) (*T, error)
}

func NewPostgreSQLCRUDRepository[T any](
	insertBuilder sq.InsertBuilder,
	selectBuilder sq.SelectBuilder,
	updateBuilder sq.UpdateBuilder,
	deleteBuilder sq.DeleteBuilder,
	scanner func(pgx.Row) (*T, error)) *PostgreSQLCRUDRepository[T] {
	return &PostgreSQLCRUDRepository[T]{
		insertBuilder: insertBuilder, selectBuilder: selectBuilder, updateBuilder: updateBuilder, deleteBuilder: deleteBuilder,
		scanner: scanner}
}

func (repo *PostgreSQLCRUDRepository[T]) Create(ctx context.Context, entity T) (int64, error) {
	builder := repo.insertBuilder.Suffix(RETURNING_ID)
	var id int64
	err := repo.dbConnection.QueryRowContextInsert(ctx, &builder).Scan(&id)
	return id, err
}

func (repo *PostgreSQLCRUDRepository[T]) GetById(ctx context.Context, id int64) (*T, error) {
	builder := repo.selectBuilder.Where(sq.Eq{idColumn: id})
	row := repo.dbConnection.QueryRowContextSelect(ctx, &builder)
	return repo.scanner(row)
}

func convertToObjects[T any](rows pgx.Rows, scanner func(pgx.Row) (*T, error)) (*[]T, error) {
	var objs []T
	for rows.Next() {
		obj, err := scanner(rows)
		if err != nil {
			return nil, err
		}
		objs = append(objs, *obj)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &objs, nil
}

func (repo *PostgreSQLCRUDRepository[T]) GetAll(ctx context.Context) (*[]T, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	rows := repo.dbConnection.QueryContextSelect(ctx, &repo.selectBuilder, nil)
	objs, err := convertToObjects(rows, repo.scanner)
	return objs, err
}

func (repo *PostgreSQLCRUDRepository[T]) Delete(ctx context.Context, id int64) error {
	panic("implement me")
}

func (repo *PostgreSQLCRUDRepository[T]) Update(ctx context.Context, id int64, entity T) error {
	panic("implement me")
}
