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