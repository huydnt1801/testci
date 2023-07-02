-- Modify "trips" table
ALTER TABLE `trips` MODIFY COLUMN `status` enum('pending','waiting','accept','done','cancel') NOT NULL DEFAULT "pending", ADD COLUMN `start_location` varchar(255) NOT NULL, ADD COLUMN `end_location` varchar(255) NOT NULL, ADD COLUMN `distance` double NOT NULL;
-- Modify "users" table
ALTER TABLE `users` ADD COLUMN `image_url` varchar(255) NULL;
