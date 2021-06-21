package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"log"
	"math/big"
	"time"

	"github.com/XSAM/otelsql"
	_ "github.com/mattn/go-sqlite3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/semconv"

	"github.com/sakti/o11ygo/api"
)

const keyLength = 20

func NewAuthService(dbPath string) (*AuthService, error) {
	driverName, err := otelsql.Register("sqlite3", semconv.DBSystemSqlite.Value.AsString())
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(driverName, dbPath)
	if err != nil {
		return nil, err
	}
	createTable, err := db.Prepare("CREATE TABLE IF NOT EXISTS tokens (id INTEGER PRIMARY KEY, token TEXT, created_at INTEGER, revoked INTEGER)")
	if err != nil {
		return nil, err
	}
	_, err = createTable.Exec()
	if err != nil {
		return nil, err
	}
	insertStmt, err := db.Prepare("INSERT INTO tokens (token, created_at, revoked) VALUES (?, ?, ?)")
	if err != nil {
		return nil, err
	}
	return &AuthService{db: db, insertStmt: insertStmt}, nil
}

type AuthService struct {
	api.UnimplementedAuthsvcServer
	db         *sql.DB
	insertStmt *sql.Stmt
}

func (a *AuthService) Enforce(ctx context.Context, in *api.EnforceRequest) (*api.EnforceReply, error) {
	tracer := otel.GetTracerProvider()
	ctx, span := tracer.Tracer("").Start(ctx, "select.token")
	defer span.End()

	var revoked bool
	err := a.db.QueryRowContext(ctx, "SELECT revoked FROM tokens where token = ?", in.Token).Scan(&revoked)
	if err != nil {
		if err == sql.ErrNoRows {
			return &api.EnforceReply{Allowed: false}, nil
		}
		log.Fatal(err)
	}
	return &api.EnforceReply{Allowed: !revoked}, nil
}

func (a *AuthService) CreateToken(ctx context.Context, in *api.CreateTokenRequest) (*api.CreateTokenReply, error) {
	tracer := otel.GetTracerProvider()
	ctx, span := tracer.Tracer("").Start(ctx, "insert.token")
	defer span.End()

	token, err := GenerateRandomString(ctx, keyLength)
	if err != nil {
		return nil, err
	}
	_, err = a.insertStmt.ExecContext(ctx, token, time.Now().Unix(), 0)
	if err != nil {
		return nil, err
	}
	return &api.CreateTokenReply{Token: token}, nil
}

// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
// ref: https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb
func GenerateRandomString(ctx context.Context, n int) (string, error) {
	tracer := otel.GetTracerProvider()
	_, span := tracer.Tracer("").Start(ctx, "GenerateRandomString")
	defer span.End()
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
