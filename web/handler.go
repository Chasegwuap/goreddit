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

	h.Get("/", h.Home())
	h.Route("/threads", func(r chi.Router) {
		r.Get("/", h.Threadslist())
		r.Get("/new", h.ThreadsCreate())
		r.Post("/", h.ThreadsStore())
		r.Get("/{id}", h.ThreadsShow())
		r.Post("/{id}/delete", h.ThreadsDelete())
		r.Get("/{id}/new", h.PostCreate())
		r.Post("/{id}", h.PostStore())
		r.Get("/{threadID}/{postID}", h.PostShow())

	})
	return h
}

type Handler struct {
	*chi.Mux

	Store goreddit.Store
}

func (h *Handler) Home() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/home.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}

}

func (h *Handler) Threadslist() http.HandlerFunc {
	type data struct {
		Threads []goreddit.Thread
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/threads.html"))
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

func (h *Handler) ThreadsCreate() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/thread_create.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}

}

func (h *Handler) ThreadsShow() http.HandlerFunc {
	type data struct {
		Thread goreddit.Thread
		Posts  []*goreddit.Post
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/thread.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.NotFound(w, r)
			return

		}

		t, err := h.Store.Thread(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pp, err := h.Store.PostByThread(t.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data{Thread: t, Posts: pp})
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

func (h *Handler) PostCreate() http.HandlerFunc {
	type data struct {
		Thread goreddit.Thread
	}
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/post_create.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		idstr := chi.URLParam(r, "id")

		id, err := uuid.Parse(idstr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		t, err := h.Store.Thread(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data{Thread: t})
	}

}
func (h *Handler) PostShow() http.HandlerFunc {
	type data struct {
		Thread goreddit.Thread
		Post   goreddit.Post
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/post_create.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		postIDSTR := chi.URLParam(r, "postID")
		threadIDSTR := chi.URLParam(r, "threadID")

		postID, err := uuid.Parse(postIDSTR)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		threadID, err := uuid.Parse(threadIDSTR)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		p, err := h.Store.Post(postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}
		t, err := h.Store.Thread(threadID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}

		tmpl.Execute(w, data{Thread: t, Post: p})
	}

}
func (h *Handler) PostStore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("title")
		content := r.FormValue("content")

		idstr := chi.URLParam(r, "id")

		id, err := uuid.Parse(idstr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		}
		t, err := h.Store.Thread(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		p := &goreddit.Post{
			ID:       uuid.New(),
			ThreadID: t.ID,
			Title:    title,
			Content:  content,
		}

		if err := h.Store.CreatePost(p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/threads/"+t.ID.String()+"/"+p.ID.String(), http.StatusFound)
	}
}
