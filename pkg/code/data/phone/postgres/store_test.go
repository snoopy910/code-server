package postgres

import (
	"database/sql"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/sirupsen/logrus"

	"github.com/code-payments/code-server/pkg/code/data/phone"
	"github.com/code-payments/code-server/pkg/code/data/phone/tests"

	postgrestest "github.com/code-payments/code-server/pkg/database/postgres/test"

	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	// Used for testing ONLY, the table and migrations are external to this repository
	tableCreate = `
		CREATE TABLE codewallet__core_phoneverification(
			id SERIAL NOT NULL PRIMARY KEY,

			phone_number TEXT NOT NULL,
			owner_account TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL,
			last_verified_at TIMESTAMP WITH TIME ZONE NOT NULL,

			CONSTRAINT codewallet__core_phoneverification__uniq__owner_account__and__phone_number UNIQUE (owner_account, phone_number)
		);

		CREATE TABLE codewallet__core_phonelinkingtoken(
			id SERIAL NOT NULL PRIMARY KEY,

			phone_number TEXT NOT NULL,
			code TEXT NOT NULL,
			current_check_count INTEGER NOT NULL,
			max_check_count INTEGER NOT NULL,
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL,

			CONSTRAINT codewallet__core_phonelinkingtoken__uniq__phone_number UNIQUE (phone_number)
		);

		CREATE TABLE codewallet__core_phonesetting(
			id SERIAL NOT NULL PRIMARY KEY,

			phone_number TEXT NOT NULL,
			owner_account TEXT NOT NULL,
			is_unlinked BOOL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL,
			last_updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

			CONSTRAINT codewallet__core_phonesetting__uniq__owner_account__and__phone_number UNIQUE (owner_account, phone_number)
		);

		CREATE TABLE codewallet__core_phoneevent(
			id SERIAL NOT NULL PRIMARY KEY,

			event_type INTEGER NOT NULL,

			verification_id TEXT NOT NULL,

			phone_number TEXT NOT NULL,
			phone_type INTEGER NULL,
			mobile_country_code INTEGER NULL,
			mobile_network_code INTEGER NULL,

			created_at TIMESTAMP WITH TIME ZONE NOT NULL
		);
	`

	// Used for testing ONLY, the table and migrations are external to this repository
	tableDestroy = `
		DROP TABLE codewallet__core_phoneverification;
		DROP TABLE codewallet__core_phonelinkingtoken;
		DROP TABLE codewallet__core_phonesetting;
		DROP TABLE codewallet__core_phoneevent;
	`
)

var (
	testStore phone.Store
	teardown  func()
)

func TestMain(m *testing.M) {
	log := logrus.StandardLogger()

	testPool, err := dockertest.NewPool("")
	if err != nil {
		log.WithError(err).Error("Error creating docker pool")
		os.Exit(1)
	}

	var cleanUpFunc func()
	db, cleanUpFunc, err := postgrestest.StartPostgresDB(testPool)
	if err != nil {
		log.WithError(err).Error("Error starting postgres image")
		os.Exit(1)
	}
	defer db.Close()

	if err := createTestTables(db); err != nil {
		logrus.StandardLogger().WithError(err).Error("Error creating test tables")
		cleanUpFunc()
		os.Exit(1)
	}

	testStore = New(db)
	teardown = func() {
		if pc := recover(); pc != nil {
			cleanUpFunc()
			panic(pc)
		}

		if err := resetTestTables(db); err != nil {
			logrus.StandardLogger().WithError(err).Error("Error resetting test tables")
			cleanUpFunc()
			os.Exit(1)
		}
	}

	code := m.Run()
	cleanUpFunc()
	os.Exit(code)
}

func TestPhonePostgresStore(t *testing.T) {
	tests.RunTests(t, testStore, teardown)
}

func createTestTables(db *sql.DB) error {
	_, err := db.Exec(tableCreate)
	if err != nil {
		logrus.StandardLogger().WithError(err).Error("could not create test tables")
		return err
	}
	return nil
}

func resetTestTables(db *sql.DB) error {
	_, err := db.Exec(tableDestroy)
	if err != nil {
		logrus.StandardLogger().WithError(err).Error("could not drop test tables")
		return err
	}

	return createTestTables(db)
}
