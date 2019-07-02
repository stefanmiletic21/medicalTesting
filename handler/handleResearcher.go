package handler

import (
	"context"
	"net/http"
	"medicalTesting/core"
	"medicalTesting/serverErr"
	"strings"
)

func handleResearcher(ctx context.Context, r *http.Request) (response interface{}, err error) {
	if strings.HasPrefix(r.URL.Path, "/filled") {
		r.URL.Path = r.URL.Path[7:]
		switch r.Method {
		case http.MethodGet:
			if strings.HasPrefix(r.URL.Path, "/") {
				return core.GetFilledTest(ctx, r.URL.Path[1:])
			}
			return core.GetFilledTests(ctx)
		}
	}
	return nil, serverErr.ErrInvalidAPICall
}
