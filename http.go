package stubby

// TODO: JWT protect meta data and creation.

import (
	"net/http"

	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
	"github.com/y0ssar1an/q"
)

func Handler(conn *Conn) http.Handler {
	router := chi.NewRouter()

	router.Get("/:id", func(w http.ResponseWriter, r *http.Request) {
		q.Q(chi.URLParam(r, "id"))
		u, err := conn.Get(chi.URLParam(r, "id"))
		if err != nil {
			render.JSON(w, r, err)
			return
		}
		// TODO: Redirect instead of JSON.
		render.JSON(w, r, u)
	})

	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		u, err := conn.Create(r.PostFormValue("url"))
		if err != nil {
			q.Q(err)
			render.JSON(w, r, err)
			return
		}
		render.JSON(w, r, u)
	})
	// TODO: /s/:jwt for secure JWT authenticated redirects? The :id is embedded in the claim.
	// TODO: /m/:jwt for secure JWT metadata for a url visits, json metadata, jwt claims etc.
	return router
}
