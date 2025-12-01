package oidc

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty string `json:"kty"`           // Key Type: "RSA", "EC"
	Kid string `json:"kid"`           // Key ID
	Use string `json:"use,omitempty"` // Usually "sig"
	Alg string `json:"alg,omitempty"` // e.g. "RS256"

	// RSA Fields
	N string `json:"n,omitempty"` // Modulus
	E string `json:"e,omitempty"` // Exponent

	// EC Fields
	Crv string `json:"crv,omitempty"`
	X   string `json:"x,omitempty"`
	Y   string `json:"y,omitempty"`

	// X.509 fields
	X5c     []string `json:"x5c,omitempty"`
	X5t     string   `json:"x5t,omitempty"`
	X5tS256 string   `json:"x5t#S256,omitempty"`
}
