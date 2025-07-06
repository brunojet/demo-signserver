package http_handlers

import (
	"context"
	"demo-signserver/internal/config"
	"demo-signserver/internal/request/services"
	db_services "demo-signserver/pkg/repository/services"
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

var (
	r *gin.Engine
)

func init() {
	os.Setenv("PROJECT_NAME", "signserver")
	os.Setenv("ENVIRONMENT", "dev")
	os.Setenv("DYNAMODB_ENDPOINT", "http://localhost:8001")
	os.Setenv("SIGN_REQUEST_TABLE", "request_handler_test")
	os.Setenv("SIGN_PROFILE_TABLE", "profile_handler_test")
	os.Setenv("SIGN_STORAGE_BUCKET", "storage_handler_test")
	os.Setenv("AWS_ACCESS_KEY_ID", "fake")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "fake")
}

func TestMain(m *testing.M) {
	db := db_services.NewDB()
	profile := config.MakeResourceName("SIGN_PROFILE_TABLE")
	db.DeleteTable(context.Background(), profile)
	db.CreateTable(context.Background(), profile, db_services.SORT_KEY)

	request := config.MakeResourceName("SIGN_REQUEST_TABLE")
	db.DeleteTable(context.Background(), request)
	db.CreateTable(context.Background(), request, db_services.NO_KEY)

	gin.SetMode(gin.TestMode)
	r = gin.Default()
	handler := NewRequestHandler(services.NewRequestService())
	handler.RegisterRoutes(r)
	handlerProfile := NewProfileHandler(services.NewProfileService())
	handlerProfile.RegisterRoutes(r)
	log.Println("HttpHandlerTest initialized")

	code := 1
	defer func() {
		db.DeleteTable(context.Background(), profile)
		db.DeleteTable(context.Background(), request)
		os.Exit(code)
	}()
	code = m.Run()
}
