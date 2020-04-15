USE `test`;
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
    `id` INT(16) NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(45) UNIQUE NOT NULL,
    `password` VARCHAR(45) NOT NULL,
    `age` TINYINT(4) default 0,
    PRIMARY KEY (`id`)
) CHARSET=UTF8, ENGINE=InnoDB;

INSERT INTO `user` VALUES (1, 'zk', '123456', 24);