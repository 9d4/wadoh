package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/9d4/wadoh/html"
	"github.com/9d4/wadoh/internal"
	"github.com/rs/zerolog/log"
)

var (
	errBadRequest = internal.NewError(
		internal.EBADINPUT,
		"Your input seems incorrect, please check before try again",
		"bad_request",
	)
)

var errorKindStatus = map[internal.ErrorKind]int{
	internal.EINTERNAL: http.StatusInternalServerError,
	internal.ENOTFOUND: http.StatusNotFound,
	internal.EBADINPUT: http.StatusBadRequest,
}

func errorKindToStatus(kind internal.ErrorKind) int {
	code, ok := errorKindStatus[kind]
	if ok {
		return code
	}
	return http.StatusInternalServerError
}

func Error(s *Server, w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
	log.Debug().Caller().Err(err).Send()

	// parse http related error
	tmpl := &html.ErrorTmpl{}
	e := parseError(err)
	if e != nil {
		tmpl.Code = e.Code()
		tmpl.Message = e.Error()
		tmpl.Status = errorKindToStatus(e.Kind())
	} else {
		kind, message, code := internal.ParseError(err)
		tmpl.Code = code
		tmpl.Message = message
		tmpl.Status = errorKindToStatus(kind)
	}

	w.WriteHeader(tmpl.Status)
	err = s.templates.R(r.Context(), w, tmpl)
	if err != nil {
		fmt.Fprintf(w, "Something went wrong and unable to render the page. Here's some messages: %s. Code: %s", tmpl.Message, tmpl.Code)
		return
	}
}

// parseError parses error in http layer
func parseError(err error) (e *internal.Error) {
	var strconvErr *strconv.NumError
	if errors.As(err, &strconvErr) {
		return errBadRequest
	}

	return
}
