package httputil

import (
	"encoding/json"
	"fmt"
	"testing"
)

type HandlerError struct {
	Client string
	Log    string
}

func TestErrorCheck(t *testing.T, response []byte, handlerErr HandlerError) {
	var responseUnmarshaled map[string]string

	if err := json.Unmarshal(response, &responseUnmarshaled); err != nil {
		t.Fatalf("failed to unmarshal error: %s", err)
	}

	if val, ok := responseUnmarshaled["client_err"]; !ok || val != handlerErr.Client {
		t.Errorf("expected client_err %q, got %q", handlerErr.Client, val)
	}

	if val, ok := responseUnmarshaled["log_err"]; !ok || val != handlerErr.Log {
		t.Errorf("expected log_err %q, got %q", handlerErr.Log, val)
	}
}

func TestInternalServerClientErr() string {
	return "internal server error"
}

func TestNonExistingRessourceErr() string {
	return "there is no ressource assosiated with provided id"
}

func TestJsonUnmarshalingClientError() string {
	return "invalid json body"
}

func TestJsonUnmarshalingLogError[T any](t *testing.T, body []byte) string {
	var target T
	err := json.Unmarshal(body, &target)
	if err == nil {
		t.Fatalf("expected unmarshal error for %q, but got nil", string(body))
	}
	return fmt.Sprintf("failed to unmarshal body: %s", err)
}
