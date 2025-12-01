package oidc

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type OIDCProvider struct {
	Issuer string
	Realm  string
	JWKS   JWKS
}

func getJWKS(url, realm string) (JWKS, error) {
	resp, err := http.Get(fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", url, realm))
	if err != nil {
		return JWKS{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var jwks JWKS
	if err = json.Unmarshal(body, &jwks); err != nil {
		return JWKS{}, err
	}

	return jwks, nil
}

func NewOIDCProvider(url, realm string) *OIDCProvider {

	jwks, err := getJWKS(url, realm)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	return &OIDCProvider{
		Issuer: url,
		Realm:  realm,
		JWKS:   jwks,
	}
}

func (p *OIDCProvider) Validate(token string) error {
	if token == "" {
		return errors.New("empty token")
	}

	// 1. Parse sin verificar para obtener el kid del header
	parser := jwt.NewParser()
	unverified, _, err := parser.ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return fmt.Errorf("invalid token format: %w", err)
	}

	kidRaw, ok := unverified.Header["kid"]
	if !ok {
		return errors.New("token missing kid")
	}

	kid, ok := kidRaw.(string)
	if !ok || kid == "" {
		return errors.New("invalid kid format")
	}

	// 2. Buscar la JWK correspondiente en el JWKS
	jwk, err := p.findKey(kid)
	if err != nil {
		return err
	}

	// 3. Construir clave RSA pública a partir del JWK
	pubKey, err := buildRSAPublicKey(jwk)
	if err != nil {
		return fmt.Errorf("invalid jwk: %w", err)
	}

	// 4. Validar el JWT con la clave
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// Validar algoritmo esperado
		if t.Method.Alg() != jwk.Alg {
			return nil, fmt.Errorf("unexpected signing algorithm: %s", t.Method.Alg())
		}
		return pubKey, nil
	})
	if err != nil {
		return fmt.Errorf("token verification failed: %w", err)
	}

	if !parsed.Valid {
		return errors.New("invalid token")
	}

	// 5. Validar claims estándar
	claims := parsed.Claims.(jwt.MapClaims)

	// Validar issuer si lo configuraste
	if p.Issuer != "" {
		iss, _ := claims["iss"].(string)
		if iss != fmt.Sprintf("%s/realms/%s", p.Issuer, p.Realm) {
			return fmt.Errorf("invalid issuer: %s | %s", iss, p.Issuer)
		}
	}

	return nil
}

// Buscar clave por kid
func (p *OIDCProvider) findKey(kid string) (*JWK, error) {
	for _, k := range p.JWKS.Keys {
		if k.Kid == kid {
			return &k, nil
		}
	}
	return nil, fmt.Errorf("kid %s not found in JWKS", kid)
}

// Construir rsa.PublicKey desde JWK
func buildRSAPublicKey(jwk *JWK) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("invalid modulus n: %w", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("invalid exponent e: %w", err)
	}

	var eInt int
	for _, b := range eBytes {
		eInt = eInt<<8 + int(b)
	}

	if eInt <= 0 {
		return nil, errors.New("invalid exponent")
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: eInt,
	}, nil
}
