package logging

import (
	"context"
	"github.com/akpor-kofi/logger/models"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var ctx = context.Background()
var ms = 5 * time.Second

type logEntryStore struct {
	coll *mgm.Collection
}

func NewLogStore() *logEntryStore {
	logEntry := new(models.LogEntry)

	return &logEntryStore{
		coll: mgm.Coll(logEntry),
	}
}

func (l logEntryStore) Insert(logEntry *models.LogEntry) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, ms)
	defer cancel()

	if err := l.coll.CreateWithCtx(ctxWithTimeout, logEntry); err != nil {
		return err
	}

	return nil
}

func (l logEntryStore) All() ([]*models.LogEntry, error) {
	logs := new([]*models.LogEntry)

	opts := options.Find().SetSort(bson.D{{"created_at", -1}})

	err := l.coll.SimpleFind(logs, bson.D{}, opts)
	if err != nil {
		return []*models.LogEntry{}, err
	}

	return *logs, nil
}

func (l logEntryStore) GetOne(id string) (*models.LogEntry, error) {
	entry := new(models.LogEntry)

	err := l.coll.FindByID(id, entry)
	if err != nil {
		return &models.LogEntry{}, err
	}

	return entry, nil
}

func (l logEntryStore) DropCollection() error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, ms)
	defer cancel()

	err := l.coll.Drop(ctxWithTimeout)
	if err != nil {
		return err
	}

	return nil
}

func (l logEntryStore) Update(logEntry *models.LogEntry) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, ms)
	defer cancel()

	err := l.coll.UpdateWithCtx(ctxWithTimeout, logEntry)
	if err != nil {
		return err
	}

	return nil
}
