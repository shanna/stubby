package stubby

// TODO: Break admin into its own handler so auth can be removed?
// TODO: Logging.

import (
	"net/http"

	"github.com/goware/jwtauth"
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
)

func Handler(shortener *Shortener, secret []byte) http.Handler {
	router := chi.NewRouter()

	router.Get("/:id", func(w http.ResponseWriter, r *http.Request) {
		record, err := shortener.Get(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Redirect(w, r, record.URL, http.StatusPermanentRedirect)
	})

	jwt := jwtauth.New("HS256", secret, nil)
	router.Group(func(router chi.Router) {
		router.Use(jwt.Verifier)
		router.Use(jwtauth.Authenticator)

		router.Post("/", func(w http.ResponseWriter, r *http.Request) {
			record, err := shortener.Create(r.PostFormValue("url"), nil)
			if err != nil {
				render.JSON(w, r, err)
				return
			}
			render.JSON(w, r, record)
		})
	})
	return router
}
