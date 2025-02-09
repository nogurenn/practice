package simpleeventworker

import (
	"encoding/json"
	"io"
)

type Service interface {
	// ParseEvents parses a list of events from an io.Reader and returns a list of events.
	// The requirements assume events are already in the correct order.
	// The requirements also assume the input JSON is an array.
	ParseEvents(r io.Reader) ([]Event, error)
	// ProcessEvents processes a list of events and returns a map of accounts reduced to their final state.
	// This function should return a map of accounts with their ID as the key.
	// The requirements assume accounts are created prior to any charges, payments, etc.
	//
	// Idempotent: If processing an event results in an error, the function should stop processing events and return the error.
	ProcessEvents(events []Event) (map[string]Account, error)
}

type EventService struct{}

func NewService() *EventService {
	return &EventService{}
}

const (
	EventTypeAccountCreated         = "AccountCreated"
	EventTypeAccountChargeReceived  = "AccountChargeReceived"
	EventTypeAccountPaymentReceived = "AccountPaymentReceived"
	EventTypeAccountRecalled        = "AccountRecalled"

	EventPayloadFieldAmount  = "Amount"
	EventPayloadFieldBalance = "Balance"
)

type Event struct {
	Type      string       `json:"Type"`
	AccountID string       `json:"AccountID"`
	Payload   EventPayload `json:"Payload"`
}

type EventPayload interface {
	json.Unmarshaler
}

type EventPayloadAccountCreated struct {
	Balance int `json:"Balance"`
}

// EventPayloadAccountTransactionReceived represents the payload for the `AccountChargeReceived` and `AccountPaymentReceived` events.
type EventPayloadAccountTransactionReceived struct {
	Amount int `json:"Amount"`
}

func (p *EventPayloadAccountCreated) UnmarshalJSON(data []byte) error {
	aux := map[string]int{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	balance, ok := aux[EventPayloadFieldBalance]
	if !ok {
		return &ErrMissingFieldInEventPayloadField{Field: EventPayloadFieldBalance}
	}

	p.Balance = balance

	return nil
}

func (p *EventPayloadAccountTransactionReceived) UnmarshalJSON(data []byte) error {
	aux := map[string]int{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	amount, ok := aux[EventPayloadFieldAmount]
	if !ok {
		return &ErrMissingFieldInEventPayloadField{Field: EventPayloadFieldAmount}
	}

	p.Amount = amount

	return nil
}

func (e *Event) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Type      string          `json:"Type"`
		AccountID string          `json:"AccountID"`
		Payload   json.RawMessage `json:"Payload"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var payload EventPayload
	switch aux.Type {
	case EventTypeAccountCreated:
		temp := &EventPayloadAccountCreated{}
		if err := json.Unmarshal(aux.Payload, temp); err != nil {
			return err
		}
		payload = temp
	case EventTypeAccountChargeReceived, EventTypeAccountPaymentReceived:
		temp := &EventPayloadAccountTransactionReceived{}
		if err := json.Unmarshal(aux.Payload, temp); err != nil {
			return err
		}
		payload = temp
	case EventTypeAccountRecalled:
		// No payload required. If network costs are a concern, we can enforce byte size limits for aux.Payload.
		payload = nil
	default:
		return &ErrUnsupportedEventType{Type: aux.Type}
	}

	e.Type = aux.Type
	e.AccountID = aux.AccountID
	e.Payload = payload

	return nil
}

func (s *EventService) ParseEvents(r io.Reader) ([]Event, error) {
	events := []Event{}

	decoder := json.NewDecoder(r)
	token, err := decoder.Token()
	if err != nil {
		return nil, err
	}

	if token != json.Delim('[') {
		return nil, ErrInputJSONIsNotArray
	}

	for decoder.More() {
		var event Event
		if err := decoder.Decode(&event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (s *EventService) ProcessEvents(events []Event) (map[string]Account, error) {
	accounts := map[string]Account{}

	for _, event := range events {
		switch event.Type {
		case EventTypeAccountCreated:
			if err := s.processEventTypeAccountCreated(event, accounts); err != nil {
				return nil, err
			}
		case EventTypeAccountChargeReceived:
			if err := s.processEventTypeAccountChargeReceived(event, accounts); err != nil {
				return nil, err
			}
		case EventTypeAccountPaymentReceived:
			if err := s.processEventTypeAccountPaymentReceived(event, accounts); err != nil {
				return nil, err
			}
		case EventTypeAccountRecalled:
			if err := s.processEventTypeAccountRecalled(event, accounts); err != nil {
				return nil, err
			}
		default:
			return nil, &ErrUnsupportedEventType{Type: event.Type}
		}
	}

	return accounts, nil
}

func (_ EventService) processEventTypeAccountCreated(event Event, accounts map[string]Account) error {
	if _, ok := accounts[event.AccountID]; ok {
		return &ErrAccountAlreadyExists{AccountID: event.AccountID}
	}

	account := &Account{
		ID: event.AccountID,
	}
	if err := account.RecordTransaction(event.Payload.(*EventPayloadAccountCreated).Balance); err != nil {
		return err
	}

	accounts[event.AccountID] = *account

	return nil
}

func (_ EventService) processEventTypeAccountRecalled(event Event, accounts map[string]Account) error {
	account, ok := accounts[event.AccountID]
	if !ok {
		return &ErrAccountDoesNotExist{AccountID: event.AccountID}
	}

	account.status = AccountStatusRecalled
	accounts[event.AccountID] = account

	return nil
}

func (_ EventService) processEventTypeAccountChargeReceived(event Event, accounts map[string]Account) error {
	account, ok := accounts[event.AccountID]
	if !ok {
		return &ErrAccountDoesNotExist{AccountID: event.AccountID}
	}

	if err := account.RecordTransaction(event.Payload.(*EventPayloadAccountTransactionReceived).Amount); err != nil {
		return err
	}

	accounts[event.AccountID] = account

	return nil
}

func (_ EventService) processEventTypeAccountPaymentReceived(event Event, accounts map[string]Account) error {
	account, ok := accounts[event.AccountID]
	if !ok {
		return &ErrAccountDoesNotExist{AccountID: event.AccountID}
	}

	if err := account.RecordTransaction(-event.Payload.(*EventPayloadAccountTransactionReceived).Amount); err != nil {
		return err
	}

	accounts[event.AccountID] = account

	return nil
}
