package rest

import (
	"crypto/tls"
	"net/http"
	"time"
)

// NewServer returns an HTTP server with all required configuration
// See https://blog.cloudflare.com/exposing-go-on-the-internet/ for details about these settings
func NewServer(mux *http.ServeMux, serverAddress string) *http.Server {
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true, // Causes servers to use Go's default cipher suite preferences, which are tuned to avoid attacks.
		CurvePreferences: []tls.CurveID{ // Only use curves which have assembly implementations
			tls.CurveP256,
		},
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	srv := &http.Server{
		Addr:         serverAddress,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSConfig:    tlsConfig,
		Handler:      mux,
	}

	return srv
}
