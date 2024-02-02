package handlers

import (
	"net/http"

	"github.com/davidhalasz/gomath/cmd/web/internal/render"
)

func AiPage(w http.ResponseWriter, r *http.Request) {
	if err := render.Template(w, r, "mi.page.gohtml", nil); err != nil {
		app.ErrorLog.Println(err)
	}
}
