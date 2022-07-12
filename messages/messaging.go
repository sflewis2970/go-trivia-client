package messages

const (
	GET_QUESTION_PATH    string = "/api/v1/trivia/question"
	ANSWER_QUESTION_PATH string = "/api/v1/trivia/answer"
)

// Message definition
type QuestionResponse struct {
	QuestionID string   `json:"questionid"`
	Question   string   `json:"question"`
	Category   string   `json:"category"`
	Choices    []string `json:"choices"`
	Timestamp  string   `json:"timestamp"`
	Warning    string   `json:"warning,omitempty"`
	Error      string   `json:"error,omitempty"`
}

type AnswerRequest struct {
	QuestionID string `json:"questionid"`
	Response   string `json:"response"`
}

type AnswerResponse struct {
	Question  string `json:"question"`
	Timestamp string `json:"timestamp"`
	Category  string `json:"category"`
	Response  string `json:"response"`
	Answer    string `json:"answer"`
	Correct   bool   `json:"correct"`
	Message   string `json:"message,omitempty"`
	Warning   string `json:"warning,omitempty"`
	Error     string `json:"error,omitempty"`
}
