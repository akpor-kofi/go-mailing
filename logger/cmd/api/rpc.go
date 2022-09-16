package main

import (
	"github.com/akpor-kofi/logger/models"
	"github.com/akpor-kofi/logger/ports"
)

type RPCServer struct {
	logStore ports.LogEntryRepository
}

type RPCPayload struct {
	Name string
	Data string
}

func NewRPCServer(logStore ports.LogEntryRepository) *RPCServer {
	return &RPCServer{logStore: logStore}
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	logEntry := models.LogEntry{
		Name: payload.Name,
		Data: payload.Data,
	}
	err := r.logStore.Insert(&logEntry)
	if err != nil {
		return err
	}

	*resp = "Processed payload via RPC: " + payload.Name
	return nil
}
