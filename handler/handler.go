package handler

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"medicalTesting/config"
	"medicalTesting/core"
	"medicalTesting/db"
	"medicalTesting/dto"
	"medicalTesting/enum"
	"medicalTesting/logger"
	"medicalTesting/serverErr"
	"strings"
	"time"
)

var (
	Login             = login
	Logout            = logout
	HandleAuthorized  = handleAuthorized
	HandleTestingHash = handleTestingHash
)

const (
	headerAuth    = "Authorization"
	headerUserUID = "UserUID"

	httpErrForbidden = "Forbidden"
	httpErrInternal  = "Internal"
)

func handleTestingHash(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Auth string
	}

	if err := shouldTestingRequestBeAllowed(r); err != nil {
		logger.Error("Will not handle request: %v", err)
		http.Error(w, httpErrForbidden, http.StatusForbidden)
		return
	}
	authorization := r.Header.Get(headerAuth)

	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	currentTime := now.Format(time.RFC3339)

	hash := md5.New()
	hash.Write([]byte(currentTime))
	hash.Write([]byte(authorization))
	h := hash.Sum(nil)
	resp := response{currentTime + "|" + hex.EncodeToString(h)}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func handleAuthorized(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get(headerAuth)
	userUID := r.Header.Get(headerUserUID)
	if authorization == "" || userUID == "" {
		http.Error(w, httpErrForbidden, http.StatusUnauthorized)
		return
	}

	ctx := context.Background()
	dbRunner := db.CreateRunner()
	ctx = context.WithValue(ctx, db.RunnerKey, dbRunner)

	userSession, err := core.GetSession(ctx, authorization, userUID)
	if err != nil {
		http.Error(w, httpErrInternal, http.StatusInternalServerError)
		return
	}
	if userSession == nil {
		http.Error(w, httpErrForbidden, http.StatusUnauthorized)
		return
	}
	ctx = context.WithValue(ctx, core.UserKey, userSession)

	// Removing /auth
	r.URL.Path = r.URL.Path[5:]

	response, err := handle(ctx, r)
	var httpResponseStatus int
	if err == nil {
		httpResponseStatus = http.StatusOK
	} else {
		switch err {
		case serverErr.ErrBadRequest:
			httpResponseStatus = http.StatusBadRequest
		case serverErr.ErrNotAuthenticated:
			httpResponseStatus = http.StatusUnauthorized
		case serverErr.ErrForbidden:
			httpResponseStatus = http.StatusForbidden
		case serverErr.ErrInvalidAPICall, serverErr.ErrResourceNotFound:
			httpResponseStatus = http.StatusNotFound
		case serverErr.ErrMethodNotAllowed:
			httpResponseStatus = http.StatusMethodNotAllowed
		default:
			httpResponseStatus = http.StatusInternalServerError
		}
	}

	w.WriteHeader(httpResponseStatus)
	json.NewEncoder(w).Encode(response)
}

func handle(ctx context.Context, r *http.Request) (response interface{}, err error) {
	userSession := ctx.Value(core.UserKey).(*dto.SessionInfo)
	switch {
	case strings.HasPrefix(r.URL.Path, "/admin"):
		if userSession.Role != enum.RoleAdmin {
			err = serverErr.ErrNotAuthenticated
			return
		}
		r.URL.Path = r.URL.Path[6:]
		response, err = handleAdmin(ctx, r)

	case strings.HasPrefix(r.URL.Path, "/doctor"):
		if userSession.Role != enum.RoleDoctor {
			err = serverErr.ErrNotAuthenticated
			return
		}
		r.URL.Path = r.URL.Path[7:]
		response, err = handleDoctor(ctx, r)

	case strings.HasPrefix(r.URL.Path, "/research"):
		if userSession.Role != enum.RoleResearch {
			err = serverErr.ErrNotAuthenticated
			return
		}
		r.URL.Path = r.URL.Path[9:]
		response, err = handleResearcher(ctx, r)

	case strings.HasPrefix(r.URL.Path, "/nurse"):
		if userSession.Role != enum.RoleNurse {
			err = serverErr.ErrNotAuthenticated
			return
		}
		r.URL.Path = r.URL.Path[6:]
		response, err = handleNurse(ctx, r)
	case strings.HasPrefix(r.URL.Path, "/pass"):
		r.URL.Path = r.URL.Path[5:]
		response, err = handlePassChange(ctx, r, userSession.UserUID)
	default:
		err = serverErr.ErrInvalidAPICall
	}
	return
}

func login(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	dbRunner := db.CreateRunner()
	ctx = context.WithValue(ctx, db.RunnerKey, dbRunner)

	type loginRequest struct {
		Username string
		Password string
	}

	type loginResponse struct {
		Username      string
		UserUID       string
		Authenticated bool
		Authorization string
		Role          enum.Role
	}

	request := &loginRequest{}
	response := &loginResponse{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		logger.Warn("Request data could not be decoded: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userSession, err := core.Login(ctx, request.Username, request.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if userSession == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := core.CreateSession(ctx, userSession.UserUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.Authenticated = true
	response.Role = userSession.Role
	response.Username = request.Username
	response.UserUID = userSession.UserUID
	response.Authorization = token
	buf := make([]byte, 0, 1000)
	responsew := bytes.NewBuffer(buf)
	json.NewEncoder(responsew).Encode(response)
	w.Write(responsew.Bytes())

}

func logout(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get(headerAuth)
	userUID := r.Header.Get(headerUserUID)
	if authorization == "" || userUID == "" {
		http.Error(w, httpErrForbidden, http.StatusUnauthorized)
		return
	}
	ctx := context.Background()
	dbRunner := db.CreateRunner()
	ctx = context.WithValue(ctx, db.RunnerKey, dbRunner)
	err := core.RemoveSession(ctx, authorization, userUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func shouldTestingRequestBeAllowed(request *http.Request) error {
	if !config.GetGeneralIsTestingMode() {
		return errors.New("Request allowed only in testing mode")
	}

	if len(request.RemoteAddr) == 0 {
		return errors.New("Request remote address is empty, which is unexpected")

	}

	portIndex := strings.LastIndex(request.RemoteAddr, ":")
	if portIndex == -1 {
		return errors.New("Request remote port is empty")

	}
	ip := request.RemoteAddr[0:portIndex]

	if ip != "127.0.0.1" && ip != "::1" {
		return errors.New("Request allowed from localhost only")
	}

	return nil
}
