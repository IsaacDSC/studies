package blocklist

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func setupTestHandler(t *testing.T) (*Handler, string) {
	tmpDir := t.TempDir()
	repo := NewFileRepository(tmpDir)
	service := NewService(repo)
	handler := NewHandler(service)
	return handler, tmpDir
}

func TestHandler_Block(t *testing.T) {
	handler, _ := setupTestHandler(t)
	userID := uuid.New()

	// Create request
	body, _ := json.Marshal(BlockRequest{UserID: userID.String()})
	req := httptest.NewRequest(http.MethodPost, "/backoffice/block-list", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Call handler
	handler.Block(rr, req)

	// Check response
	if rr.Code != http.StatusCreated {
		t.Errorf("Block() status = %d, want %d", rr.Code, http.StatusCreated)
	}

	var resp BlockResponse
	json.NewDecoder(rr.Body).Decode(&resp)
	if resp.UserID != userID.String() {
		t.Errorf("Block() resp.UserID = %v, want %v", resp.UserID, userID.String())
	}
	if resp.Message != "user blocked successfully" {
		t.Errorf("Block() resp.Message = %v", resp.Message)
	}
}

func TestHandler_Block_InvalidUUID(t *testing.T) {
	handler, _ := setupTestHandler(t)

	// Create request with invalid UUID
	body, _ := json.Marshal(BlockRequest{UserID: "invalid-uuid"})
	req := httptest.NewRequest(http.MethodPost, "/backoffice/block-list", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Call handler
	handler.Block(rr, req)

	// Check response
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Block() status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestHandler_Block_InvalidJSON(t *testing.T) {
	handler, _ := setupTestHandler(t)

	// Create request with invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/backoffice/block-list", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Call handler
	handler.Block(rr, req)

	// Check response
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Block() status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestHandler_Block_AlreadyBlocked(t *testing.T) {
	handler, _ := setupTestHandler(t)
	userID := uuid.New()

	// Block user first
	body, _ := json.Marshal(BlockRequest{UserID: userID.String()})
	req := httptest.NewRequest(http.MethodPost, "/backoffice/block-list", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.Block(rr, req)

	// Try to block again
	body, _ = json.Marshal(BlockRequest{UserID: userID.String()})
	req = httptest.NewRequest(http.MethodPost, "/backoffice/block-list", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	handler.Block(rr, req)

	// Check response
	if rr.Code != http.StatusConflict {
		t.Errorf("Block() status = %d, want %d", rr.Code, http.StatusConflict)
	}
}

func TestHandler_Unblock(t *testing.T) {
	handler, _ := setupTestHandler(t)
	userID := uuid.New()

	// Block user first
	body, _ := json.Marshal(BlockRequest{UserID: userID.String()})
	req := httptest.NewRequest(http.MethodPost, "/backoffice/block-list", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.Block(rr, req)

	// Unblock user
	req = httptest.NewRequest(http.MethodDelete, "/backoffice/block-list/"+userID.String(), nil)
	req.SetPathValue("user_id", userID.String())
	rr = httptest.NewRecorder()
	handler.Unblock(rr, req)

	// Check response
	if rr.Code != http.StatusOK {
		t.Errorf("Unblock() status = %d, want %d", rr.Code, http.StatusOK)
	}

	var resp UnblockResponse
	json.NewDecoder(rr.Body).Decode(&resp)
	if resp.UserID != userID.String() {
		t.Errorf("Unblock() resp.UserID = %v, want %v", resp.UserID, userID.String())
	}
}

func TestHandler_Unblock_NotBlocked(t *testing.T) {
	handler, _ := setupTestHandler(t)
	userID := uuid.New()

	// Unblock non-blocked user
	req := httptest.NewRequest(http.MethodDelete, "/backoffice/block-list/"+userID.String(), nil)
	req.SetPathValue("user_id", userID.String())
	rr := httptest.NewRecorder()
	handler.Unblock(rr, req)

	// Check response
	if rr.Code != http.StatusNotFound {
		t.Errorf("Unblock() status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}

func TestHandler_Unblock_InvalidUUID(t *testing.T) {
	handler, _ := setupTestHandler(t)

	// Unblock with invalid UUID
	req := httptest.NewRequest(http.MethodDelete, "/backoffice/block-list/invalid-uuid", nil)
	req.SetPathValue("user_id", "invalid-uuid")
	rr := httptest.NewRecorder()
	handler.Unblock(rr, req)

	// Check response
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Unblock() status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}


