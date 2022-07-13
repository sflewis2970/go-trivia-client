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

type MessageResponse struct {
	QuestionID          string
	QuestionIDInputCtrl template.HTML
	Question            string
	QuestionInputCtrl   template.HTML
	Category            string
	CategoryInputCtrl   template.HTML
	Choices             []string
	ChoicesInputCtrl    template.HTML
	DisplayAlert        bool
	CorrectAnswer       bool
	Answer              string
	Response            string
	Message             string
}

type Model struct {
	CfgData         *config.ConfigData
	Action          string
	MessageResponse MessageResponse
}

func (m *Model) NewQuestion() {
	// Call Trivia Service to get new question
	log.Print("Retrieve new question from trivia service...")
	getQuestionURL := m.CfgData.TriviaServiceName + ":" + m.CfgData.TriviaServicePort + messages.GET_QUESTION_PATH
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

	// Reset alert
	m.MessageResponse.DisplayAlert = false
	m.MessageResponse.CorrectAnswer = false

	// Update action
	m.Action = NEW_QUESTION

	// Question ID will come from Trivia Service
	m.UpdateModelData(qResponse)
	m.UpdateModelControls()
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

	answerQuestionURL := m.CfgData.TriviaServiceName + ":" + m.CfgData.TriviaServicePort + messages.ANSWER_QUESTION_PATH
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

	// Update results message field
	m.MessageResponse.DisplayAlert = true
	m.MessageResponse.CorrectAnswer = aResponse.Correct
	m.MessageResponse.Answer = aResponse.Answer
	m.MessageResponse.Response = aResponse.Response

	log.Print("Message: ", aResponse.Message)
	m.MessageResponse.Message = aResponse.Message

	// Update action
	m.Action = ANSWER_QUESTION

	// Question ID will come from Trivia Service
	m.UpdateModelData(aResponse)
	m.UpdateModelControls()
}

func (m *Model) UpdateModelData(msg any) {
	if m.Action == RESET_MODEL_DATA || m.Action == NEW_QUESTION {
		response := msg.(messages.QuestionResponse)

		m.MessageResponse.QuestionID = response.QuestionID
		m.MessageResponse.Question = response.Question
		if len(response.Category) > 0 {
			m.MessageResponse.Category = response.Category
		} else {
			m.MessageResponse.Category = "<unclassified>"
		}

		m.MessageResponse.Choices = response.Choices
	} else if m.Action == ANSWER_QUESTION {
		response := msg.(messages.AnswerResponse)

		m.MessageResponse.Question = response.Question
		m.MessageResponse.Category = response.Category
		if len(response.Category) > 0 {
			m.MessageResponse.Category = response.Category
		} else {
			m.MessageResponse.Category = "<unclassified>"
		}
	}
}

func (m *Model) UpdateModelControls() {
	// Question ID control
	questionIDInputCtrl := "<input type='hidden' class='form-control' id='inputHiddenQuestionID' name='questionID' value='" + m.MessageResponse.QuestionID + "'>"
	m.MessageResponse.QuestionIDInputCtrl = template.HTML(questionIDInputCtrl)

	// Question control
	questionInputCtrl := "<input type='text' readonly class='form-control-plaintext' id='staticQuestion' value='" + m.MessageResponse.Question + "'> "
	m.MessageResponse.QuestionInputCtrl = template.HTML(questionInputCtrl)

	// Category control
	categoryInputCtrl := "<input type='text' readonly class='form-control-plaintext' id='staticCategory' value='" + m.MessageResponse.Category + "'> "
	m.MessageResponse.CategoryInputCtrl = template.HTML(categoryInputCtrl)

	// Choices control
	choicesInputCtrl := "<select id='idResponse' name='response'>"
	for idx, choice := range m.MessageResponse.Choices {
		if idx == 0 {
			choicesInputCtrl = choicesInputCtrl + "<option selected>" + choice + "</option>"
		} else {
			choicesInputCtrl = choicesInputCtrl + "<option value='" + choice + "'>" + choice + "</option>"
		}
	}
	choicesInputCtrl = choicesInputCtrl + "</select>"
	m.MessageResponse.ChoicesInputCtrl = template.HTML(choicesInputCtrl)
}

func (m *Model) ResetModel() {
	var qResponse messages.QuestionResponse

	// Reset alert
	m.MessageResponse.DisplayAlert = false
	m.MessageResponse.CorrectAnswer = false

	m.Action = RESET_MODEL_DATA
	qResponse.Question = "Question will be placed here"
	qResponse.Category = "Category will be placed here"
	m.UpdateModelData(qResponse)
	m.UpdateModelControls()
}

func (m *Model) GetConfigData() {
	var cfgDataErr error
	m.CfgData, cfgDataErr = config.Get().GetData()
	if cfgDataErr != nil {
		log.Print("Error getting response data: ", cfgDataErr)
		return
	}
}

func New() *Model {
	model := new(Model)

	model.GetConfigData()

	return model
}
