package service

import (
	pb_prf "github.com/msqtt/sevencowcloud-shortvideo-service/api/pb/v1/profile"
	pb_usr "github.com/msqtt/sevencowcloud-shortvideo-service/api/pb/v1/user"
	db "github.com/msqtt/sevencowcloud-shortvideo-service/internal/repo/sqlc"
)

// convert database User type to protobuf User one.
func db2pbUser(u db.User, prof *db.Profile) *pb_usr.User {
	return &pb_usr.User{
		Id: u.ID,
		Nickname:  u.Nickname,
		Email:     u.Email,
		Profile:   db2pbProfile(prof),
		CreatedAt: u.CreatedAt.Unix(),
	}
}

// convert database Profile type to protobuf Profile one.
func db2pbProfile(prof *db.Profile) *pb_prf.Profile {
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
		AvatarLink:   prof.AvatarLink.String,
		UpdatedAt:    updateat,
	}
}

// convert int32 to bool.
func int2Bool(i int32) (res bool) {
	if i == 1 {
		res = true
	}
	return
}
