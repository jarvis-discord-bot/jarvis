package applies

type Apply struct {
	UserID    string      `json:"user_id"`
	Title     string      `json:"title"`
	URL       string      `json:"url"`
	Questions []*Question `json:"questions"`
}

type Question struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}
