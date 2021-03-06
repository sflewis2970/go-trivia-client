package controllers

import (
	"log"
	"net/http"

	"github.com/sflewis2970/go-trivia-client/models"
	"github.com/sflewis2970/go-trivia-client/views"
)

var viewsDir string = "views"

type Controller struct {
	triviaView     *views.View
	triviaFAQView  *views.View
	triviaModel    *models.Model
	triviaFAQModel *models.Model
}

var controller *Controller

func New() *Controller {
	controller = new(Controller)

	controller.createModels()
	controller.createViews()
	controller.defineRoutes()

	return controller
}

func (c *Controller) createModels() {
	controller.triviaModel = models.New()
	controller.triviaFAQModel = models.New()
}

func (c *Controller) createViews() {
	c.triviaView = views.New(viewsDir, "base", viewsDir+"/partial/trivia.gohtml", viewsDir+"/partial/select.gohtml", viewsDir+"/partial/alert.gohtml")
	c.triviaFAQView = views.New(viewsDir, "base", viewsDir+"/partial/faq.gohtml")
}

func (c *Controller) defineRoutes() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/trivia", TriviaHandler)
	http.HandleFunc("/trivia/faq", FaqHandler)
	http.HandleFunc("/trivia/newquestion", NewQuestionHandler)
	http.HandleFunc("/trivia/answerquestion", AnswerQuestionHandler)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Redirect request
	http.Redirect(w, r, "/trivia", http.StatusMovedPermanently)
}

func TriviaHandler(w http.ResponseWriter, r *http.Request) {
	// Reset data model
	controller.triviaModel.ResetModel()

	// Render view with updated data
	controller.triviaView.Render(w, controller.triviaModel)
}

func FaqHandler(w http.ResponseWriter, r *http.Request) {
	// Render view with initially created data model (no update required)
	controller.triviaFAQView.Render(w, controller.triviaFAQModel)
}

func NewQuestionHandler(w http.ResponseWriter, r *http.Request) {
	// Get new question from trivia service
	controller.triviaModel.NewQuestion()

	// Render view with updated data
	controller.triviaView.Render(w, controller.triviaModel)
}

func AnswerQuestionHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	parseErr := r.ParseForm()
	if parseErr != nil {
		log.Print("Parsing error: ", parseErr)
		return
	}

	// Get data from form
	questionID := r.PostForm.Get("questionID")
	response := r.PostForm.Get("response")

	// Check answer with trivia service
	controller.triviaModel.AnswerQuestion(questionID, response)

	// Render view with updated data
	controller.triviaView.Render(w, controller.triviaModel)
}
