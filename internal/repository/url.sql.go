// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: url.sql

package repository

import (
	"context"

	"github.com/google/uuid"
)

const deleteUrlByShortUrl = `-- name: DeleteUrlByShortUrl :one
delete from urls where short_url = $1 returning id, url, short_url, created_at, updated_at, visited_count
`

func (q *Queries) DeleteUrlByShortUrl(ctx context.Context, shortUrl string) (Url, error) {
	row := q.queryRow(ctx, q.deleteUrlByShortUrlStmt, deleteUrlByShortUrl, shortUrl)
	var i Url
	err := row.Scan(
		&i.ID,
		&i.Url,
		&i.ShortUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.VisitedCount,
	)
	return i, err
}

const findUrlByShortUrl = `-- name: FindUrlByShortUrl :one
select id, url, short_url, created_at, updated_at, visited_count from urls where short_url = $1
`

func (q *Queries) FindUrlByShortUrl(ctx context.Context, shortUrl string) (Url, error) {
	row := q.queryRow(ctx, q.findUrlByShortUrlStmt, findUrlByShortUrl, shortUrl)
	var i Url
	err := row.Scan(
		&i.ID,
		&i.Url,
		&i.ShortUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.VisitedCount,
	)
	return i, err
}

const insertUrl = `-- name: InsertUrl :one
insert into urls(id, url, short_url) values($1, $2, $3) returning id, url, short_url, created_at, updated_at, visited_count
`

type InsertUrlParams struct {
	ID       uuid.UUID `json:"id"`
	Url      string    `json:"url"`
	ShortUrl string    `json:"short_url"`
}

func (q *Queries) InsertUrl(ctx context.Context, arg InsertUrlParams) (Url, error) {
	row := q.queryRow(ctx, q.insertUrlStmt, insertUrl, arg.ID, arg.Url, arg.ShortUrl)
	var i Url
	err := row.Scan(
		&i.ID,
		&i.Url,
		&i.ShortUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.VisitedCount,
	)
	return i, err
}

const updateUrl = `-- name: UpdateUrl :one
update urls set url = $2 where short_url = $1 returning id, url, short_url, created_at, updated_at, visited_count
`

type UpdateUrlParams struct {
	ShortUrl string `json:"short_url"`
	Url      string `json:"url"`
}

func (q *Queries) UpdateUrl(ctx context.Context, arg UpdateUrlParams) (Url, error) {
	row := q.queryRow(ctx, q.updateUrlStmt, updateUrl, arg.ShortUrl, arg.Url)
	var i Url
	err := row.Scan(
		&i.ID,
		&i.Url,
		&i.ShortUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.VisitedCount,
	)
	return i, err
}

const updateVisitedCountUrl = `-- name: UpdateVisitedCountUrl :one
update urls set visited_count = $2 where id = $1 returning id, url, short_url, created_at, updated_at, visited_count
`

type UpdateVisitedCountUrlParams struct {
	ID           uuid.UUID `json:"id"`
	VisitedCount int32     `json:"visited_count"`
}

func (q *Queries) UpdateVisitedCountUrl(ctx context.Context, arg UpdateVisitedCountUrlParams) (Url, error) {
	row := q.queryRow(ctx, q.updateVisitedCountUrlStmt, updateVisitedCountUrl, arg.ID, arg.VisitedCount)
	var i Url
	err := row.Scan(
		&i.ID,
		&i.Url,
		&i.ShortUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.VisitedCount,
	)
	return i, err
}