#define model here for query linter, plz refer to: internal/adapters/repository/sqlc/schema/table_account.sql
CREATE TABLE `account`
(
    `id`           BIGINT       NOT NULL AUTO_INCREMENT COMMENT 'is identify account with primary key auto increment',
    `user_id`      BIGINT       NOT NULL COMMENT 'is identify user with foreign key with table user.id',
    `account_type` enum('bank_account', 'ewallet', 'bank_token') NOT NULL COMMENT 'type of account sof',
    `account_data` JSON         NOT NULL COMMENT 'data of account',
    `created_at`   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY pk__account (id),
    INDEX uq__account__user_id (user_id)
);


-- name: InsertAccount :execresult
INSERT INTO account(user_id, account_type, account_data)
VALUES (?, ?, ?);

-- name: GetAccountsByUserID :many
SELECT *
FROM account
WHERE user_id = ?;

-- name: DeleteAccountByUserID :execresult
DELETE
FROM account
WHERE user_id = ?;

-- name: GetAccountByID :one
SELECT *
FROM account
WHERE id = ?
LIMIT 1;