create table `settings` (
    `key` varchar(255) collate utf8_unicode_ci not null,
    `value` text collate utf8_unicode_ci not null,
    primary key (`key`)
) engine=InnoDB default charset=utf8 collate=utf8_unicode_ci;

insert into `settings` (`key`, `value`) values ('username', 'admin');
insert into `settings` (`key`, `value`) values ('password', '$2a$10$LpbHiC5IKXDKIwi33gmj9uipd33nMsLov0rIL9kCFw45zhf72fHme');
