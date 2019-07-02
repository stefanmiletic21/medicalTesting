package core

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"medicalTesting/data"
	"medicalTesting/dto"
	"medicalTesting/logger"
	"medicalTesting/serverErr"
	"medicalTesting/utils"
	"strings"
	"time"
)

var (
	Login         = login
	CreateSession = createSession
	GetSession    = getSession
	RemoveSession = removeSession
	generateToken = _generateToken
)

const (
	UserKey = "UserKey"

	allowedTimeDifference = 5 * time.Minute
)

func login(ctx context.Context, name, pass string) (*dto.SessionInfo, error) {
	pass = utils.GetPasswordHash(pass)
	return data.Login(ctx, name, pass)
}

func createSession(ctx context.Context, userUID string) (token string, err error) {
	token = generateToken()
	err = data.CreateSession(ctx, userUID, token)
	if err != nil {
		logger.Error("Couldn't create session: %v", err)
		return
	}
	return
}

func getSession(ctx context.Context, authorizationString, userUID string) (*dto.SessionInfo, error) {
	sessions, err := data.GetSessionsForUser(ctx, userUID)
	if err != nil {
		logger.Error("Couldn't get sessions for user with uid %v: %v", userUID, err)
		return nil, err
	}
	splitAuthorization := strings.Split(authorizationString, "|")
	if len(splitAuthorization) != 2 {
		logger.Warn("User with uid %v sent a bad header authorization", userUID)
		return nil, serverErr.ErrBadRequest
	}
	sentTime := splitAuthorization[0]
	hashedTime := splitAuthorization[1]
	t, err := time.Parse(time.RFC3339, sentTime)
	if err != nil {
		logger.Warn("User with uid %v sent a bad header authorization: %v", userUID, err)
	}

	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	if t.Before(now.Add(-allowedTimeDifference)) || t.After(now.Add(allowedTimeDifference)) {
		logger.Warn("User with uid %v is using invalid time in header", userUID)
		return nil, serverErr.ErrBadRequest
	}

	for _, session := range sessions {
		hash := md5.New()
		hash.Write([]byte(sentTime))
		hash.Write([]byte(session.Token))
		h := hash.Sum(nil)
		if hex.EncodeToString(h) == hashedTime {
			return session, nil
		}
	}
	logger.Warn("User with uid %v is using bad authentication token", userUID)
	return nil, serverErr.ErrNotAuthenticated
}

func removeSession(ctx context.Context, authorizationString, userUID string) (err error) {
	session, err := getSession(ctx, authorizationString, userUID)
	if err != nil {
		logger.Error("Couldn't get session for user with uid %v: %v", userUID, err)
		return
	}
	err = data.RemoveSession(ctx, session.Token)
	if err != nil {
		logger.Error("Couldn't remove session for user with uid %v: %v", userUID, err)
		return
	}
	return
}

func _generateToken() string {
	raw := make([]byte, 18, 18)
	enc := make([]byte, 24, 24)

	rand.Read(raw)
	base64.RawURLEncoding.Encode(enc, raw)

	return string(enc)
}
