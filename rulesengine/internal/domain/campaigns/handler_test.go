package campaigns

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestHandler(t *testing.T) *Handler {
	tmpDir := t.TempDir()
	repo := NewFileRepository(tmpDir)
	service := NewService(repo)
	return NewHandler(service)
}

func TestHandler_Upsert(t *testing.T) {
	handler := setupTestHandler(t)

	req := httptest.NewRequest(http.MethodPatch, "/backoffice/campaigns/test_campaign?operations=credit", bytes.NewReader(ValidGraph))
	req.SetPathValue("campaignName", "test_campaign")
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Upsert(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Upsert() status = %d, want %d, body: %s", rr.Code, http.StatusOK, rr.Body.String())
	}

	var resp UpsertResponse
	json.NewDecoder(rr.Body).Decode(&resp)
	if resp.Campaign != "test_campaign" {
		t.Errorf("Upsert() campaign = %v, want test_campaign", resp.Campaign)
	}
	if resp.Message != "rules updated successfully" {
		t.Errorf("Upsert() message = %v", resp.Message)
	}
}

func TestHandler_Upsert_MissingOperations(t *testing.T) {
	handler := setupTestHandler(t)

	req := httptest.NewRequest(http.MethodPatch, "/backoffice/campaigns/test_campaign", bytes.NewReader(ValidGraph))
	req.SetPathValue("campaignName", "test_campaign")
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Upsert(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Upsert() status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestHandler_Upsert_InvalidOperation(t *testing.T) {
	handler := setupTestHandler(t)

	req := httptest.NewRequest(http.MethodPatch, "/backoffice/campaigns/test_campaign?operations=invalid", bytes.NewReader(ValidGraph))
	req.SetPathValue("campaignName", "test_campaign")
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Upsert(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Upsert() status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestHandler_Upsert_InvalidGraph(t *testing.T) {
	handler := setupTestHandler(t)

	req := httptest.NewRequest(http.MethodPatch, "/backoffice/campaigns/test_campaign?operations=credit", bytes.NewReader(InvalidGraph))
	req.SetPathValue("campaignName", "test_campaign")
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Upsert(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("Upsert() status = %d, want %d", rr.Code, http.StatusUnprocessableEntity)
	}
}

func TestHandler_List(t *testing.T) {
	handler := setupTestHandler(t)

	// Create some campaigns
	for _, name := range []string{"campaign1", "campaign2"} {
		req := httptest.NewRequest(http.MethodPatch, "/backoffice/campaigns/"+name+"?operations=credit", bytes.NewReader(ValidGraph))
		req.SetPathValue("campaignName", name)
		rr := httptest.NewRecorder()
		handler.Upsert(rr, req)
	}

	// List
	req := httptest.NewRequest(http.MethodGet, "/backoffice/campaigns/list", nil)
	rr := httptest.NewRecorder()
	handler.List(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("List() status = %d, want %d", rr.Code, http.StatusOK)
	}

	var resp ListResponse
	json.NewDecoder(rr.Body).Decode(&resp)
	if resp.Total != 2 {
		t.Errorf("List() total = %v, want 2", resp.Total)
	}
}

func TestHandler_Delete(t *testing.T) {
	handler := setupTestHandler(t)

	// Create campaign
	req := httptest.NewRequest(http.MethodPatch, "/backoffice/campaigns/test_campaign?operations=credit", bytes.NewReader(ValidGraph))
	req.SetPathValue("campaignName", "test_campaign")
	rr := httptest.NewRecorder()
	handler.Upsert(rr, req)

	// Delete
	req = httptest.NewRequest(http.MethodDelete, "/backoffice/campaigns/test_campaign/credit", nil)
	req.SetPathValue("campaignName", "test_campaign")
	req.SetPathValue("operation", "credit")
	rr = httptest.NewRecorder()
	handler.Delete(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Delete() status = %d, want %d", rr.Code, http.StatusOK)
	}

	var resp DeleteResponse
	json.NewDecoder(rr.Body).Decode(&resp)
	if resp.Campaign != "test_campaign" {
		t.Errorf("Delete() campaign = %v, want test_campaign", resp.Campaign)
	}
}

func TestHandler_Delete_NotFound(t *testing.T) {
	handler := setupTestHandler(t)

	req := httptest.NewRequest(http.MethodDelete, "/backoffice/campaigns/nonexistent/credit", nil)
	req.SetPathValue("campaignName", "nonexistent")
	req.SetPathValue("operation", "credit")
	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Delete() status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}

func TestHandler_Delete_InvalidOperation(t *testing.T) {
	handler := setupTestHandler(t)

	req := httptest.NewRequest(http.MethodDelete, "/backoffice/campaigns/test_campaign/invalid", nil)
	req.SetPathValue("campaignName", "test_campaign")
	req.SetPathValue("operation", "invalid")
	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Delete() status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}


