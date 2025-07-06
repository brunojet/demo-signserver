package dtos

import (
	"demo-signserver/internal/repository/domain"
	"reflect"
	"testing"
)

func TestCreateSignRequestDTO_GetDomainCreateSignRequest(t *testing.T) {
	webhook := "https://callback.com/webhook"
	dto := CreateSignRequestDTO{
		ProfileId:  "123",
		WebhookURL: &webhook,
	}
	domainObj := dto.GetDomainCreateSignRequest()

	if domainObj.SignerProfileId == nil || *domainObj.SignerProfileId != dto.ProfileId {
		t.Errorf("Expected SignerProfileId %s, got %v", dto.ProfileId, domainObj.SignerProfileId)
	}
	if domainObj.SignerStatus == nil || *domainObj.SignerStatus != domain.SignerStatusCreated {
		t.Errorf("Expected SignerStatus %s, got %v", domain.SignerStatusCreated, domainObj.SignerStatus)
	}
	if domainObj.WebhookURL == nil || *domainObj.WebhookURL != webhook {
		t.Errorf("Expected WebhookURL %s, got %v", webhook, domainObj.WebhookURL)
	}
	if domainObj.History == nil || len(*domainObj.History) != 1 {
		t.Errorf("Expected History with 1 entry, got %v", domainObj.History)
	} else {
		entry := (*domainObj.History)[0]
		if entry.SignerStatus == nil || *entry.SignerStatus != domain.SignerStatusCreated {
			t.Errorf("Expected first history step to be %s, got %v", domain.SignerStatusCreated, entry.SignerStatus)
		}
		if entry.Error != nil {
			t.Errorf("Expected first history error to be nil, got %v", entry.Error)
		}
		if entry.CreatedAt == "" {
			t.Errorf("Expected first history timestamp to be set, got empty")
		}
	}
}

func TestCreateSignResponseDTO_JSONTags(t *testing.T) {
	resp := CreateSignResponseDTO{ID: "abc", UploadURL: "https://upload.com/file"}
	respType := reflect.TypeOf(resp)
	idField, _ := respType.FieldByName("ID")
	if idField.Tag.Get("json") != "id" {
		t.Errorf("Expected json tag 'id' for ID field, got %s", idField.Tag.Get("json"))
	}
	uploadField, _ := respType.FieldByName("UploadURL")
	if uploadField.Tag.Get("json") != "pre_signed_url" {
		t.Errorf("Expected json tag 'pre_signed_url' for UploadURL field, got %s", uploadField.Tag.Get("json"))
	}
}

func TestGetResponseDTO_JSONTags(t *testing.T) {
	resp := GetResponseDTO{}
	respType := reflect.TypeOf(resp)
	idField, _ := respType.FieldByName("ID")
	if idField.Tag.Get("json") != "id" {
		t.Errorf("Expected json tag 'id' for ID field, got %s", idField.Tag.Get("json"))
	}
	downloadField, _ := respType.FieldByName("DownloadURL")
	if downloadField.Tag.Get("json") != "pre_signed_url" {
		t.Errorf("Expected json tag 'pre_signed_url' for DownloadURL field, got %s", downloadField.Tag.Get("json"))
	}
}
