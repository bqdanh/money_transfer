#define model here for query linter, plz refer to: internal/adapters/repository/sqlc/schema/table_account.sql

CREATE TABLE `transaction`
(
    `id`                         BIGINT AUTO_INCREMENT NOT NULL COMMENT 'is identify transaction',
    `account_id`                 BIGINT                NOT NULL COMMENT 'is identify account',
    `amount`                     DECIMAL(18, 2)        NOT NULL COMMENT 'amount of transaction',
    `version`                    INT                   NOT NULL DEFAULT 0 COMMENT 'version of transaction',
    `request_id`                 NVARCHAR(255)         NOT NULL DEFAULT '' COMMENT 'request id of transaction',
    `description`                NVARCHAR(255)         NOT NULL DEFAULT '' COMMENT 'description of transaction',
    `partner_ref_transaction_id` NVARCHAR(255)         NOT NULL DEFAULT '' COMMENT 'partner reference transaction id',
    `status`                     NVARCHAR(255)         NOT NULL DEFAULT 'PENDING' COMMENT 'status of transaction',
    `type`                       NVARCHAR(255)         NOT NULL COMMENT 'type of transaction',
    `data`                       JSON                  NOT NULL COMMENT 'transaction extra data',
    `created_at`                 TIMESTAMP             NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`                 TIMESTAMP             NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY pk__transaction__created_at (id, created_at),
    INDEX idx__transaction__account_id__created_at (account_id, request_id, created_at)
) PARTITION BY RANGE (TO_DAYS(`created_at`)) (
    PARTITION start VALUES LESS THAN (0),
    PARTITION before_2024_04 VALUES LESS THAN (TO_DAYS('2024-04-01')),
    PARTITION before_2024_06 VALUES LESS THAN (TO_DAYS('2024-06-01')),
    PARTITION before_2024_08 VALUES LESS THAN (TO_DAYS('2024-08-01')),
    PARTITION before_2024_10 VALUES LESS THAN (TO_DAYS('2024-10-01')),
    PARTITION before_2024_12 VALUES LESS THAN (TO_DAYS('2024-12-01')),
    PARTITION before_2025_02 VALUES LESS THAN (TO_DAYS('2025-02-01')),
    PARTITION before_2025_12 VALUES LESS THAN MAXVALUE
    );

-- name: GetTransactionByID :one
SELECT *
FROM `transaction`
WHERE `id` = ?
limit 1;

-- name: CreateTransaction :execresult
INSERT INTO `transaction` (`account_id`, `amount`, `version`, `request_id`, `description`,
                           `partner_ref_transaction_id`, `status`, `type`, `data`)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetTransactionByRequestID :one
SELECT *
FROM `transaction`
WHERE `account_id` = ?
  AND `request_id` = ?
limit 1;

-- name: UpdateTransaction :exec
UPDATE `transaction`
SET `amount`                     = ?,
    `version`                    = ?,
    `request_id`                 = ?,
    `description`                = ?,
    `partner_ref_transaction_id` = ?,
    `status`                     = ?,
    `type`                       = ?,
    `data`                       = ?
WHERE `id` = ?;