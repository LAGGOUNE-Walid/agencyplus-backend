package gorseclient

import "time"

type Feedback struct {
	FeedbackType string    `json:"FeedbackType"`
	UserId       string    `json:"UserId"`
	ItemId       string    `json:"ItemId"`
	Timestamp    time.Time `json:"Timestamp"`
}

type ErrorMessage string

func (e ErrorMessage) Error() string {
	return string(e)
}

type RowAffected struct {
	RowAffected int `json:"RowAffected"`
}

type Score struct {
	Id    string  `json:"Id"`
	Score float64 `json:"Score"`
}

type UserIterator struct {
	Cursor string
	Users  []User
}

type User struct {
	UserId    string   `json:"UserId"`
	Labels    []string `json:"Labels"`
	Subscribe []string `json:"Subscribe"`
	Comment   string   `json:"Comment"`
}

type UserPatch struct {
	Labels    []string
	Subscribe []string
	Comment   *string
}

type ItemIterator struct {
	Cursor string
	Items  []Item
}

type Item struct {
	ItemId     string    `json:"ItemId"`
	IsHidden   bool      `json:"IsHidden"`
	Labels     any       `json:"Labels"`
	Categories []string  `json:"Categories"`
	Timestamp  time.Time `json:"Timestamp"`
	Comment    string    `json:"Comment"`
}

type ItemPatch struct {
	IsHidden   *bool
	Categories []string
	Timestamp  *time.Time
	Labels     any
	Comment    *string
}
