CREATE DATABASE IF NOT EXISTS `id_alloc_db`;

USE `id_alloc_db`;

CREATE TABLE `segments`
(
    `app_tag`     VARCHAR(32) NOT NULL,
    `max_id`      BIGINT      NOT NULL,
    `step`        BIGINT      NOT NULL,
    `update_time` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`app_tag`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
    COMMENT ='test业务ID池';

INSERT INTO segments(`app_tag`, `max_id`, `step`)
VALUES ('test_business', 0, 100000);
