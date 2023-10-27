CREATE TABLE `users` (
  `id` bigint PRIMARY KEY NOT NULL,
  `nickname` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `salt` int NOT NULL COMMENT '哈希加盐密码',
  `profile_id` bigint NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  `is_deleted` tinyint NOT NULL COMMENT '0 表示否 1 表示是'
);

CREATE TABLE `profiles` (
  `id` bigint PRIMARY KEY NOT NULL,
  `real_name` varchar(255),
  `mood` varchar(255),
  `gender` ENUM ('male', 'female', 'unknown') COMMENT '性别：男、女、未知',
  `brithday` date,
  `introduction` varchar(255)
);

CREATE TABLE `follows` (
  `following_user_id` bigint NOT NULL,
  `followed_user_id` bigint NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `comments` (
  `id` bigint PRIMARY KEY NOT NULL,
  `comment_id` bigint NOT NULL COMMENT '为空表示属于视频的一级评论，不为空表示为二级评论，最多允许三极评论',
  `user_id` bigint NOT NULL,
  `post_id` bigint NOT NULL,
  `content` varchar(255)
);

CREATE TABLE `likes` (
  `id` bigint PRIMARY KEY NOT NULL,
  `user_id` bigint NOT NULL,
  `post_id` bigint NOT NULL,
  `comment_id` bigint COMMENT '用于标记点赞的是评论还是视频'
);

CREATE TABLE `shares` (
  `id` bigint PRIMARY KEY NOT NULL,
  `user_id` bigint NOT NULL,
  `post_id` bigint NOT NULL
);

CREATE TABLE `posts` (
  `id` bigint PRIMARY KEY NOT NULL,
  `title` varchar(255) NOT NULL,
  `description` varchar(255) NOT NULL,
  `user_id` bigint NOT NULL,
  `video_id` bigint NOT NULL,
  `video_class_id` int NOT NULL,
  `is_deleted` tinyint NOT NULL COMMENT '0 表示否 1 表示是',
  `updated_at` timestamp NOT NULL DEFAULT (now()),
  `created_at` timestamp NOT NULL
);

CREATE TABLE `user_collect_folders` (
  `id` bigint PRIMARY KEY NOT NULL,
  `user_id` bigint NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` varchar(255)
);

CREATE TABLE `collections` (
  `id` bigint PRIMARY KEY NOT NULL,
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
  `id` bigint PRIMARY KEY NOT NULL,
  `user_id` bigint NOT NULL,
  `video_id` bigint NOT NULL,
  `is_deleted` tinyint NOT NULL COMMENT '0 表示否 1 表示是'
);

CREATE TABLE `video` (
  `id` bigint PRIMARY KEY NOT NULL,
  `content_hash` char NOT NULL,
  `updated_at` timestamp NOT NULL
);

CREATE TABLE `video_class` (
  `id` int PRIMARY KEY NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` varchar(255) NOT NULL,
  `is_enabled` tinyint NOT NULL COMMENT '0 表示否 1 表示是'
);

ALTER TABLE `users` ADD FOREIGN KEY (`profile_id`) REFERENCES `profiles` (`id`);

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

ALTER TABLE `posts` ADD FOREIGN KEY (`video_id`) REFERENCES `video` (`id`);

ALTER TABLE `posts` ADD FOREIGN KEY (`video_class_id`) REFERENCES `video_class` (`id`);

ALTER TABLE `user_collect_folders` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `collections` ADD FOREIGN KEY (`folder_id`) REFERENCES `user_collect_folders` (`id`);

ALTER TABLE `collections` ADD FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`);

ALTER TABLE `post_tag` ADD FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`);

ALTER TABLE `user_upload` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `user_upload` ADD FOREIGN KEY (`video_id`) REFERENCES `video` (`id`);