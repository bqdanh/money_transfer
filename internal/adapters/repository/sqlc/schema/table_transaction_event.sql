CREATE TABLE `transaction_event`
(
    `id`             BIGINT    NOT NULL AUTO_INCREMENT,
    `transaction_id` BIGINT    NOT NULL,
    `version`        INT       NOT NULL,
    `data`           JSON      NOT NULL,
    `created_date`   DATE      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_at`     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY pk__transaction_event (`id`),
    INDEX `idx__transaction_event__transaction_id__version` (`transaction_id`, `version`)
);
