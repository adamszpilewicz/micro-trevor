-- name: ResetPassword :exec
update users set password = $1 where id = $2;