package server

import (
	"embed"

	"github.com/google/go-safeweb/safehttp"
	"github.com/google/go-safeweb/safehttp/plugins/htmlinject"
	"github.com/google/safehtml/template"

	"github.com/aiit2022-pbl-okuhara/traning/micchie/sample-web-app/secure/auth"
	"github.com/aiit2022-pbl-okuhara/traning/micchie/sample-web-app/secure/responses"
	"github.com/aiit2022-pbl-okuhara/traning/micchie/sample-web-app/storage"
)

//go:embed static
var staticFiles embed.FS

//go:embed templates
var templatesFS embed.FS

var templates *template.Template

func init() {
	tplSrc := template.TrustedSourceFromConstant("templates/*.tpl.html")
	var err error
	// htmlinject automatically injects CSP nonces and XSRF tokens placeholders.
	templates, err = htmlinject.LoadGlobEmbed(nil, htmlinject.LoadConfig{}, tplSrc, templatesFS)
	if err != nil {
		panic(err)
	}
}

type serverDeps struct {
	db *storage.DB
}

func Load(db *storage.DB, mux *safehttp.ServeMux) {
	deps := &serverDeps{
		db: db,
	}

	// Private endpoints, only accessible to authenticated users (default).
	mux.Handle("/notes/", "GET", getNotesHandler(deps))
	mux.Handle("/notes", "POST", postNotesHandler(deps))
	mux.Handle("/logout", "POST", logoutHandler(deps))

	// Public enpoints, no auth checks performed.
	mux.Handle("/login", "POST", postLoginHandler(deps), auth.Skip{})
	mux.Handle("/static/", "GET", safehttp.FileServerEmbed(staticFiles), auth.Skip{})
	mux.Handle("/", "GET", indexHandler(deps), auth.Skip{})
}

func getNotesHandler(deps *serverDeps) safehttp.Handler {
	return safehttp.HandlerFunc(func(rw safehttp.ResponseWriter, r *safehttp.IncomingRequest) safehttp.Result {
		user := auth.User(r)
		notes := deps.db.GetNotes(user)
		return safehttp.ExecuteNamedTemplate(rw, templates, "notes.tpl.html", map[string]interface{}{
			"notes": notes,
			"user":  user,
		})
	})
}

func postNotesHandler(deps *serverDeps) safehttp.Handler {
	noFormErr := responses.NewError(
		safehttp.StatusBadRequest,
		template.MustParseAndExecuteToHTML(`Please submit a valid form with "title" and "text" parameters.`),
	)
	noFieldsErr := responses.NewError(
		safehttp.StatusBadRequest,
		template.MustParseAndExecuteToHTML("Both title and text must be specified."),
	)

	return safehttp.HandlerFunc(func(rw safehttp.ResponseWriter, r *safehttp.IncomingRequest) safehttp.Result {
		form, err := r.PostForm()
		if err != nil {
			return rw.WriteError(noFormErr)
		}
		title := form.String("title", "")
		body := form.String("text", "")
		if title == "" || body == "" {
			return rw.WriteError(noFieldsErr)
		}
		user := auth.User(r)
		deps.db.AddOrEditNote(user, storage.Note{Title: title, Text: body})

		notes := deps.db.GetNotes(user)
		return safehttp.ExecuteNamedTemplate(rw, templates, "notes.tpl.html", map[string]interface{}{
			"notes": notes,
			"user":  user,
		})
	})
}

func indexHandler(deps *serverDeps) safehttp.Handler {
	return safehttp.HandlerFunc(func(rw safehttp.ResponseWriter, r *safehttp.IncomingRequest) safehttp.Result {
		user := auth.User(r)
		if user != "" {
			return safehttp.Redirect(rw, r, "/notes/", safehttp.StatusTemporaryRedirect)
		}
		return safehttp.ExecuteNamedTemplate(rw, templates, "index.tpl.html", nil)
	})
}

// Logout and Login handlers would normally be centralized and provided by a separate package owned by the security team.
// Since this is a simple example application they are here together with the rest.
func logoutHandler(deps *serverDeps) safehttp.Handler {
	return safehttp.HandlerFunc(func(rw safehttp.ResponseWriter, r *safehttp.IncomingRequest) safehttp.Result {
		auth.ClearSession(r)
		return safehttp.Redirect(rw, r, "/", safehttp.StatusSeeOther)
	})
}

func postLoginHandler(deps *serverDeps) safehttp.Handler {
	invalidAuthErr := responses.NewError(
		safehttp.StatusBadRequest,
		template.MustParseAndExecuteToHTML("Please specify a username and a password, both must be non-empty and your password must match the one you use to register."),
	)

	// Always return the same error to not leak the existence of a user.
	return safehttp.HandlerFunc(func(rw safehttp.ResponseWriter, r *safehttp.IncomingRequest) safehttp.Result {
		form, err := r.PostForm()
		if err != nil {
			return rw.WriteError(invalidAuthErr)
		}
		username := form.String("username", "")
		password := form.String("password", "")
		if username == "" || password == "" {
			return rw.WriteError(invalidAuthErr)
		}
		if err := deps.db.AddOrAuthUser(username, password); err != nil {
			return rw.WriteError(invalidAuthErr)
		}
		auth.CreateSession(r, username)
		return safehttp.Redirect(rw, r, "/notes/", safehttp.StatusSeeOther)
	})
}
