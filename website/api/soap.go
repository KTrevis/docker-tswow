package api

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// callSOAPAccountCreate calls TrinityCore SOAP to create an account.
// Requires SOAP host/port and basic auth user/pass.
// Env vars with defaults:
//
//	SOAP_HOST (default: tswow)
//	SOAP_PORT (default: 7878)
//	SOAP_USER (no default)
//	SOAP_PASS (no default)
func callSOAPAccountCreate(usernameUpper, password, email string) error {
	host := getEnv("SOAP_HOST", "tswow")
	port := getEnv("SOAP_PORT", "7878")
	user := getEnv("SOAP_USER", "admin")
	pass := getEnv("SOAP_PASS", "Galamerde_!")

	if user == "" || pass == "" {
		return fmt.Errorf("missing SOAP_USER or SOAP_PASS env var")
	}

	// TrinityCore SOAP endpoint
	url := fmt.Sprintf("http://%s:%s/", host, port)

	command := fmt.Sprintf("account create %s %s %s", usernameUpper, password, email)
	envelope := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
  <SOAP-ENV:Body>
    <ns1:executeCommand xmlns:ns1="urn:TC">
      <command>%s</command>
    </ns1:executeCommand>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`, command)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(envelope))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	// Basic auth header
	auth := base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
	req.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("soap error: status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}
