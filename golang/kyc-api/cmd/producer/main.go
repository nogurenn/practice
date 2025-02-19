package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/civil"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

func main() {
	slog.Info("Hello, Producer!")

	// to create topics when auto.create.topics.enable='true'
	time.Sleep(5 * time.Second)
	conn, err := kafka.DialLeader(context.Background(), "tcp", "kafka:9092", "kyc-verification-requests", 0)
	if err != nil {
		panic(err.Error())
	}
	conn.Close()

	kafkaWriter := &kafka.Writer{
		Addr:                   kafka.TCP("kafka:9092"),
		Topic:                  "kyc-verification-requests",
		RequiredAcks:           kafka.RequireOne,
		MaxAttempts:            3,
		AllowAutoTopicCreation: true,
	}

	// Create a list of static users to seed the KYC verification requests from consistent user data.
	users := getUsers()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	newVerificationRequestTicker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-shutdown:
			slog.Warn("Shutting down...")
			kafkaWriter.Close()
			os.Exit(0)
		case <-newVerificationRequestTicker.C:
			pickedUser := users[rand.IntN(len(users))]
			newRequest := generateRandomKYCVerificationRequest(pickedUser)
			// slog.Info(fmt.Sprintf("Generated new KYC verification request: %#v\n", newRequest))

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(newRequest)
			if err != nil {
				slog.Error("Failed to encode KYC verification request", slog.Any("error", err))
				continue
			}
			err = kafkaWriter.WriteMessages(
				context.Background(),
				kafka.Message{
					Key:   []byte(newRequest.ID.String()),
					Value: buf.Bytes(),
				},
			)
			if err != nil {
				slog.Error("Failed to write KYC verification request to Kafka", slog.Any("error", err))
				continue
			}
			slog.Info(fmt.Sprintf("Sent KYC verification request: %#v\n", newRequest))
		}
	}
}

func getUsers() []User {
	return []User{
		{
			ID:        uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479"),
			FirstName: "John",
			LastName:  "Doe",
			Suffix:    "Jr.",
			BirthDate: civil.Date{Year: 1990, Month: 1, Day: 5},
			Email:     "john.doe@gmail.com",
			Phone:     "09123456789",
		},

		{
			ID:        uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d480"),
			FirstName: "Jane",
			LastName:  "Reyes",
			BirthDate: civil.Date{Year: 1995, Month: 2, Day: 10},
			Email:     "janeeeee@gmail.com",
			Phone:     "09123456788",
		},
		{
			ID:         uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d481"),
			FirstName:  "Maris",
			LastName:   "Lopez",
			MiddleName: "Mendoza",
			BirthDate:  civil.Date{Year: 1992, Month: 3, Day: 15},
			Email:      "mendozamaris@yahoo.com",
			Phone:      "09123456787",
		},
		{
			ID:         uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d482"),
			FirstName:  "Juan",
			LastName:   "Dela Cruz",
			MiddleName: "Santos",
			Suffix:     "III",
			BirthDate:  civil.Date{Year: 1993, Month: 4, Day: 20},
			Email:      "juantos@gmail.com",
			Phone:      "09123456786",
		},
		{
			ID:         uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d483"),
			FirstName:  "Maria",
			LastName:   "Jimenez",
			MiddleName: "De Guzman",
			BirthDate:  civil.Date{Year: 1994, Month: 5, Day: 25},
			Email:      "mjdeguzman@yahoo.com",
			Phone:      "09123456785",
		},
	}
}

type User struct {
	ID         uuid.UUID
	FirstName  string
	LastName   string
	MiddleName string
	Suffix     string
	BirthDate  civil.Date
	Email      string
	Phone      string
}

type KYCVerificationRequest struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	MiddleName  string     `json:"middle_name"`
	Suffix      string     `json:"suffix"`
	BirthDate   civil.Date `json:"birth_date"`
	PEPStatus   PEPStatus  `json:"pep_status"` // We enforce that PEPStatus is never null.
	SubmittedAt time.Time  `json:"submitted_at"`
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
	// true  | true  | Ambiguous from a compliance pov, so we reject the incoming request early.
	// false | false | Non-PEP
	// false | true  | RCA only
	IsPep bool `json:"is_pep"`
	IsRca bool `json:"is_rca"`
	// Affiliation is the name of the relevant individual or organization that the user is associated with.
	PepAffiliation string `json:"pep_affiliation"`
	RcaAffiliation string `json:"rca_affiliation"`
}

func generateRandomKYCVerificationRequest(user User) *KYCVerificationRequest {
	pepStatus := PEPStatus{
		IsPep: gofakeit.Bool(),
		IsRca: gofakeit.Bool(),
	}
	var affiliations []string
	if pepStatus.IsPep {
		affiliations = []string{
			"",
			gofakeit.Name(),
			gofakeit.Company(),
		}
		pepStatus.PepAffiliation = affiliations[rand.IntN(len(affiliations))]
	}
	if pepStatus.IsRca {
		// RCAs must have an affiliation somehow.
		affiliations = []string{
			gofakeit.Name(),
			gofakeit.Company(),
		}
		pepStatus.RcaAffiliation = affiliations[rand.IntN(len(affiliations))]
	}

	return &KYCVerificationRequest{
		ID:          uuid.New(),
		UserID:      user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		MiddleName:  user.MiddleName,
		Suffix:      user.Suffix,
		BirthDate:   user.BirthDate,
		PEPStatus:   pepStatus,
		SubmittedAt: time.Now(),
	}
}
