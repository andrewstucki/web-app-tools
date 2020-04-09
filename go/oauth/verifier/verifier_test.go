// +build integration

package verifier

import (
	"testing"
	"time"
)

var (
	testToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjUzYzY2YWFiNTBjZmRkOTFhMTQzNTBhNjY0ODJkYjM4MDBjODNjNjMiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhenAiOiI0MDc0MDg3MTgxOTIuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJhdWQiOiI0MDc0MDg3MTgxOTIuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMTM3OTI0NDk4NzQ1ODEyNTAxMDUiLCJlbWFpbCI6ImdwLm9wcy5ib3RAZ21haWwuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImF0X2hhc2giOiJrT1ZnNEdJRW13dUF0d0xUODM2Y0RBIiwibmFtZSI6IkdyYWNlcG9pbnQgQm90IiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS8tejFRaXVaY2k4SzAvQUFBQUFBQUFBQUkvQUFBQUFBQUFBQUEvQUtGMDVuQUNzeEZTUUlWOWgyRE00RllCSWlRWUxEUFhkQS9zOTYtYy9waG90by5qcGciLCJnaXZlbl9uYW1lIjoiR3JhY2Vwb2ludCIsImZhbWlseV9uYW1lIjoiQm90IiwibG9jYWxlIjoiZW4iLCJpYXQiOjE1ODUzMzk2MDQsImV4cCI6MTU4NTM0MzIwNH0.BMk59OietkQxNXZN36dTQ6Hfj2IH_wis4o2Ay15SZt9wzQuuImNS_Uqt44mo198ss5jrn6NnrOyWli1JZDQ-HCwcJXi7py845TQT20qZTUknTQQZJb9gseD9350-_u1Lx_qphoMMsSljTlygSWmvWGFwmJEweOx1v7HmwX2wZkS_Vir7jN3eMr1AqOi3j8A8iY0e3mg3u_q_X8a4-7-eJHc-3VwlY_ITvVfYGMMQoo14Pb6sWjzbZPJnpoXkjc40arORK3Kn_F1LpKT0QS8y4S5-0eZ5DdcOMvOeKVt2gur9bm9snD8emnK4jTWz92nYl4eqEcNHdsM3jI6ll4EUPg"
	issued    = time.Unix(1585339604, 0)
	expires   = time.Unix(1585343204, 0)
)

func assertErrorEquality(t *testing.T, expected, actual error) {
	if actual == nil || expected.Error() != actual.Error() {
		t.Fatalf("expected '%v', got '%v'", expected, actual)
	}
}

func assertNoError(t *testing.T, actual error) {
	if actual != nil {
		t.Fatalf("unexpected error '%v'", actual)
	}
}

const ()

func TestVerifier(t *testing.T) {
	TimeFunc = func() time.Time {
		return issued
	}

	claims := &StandardClaims{}

	verifier := NewVerifier().WithAudiences("407408718192.apps.googleusercontent.com")
	assertNoError(t, verifier.VerifyIDToken(testToken, claims))

	verifier = NewVerifier().WithAudiences("google.com")
	assertErrorEquality(t, ErrInvalidAudience, verifier.VerifyIDToken(testToken, claims))

	verifier = NewVerifier().WithDomains("gpmail.org")
	assertErrorEquality(t, ErrInvalidDomain, verifier.VerifyIDToken(testToken, claims))

	verifier = NewVerifier()

	TimeFunc = func() time.Time {
		return issued.Add(-1 * time.Minute)
	}
	claims.SetAllowedSkew(0)
	assertErrorEquality(t, ErrIssuedAt, verifier.VerifyIDToken(testToken, claims))
	claims.SetAllowedSkew(1 * time.Minute)
	assertNoError(t, verifier.VerifyIDToken(testToken, claims))

	TimeFunc = func() time.Time {
		return expires
	}
	claims.SetAllowedSkew(0)
	assertErrorEquality(t, ErrExpired, verifier.VerifyIDToken(testToken, claims))
	claims.SetAllowedSkew(1 * time.Minute)
	assertNoError(t, verifier.VerifyIDToken(testToken, claims))
}
