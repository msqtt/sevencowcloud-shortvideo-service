package db

import (
	"context"
	"database/sql"
	"math/rand"
	"testing"

	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/pkg/sample"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func testCreateProfile(t *testing.T) Profile {
	params := AddProfileParams{
		RealName:     sql.NullString{String: sample.RandomStr(rand.Intn(20)), Valid: true},
		Mood:         sql.NullString{String: sample.RandomStr(rand.Intn(50)), Valid: true},
		Gender:       NullProfilesGender{ProfilesGender: ProfilesGender(sample.RandomGender()), Valid: true},
		BirthDate:    sql.NullTime{Time: sample.RandomBirthDate(), Valid: true},
		Introduction: sql.NullString{String: sample.RandomStr(rand.Intn(200)), Valid: true},
	}
	res, err := testQueries.AddProfile(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, res)
	id, err := res.LastInsertId()
	require.NoError(t, err)
	require.NotZero(t, id)

	p, err2 := testQueries.GetProfile(context.Background(), id)
	require.NoError(t, err2)
	require.NotEmpty(t, p)

	require.Equal(t, p.RealName, params.RealName)
	require.Equal(t, p.Mood, params.Mood)
	require.Equal(t, p.Gender, params.Gender)
	require.Equal(t, p.BirthDate, params.BirthDate)
	require.Equal(t, p.Introduction, params.Introduction)

	return p
}

func TestCreateProfile(t *testing.T) {
	testCreateProfile(t)
}

func testCreateUser(t *testing.T) User {
	p := testCreateProfile(t)

	oriPasswd := sample.RandomStr(int(sample.RandomInt(8, 20)))
	passwd, err := bcrypt.GenerateFromPassword([]byte(oriPasswd), bcrypt.DefaultCost)
	require.NoError(t, err)
	require.Len(t, passwd, 60)

	params := AddUserParams{
		Nickname:  sample.RandomStr(rand.Intn(20)),
		Email:     sample.RandomEmail(),
		Password:  string(passwd),
		ProfileID: p.ID,
	}
	res, err := testQueries.AddUser(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, res)

	i, err2 := res.LastInsertId()
	require.NoError(t, err2)
	require.NotZero(t, i)
	
	u, err3 := testQueries.GetUser(context.Background(), i)
	require.NoError(t, err3)
	require.NotEmpty(t, u)
	require.Equal(t, u.Nickname, params.Nickname)
	require.Equal(t, u.Email, params.Email)
	require.Equal(t, u.Password, params.Password)
	require.Equal(t, u.ProfileID, params.ProfileID)

	// check for password
	err4 := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(oriPasswd))
	require.NoError(t, err4)

	return u
}

func TestCreateUser(t *testing.T) {
	testCreateUser(t)
}
