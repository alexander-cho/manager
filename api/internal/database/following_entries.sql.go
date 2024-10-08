// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: following_entries.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFollowingEntry = `-- name: CreateFollowingEntry :one
INSERT INTO following_entries (id, created_at, updated_at, user_id, entry_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at, user_id, entry_id
`

type CreateFollowingEntryParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	EntryID   uuid.UUID
}

func (q *Queries) CreateFollowingEntry(ctx context.Context, arg CreateFollowingEntryParams) (FollowingEntry, error) {
	row := q.db.QueryRowContext(ctx, createFollowingEntry,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.EntryID,
	)
	var i FollowingEntry
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.EntryID,
	)
	return i, err
}

const deleteFollowingEntry = `-- name: DeleteFollowingEntry :exec
DELETE FROM following_entries WHERE id=$1 AND user_id=$2
`

type DeleteFollowingEntryParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) DeleteFollowingEntry(ctx context.Context, arg DeleteFollowingEntryParams) error {
	_, err := q.db.ExecContext(ctx, deleteFollowingEntry, arg.ID, arg.UserID)
	return err
}

const getFollowingEntries = `-- name: GetFollowingEntries :many
SELECT id, created_at, updated_at, user_id, entry_id FROM following_entries WHERE user_id=$1
`

func (q *Queries) GetFollowingEntries(ctx context.Context, userID uuid.UUID) ([]FollowingEntry, error) {
	rows, err := q.db.QueryContext(ctx, getFollowingEntries, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FollowingEntry
	for rows.Next() {
		var i FollowingEntry
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.EntryID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
