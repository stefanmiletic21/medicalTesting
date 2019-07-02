package handler

import (
	"context"
	"net/http"
	"medicalTesting/core"
	"medicalTesting/serverErr"
	"strings"
)

func handleNurse(ctx context.Context, r *http.Request) (response interface{}, err error) {
	if strings.HasPrefix(r.URL.Path, "/person") {
		r.URL.Path = r.URL.Path[7:]
		return handlePerson(ctx, r)
	}
	if strings.HasPrefix(r.URL.Path, "/patient") {
		r.URL.Path = r.URL.Path[8:]
		switch r.Method {
		case http.MethodPost:
			return core.CreatePatient(ctx, r.Body)
		case http.MethodPatch:
			if strings.HasPrefix(r.URL.Path, "/") {
				return nil, core.UpdatePatient(ctx, r.URL.Path[1:], r.Body)
			}
			return nil, serverErr.ErrBadRequest
		case http.MethodGet:
			if strings.HasPrefix(r.URL.Path, "/") {
				return core.GetPatient(ctx, r.URL.Path[1:])
			}
			return core.GetPatients(ctx)
		case http.MethodDelete:
			if strings.HasPrefix(r.URL.Path, "/") {
				return nil, core.RemovePatient(ctx, r.URL.Path[1:])
			}
			return nil, serverErr.ErrBadRequest
		}
	}
	if strings.HasPrefix(r.URL.Path, "/examination") {
		switch r.Method {
		case http.MethodPost:
			return core.CreateExamination(ctx, r.Body)
		case http.MethodGet:
			return core.GetExaminations(ctx)
		case http.MethodDelete:
			r.URL.Path = r.URL.Path[12:]
			if strings.HasPrefix(r.URL.Path, "/") {
				return nil, core.RemoveExamination(ctx, r.URL.Path[1:])
			}
			return nil, serverErr.ErrBadRequest
		}
	}
	if strings.HasPrefix(r.URL.Path, "/doctor") {
		switch r.Method {
		case http.MethodGet:
			return core.GetDoctors(ctx)
		}
	}

	return nil, serverErr.ErrInvalidAPICall
}
