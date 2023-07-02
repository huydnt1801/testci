-- Create "otps" table
CREATE TABLE `otps` (`id` bigint NOT NULL AUTO_INCREMENT, `phone_number` varchar(255) NOT NULL, `otp` varchar(255) NOT NULL, `created_at` timestamp NOT NULL, PRIMARY KEY (`id`), UNIQUE INDEX `phone_number` (`phone_number`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "users" table
CREATE TABLE `users` (`id` bigint NOT NULL AUTO_INCREMENT, `created_at` timestamp NOT NULL, `updated_at` timestamp NOT NULL, `phone_number` varchar(15) NOT NULL, `confirmed` bool NOT NULL DEFAULT 0, `full_name` varchar(255) NOT NULL, `password` varchar(255) NOT NULL, PRIMARY KEY (`id`), UNIQUE INDEX `phone_number` (`phone_number`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
