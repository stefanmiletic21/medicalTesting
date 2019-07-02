package handler

import (
	"context"
	"net/http"
	"medicalTesting/core"
	"medicalTesting/serverErr"
	"strings"
)

func handlePerson(ctx context.Context, r *http.Request) (response interface{}, err error) {
	switch r.Method {
	case http.MethodPost:
		return core.CreatePerson(ctx, r.Body)
	case http.MethodPatch:
		if strings.HasPrefix(r.URL.Path, "/") {
			return nil, core.UpdatePerson(ctx, r.URL.Path[1:], r.Body)
		}
		return nil, serverErr.ErrBadRequest
	case http.MethodGet:
		if strings.HasPrefix(r.URL.Path, "/") {
			return core.GetPerson(ctx, r.URL.Path[1:])
		}
		return core.GetPersons(ctx)
	case http.MethodDelete:
		if strings.HasPrefix(r.URL.Path, "/") {
			return nil, core.RemovePerson(ctx, r.URL.Path[1:])
		}
		return nil, serverErr.ErrBadRequest
	}

	return nil, serverErr.ErrInvalidAPICall
}
