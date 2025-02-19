package main

import (
	"fmt"
	"time"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
)

func main() {
	fmt.Println("Hello, Producer!")

}

type KYCVerificationRequest struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	FirstName   string
	LastName    string
	MiddleName  string
	Suffix      string
	BirthDate   civil.Date
	PEPStatus   PEPStatus // We enforce that PEPStatus is never null.
	SubmittedAt time.Time
}

type PEPStatus struct {
	// IsPEP is true when the user is a Politically Exposed Person.
	// IsRCA is true when the user is a Relative or Close Associate of a PEP.
	// As far as we know, a user cannot be both a PEP and an RCA at the same time.
	// That would just mean that the user is a PEP.
	//
	// IsPEP XOR IsRCA
	//
	// IsPEP | IsRCA |
	// ------|-------|
	// true  | false | PEP
	// true  | true  | Ambiguous, so we fail the input.
	// false | false | Non-PEP
	// false | true  | RCA only
	IsPEP bool
	IsRCA bool
	// Affiliation is the name of the relevant individual or organization that the user is associated with.
	// If Affiliation is empty, then the PEP is the user themselves, and we let Compliance fill in the details if necessary.
	Affiliation string
}

func generateRandomKYCVerification() *KYCVerificationRequest {
	return nil
}
