package web

import (
	"html/template"
	"net/http"

	"github.com/google/uuid"

	"github.com/Chasegwuap/goreddit"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func NewHandler(store goreddit.Store) *Handler {
	h := &Handler{

		Mux:   chi.NewMux(),
		Store: store,
	}
	h.Use(middleware.Logger)
	h.Route("/threads", func(r chi.Router) {
		r.Get("/", h.Threadslist())
		r.Get("/new", h.ThreadsCreate())
		r.Post("/", h.ThreadsStore())
		r.Post("/{id}/delete", h.ThreadsDelete())

	})

	h.Get("/html", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.ParseFiles("templates/layout.html"))

		type params struct {
			Title   string
			Text    string
			Lines   []string
			Number1 int
			Number2 int
		}

		t.Execute(w, params{
			Title: "Reddit clone",
			Text:  "Welcome to our 123 reddit clone ",
			Lines: []string{
				"Line1",
				"Line2",
				"Line3",
			},
			Number1: 421,
			Number2: 421,
		})

	})

	return h
}

type Handler struct {
	*chi.Mux

	Store goreddit.Store
}

const threadsListHTML = `
<h1>Threads</h1>
{{range .Threads}}
<dl>
  <dt><strong>{{.Title}}</strong></dt>
  <dd>{{.Description}}</dd> 
  <dd>
	<form action = "/threads/{{.ID}}/delete" method="POST">
		<button type="submit">Delete</button>
	</form>
	</dd>
{{end}}
</dl>
<a href="/threads/new">Create thread</a>s
`

func (h *Handler) Threadslist() http.HandlerFunc {
	type data struct {
		Threads []goreddit.Thread
	}

	tmpl := template.Must(template.New("").Parse(threadsListHTML))
	return func(w http.ResponseWriter, r *http.Request) {
		ttPtrs, err := h.Store.Threads()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert []*Thread → []Thread for the template
		tt := make([]goreddit.Thread, len(ttPtrs))
		for i, t := range ttPtrs {
			tt[i] = *t
		}

		tmpl.Execute(w, data{Threads: tt})
	}
}

const ThreadsCreateHTML = `
<h1>Create New Thread</h1>
<form action="/threads" method="POST">
  Title:<br>
  <input type="text" name="title"><br><br>

  Description:<br>
  <textarea name="description"></textarea><br><br>

  <button type="submit">Create Thread</button>
</form>


`

func (h *Handler) ThreadsCreate() http.HandlerFunc {
	tmpl := template.Must(template.New("").Parse(ThreadsCreateHTML))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}

}

func (h *Handler) ThreadsStore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("title")
		description := r.FormValue("description")

		if err := h.Store.CreateThread(&goreddit.Thread{
			ID:          uuid.New(),
			Title:       title,
			Description: description,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/threads", http.StatusFound)
	}
}

func (h *Handler) ThreadsDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idstr := chi.URLParam(r, "id")

		id, err := uuid.Parse(idstr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		}
		if err := h.Store.DeleteThread(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}

		http.Redirect(w, r, "/threads", http.StatusFound)
	}

}
