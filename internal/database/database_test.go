package database

import (
	"database/sql/driver"
	"net/url"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/guvanch999/abrayly-mebel-backend/internal/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestDSN(t *testing.T) {
	cfg := &config.Schema{
		PostgresUser:     "abrayly",
		PostgresPassword: "p@ss:w/rd?#",
		PostgresHost:     "db",
		PostgresPort:     "5432",
		PostgresDb:       "abrayly",
		PostgresSSLMode:  "require",
	}

	dsn := DSN(cfg)

	parsed, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("DSN is not a valid URL: %v", err)
	}
	if parsed.Scheme != "postgres" || parsed.Host != "db:5432" || parsed.Path != "/abrayly" {
		t.Errorf("unexpected DSN parts: scheme=%q host=%q path=%q", parsed.Scheme, parsed.Host, parsed.Path)
	}
	if parsed.User.Username() != "abrayly" {
		t.Errorf("username = %q", parsed.User.Username())
	}
	if password, _ := parsed.User.Password(); password != cfg.PostgresPassword {
		t.Error("password with special characters did not survive escaping")
	}
	if got := parsed.Query().Get("sslmode"); got != "require" {
		t.Errorf("sslmode = %q, want %q", got, "require")
	}
}

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Silent),
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}
	return db, mock
}

// bcryptOf matches an argument that is a bcrypt hash of the given password.
type bcryptOf string

func (p bcryptOf) Match(v driver.Value) bool {
	hash, ok := v.(string)
	return ok && bcrypt.CompareHashAndPassword([]byte(hash), []byte(p)) == nil
}

func TestSeedAdminSkipsWithoutPassword(t *testing.T) {
	db, mock := newMockDB(t)

	if err := SeedAdmin(db, AdminSeed{Username: "admin"}); err != nil {
		t.Fatalf("SeedAdmin: %v", err)
	}

	// No query may run: without ADMIN_PASSWORD nothing is created.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestSeedAdminCreatesFirstAdmin(t *testing.T) {
	db, mock := newMockDB(t)

	mock.ExpectQuery(`SELECT count\(\*\) FROM "admins"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectExec(`INSERT INTO "admins"`).
		WithArgs(sqlmock.AnyArg(), "Store Owner", "owner", bcryptOf("s3cret-from-env"), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := SeedAdmin(db, AdminSeed{Name: "Store Owner", Username: "owner", Password: "s3cret-from-env"})
	if err != nil {
		t.Fatalf("SeedAdmin: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestSeedAdminDoesNothingWhenAdminsExist(t *testing.T) {
	db, mock := newMockDB(t)

	mock.ExpectQuery(`SELECT count\(\*\) FROM "admins"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	if err := SeedAdmin(db, AdminSeed{Username: "admin", Password: "whatever"}); err != nil {
		t.Fatalf("SeedAdmin: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
