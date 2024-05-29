package html

import (
	"context"
	"fmt"
	"html/template"
	"io"

	"github.com/rs/zerolog/log"
)

type ErrorTmpl struct {
	Message string
	Code    string
}

func (t *ErrorTmpl) Render(ctx context.Context, w io.Writer) {
	sendText := func() {
		fmt.Fprintf(w, "Something went wrong and unable to render the page. Here's some messages: %s. Code: %s", t.Message, t.Code)
	}
	tmp, err := template.ParseFS(TemplatesFS(), "layouts/base.html", "pages/error.html", "templates/*.html")
	if err != nil {
		sendText()
		log.Debug().Caller().Err(err).Send()
		return
	}

	err = tmp.Execute(w, t)
	if err != nil {
		sendText()
		log.Debug().Caller().Err(err).Send()
		return
	}
}
