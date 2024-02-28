package endpoint

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
)

func HFUnmarshal[T any](
	w http.ResponseWriter,
	r *http.Request,
) (
	v T,
	cont bool,
) {
	if w == nil || r == nil {
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return v, false
	}

	defer r.Body.Close()

	err = json.Unmarshal(b, &v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return v, false
	}

	return v, true
}

func HFMarshal(w http.ResponseWriter, v any, ssc int) {
	if w == nil || v == nil {
		if ssc != 0 {
			w.WriteHeader(ssc)
		}
		return
	}

	b, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if ssc != 0 {
		w.WriteHeader(ssc)
	}

	w.Write(b)
}

func HFRespond(w http.ResponseWriter, v any, ssc int, err error) {
	switch {
	case err != nil && ssc == 0:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	case err != nil && ssc != 0:
		http.Error(w, err.Error(), ssc)
	case err == nil && ssc == 0:
		HFMarshal(w, v, http.StatusOK)
	case err == nil && ssc != 0:
		HFMarshal(w, v, ssc)
	}
}

func HFTryParseURLParam(r *http.Request, accessor string) string {
	ns := regexp.MustCompile("{(.+?)}").FindString(r.URL.String())
	return ns
}
