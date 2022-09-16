package models

import "github.com/kamva/mgm/v3"

type LogEntry struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `json:"name" bson:"name"`
	Data             string `json:"data" bson:"data"`
}
