package database

import (
	"database/sql"
	"errors"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/pedromsmoreira/jarvis/internal/configuration"
	"github.com/sirupsen/logrus"
)

func DoMigrations(targetVersion uint, settings *configuration.Sql, db *sql.DB) bool {
	migrationDriver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		logrus.WithError(err).
			Error("Error initializing sql migration driver.")
		return false
	}
	migrationDb, err := migrate.NewWithDatabaseInstance(
		"file://"+settings.MigrationPath,
		"mysql",
		migrationDriver,
	)
	if err != nil {
		logrus.WithError(err).
			Error("Error creating sql migration db instance.")
		return false
	}

	currentVersion, dirty, err := migrationDb.Version()
	if dirty {
		logrus.WithError(err).
			Error("Database state is dirty! Can not start bot... Resolve the issue with the database state and restart the bot.")
		return false
	}
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			err = nil
			currentVersion = 0
		} else {
			logrus.WithError(err).
				Error("Error checking database migration version.")
			return false
		}
	}

	if currentVersion > targetVersion {
		logrus.
			WithField("target_version", targetVersion).
			WithField("current_version", currentVersion).
			Error("Database migration version is higher than target version in the code, can not proceed. Maybe you forgot to increase the target version?")
		return false
	} else if currentVersion == targetVersion {
		logrus.
			WithField("target_version", targetVersion).
			WithField("current_version", currentVersion).
			Info("Database migration version is up to date. No migration performed.")
	} else {
		err = migrationDb.Migrate(targetVersion)
		if err != nil {
			logrus.WithError(err).
				WithField("target_version", targetVersion).
				WithField("current_version", currentVersion).
				Error("Error performing database migrations.")
			return false
		} else {
			logrus.
				WithField("target_version", targetVersion).
				WithField("current_version", currentVersion).
				Info("Performed database migration successfully.")
		}
	}
	return true
}
