package routes

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/learnToCrypto/lakoposlati/internal/user"
)

// Convenience function to redirect to the error message page

// Checks if the user is logged in and has a session, if not err is not nil
func session(writer http.ResponseWriter, request *http.Request) (sess user.Session, err error) {
	cookie, err := request.Cookie("_cookie")
	if err == nil {
		sess = user.Session{Uuid: cookie.Value}
		if ok, _ := sess.Check(); !ok {
			err = errors.New("Invalid session")
		}
	}
	return
}

// parse HTML templates
// pass in a list of file names, and get a template
func parseTemplateFiles(filenames ...string) (t *template.Template) {
	var files []string
	t = template.New("layout")
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	t = template.Must(t.ParseFiles(files...))
	return
}

func generateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(writer, "layout", data)
}

func generateHTMLwithFunc(writer http.ResponseWriter, data interface{}, funcMap template.FuncMap, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	templates, err := template.New("myhtml").Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}
	templates.ExecuteTemplate(writer, "layout", data)
}
