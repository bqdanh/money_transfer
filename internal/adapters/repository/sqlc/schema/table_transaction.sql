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
);