package repositories

import (
	"context"
	db_services "demo-signserver/pkg/repository/services"
	"fmt"
	"os"
	"testing"
)

func initConfig() {
	os.Setenv("PROJECT_NAME", "signserver")
	os.Setenv("ENVIRONMENT", "dev")
	os.Setenv("DYNAMODB_ENDPOINT", "http://localhost:8001")
	os.Setenv("SIGN_REQUEST_TABLE", "request_test")
	os.Setenv("SIGN_PROFILE_TABLE", "profile_test")
	os.Setenv("SIGN_STORAGE_BUCKET", "storage_test")
	os.Setenv("AWS_ACCESS_KEY_ID", "fake")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "fake")
}

func TestMain(m *testing.M) {
	initConfig()
	db := db_services.NewDB()
	profile := fmt.Sprintf("%s-%s-%s", os.Getenv("PROJECT_NAME"), os.Getenv("ENVIRONMENT"), os.Getenv("SIGN_PROFILE_TABLE"))
	db.DeleteTable(context.Background(), profile)
	db.CreateTable(context.Background(), profile, db_services.SORT_KEY)

	request := fmt.Sprintf("%s-%s-%s", os.Getenv("PROJECT_NAME"), os.Getenv("ENVIRONMENT"), os.Getenv("SIGN_REQUEST_TABLE"))
	db.DeleteTable(context.Background(), request)
	db.CreateTable(context.Background(), request, db_services.NO_KEY)

	code := 1
	defer func() {
		db.DeleteTable(context.Background(), profile)
		db.DeleteTable(context.Background(), request)
		os.Exit(code)
	}()
	code = m.Run()
}
