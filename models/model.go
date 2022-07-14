package models

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/sflewis2970/go-trivia-client/config"
	"github.com/sflewis2970/go-trivia-client/messages"
)

const (
	RESET_MODEL_DATA string = "RESET_MODEL_DATA"
	NEW_QUESTION     string = "NEW_QUESTION"
	ANSWER_QUESTION  string = "ANSWER_QUESTION"
)

type Responses struct {
	QuestionID    string
	Question      string
	Category      string
	Choices       []string
	RecordFound   bool
	CorrectAnswer bool
	Answer        string
	UserResponse  string
	Message       string
}

type WebControls struct {
	QuestionIDCtrl template.HTML
	CategoryCtrl   template.HTML
	ChoicesCtrl    template.HTML
	DisplayAlert   bool
}

type Model struct {
	CfgData           *config.ConfigData
	IsAppInProduction bool
	Action            string
	Response          Responses
	WebControl        WebControls
}

// Exported type methods
func (m *Model) NewQuestion() {
	// Call Trivia Service to get new question
	log.Print("Retrieve new question from trivia service...")
	var getQuestionURL string
	if m.CfgData.Env == config.PRODUCTION {
		getQuestionURL = m.CfgData.TriviaServiceName + messages.GET_QUESTION_PATH
	} else {
		getQuestionURL = m.CfgData.TriviaServiceName + ":" + m.CfgData.TriviaServicePort + messages.GET_QUESTION_PATH
	}
	log.Print("getQuestionURL: ", getQuestionURL)

	response, getErr := http.Get(getQuestionURL)
	if getErr != nil {
		log.Print("Error getting response data: ", getErr)
		return
	}
	defer response.Body.Close()

	// Handle QuestionRresponse
	var qResponse messages.QuestionResponse

	// Read response stream into JSON
	json.NewDecoder(response.Body).Decode(&qResponse)

	// Update action
	m.Action = NEW_QUESTION

	// Question ID will come from Trivia Service
	m.updateModelData(qResponse)
	m.updateModelControls()
}

func (m *Model) AnswerQuestion(questionID string, userResponse string) {
	// Call Trivia Service to answer question
	log.Print("Send answer to trivia service...")
	var aRequest messages.AnswerRequest

	// Build add question request
	aRequest.QuestionID = questionID
	aRequest.Response = userResponse

	// Convert struct to byte array
	requestBody, marshalErr := json.Marshal(aRequest)
	if marshalErr != nil {
		log.Print("marshaling error: ", marshalErr)
		return
	}

	var answerQuestionURL string
	if m.CfgData.Env == config.PRODUCTION {
		answerQuestionURL = m.CfgData.TriviaServiceName + messages.ANSWER_QUESTION_PATH
	} else {
		answerQuestionURL = m.CfgData.TriviaServiceName + ":" + m.CfgData.TriviaServicePort + messages.ANSWER_QUESTION_PATH
	}
	log.Print("answerQuestionURL: ", answerQuestionURL)

	response, postErr := http.Post(answerQuestionURL, "application/json", bytes.NewBuffer(requestBody))
	if postErr != nil {
		log.Print("Error posting message: ", postErr)
		return
	}
	defer response.Body.Close()

	// Handle AnswerResponse
	var aResponse messages.AnswerResponse

	// Read response stream into JSON
	json.NewDecoder(response.Body).Decode(&aResponse)

	// Update action
	m.Action = ANSWER_QUESTION

	// Question ID will come from Trivia Service
	m.updateModelData(aResponse)
	m.updateModelControls()
}

func (m *Model) ResetModel() {
	var qResponse messages.QuestionResponse

	// Update action
	m.Action = RESET_MODEL_DATA

	// Update response fields
	qResponse.Question = "Question will be placed here"
	qResponse.Category = "Category will be placed here"
	m.updateModelData(qResponse)
	m.updateModelControls()
}

// Unexported type methods
func (m *Model) updateModelData(msg any) {
	if m.Action == RESET_MODEL_DATA || m.Action == NEW_QUESTION {
		// Set Alert flag
		m.WebControl.DisplayAlert = false

		// Convert message to QuestionResponse
		response := msg.(messages.QuestionResponse)

		// Set Response fields
		m.Response.CorrectAnswer = false
		m.Response.RecordFound = false
		m.Response.QuestionID = response.QuestionID
		m.Response.Question = response.Question
		if len(response.Category) > 0 {
			m.Response.Category = response.Category
		} else {
			m.Response.Category = "<unclassified>"
		}

		m.Response.Choices = response.Choices
	} else if m.Action == ANSWER_QUESTION {
		// Set Alert flag
		m.WebControl.DisplayAlert = true

		// Convert message to AnswerResponse
		response := msg.(messages.AnswerResponse)

		// Set Response fields
		m.Response.CorrectAnswer = response.Correct
		m.Response.Answer = response.Answer
		m.Response.Message = response.Message

		m.Response.Question = response.Question
		if len(m.Response.Question) > 0 {
			m.Response.RecordFound = true
			m.Response.UserResponse = response.Response
		} else {
			m.Response.RecordFound = false
			m.Response.UserResponse = "Record not found."
		}

		m.Response.Category = response.Category
		if len(response.Category) > 0 {
			m.Response.Category = response.Category
		} else {
			m.Response.Category = "<unclassified>"
		}
	}
}

func (m *Model) updateModelControls() {
	// Question ID control
	questionIDCtrl := "<input type='hidden' class='form-control' id='inputHiddenQuestionID' name='questionID' value='" + m.Response.QuestionID + "'>"
	m.WebControl.QuestionIDCtrl = template.HTML(questionIDCtrl)

	// Category control
	categoryCtrl := "<input type='text' readonly class='form-control-plaintext' id='staticCategory' value='" + m.Response.Category + "'> "
	m.WebControl.CategoryCtrl = template.HTML(categoryCtrl)

	// Choices control
	choicesCtrl := "<select id='idResponse' name='response'>"
	for idx, choice := range m.Response.Choices {
		if idx == 0 {
			choicesCtrl = choicesCtrl + "<option selected>" + choice + "</option>"
		} else {
			choicesCtrl = choicesCtrl + "<option value='" + choice + "'>" + choice + "</option>"
		}
	}
	choicesCtrl = choicesCtrl + "</select>"
	m.WebControl.ChoicesCtrl = template.HTML(choicesCtrl)
}

func (m *Model) GetConfigData() {
	var cfgDataErr error
	m.CfgData, cfgDataErr = config.Get().GetData()
	if cfgDataErr != nil {
		log.Print("Error getting response data: ", cfgDataErr)
		return
	}

	if m.CfgData.Env == config.PRODUCTION {
		m.IsAppInProduction = true
	}
}

func New() *Model {
	model := new(Model)

	model.GetConfigData()

	return model
}
