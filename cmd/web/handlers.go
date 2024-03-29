package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

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
    app.render(w, r, "create.page.tmpl", nil)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

    err := r.ParseForm()
    if err != nil { 
        app.clientError(w, http.StatusBadRequest)
        return
    }

    title := r.PostForm.Get("title")
    content := r.PostForm.Get("content")
    expires := r.PostForm.Get("expires")

    errors := make(map[string]string)

    if strings.TrimSpace(title) == "" { 
        errors["title"] = "This feild cannot be blank"
    } else if utf8.RuneCountInString(title) > 100 { 
        errors["title"] = "This field is too long"
    }

    if strings.TrimSpace(content) == "" { 
        errors["content"] = "This feild cannot be blank"
    }

    if strings.TrimSpace(expires) == "" { 
        errors["expires"] = "This feild cannot be blank"
    } else if expires != "365" && expires != "7" && expires != "1" { 
        errors["expires"] = "This field is invalid"
    }

    if len(errors) > 0 { 
        app.render(w, r, "create.page.tmpl", &templateData{
            FormErrors: errors,
            FormData: r.PostForm,
        })
        return
    }

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
