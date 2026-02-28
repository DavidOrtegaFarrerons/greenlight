package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

type envelop map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelop, headers http.Header) error {
	//MarshalIndent is 65% slower than Marshal, we use it here as it is a toy project, in a real project Marshal is better.
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	//Go does not throw an error if the headers map is nil and we try to loop over it
	for k, v := range headers {
		w.Header()[k] = v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
