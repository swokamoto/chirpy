package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-key"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() failed: %v", err)
	}

	if token == "" {
		t.Error("MakeJWT() returned empty token")
	}

	// Validate that the token can be parsed back
	parsedUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT() failed for token created by MakeJWT(): %v", err)
	}

	if parsedUserID != userID {
		t.Errorf("ValidateJWT() returned userID = %v, want %v", parsedUserID, userID)
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-key"
	wrongSecret := "wrong-secret-key"

	tests := []struct {
		name       string
		setupToken func() string
		secret     string
		wantUserID uuid.UUID
		wantErr    bool
	}{
		{
			name: "Valid token",
			setupToken: func() string {
				token, _ := MakeJWT(userID, secret, time.Hour)
				return token
			},
			secret:     secret,
			wantUserID: userID,
			wantErr:    false,
		},
		{
			name: "Token signed with wrong secret",
			setupToken: func() string {
				token, _ := MakeJWT(userID, wrongSecret, time.Hour)
				return token
			},
			secret:     secret,
			wantUserID: uuid.Nil,
			wantErr:    true,
		},
		{
			name: "Expired token",
			setupToken: func() string {
				// Create token that expires immediately
				token, _ := MakeJWT(userID, secret, -time.Hour)
				return token
			},
			secret:     secret,
			wantUserID: uuid.Nil,
			wantErr:    true,
		},
		{
			name: "Malformed token",
			setupToken: func() string {
				return "not.a.valid.jwt"
			},
			secret:     secret,
			wantUserID: uuid.Nil,
			wantErr:    true,
		},
		{
			name: "Empty token",
			setupToken: func() string {
				return ""
			},
			secret:     secret,
			wantUserID: uuid.Nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.setupToken()
			gotUserID, err := ValidateJWT(token, tt.secret)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestJWTRoundTrip(t *testing.T) {
	// Test multiple round trips with different users and secrets
	testCases := []struct {
		userID    uuid.UUID
		secret    string
		expiresIn time.Duration
	}{
		{uuid.New(), "secret1", time.Hour},
		{uuid.New(), "different-secret", 30 * time.Minute},
		{uuid.New(), "very-long-secret-key-for-testing", 2 * time.Hour},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			// Create token
			token, err := MakeJWT(tc.userID, tc.secret, tc.expiresIn)
			if err != nil {
				t.Fatalf("MakeJWT() failed: %v", err)
			}

			// Validate token
			parsedUserID, err := ValidateJWT(token, tc.secret)
			if err != nil {
				t.Fatalf("ValidateJWT() failed: %v", err)
			}

			if parsedUserID != tc.userID {
				t.Errorf("Round trip failed: got userID %v, want %v", parsedUserID, tc.userID)
			}
		})
	}
}
