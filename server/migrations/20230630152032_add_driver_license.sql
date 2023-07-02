-- Modify "trips" table
ALTER TABLE `trips` MODIFY COLUMN `status` enum('waiting','accept','done','cancel') NOT NULL DEFAULT "waiting";
-- Modify "vehicle_drivers" table
ALTER TABLE `vehicle_drivers` ADD COLUMN `license` enum('motor','car') NOT NULL;
