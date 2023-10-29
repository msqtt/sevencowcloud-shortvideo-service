// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0

package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type ProfilesGender string

const (
	ProfilesGenderMale    ProfilesGender = "male"
	ProfilesGenderFemale  ProfilesGender = "female"
	ProfilesGenderUnknown ProfilesGender = "unknown"
)

func (e *ProfilesGender) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ProfilesGender(s)
	case string:
		*e = ProfilesGender(s)
	default:
		return fmt.Errorf("unsupported scan type for ProfilesGender: %T", src)
	}
	return nil
}

type NullProfilesGender struct {
	ProfilesGender ProfilesGender `json:"profiles_gender"`
	Valid          bool           `json:"valid"` // Valid is true if ProfilesGender is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullProfilesGender) Scan(value interface{}) error {
	if value == nil {
		ns.ProfilesGender, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.ProfilesGender.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullProfilesGender) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.ProfilesGender), nil
}

type Collection struct {
	ID        int64     `json:"id"`
	FolderID  int64     `json:"folder_id"`
	PostID    int64     `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	ID int64 `json:"id"`
	// 为空表示属于视频的一级评论，不为空表示为二级评论，最多允许三极评论
	CommentID int64          `json:"comment_id"`
	UserID    int64          `json:"user_id"`
	PostID    int64          `json:"post_id"`
	Content   sql.NullString `json:"content"`
}

type Follow struct {
	FollowingUserID int64     `json:"following_user_id"`
	FollowedUserID  int64     `json:"followed_user_id"`
	CreatedAt       time.Time `json:"created_at"`
}

type Like struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
	PostID int64 `json:"post_id"`
	// 用于标记点赞的是评论还是视频
	CommentID sql.NullInt64 `json:"comment_id"`
}

type Post struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	UserID       int64  `json:"user_id"`
	VideoID      int64  `json:"video_id"`
	VideoClassID int32  `json:"video_class_id"`
	// 0 表示否 1 表示是
	IsDeleted int32     `json:"is_deleted"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type PostTag struct {
	ID         sql.NullInt64  `json:"id"`
	PostID     int64          `json:"post_id"`
	TagContent sql.NullString `json:"tag_content"`
}

type Profile struct {
	ID       int64          `json:"id"`
	RealName sql.NullString `json:"real_name"`
	Mood     sql.NullString `json:"mood"`
	// 性别：男、女、未知
	Gender       NullProfilesGender `json:"gender"`
	BirthDate    sql.NullTime       `json:"birth_date"`
	Introduction sql.NullString     `json:"introduction"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

type Share struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
	PostID int64 `json:"post_id"`
}

type User struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	// 哈希加盐密码
	Password  string    `json:"password"`
	ProfileID int64     `json:"profile_id"`
	CreatedAt time.Time `json:"created_at"`
	// 账户是否被删除, 0 表示否 1 表示是
	IsDeleted int32 `json:"is_deleted"`
	// 账户是否激活, 0 表示否 1 表示是
	IsActivated int32 `json:"is_activated"`
}

type UserActivation struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"user_id"`
	ActivateCode string `json:"activate_code"`
	// 激活码是否被删除, 0 表示否 1 表示是
	IsDeleted int32     `json:"is_deleted"`
	ExpiredAt time.Time `json:"expired_at"`
}

type UserCollectFolder struct {
	ID          int64          `json:"id"`
	UserID      int64          `json:"user_id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
}

type UserUpload struct {
	ID      int64 `json:"id"`
	UserID  int64 `json:"user_id"`
	VideoID int64 `json:"video_id"`
	// 0 表示否 1 表示是
	IsDeleted int32 `json:"is_deleted"`
}

type Video struct {
	ID          int64     `json:"id"`
	ContentHash string    `json:"content_hash"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type VideoClass struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// 0 表示否 1 表示是
	IsEnabled int32 `json:"is_enabled"`
}
