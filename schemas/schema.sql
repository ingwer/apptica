CREATE DATABASE IF NOT EXISTS `apptica`
  COLLATE 'utf8mb4_unicode_ci'
  DEFAULT CHARSET 'utf8mb4';

USE apptica;

CREATE TABLE `top_positions` (
  `app_id` BIGINT(20) UNSIGNED NOT NULL,
  `date` DATETIME NOT NULL,
  `country_code` INT(11) NOT NULL,
  `category_id` BIGINT(20) NOT NULL,
  `position` BIGINT(20) NOT NULL,
  UNIQUE KEY `app_in_category_date` (`app_id`, `category_id`, `date`),
  KEY `app_id_date` (`app_id`, `date`)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
