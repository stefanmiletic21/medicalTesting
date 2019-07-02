package handler

import (
	"context"
	"net/http"
	"medicalTesting/core"
	"medicalTesting/serverErr"
	"strings"
)

func handleAdmin(ctx context.Context, r *http.Request) (response interface{}, err error) {
	if strings.HasPrefix(r.URL.Path, "/person") {
		r.URL.Path = r.URL.Path[7:]
		return handlePerson(ctx, r)
	}
	if strings.HasPrefix(r.URL.Path, "/employee") {
		r.URL.Path = r.URL.Path[9:]
		switch r.Method {
		case http.MethodPost:
			return core.CreateEmployee(ctx, r.Body)
		case http.MethodPatch:
			if strings.HasPrefix(r.URL.Path, "/") {
				return nil, core.UpdateEmployee(ctx, r.URL.Path[1:], r.Body)
			}
			return nil, serverErr.ErrBadRequest

		case http.MethodGet:
			if strings.HasPrefix(r.URL.Path, "/") {
				return core.GetEmployee(ctx, r.URL.Path[1:])
			}
			return core.GetEmployees(ctx)
		case http.MethodDelete:
			if strings.HasPrefix(r.URL.Path, "/") {
				return nil, core.RemoveEmployee(ctx, r.URL.Path[1:])
			}
			return nil, serverErr.ErrBadRequest
		}
	}
	if strings.HasPrefix(r.URL.Path, "/user") {
		r.URL.Path = r.URL.Path[5:]

		switch r.Method {
		case http.MethodPost:
			return core.CreateUser(ctx, r.Body)
		case http.MethodPatch:
			if strings.HasPrefix(r.URL.Path, "/") {
				return nil, core.UpdateUser(ctx, r.URL.Path[1:], r.Body)
			}
			return nil, serverErr.ErrBadRequest
		case http.MethodGet:
			if strings.HasPrefix(r.URL.Path, "/") {
				return core.GetUser(ctx, r.URL.Path[1:])
			}
			return core.GetUsers(ctx)
		case http.MethodDelete:
			if strings.HasPrefix(r.URL.Path, "/") {
				return nil, core.RemoveUser(ctx, r.URL.Path[1:])
			}
			return nil, serverErr.ErrBadRequest
		}
	}

	return nil, serverErr.ErrInvalidAPICall
}
