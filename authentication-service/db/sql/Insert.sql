-- name: Insert :exec
insert into users (email, first_name, last_name, password, user_active, created_at, updated_at)
values ($1, $2, $3, $4, $5, $6, $7) returning id;