-- name: InsertUser :execresult
INSERT INTO user(user_name,
                 password,
                 full_name,
                 phone)
VALUES (?, ?, ?, ?);