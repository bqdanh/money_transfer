CREATE TABLE `user`
(
    `id`         BIGINT        NOT NULL AUTO_INCREMENT COMMENT 'is identify user with primary key auto increment',
    `user_name`  NVARCHAR(255) NOT NULL COMMENT 'user name for login, unique',
    `password`   NVARCHAR(255) NOT NULL DEFAULT '' COMMENT 'password hashed',
    `full_name`  NVARCHAR(255) NOT NULL DEFAULT '' COMMENT 'user full name',
    `phone`      NVARCHAR(255)  NOT NULL DEFAULT '' COMMENT 'user phone number',
    `created_at` TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY pk__user (id),
    UNIQUE uq__user__user_name (user_name),
    INDEX ix__user__created_at (created_at)
)