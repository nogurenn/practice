package simpleeventworker

const (
	AccountStatusOutstanding = "Outstanding"
	AccountStatusRecalled    = "Recalled"
	AccountStatusSettled     = "Settled"
	AccountStatusOverpaid    = "Overpaid"
)

type Account struct {
	ID      string
	status  string
	balance int
}

func NewAccount(id string, balance int) *Account {
	account := &Account{
		ID:     id,
		status: AccountStatusSettled,
	}

	_ = account.RecordTransaction(balance)

	return account
}

// RecordTransaction records a transaction for the account.
// At the time of writing, both charges and payments are positive integers from input.
// Ensure `amount` is positive for charges and negative for payments.
// The account's status is updated based on the new balance.
// Returns an error if the account is already recalled.
func (a *Account) RecordTransaction(amount int) error {
	if a.IsRecalled() {
		return &ErrCannotTransactWithRecalledAccount{AccountID: a.ID}
	}

	a.balance += amount

	if a.balance == 0 {
		a.status = AccountStatusSettled
	} else if a.balance < 0 {
		a.status = AccountStatusOverpaid
	} else {
		a.status = AccountStatusOutstanding
	}

	return nil
}

func (a *Account) Status() string {
	return a.status
}

func (a *Account) Balance() int {
	return a.balance
}

func (a *Account) IsRecalled() bool {
	return a.status == AccountStatusRecalled
}

func (a *Account) Recall() {
	a.status = AccountStatusRecalled
}
