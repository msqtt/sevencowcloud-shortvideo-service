ALTER TABLE videos MODIFY COLUMN content_hash varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL;

ALTER TABLE posts MODIFY COLUMN is_deleted tinyint DEFAULT 0 NOT NULL COMMENT '0 表示否 1 表示是';

alter table videos
    modify content_hash varchar(255) null;
