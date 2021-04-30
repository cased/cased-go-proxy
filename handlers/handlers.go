package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cased/cased-go"
)

type requestError struct {
	status  int
	Message string `json:"error"`
}

func (re requestError) Error() string {
	return re.Message
}

func AuditEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	auditEvent := cased.AuditEvent{}
	if err := loadJSON(r, auditEvent); err != nil {
		if err := writeError(w, err); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}

	if err := cased.Publish(auditEvent); err != nil {
		if err := writeError(w, err); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}

	if err := json.NewEncoder(w).Encode(auditEvent); err != nil {
		if err := writeError(w, err); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}
}

func writeError(w http.ResponseWriter, err error) error {
	var re *requestError
	if errors.As(err, &re) {
		d, err := json.Marshal(re)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		} else {
			http.Error(w, string(d), re.status)
		}
	} else {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	return nil
}

func loadJSON(r *http.Request, v map[string]interface{}) error {
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		return &requestError{http.StatusUnsupportedMediaType, "Content-Type header is not application/json"}
	}

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		var syntaxErr *json.SyntaxError

		switch {
		case errors.As(err, &syntaxErr):
			msg := fmt.Sprintf("request body contains malformed JSON (at position %d)", syntaxErr.Offset)
			return &requestError{http.StatusBadRequest, msg}
		case errors.Is(err, io.ErrUnexpectedEOF):
			return &requestError{http.StatusBadRequest, "request body contains invalid JSON"}
		case errors.Is(err, io.EOF):
			return &requestError{http.StatusBadRequest, "request body was empty"}
		default:
			return &requestError{http.StatusInternalServerError, "internal server error"}
		}
	}

	return nil
}
