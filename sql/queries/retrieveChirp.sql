-- name: RetrieveChirp :one
SELECT * FROM chirps 
where id = $1;