package model

type Item struct {
	ID        string `json:"-"`
	Title     string `json:"title,omitempty"`
	IsDeleted bool   `json:"-"`
	CreatedAt string `json:"-"`
}
