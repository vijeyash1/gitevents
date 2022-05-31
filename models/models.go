package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type Models struct {
	DB DBModel
}

// NewModels returns models with db and nats jetsStream pool
func NewModels(js nats.JetStreamContext) Models {
	return Models{
		DB: DBModel{
			JS: js},
	}
}

type Gitevent struct {
	Uuid              uuid.UUID
	CommitedBy        string
	CommitedAt        time.Time
	Repository        string
	Commitstat        string
	Availablebranches string
	Commitmessage     string
}
