package services

import (
	"database/sql"
	"events/pkg/events"
	"events/pkg/storage/mock"
	"testing"
)

const (
	uid = 1
	eventID = 1
)

func TestEventService_CreateEvent(t *testing.T) {
	tests := []struct{
		name string
		event *events.Event
		wantedErr error
	} {
		{
			name:      "valid title only",
			event: &events.Event{Title: "Some valid title"},
			wantedErr: nil,
		},
		{
			name:      "long title",
			event: &events.Event{Title: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
			wantedErr: nil,
		},
	}

	for _, tt := range tests {
		db, err := mock.NewDB()
		if err != nil {
			t.Fatal(err)
		}

		tx, err := db.Begin()
		if err != nil {
			t.Fatal(err)
		}

		eventRepo := mock.NewEventStorage()

		eventRepo.TxFn = func() (*sql.Tx, error) {
			return tx, nil
		}

		eventRepo.SaveEventTxFn = func(tx *sql.Tx, title string, uid int) (int, error) {
			return eventID, nil
		}

		eventRepo.UpdateEventTxFn = func(tx *sql.Tx, i *events.Event) error {
			return nil
		}

		service := NewEventService(eventRepo)
		t.Run(tt.name, func(t *testing.T) {
			id, err := service.CreateEvent(tt.event, uid)
			if id != 1 {
				t.Fatalf("wanted %d, got %v", 1, id)
			}

			if err != tt.wantedErr {
				t.Fatalf("wanted %v, got %v", tt.wantedErr, err)
			}

			if !eventRepo.SaveEventTxInvoked {
				t.Fatalf("SaveEventTx method was not invoked")
			}

			if !eventRepo.UpdateEventTxInvoked {
				t.Fatalf("UpdateEventTx method was not invoked")
			}
		})
	}

}
