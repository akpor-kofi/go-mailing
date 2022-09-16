package ports

import "github.com/akpor-kofi/logger/models"

type LogEntryRepository interface {
	Insert(logEntry *models.LogEntry) error
	All() ([]*models.LogEntry, error)
	GetOne(id string) (*models.LogEntry, error)
	DropCollection() error
	Update(logEntry *models.LogEntry) error
}
