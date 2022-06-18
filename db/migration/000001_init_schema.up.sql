CREATE TABLE `accounts` (
                           `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
                           `owner` varchar(255) NOT NULL,
                           `balance` bigint NOT NULL,
                           `currency` varchar(255) NOT NULL,
                           `created_at` TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE `entries` (
                           `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
                           `account_id` bigint,
                           `amount` bigint NOT NULL COMMENT 'can be nagative or positive',
                           `created_at` TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE `transfers` (
                             `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
                             `from_account_id` bigint,
                             `to_account_id` bigint,
                             `amount` bigint NOT NULL COMMENT 'must be positive',
                             `created_at` TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX `account_index_0` ON `accounts` (`owner`);

CREATE INDEX `entries_index_1` ON `entries` (`account_id`);

CREATE INDEX `transfers_index_2` ON `transfers` (`from_account_id`);

CREATE INDEX `transfers_index_3` ON `transfers` (`to_account_id`);

CREATE INDEX `transfers_index_4` ON `transfers` (`from_account_id`, `to_account_id`);

-- ALTER TABLE `entries` ADD FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`);
--
-- ALTER TABLE `transfers` ADD FOREIGN KEY (`from_account_id`) REFERENCES `accounts` (`id`);
--
-- ALTER TABLE `transfers` ADD FOREIGN KEY (`to_account_id`) REFERENCES `accounts` (`id`);
