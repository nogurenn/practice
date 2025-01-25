package simpleeventworker_test

import (
	"io"
	"strings"
	"testing"

	event "github.com/nogurenn/assorted-programs/simple-event-worker"
	"github.com/stretchr/testify/assert"
)

func TestEvent_ParseEvents_Success(t *testing.T) {
	s := event.NewService()

	subtests := []struct {
		name  string
		input io.Reader
		want  []event.Event
	}{
		{
			name:  "AccountCreated",
			input: strings.NewReader(`[{"Type":"AccountCreated","AccountID":"Jack","Payload":{"Balance":50}}]`),
			want:  []event.Event{{Type: event.EventTypeAccountCreated, AccountID: "Jack", Payload: &event.EventPayloadAccountCreated{Balance: 50}}},
		},
		{
			name:  "AccountChargeReceived",
			input: strings.NewReader(`[{"Type":"AccountChargeReceived","AccountID":"Jack","Payload":{"Amount":25}}]`),
			want:  []event.Event{{Type: event.EventTypeAccountChargeReceived, AccountID: "Jack", Payload: &event.EventPayloadAccountTransactionReceived{Amount: 25}}},
		},
		{
			name:  "AccountPaymentReceived",
			input: strings.NewReader(`[{"Type":"AccountPaymentReceived","AccountID":"Jack","Payload":{"Amount":25}}]`),
			want:  []event.Event{{Type: event.EventTypeAccountPaymentReceived, AccountID: "Jack", Payload: &event.EventPayloadAccountTransactionReceived{Amount: 25}}},
		},
		{
			name:  "AccountRecalled",
			input: strings.NewReader(`[{"Type":"AccountRecalled","AccountID":"Jack","Payload":{}}]`),
			want:  []event.Event{{Type: event.EventTypeAccountRecalled, AccountID: "Jack", Payload: nil}},
		},
		{
			name:  "AccountRecalled_PayloadIgnored",
			input: strings.NewReader(`[{"Type":"AccountRecalled","AccountID":"Jack","Payload":{"Amount":25}}]`),
			want:  []event.Event{{Type: event.EventTypeAccountRecalled, AccountID: "Jack", Payload: nil}},
		},
		{
			name: "MultipleEvents",
			input: strings.NewReader(`[
				{"Type":"AccountCreated","AccountID":"Jack","Payload":{"Balance":50}},
				{"Type":"AccountChargeReceived","AccountID":"Jack","Payload":{"Amount":25}},
				{"Type":"AccountPaymentReceived","AccountID":"Jack","Payload":{"Amount":25}}
			]`),
			want: []event.Event{
				{Type: event.EventTypeAccountCreated, AccountID: "Jack", Payload: &event.EventPayloadAccountCreated{Balance: 50}},
				{Type: event.EventTypeAccountChargeReceived, AccountID: "Jack", Payload: &event.EventPayloadAccountTransactionReceived{Amount: 25}},
				{Type: event.EventTypeAccountPaymentReceived, AccountID: "Jack", Payload: &event.EventPayloadAccountTransactionReceived{Amount: 25}},
			},
		},
		{
			name:  "NoEvents",
			input: strings.NewReader(`[]`),
			want:  []event.Event{},
		},
	}

	for _, tt := range subtests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.ParseEvents(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEvent_ParseEvents_CustomErrors(t *testing.T) {
	subtests := []struct {
		name  string
		input io.Reader
		want  error
	}{
		{
			name:  "ErrInputJSONIsNotArray",
			input: strings.NewReader(`{"Type":"AccountCreated","AccountID":"Jack","Payload":{"Balance":50}}`),
			want:  event.ErrInputJSONIsNotArray,
		},
		{
			name:  "ErrUnsupportedEventType",
			input: strings.NewReader(`[{"Type":"AccountCreate","AccountID":"Jack","Payload":{"Balance":50}}]`),
			want:  &event.ErrUnsupportedEventType{Type: "AccountCreate"},
		},
		{
			name:  "ErrMissingFieldInEventPayloadField AccountCreated",
			input: strings.NewReader(`[{"Type":"AccountCreated","AccountID":"Jack","Payload":{"Amount":50}}]`),
			want:  &event.ErrMissingFieldInEventPayloadField{Field: event.EventPayloadFieldBalance},
		},
		{
			name:  "ErrMissingFieldInEventPayloadField AccountChargeReceived",
			input: strings.NewReader(`[{"Type":"AccountChargeReceived","AccountID":"Jack","Payload":{"Balance":50}}]`),
			want:  &event.ErrMissingFieldInEventPayloadField{Field: event.EventPayloadFieldAmount},
		},
		{
			name:  "ErrMissingFieldInEventPayloadField AccountPaymentReceived",
			input: strings.NewReader(`[{"Type":"AccountPaymentReceived","AccountID":"Jack","Payload":{"Balance":50}}]`),
			want:  &event.ErrMissingFieldInEventPayloadField{Field: event.EventPayloadFieldAmount},
		},
	}

	for _, tt := range subtests {
		t.Run(tt.name, func(t *testing.T) {
			s := event.NewService()
			_, err := s.ParseEvents(tt.input)
			assert.EqualError(t, err, tt.want.Error())
		})
	}
}

func TestEvent_ProcessEvents_Success(t *testing.T) {
	s := event.NewService()

	subtests := []struct {
		name   string
		events []event.Event
		want   map[string]event.Account
	}{
		{
			name: "AccountCreated",
			events: []event.Event{
				{Type: event.EventTypeAccountCreated, AccountID: "Jack", Payload: &event.EventPayloadAccountCreated{Balance: 50}},
			},
			want: func() map[string]event.Account {
				account := event.NewAccount("Jack", 50)
				return map[string]event.Account{"Jack": *account}
			}(),
		},
		{
			name: "AccountChargeReceived",
			events: []event.Event{
				{Type: event.EventTypeAccountCreated, AccountID: "Jack", Payload: &event.EventPayloadAccountCreated{Balance: 50}},
				{Type: event.EventTypeAccountChargeReceived, AccountID: "Jack", Payload: &event.EventPayloadAccountTransactionReceived{Amount: 25}},
			},
			want: func() map[string]event.Account {
				account := event.NewAccount("Jack", 75)
				return map[string]event.Account{account.ID: *account}
			}(),
		},
		{
			name: "AccountPaymentReceived",
			events: []event.Event{
				{Type: event.EventTypeAccountCreated, AccountID: "Jack", Payload: &event.EventPayloadAccountCreated{Balance: 50}},
				{Type: event.EventTypeAccountPaymentReceived, AccountID: "Jack", Payload: &event.EventPayloadAccountTransactionReceived{Amount: 25}},
			},
			want: func() map[string]event.Account {
				account := event.NewAccount("Jack", 25)
				return map[string]event.Account{account.ID: *account}
			}(),
		},
		{
			name: "AccountRecalled",
			events: []event.Event{
				{Type: event.EventTypeAccountCreated, AccountID: "Jack", Payload: &event.EventPayloadAccountCreated{Balance: 50}},
				{Type: event.EventTypeAccountRecalled, AccountID: "Jack", Payload: nil},
			},
			want: func() map[string]event.Account {
				account := event.NewAccount("Jack", 50)
				account.Recall()
				return map[string]event.Account{account.ID: *account}
			}(),
		},
		{
			name: "AccountRecalled PayloadIgnored",
			events: []event.Event{
				{Type: event.EventTypeAccountCreated, AccountID: "Jack", Payload: &event.EventPayloadAccountCreated{Balance: 50}},
				{Type: event.EventTypeAccountRecalled, AccountID: "Jack", Payload: &event.EventPayloadAccountTransactionReceived{Amount: 25}},
			},
			want: func() map[string]event.Account {
				account := event.NewAccount("Jack", 50)
				account.Recall()
				return map[string]event.Account{account.ID: *account}
			}(),
		},
		{
			name:   "NoEvents",
			events: []event.Event{},
			want:   map[string]event.Account{},
		},
		{
			name: "MultipleEvents",
			events: []event.Event{
				{
					Type:      event.EventTypeAccountCreated,
					AccountID: "Jack",
					Payload:   &event.EventPayloadAccountCreated{Balance: 50},
				},
				{
					Type:      event.EventTypeAccountCreated,
					AccountID: "Jen",
					Payload:   &event.EventPayloadAccountCreated{Balance: 100},
				},
				{
					Type:      event.EventTypeAccountChargeReceived,
					AccountID: "Jack",
					Payload:   &event.EventPayloadAccountTransactionReceived{Amount: 25},
				},
				{
					Type:      event.EventTypeAccountCreated,
					AccountID: "Robert",
					Payload:   &event.EventPayloadAccountCreated{Balance: 100},
				},
				{
					Type:      event.EventTypeAccountChargeReceived,
					AccountID: "Jack",
					Payload:   &event.EventPayloadAccountTransactionReceived{Amount: 25},
				},
				{
					Type:      event.EventTypeAccountCreated,
					AccountID: "Olivia",
					Payload:   &event.EventPayloadAccountCreated{Balance: 50},
				},
				{
					Type:      event.EventTypeAccountPaymentReceived,
					AccountID: "Robert",
					Payload:   &event.EventPayloadAccountTransactionReceived{Amount: 50},
				},
				{
					Type:      event.EventTypeAccountPaymentReceived,
					AccountID: "Jen",
					Payload:   &event.EventPayloadAccountTransactionReceived{Amount: 50},
				},
				{
					Type:      event.EventTypeAccountChargeReceived,
					AccountID: "Robert",
					Payload:   &event.EventPayloadAccountTransactionReceived{Amount: 25},
				},
				{
					Type:      event.EventTypeAccountPaymentReceived,
					AccountID: "Jack",
					Payload:   &event.EventPayloadAccountTransactionReceived{Amount: 100},
				},
				{
					Type:      event.EventTypeAccountPaymentReceived,
					AccountID: "Jen",
					Payload:   &event.EventPayloadAccountTransactionReceived{Amount: 60},
				},
				{
					Type:      event.EventTypeAccountRecalled,
					AccountID: "Olivia",
					Payload:   nil,
				},
				{
					Type:      event.EventTypeAccountPaymentReceived,
					AccountID: "Robert",
					Payload:   &event.EventPayloadAccountTransactionReceived{Amount: 50},
				},
			},
			want: func() map[string]event.Account {
				targetValues := []struct {
					id      string
					balance int
				}{
					{"Jack", 0},
					{"Jen", -10},
					{"Robert", 25},
					{"Olivia", 50}, // Recalled
				}

				m := make(map[string]event.Account)
				for _, v := range targetValues {
					account := event.NewAccount(v.id, v.balance)
					m[v.id] = *account
				}

				forRecall := m["Olivia"]
				forRecall.Recall()
				m["Olivia"] = forRecall

				return m
			}(),
		},
	}

	for _, tt := range subtests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.ProcessEvents(tt.events)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEvent_ProcessEvents_CustomErrors(t *testing.T) {
	s := event.NewService()

	subtests := []struct {
		name   string
		events []event.Event
		want   error
	}{
		{
			name: "ErrAccountAlreadyExists",
			events: []event.Event{
				{Type: event.EventTypeAccountCreated, AccountID: "Jack", Payload: &event.EventPayloadAccountCreated{Balance: 50}},
				{Type: event.EventTypeAccountCreated, AccountID: "Jack", Payload: &event.EventPayloadAccountCreated{Balance: 50}},
			},
			want: &event.ErrAccountAlreadyExists{AccountID: "Jack"},
		},
		{
			name: "ErrAccountDoesNotExist",
			events: []event.Event{
				{Type: event.EventTypeAccountChargeReceived, AccountID: "Jack", Payload: &event.EventPayloadAccountTransactionReceived{Amount: 25}},
			},
			want: &event.ErrAccountDoesNotExist{AccountID: "Jack"},
		},
		{
			name: "ErrCannotTransactWithRecalledAccount",
			events: []event.Event{
				{Type: event.EventTypeAccountCreated, AccountID: "Jack", Payload: &event.EventPayloadAccountCreated{Balance: 50}},
				{Type: event.EventTypeAccountRecalled, AccountID: "Jack", Payload: nil},
				{Type: event.EventTypeAccountChargeReceived, AccountID: "Jack", Payload: &event.EventPayloadAccountTransactionReceived{Amount: 25}},
			},
			want: &event.ErrCannotTransactWithRecalledAccount{AccountID: "Jack"},
		},
	}

	for _, tt := range subtests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.ProcessEvents(tt.events)
			assert.EqualError(t, err, tt.want.Error())
		})
	}
}
