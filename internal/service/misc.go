package service

import (
	"fmt"
	"net/url"
	"strings"

	pb_pst "github.com/msqtt/sevencowcloud-shortvideo-service/api/pb/v1/post"
	pb_prf "github.com/msqtt/sevencowcloud-shortvideo-service/api/pb/v1/profile"
	pb_usr "github.com/msqtt/sevencowcloud-shortvideo-service/api/pb/v1/user"
	pb_vid "github.com/msqtt/sevencowcloud-shortvideo-service/api/pb/v1/video"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/pkg/config"
	db "github.com/msqtt/sevencowcloud-shortvideo-service/internal/repo/sqlc"
)

// convert database User type to protobuf User one.
func db2pbUser(conf config.Config, u db.User, prof db.Profile) *pb_usr.User {
	return &pb_usr.User{
		UserId:        u.ID,
		Nickname:  u.Nickname,
		Email:     u.Email,
		Profile:   db2pbProfile(conf, prof),
		CreatedAt: u.CreatedAt.Unix(),
	}
}

// convert database Profile type to protobuf Profile one.
func db2pbProfile(conf config.Config, prof db.Profile) *pb_prf.Profile {
	birthdate := prof.BirthDate.Time.Unix()
	if prof.BirthDate.Time.IsZero() {
		birthdate = 0
	}
	updateat := prof.UpdatedAt.Unix()
	if prof.UpdatedAt.IsZero() {
		updateat = 0
	}
	return &pb_prf.Profile{
		RealName:     prof.RealName.String,
		Mood:         prof.Mood.String,
		Gender:       string(prof.Gender.ProfilesGender),
		BirthDate:    birthdate,
		Introduction: prof.Introduction.String,
		AvatarLink: uriPath(conf.KodoLink,
			prof.AvatarLink.String),
		UpdatedAt: updateat,
	}
}

func db2pbVideo(conf config.Config, video db.Video) *pb_vid.Video {
	return &pb_vid.Video{
		CoverLink: uriPath(conf.KodoLink, video.CoverLink),
		SrcLink: uriPath(conf.KodoLink, video.SrcLink),
	}
}

func db2pbTag(vc db.Tag) *pb_vid.Tag {
	return &pb_vid.Tag{
		TagId: vc.ID,
		Name: vc.Name,
		Description: vc.Description,
	}
}

func db2pbPost(conf config.Config, post db.Post, video db.Video, dbtags []db.Tag) *pb_pst.Post {
	tags := make([]*pb_vid.Tag, len(dbtags))

	for i, t := range dbtags {
		tags[i] = &pb_vid.Tag{
			TagId: t.ID,
			Name: t.Name,
			Description: t.Description,
		}
	}
	return &pb_pst.Post{
		PostId: post.ID,
		Title: post.Title,
		Description: post.Description,
		UserId: post.UserID,
		Tags: tags,
		Video: db2pbVideo(conf, video),
		UpdatedAt: post.UpdatedAt.Unix(),
		CreatedAt: post.CreatedAt.Unix(),
	}
}

func uriPath(prefix, path string) string {
	s := strings.Split(path, "/")
	s[len(s)-1] = url.QueryEscape(s[len(s)-1])
	return fmt.Sprintf("%s/%s", prefix, strings.Join(s, "/"))
}

// convert int32 to bool.
func int2Bool(i int32) (res bool) {
	if i == 1 {
		res = true
	}
	return
}
