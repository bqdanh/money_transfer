CREATE TABLE `transaction`
(
    `id`                         BIGINT         NOT NULL COMMENT 'is identify transaction',
    `account_id`                 BIGINT         NOT NULL COMMENT 'is identify account',
    `amount`                     DECIMAL(18, 2) NOT NULL COMMENT 'amount of transaction',
    `version`                    INT            NOT NULL DEFAULT 0 COMMENT 'version of transaction',
    `request_id`                 NVARCHAR(255)  NOT NULL DEFAULT '' COMMENT 'request id of transaction',
    `description`                NVARCHAR(255)  NOT NULL DEFAULT '' COMMENT 'description of transaction',
    `partner_ref_transaction_id` NVARCHAR(255)  NOT NULL DEFAULT '' COMMENT 'partner reference transaction id',
    `status`                     NVARCHAR(255)  NOT NULL DEFAULT 'PENDING' COMMENT 'status of transaction',
    `type`                       NVARCHAR(255)  NOT NULL COMMENT 'type of transaction',
    `data`                       JSON           NOT NULL COMMENT 'data of transaction',
    `created_date`               DATE       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_at`                 TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`                 TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY pk__transaction__created_at (id, created_date),
    INDEX idx__transaction__account_id__created_at (account_id, request_id, created_date)
) PARTITION BY RANGE (TO_DAYS(`created_date`)) (
    PARTITION start VALUES LESS THAN (0),
    PARTITION before_2024_04 VALUES LESS THAN (TO_DAYS('2024-04-01')),
    PARTITION before_2024_06 VALUES LESS THAN (TO_DAYS('2024-06-01')),
    PARTITION before_2024_08 VALUES LESS THAN (TO_DAYS('2024-08-01')),
    PARTITION before_2024_10 VALUES LESS THAN (TO_DAYS('2024-10-01')),
    PARTITION before_2024_12 VALUES LESS THAN (TO_DAYS('2024-12-01')),
    PARTITION before_2025_02 VALUES LESS THAN (TO_DAYS('2025-02-01')),
    PARTITION before_2025_12 VALUES LESS THAN MAXVALUE
    );