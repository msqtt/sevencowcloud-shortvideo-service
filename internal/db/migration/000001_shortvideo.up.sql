CREATE TABLE `users` (
  `id` bigint PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `nickname` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password` char(60) NOT NULL COMMENT '哈希加盐密码',
  `profile_id` bigint NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  `is_deleted` tinyint NOT NULL DEFAULT 0 COMMENT '账户是否被删除, 0 表示否 1 表示是',
  `is_activated` tinyint NOT NULL DEFAULT 0 COMMENT '账户是否激活, 0 表示否 1 表示是'
);

CREATE TABLE `user_activation` (
  `id` bigint PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `activate_code` char(8) NOT NULL,
  `is_deleted` tinyint NOT NULL DEFAULT 0 COMMENT '激活码是否被删除, 0 表示否 1 表示是',
  `expired_at` timestamp NOT NULL
);

CREATE TABLE `profiles` (
  `id` bigint PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `real_name` varchar(255),
  `mood` varchar(255),
  `gender` ENUM ('male', 'female', 'unknown') COMMENT '性别：男、女、未知',
  `birth_date` date,
  `introduction` varchar(255),
  `avatar_link` varchar(255) COMMENT '存放头像文件的相对路径',
  `updated_at` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `follows` (
  `following_user_id` bigint NOT NULL,
  `followed_user_id` bigint NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY (`followed_user_id`, `following_user_id`)
);

CREATE TABLE `comments` (
  `id` bigint PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `comment_id` bigint NOT NULL COMMENT '为空表示属于视频的一级评论，不为空表示为二级评论，最多允许三极评论',
  `user_id` bigint NOT NULL,
  `post_id` bigint NOT NULL,
  `content` varchar(255)
);

CREATE TABLE `likes` (
  `id` bigint PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `post_id` bigint NOT NULL,
  `comment_id` bigint COMMENT '用于标记点赞的是评论还是视频'
);

CREATE TABLE `shares` (
  `id` bigint PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `post_id` bigint NOT NULL
);

CREATE TABLE `posts` (
  `id` bigint PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `description` varchar(255) NOT NULL,
  `user_id` bigint NOT NULL,
  `video_id` bigint NOT NULL,
  `video_class_id` int NOT NULL,
  `is_deleted` tinyint NOT NULL COMMENT '0 表示否 1 表示是',
  `updated_at` timestamp NOT NULL DEFAULT (now()),
  `created_at` timestamp NOT NULL
);

CREATE TABLE `user_collect_folder` (
  `id` bigint PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` varchar(255)
);

CREATE TABLE `collections` (
  `id` bigint PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `folder_id` bigint NOT NULL,
  `post_id` bigint NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `post_tag` (
  `id` bigint,
  `post_id` bigint NOT NULL,
  `tag_content` varchar(255)
);

CREATE TABLE `user_upload` (
  `id` bigint PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `video_id` bigint NOT NULL,
  `is_deleted` tinyint NOT NULL COMMENT '0 表示否 1 表示是'
);

CREATE TABLE `videos` (
  `id` bigint PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `content_hash` char NOT NULL,
  `updated_at` timestamp NOT NULL
);

CREATE TABLE `video_class` (
  `id` int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` varchar(255) NOT NULL,
  `is_enabled` tinyint NOT NULL COMMENT '0 表示否 1 表示是'
);

CREATE INDEX `users_index_0` ON `users` (`email`);

ALTER TABLE `users` ADD FOREIGN KEY (`profile_id`) REFERENCES `profiles` (`id`);

ALTER TABLE `user_activation` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `follows` ADD FOREIGN KEY (`following_user_id`) REFERENCES `users` (`id`);

ALTER TABLE `follows` ADD FOREIGN KEY (`followed_user_id`) REFERENCES `users` (`id`);

ALTER TABLE `comments` ADD FOREIGN KEY (`comment_id`) REFERENCES `comments` (`id`);

ALTER TABLE `comments` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `comments` ADD FOREIGN KEY (`post_id`) REFERENCES `profiles` (`id`);

ALTER TABLE `likes` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `likes` ADD FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`);

ALTER TABLE `likes` ADD FOREIGN KEY (`comment_id`) REFERENCES `comments` (`id`);

ALTER TABLE `shares` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `shares` ADD FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`);

ALTER TABLE `posts` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `posts` ADD FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`);

ALTER TABLE `posts` ADD FOREIGN KEY (`video_class_id`) REFERENCES `video_class` (`id`);

ALTER TABLE `user_collect_folder` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `collections` ADD FOREIGN KEY (`folder_id`) REFERENCES `user_collect_folder` (`id`);

ALTER TABLE `collections` ADD FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`);

ALTER TABLE `post_tag` ADD FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`);

ALTER TABLE `user_upload` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `user_upload` ADD FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`);
