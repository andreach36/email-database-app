package emails

import "github.com/go-chi/chi"

func EmailsRoutes() chi.Router {
	r := chi.NewRouter()
	emailHandler := EmailHandler{}

	r.Get("/", emailHandler.GetAllEmails)
	r.Get("/search", emailHandler.SearchEmails)

	return r
}
