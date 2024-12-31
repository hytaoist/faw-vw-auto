package http

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"

	. "github.com/hytaoist/autosignin/domain"
	"github.com/hytaoist/autosignin/internal/log"
)

type api struct {
	use Usecaser
}

func newAPI(use Usecaser) *api {
	return &api{use}
}

func (a *api) versions() httprouter.Handle {
	type request struct {
		Product string
	}
	type response struct {
		Versions []string `json:"versions"`
	}
	input := func(r *http.Request) (*request, error) {
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			err = errors.WithStack(err)
			return nil, err
		}
		if req.Product == "" {
			return nil, errors.New("product is missing")
		}
		return req, nil
	}
	process := func(req *request) ([]string, error) {
		data, err := a.use.Versions(req.Product)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	output := func(w http.ResponseWriter, data []string) error {
		resp := &response{}
		resp.Versions = data
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			err = errors.WithStack(err)
			return err
		}
		return nil
	}
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		req, err := input(r)
		if err != nil {
			err = errors.WithMessage(err, "bad request")
			log.Debug(err)
			http.Error(w, "Bad Request", 400)
			return
		}
		data, err := process(req)
		if err != nil {
			log.Error(err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
		err = output(w, data)
		if err != nil {
			log.Error(err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
	}
}