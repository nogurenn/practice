package simpleeventworker_test

import (
	"testing"

	event "github.com/nogurenn/assorted-programs/simple-event-worker"
	"github.com/stretchr/testify/assert"
)

func TestAccount_NewAccount(t *testing.T) {
	subtests := []struct {
		name    string
		id      string
		balance int
		want    *event.Account
	}{
		{
			name:    "WithOutstanding",
			id:      "Jack",
			balance: 100,
			want: func() *event.Account {
				return event.NewAccount("Jack", 100)
			}(),
		},
		{
			name:    "Settled",
			id:      "Jack",
			balance: 0,
			want: func() *event.Account {
				return event.NewAccount("Jack", 0)
			}(),
		},
		{
			name:    "Overpaid",
			id:      "Jack",
			balance: -100,
			want: func() *event.Account {
				return event.NewAccount("Jack", -100)
			}(),
		},
	}

	for _, tt := range subtests {
		t.Run(tt.name, func(t *testing.T) {
			account := event.NewAccount(tt.id, tt.balance)
			assert.Equal(t, tt.want, account)
			assert.Equal(t, tt.want.Balance(), account.Balance())
			assert.Equal(t, tt.want.Status(), account.Status())
		})
	}
}

func TestAccount_RecordTransaction_Success(t *testing.T) {
	subtests := []struct {
		name         string
		startBalance int
		amount       int
		want         *event.Account
	}{
		{
			name:         "OutstandingToSettled",
			startBalance: 100,
			amount:       -100,
			want: func() *event.Account {
				return event.NewAccount("Jack", 0)
			}(),
		},
		{
			name:         "SettledToOutstanding",
			startBalance: 0,
			amount:       100,
			want: func() *event.Account {
				return event.NewAccount("Jack", 100)
			}(),
		},
		{
			name:         "OutstandingToOverpaid",
			startBalance: 100,
			amount:       -101,
			want: func() *event.Account {
				return event.NewAccount("Jack", -1)
			}(),
		},
		{
			name:         "OverpaidToOutstanding",
			startBalance: -100,
			amount:       101,
			want: func() *event.Account {
				return event.NewAccount("Jack", 1)
			}(),
		},
		{
			name:         "OverpaidToSettled",
			startBalance: -100,
			amount:       100,
			want: func() *event.Account {
				return event.NewAccount("Jack", 0)
			}(),
		},
		{
			name:         "OutstandingToOutstanding",
			startBalance: 100,
			amount:       50,
			want: func() *event.Account {
				return event.NewAccount("Jack", 150)
			}(),
		},
		{
			name:         "SettledToSettled",
			startBalance: 0,
			amount:       0,
			want: func() *event.Account {
				return event.NewAccount("Jack", 0)
			}(),
		},
	}

	for _, tt := range subtests {
		t.Run(tt.name, func(t *testing.T) {
			account := event.NewAccount("Jack", tt.startBalance)
			err := account.RecordTransaction(tt.amount)

			assert.Nil(t, err)
			assert.Equal(t, tt.want, account)
			assert.Equal(t, tt.want.Balance(), account.Balance())
			assert.Equal(t, tt.want.Status(), account.Status())
		})
	}
}

func TestAccount_RecordTransaction_CustomErrors(t *testing.T) {
	subtests := []struct {
		name         string
		startBalance int
		amount       int
		want         error
	}{
		{
			name:         "AlreadyRecalled",
			startBalance: 100,
			amount:       -100,
			want:         &event.ErrCannotTransactWithRecalledAccount{AccountID: "Jack"},
		},
	}

	for _, tt := range subtests {
		t.Run(tt.name, func(t *testing.T) {
			account := event.NewAccount("Jack", tt.startBalance)
			account.Recall()
			err := account.RecordTransaction(tt.amount)
			assert.Equal(t, tt.want, err)
		})
	}
}

func TestAccount_Recall(t *testing.T) {
	subtests := []struct {
		name string
		want *event.Account
	}{
		{
			name: "Recalled",
			want: func() *event.Account {
				account := event.NewAccount("Jack", 100)
				account.Recall()
				return account
			}(),
		},
	}

	for _, tt := range subtests {
		t.Run(tt.name, func(t *testing.T) {
			account := event.NewAccount("Jack", 100)
			account.Recall()

			assert.Equal(t, tt.want, account)
			assert.Equal(t, tt.want.Balance(), account.Balance())
			assert.Equal(t, tt.want.Status(), account.Status())
		})
	}
}
