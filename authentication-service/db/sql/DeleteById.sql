-- name: DeleteById :exec
delete from users where id = $1;