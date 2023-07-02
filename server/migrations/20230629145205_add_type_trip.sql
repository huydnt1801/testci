-- Modify "trips" table
ALTER TABLE `trips` ADD COLUMN `type` enum('motor','car') NOT NULL;
