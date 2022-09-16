package logging

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	err := mgm.SetDefaultConfig(nil, "logger", options.Client().ApplyURI("mongodb://root:example@logger-mongo-srv"))
	must(err)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
