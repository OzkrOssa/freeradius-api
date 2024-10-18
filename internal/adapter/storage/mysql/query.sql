-- name: CreateUser :execresult
INSERT INTO radcheck(username, attribute, op, value) VALUES (?, ?, ?, ?);

-- name: GetUserByID :one
SELECT * FROM radcheck WHERE id = ? AND attribute = "User-Profile" LIMIT 1;

-- name: ListUsers :many
SELECT * FROM radcheck WHERE attribute = "User-Profile";

-- name: UpdateUserByUserName :execresult
UPDATE radcheck SET username = ? WHERE id = ?;

-- name: UpdateUserByProfile :execresult
UPDATE radcheck SET value = ? WHERE id = ? AND attribute = "User-Profile";

-- name: DeleteUser :exec
DELETE FROM radcheck WHERE username = ?;

-- name: CreateProfile :execresult
INSERT INTO radusergroup (username, groupname, priority) VALUES (?, ?, ?);

-- name: ListProfiles :many
SELECT * FROM radusergroup;

-- name: GetProfileByID :one
SELECT * FROM radusergroup WHERE id = ? LIMIT 1;

-- name: UpdateProfile :execresult
UPDATE radusergroup SET username = ?, groupname = ?, priority = ? WHERE id = ?;

-- name: DeleteProfile :exec
DELETE FROM radusergroup WHERE id = ?;

-- name: CreateGroupName :execresult
INSERT INTO radgroupreply (groupname, attribute, op, value) VALUES (?, ?, "=", ?);

-- name: ListGroupNames :many
SELECT * FROM radgroupreply;

-- name: UpdateGroupName :execresult
UPDATE radgroupreply SET groupname = ?, attribute = ?, value = ? WHERE id = ?;

-- name: GetGroupNameByID :one
SELECT * FROM radgroupreply WHERE id = ? LIMIT 1;

-- name: DeleteGroupName :exec
DELETE FROM radgroupreply WHERE id = ?;

