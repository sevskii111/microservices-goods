package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")
)

type Good struct {
	ID          int
	Name        string
	Description string
	Price       int
	InStock     bool
	Created     time.Time
}

type GoodUpdate struct {
	Name        *string
	Description *string
	Price       *int
	InStock     *bool
}

func (u *GoodUpdate) IsEmpty() bool {
	return u.Name == nil && u.Description == nil && u.Price == nil && u.InStock == nil
}
