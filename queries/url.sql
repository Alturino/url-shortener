-- name: InsertUrl :one
insert into urls(id, url, short_url) values($1, $2, $3) returning *;

-- name: UpdateUrl :one
update urls set url = $2 where short_url = $1 returning *;

-- name: UpdateVisitedCountUrl :one
update urls set visited_count = $2 where id = $1 returning *;

-- name: FindUrlByShortUrl :one
select * from urls where short_url = $1;

-- name: DeleteUrlByShortUrl :one
delete from urls where short_url = $1 returning *;
