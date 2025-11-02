package entity

type Book struct {
	Id      string   `json:"id"`
	Title   string   `json:"title"`
	Authors []string `json:"authors"`
	Text    string   `json:"text"`
}
