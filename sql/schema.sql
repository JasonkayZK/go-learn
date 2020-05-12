use `test`;

DROP TABLE IF EXISTS `dbgrpc`;

CREATE TABLE `dbgrpc` (
    `id` INT(10) NOT NULL AUTO_INCREMENT COMMENT '学生Id信息',
    `name` VARCHAR(20) NOT NULL COMMENT '学生姓名',
    `grade` INT(10) NOT NULL COMMENT '学生年级',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB, CHARSET=utf8;

INSERT INTO `dbgrpc` VALUES (1, 'tester1', '1');
INSERT INTO `dbgrpc` VALUES (2, 'tester2', '1');
INSERT INTO `dbgrpc` VALUES (3, 'tester3', '2');
INSERT INTO `dbgrpc` VALUES (4, 'tester4', '2');
INSERT INTO `dbgrpc` VALUES (5, 'tester5', '3');
INSERT INTO `dbgrpc` VALUES (6, 'tester6', '3');

