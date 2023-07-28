package main

import (
	"fmt"
	"net/http"
	"strconv"

	"alire.me/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    // pat takes care of the url matching 

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

    app.render(w, r, "home.page.tmpl", &templateData{
        Snippets: s,
    })
}
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {

    id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err == models.ErrorNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

    app.render(w, r, "show.page.tmpl", &templateData{
        Snippet: s,
    })

}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Create a new snippet..."))
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		// WriteHeader is only possible to call once per response
		// the default header status code is always 200 OK
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O Snail"
	content := "O Snail\nClimb mountain fuji\n but slowly slowly"
	expires := "7"
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
