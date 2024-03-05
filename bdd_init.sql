DROP DATABASE IF EXISTS `DATA`;

CREATE DATABASE IF NOT EXISTS `DATA`;

use DATA;

CREATE TABLE IF NOT EXISTS `USER` (
    `id` int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `username` varchar(255) NOT NULL,
    `password` varchar(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS `EVENT_TYPE` (
    `id` int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `name` varchar(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS `EVENT` (
    `id` int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `title` varchar(255) NOT NULL,
    `description` varchar(255) DEFAULT NULL,
    `localisation` varchar(255) DEFAULT NULL,
    `start_date` datetime DEFAULT NULL,
    `end_date` datetime DEFAULT NULL,
    `event_type` int(11) NOT NULL REFERENCES EVENT_TYPE(id)
);

CREATE TABLE IF NOT EXISTS `USER_EVENT` (
    `user` int(11),
    `event` int(11),
    PRIMARY KEY (`user`,`event`),
    FOREIGN KEY (`user`) REFERENCES USER(`id`),
    FOREIGN KEY (`event`) REFERENCES EVENT(`id`)
);

INSERT INTO `USER` (`username`, `password`) VALUES 
('JeanDupont', 'motdepasse123'),
('MarieCurie', 'radioactivite234'),
('LucBesson', 'cinema567');

INSERT INTO `EVENT_TYPE` (`name`) VALUES 
('Conférence'),
('Concert'),
('Exposition');

INSERT INTO `EVENT` (`title`, `description`, `localisation`, `start_date`, `end_date`, `event_type`) VALUES 
('Conférence sur l\'environnement', 'Discussion sur les changements climatiques', 'Paris, France', '2024-01-03 09:00:00', '2024-02-15 17:00:00', 1),
('Concert de Rock', 'Concert de rock avec des groupes locaux', 'Lyon, France', '2024-03-20 20:00:00', '2024-03-21 01:00:00', 2),
('Exposition d\'art moderne', 'Exposition des nouvelles œuvres d\'artistes contemporains', 'Marseille, France', '2024-04-10 10:00:00', '2024-05-10 18:00:00', 3);

INSERT INTO `USER_EVENT` (`user`, `event`) VALUES 
(1, 1),
(2, 2),
(3, 3);