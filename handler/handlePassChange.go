package handler

import (
	"context"
	"net/http"
	"medicalTesting/core"
	"medicalTesting/serverErr"
)

func handlePassChange(ctx context.Context, r *http.Request, userUID string) (response interface{}, err error) {
	if r.Method == http.MethodPost {
		return nil, core.ChangePass(ctx, r.Body, userUID)
	}

	return nil, serverErr.ErrInvalidAPICall
}
