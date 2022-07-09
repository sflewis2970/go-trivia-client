package controllers

import (
	"net/http"

	"github.com/sflewis2970/go-trivia-client/views"
)

var viewsDir string = "views"

type Controller struct {
	trivia *views.View
}

var controller *Controller

func New() *Controller {
	controller = &Controller{}

	controller.defineRoutes()

	return controller
}

func (c *Controller) defineRoutes() {
	c.trivia = views.New(viewsDir, "base", viewsDir+"/partial/trivia.gohtml")

	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/trivia", TriviaHandler)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/trivia", http.StatusMovedPermanently)
}

func TriviaHandler(w http.ResponseWriter, r *http.Request) {
	controller.trivia.Render(w, nil)
}
