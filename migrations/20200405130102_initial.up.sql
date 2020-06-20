CREATE TABLE `category` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `title` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `order` int(11) NOT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime DEFAULT NULL,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `feed` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `feed_title` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `feed_url` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `feed_image` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `feed_subtitle` longtext COLLATE utf8mb4_unicode_ci,
    `site_url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `site_favicon_url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `site_favicon` text COLLATE utf8mb4_unicode_ci,
    `category_id` int(11) DEFAULT NULL,
    `last_read_at` datetime NOT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime DEFAULT NULL,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `feed_category_id` (`category_id`),
    KEY `feed__deleted_at__category_id` (`deleted_at`,`category_id`),
    CONSTRAINT `feed_ibfk_1` FOREIGN KEY (`category_id`) REFERENCES `category` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `entry` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `title` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL,
    `author` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `summary` longtext COLLATE utf8mb4_unicode_ci,
    `link` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `published_at` datetime NOT NULL,
    `feed_id` int(11) NOT NULL,
    `read_at` datetime DEFAULT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime DEFAULT NULL,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `entry_feed_id` (`feed_id`),
    KEY `entry__feed_id__read_at__created_at_desc` (`feed_id`,`read_at`,`created_at`),
    KEY `entry_feed_id_published_at` (`feed_id`,`published_at`),
    KEY `deleted_at_idx` (`deleted_at`),
    KEY `feed_read_at_idx` (`read_at`,`feed_id`,`deleted_at`),
    KEY `entry_link_idx` (`link`), -- used when upserting entry
    CONSTRAINT `entry_ibfk_1` FOREIGN KEY (`feed_id`) REFERENCES `feed` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- CREATE DATABASE `webrss` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;