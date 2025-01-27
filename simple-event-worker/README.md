# simple-event-worker

## Exercise

Given an array of events as JSON objects, get the final state of the accounts associated to the events.

## Details

1. The events are processed in order (0 to n).
2. Events are processed in a single transaction. If one of the events fails processing, the program should fail so the processing could be restarted later.
3. The final state of the accounts should be produced to stdout.
4. An event is always associated to an `AccountID`.
5. The payload of each event depends on the event type.
6. There are four event types:
    1. `AccountCreated`: carries an initial balance, assume to be non-negative always.
    2. `AccountChargeReceived`
    3. `AccountPaymentReceived`
    4. `AccountRecalled`
7. The balance of an account tracks how much debt is left unpaid. All monetary representations are integers in this exercise.
8. Both charges and payments provide positive `Amount`. Handle the calculations according to the following:
    * A charge is debt added to the balance: `balance += amount`
    * A payment is a reduction in debt: `balance -= amount`
9. Assume events are ordered correctly:
    * The `AccountCreated` event of an `AccountID` always comes first before the rest of its associated events.
    * Accounts cannot be recreated after its first `AccountCreated` event. Assume it is an error.
    * Payment/Charge events cannot be processed after an account gets recalled (`AccountRecalled`). Assume it is frozen -- no changes to the account could be made after it is recalled.
10. There are four Account states:
    1. `Outstanding`: balance is positive
    2. `Settled`: balance is zero
    3. `Overpaid`: balance is negative
    4. `Recalled`: account is recalled and frozen, regardless of balance
11. The final state of an account should be printed as: `Jack: {Status: outstanding, Balance: 50}`.
12. If processing fails at any point, exit, and calculate the state from the beginning upon the next run. Automatic restart/recovery is outside the scope.

## Design

We stream the input one JSON object at a time, and process one at a time folding each event into the output collection. The strategy is essentially an input-streamed `foldLeft` process.

```pseudocode
func foldLeft(input stream) map[id]account {
    for each json object in input stream {
        get event from json object
        process event and update corresponding accounts[id]
    }
    return accounts
}
```

1. We represent an `Event` object using a struct, with its payload abstracted using an `EventPayload` interface.
2. We parse the input via an `io.Reader` stream, and handle one `Event` object at a time to prevent loading the entire input into memory.
3. We process one `Event` object at a time.
4. If something fails at any step, we exit the program.
