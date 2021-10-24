package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"golang.org/x/xerrors"
)

type selfValidator interface {
	validate() (code int, err error)
}

// bind function unmarshall http.Request.Body with json to data structure 'm'
// and validate received message.
func bind(r *http.Request, m selfValidator) (code int, err error) {
	if err = json.NewDecoder(r.Body).Decode(m); err != nil {
		return http.StatusBadRequest, xerrors.Errorf("request json decoding: %w", err)
	}

	log.Debug().
		Str("request_id", rID(r)).
		Interface("message", m).
		Msg(fmt.Sprintf("message '%T' received", m))

	return m.validate()
}

func writeResp(w http.ResponseWriter, reqID string, ret interface{}, code int, err error) {
		w.WriteHeader(code)

	lg := log.With().Str("component", "handlers").Str("request_id", reqID).Logger()
	resp := make(map[string]interface{})

	if ret != nil {
		resp["result"] = ret
	}

	if err != nil {
		lg.Err(err).Str("errors_stack", fmt.Sprintf("%v+", err)).Msg("request processing")

		resp["error"] = map[string]interface{}{
			"code":    code,
			"message": err.Error(),
		}
		resp["result"] = "error"
	} else if ret == nil {
		resp["result"] = "success"
	}

	b, err := json.Marshal(resp)
	if err != nil {
		b = []byte(`{"result":null,"error":{"code":599,"message":"marshal resp err","details":null}}`)

		lg.Err(err).Interface("resp", resp).Msg("response marshall err")
	}

	if _, err = w.Write(b); err != nil {
		lg.Error().Bytes("body", b).Err(err).Msg("body write err")
	}
}

func rID(r *http.Request) string {
	return r.Header.Get("X-RequestID")
}
