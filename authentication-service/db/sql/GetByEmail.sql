-- name: GetByEmail :one
select id,
       email,
       first_name,
       last_name,
       password,
       user_active,
       created_at,
       updated_at
from users
where email = $1 ;
