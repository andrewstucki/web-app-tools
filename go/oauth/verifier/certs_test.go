// +build integration

package verifier

import (
	"testing"
	"time"
)

const currentCert = "53c66aab50cfdd91a14350a66482db3800c83c63"

func TestGetFederatedSignonCerts(t *testing.T) {
	certs, err := getFederatedSignonCerts()
	if err != nil {
		t.Fatal(err)
	}

	cacheAge := certs.Expiry.Sub(time.Now()).Seconds()
	if cacheAge <= 7200 {
		t.Fatal("max-age not found")
	}

	if certs.Key(currentCert) == nil {
		t.Fatal("key should exists")
	}
}

func TestGetFederatedSignonCertsCache(t *testing.T) {
	certs = &Certs{
		Expiry: time.Now(),
	}
	certs, err := getFederatedSignonCerts() // trigger update
	if err != nil {
		t.Fatal(err)
	}
	if certs.Key(currentCert) == nil {
		t.Fatal("key should exists")
	}
}
