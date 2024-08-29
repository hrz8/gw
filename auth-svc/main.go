package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	name = "auth-service"
	port = 3001
)

type JWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, name)
	})
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})
	mux.HandleFunc("POST /api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		privKey, err := getPrivateKey()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		aud := r.Header.Get("X-Audience-Id")
		if aud == "" {
			aud = "tb-app"
		}

		token := jwt.New(jwt.SigningMethodRS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["iss"] = "tb-auth-svc"
		claims["iat"] = time.Now().Unix()
		claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
		claims["aud"] = aud
		claims["sub"] = "123"

		tokenString, err := token.SignedString(privKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		json.NewEncoder(w).Encode(struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}{
			AccessToken:  tokenString,
			RefreshToken: "refresh_token",
		})
	})
	mux.HandleFunc("POST /api/v1/refresh", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(struct {
			AccessToken string `json:"access_token"`
		}{
			AccessToken: "accesstokenbaru",
		})
	})
	mux.HandleFunc("GET /.well-known/jwks.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		jwks, _ := getJwks()
		json.NewEncoder(w).Encode(jwks)
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	fmt.Println("serve http server on port:", port)
	log.Fatalf("cannot start http server: %v\n", srv.ListenAndServe())
}

func getJwks() (*JWKS, error) {
	kid := os.Getenv("JWKS_KID")
	if kid == "" {
		return nil, fmt.Errorf("missing: %v", "kid")
	}

	n := os.Getenv("JWKS_MODULUS")
	if n == "" {
		return nil, fmt.Errorf("missing: %v", "n")
	}

	e := os.Getenv("JWKS_EXPONENT")
	if e == "" {
		return nil, fmt.Errorf("missing: %v", "e")
	}

	jwk := JWK{
		Kty: "RSA",
		Kid: kid,
		Use: "sig",
		N:   n,
		E:   e,
	}

	jwks := &JWKS{
		Keys: []JWK{jwk},
	}

	return jwks, nil
}

func getPrivateKey() (any, error) {
	privateKeyPEM := os.Getenv("JWT_SECRET")
	privateKeyPEM = restoreNewlines(privateKeyPEM)

	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("failed to decode private key")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %v", err)
	}

	return privateKey, nil
}

func restoreNewlines(s string) string {
	return strings.ReplaceAll(s, "\\n", "\n")
}
