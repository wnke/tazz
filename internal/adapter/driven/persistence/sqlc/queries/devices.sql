-- name: CreateDevice :one
INSERT INTO devices (url, username, password)
VALUES (?, ?, ?)
RETURNING id, url, username, password;

-- name: GetDeviceByID :one
SELECT id, url, username, password
FROM devices
WHERE id = ?;

-- name: ListDevices :many
SELECT id, url, username, password
FROM devices
ORDER BY id;

-- name: UpdateDevice :one
UPDATE devices
SET url = ?,
    username = ?,
    password = ?
WHERE id = ?
RETURNING id, url, username, password;

-- name: DeleteDevice :execrows
DELETE FROM devices
WHERE id = ?;
