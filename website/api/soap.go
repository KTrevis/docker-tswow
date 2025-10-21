package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/tiaguinho/gosoap"
)

type roundTripperWithAuth struct {
	underlying http.RoundTripper
	authHeader string
}

func (r roundTripperWithAuth) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := req.Clone(req.Context())
	req2.Header.Set("Authorization", r.authHeader)
	return r.underlying.RoundTrip(req2)
}

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

	// Create SOAP client with timeout
	httpClient := &http.Client{Timeout: 10 * time.Second}
	soapClient, err := gosoap.SoapClient(url, httpClient)
	if err != nil {
		return err
	}
	// Add Basic Auth header via client transport RoundTripper
	auth := base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
	// gosoap.Client exposes HttpClient; attach header via Transport wrapper
	soapClient.HTTPClient.Transport = roundTripperWithAuth{underlying: http.DefaultTransport, authHeader: "Basic " + auth}

	// TrinityCore method and params
	params := gosoap.Params{
		"command": fmt.Sprintf("account create %s %s %s", usernameUpper, password, email),
	}

	// Note: TrinityCore uses namespace urn:TC with method executeCommand
	// gosoap will build the envelope; we rely on server accepting basic form
	_, err = soapClient.Call("executeCommand", params)
	return err
}
