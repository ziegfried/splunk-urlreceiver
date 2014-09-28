package server

import (
	"fmt"
	"io/ioutil"
	. "modinputs"
	"net/http"
	"os"
	"path/filepath"
	. "scheme"
	"strings"
	. "utils"
)

var servers map[int]*Mux = make(map[int]*Mux)

func Get(port int, stanza InputStanza, config InputConfig) (*Mux, error) {
	ssl := NormalizeBoolean(stanza.GetParamValue(PARAM_SSL, "false"), false)
	server, found := servers[port]
	if !found {
		Log.Debug("Creating new HTTP server on port=%d", port)
		server = new(Mux)
		server.port = port
		server.m = make(map[string]http.Handler)
		server.ssl = ssl
		if ssl {
			if err := server.setupSSL(stanza, config); err != nil {
				return server, fmt.Errorf("Error setting up SSL: %s", err)
			}
		}
		servers[port] = server
	}
	if server.ssl != ssl {
		Log.Error("Inputs need server on port=%d both with and without SSL. This is not possible.", port)
		return server, fmt.Errorf("Server on port %d does not match SSL setting", port)
	}
	return server, nil
}

type Mux struct {
	port        int
	ssl         bool
	sslKeyPath  string
	sslCertPath string
	m           map[string]http.Handler
}

func (mux *Mux) Register(path string, h http.Handler) {
	if _, exists := mux.m[path]; exists {
		Log.Warn("Handler for path %s already exists (duplicate)", path)
	} else {
		mux.m[path] = h
	}
}

func (mux *Mux) setupSSL(stanza InputStanza, config InputConfig) error {
	debug := NormalizeBoolean(stanza.GetParamValue(PARAM_DEBUG, "false"), false)
	if debug {
		Log.Info("[%s] Setting up SSL support", stanza.Name)
	}

	certPath := os.ExpandEnv(stanza.GetParamValue(PARAM_SSL_CERT_PATH, ""))
	keyPath := os.ExpandEnv(stanza.GetParamValue(PARAM_SSL_KEY_PATH, ""))

	if certPath == "" {
		if NormalizeBoolean(stanza.GetParamValue(PARAM_SSL_SELF_SIGNED, "false"), false) {
			if debug {
				Log.Info("[%s] Generating self-signed SSL certificate", stanza.Name)
			}
			host := stanza.GetParamValue(PARAM_SSL_HOST, config.ServerHost)
			bits := NormalizeInt(stanza.GetParamValue(PARAM_SSL_KEY_BITS, "2048"), 2048)
			prefix := SanitizeFilename(strings.Replace(stanza.Name, "urlreceiver://", "", 1))
			certPath, keyPath, err := generateSelfSignedCertificate(host, bits, os.TempDir(), prefix)

			if err != nil {
				return fmt.Errorf("Error creating self-sigend certificate: %s", err)
			} else {
				mux.sslCertPath = certPath
				mux.sslKeyPath = keyPath
				if debug {
					Log.Info("[%s] Using generated SSL cert=%s and key=%s", stanza.Name, certPath, keyPath)
				}
			}
		} else {
			certData := stanza.GetParamValue(PARAM_SSL_CERT, "")
			keyData := stanza.GetParamValue(PARAM_SSL_KEY, "")
			if certData != "" && keyData != "" {
				prefix := filepath.Join(os.TempDir(), SanitizeFilename(strings.Replace(stanza.Name, "urlreceiver://", "", 1)))
				certPath := fmt.Sprintf("%s_cert.pem", prefix)
				keyPath := fmt.Sprintf("%s_key.pem", prefix)
				if err := ioutil.WriteFile(certPath, []byte(certData), 0600); err != nil {
					return fmt.Errorf("[%s] Error writing SSL certificate to file=%s", stanza.Name, certPath)
				}
				if err := ioutil.WriteFile(keyPath, []byte(keyData), 0600); err != nil {
					return fmt.Errorf("[%s] Error writing SSL private key to file=%s", stanza.Name, certPath)
				}
				mux.sslCertPath = certPath
				mux.sslKeyPath = keyPath
				if debug {
					Log.Info("[%s] Using just written SSL cert=%s and key=%s", stanza.Name, certPath, keyPath)
				}
			} else {
				return fmt.Errorf("[%s] No SSL certificate provided", stanza.Name)
			}
		}
	} else {
		if debug {
			Log.Info("[%s] Using configured SSL cert=%s and key=%s", stanza.Name, certPath, keyPath)
		}
		mux.sslCertPath = certPath
		mux.sslKeyPath = keyPath
	}

	return nil
}

func (mux *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Log.Debug("Incoming request on port=%d -> method=%s uri=%s length=%d from clientip=%s", mux.port, r.Method, r.URL.Path, r.ContentLength, r.RemoteAddr)

	if handler, found := mux.m[r.URL.Path]; found {
		handler.ServeHTTP(w, r)
	} else {
		Log.Warn("Unhandled request method=%s path=%s on port=%d from clientip=%s", r.Method, r.URL.Path, mux.port, r.RemoteAddr)
		http.Error(w, "404", 404)
	}
}

func StartListening() {
	for port, server := range servers {
		if server.ssl {
			go startServerTLS(port, server)
		} else {
			go startServer(port, server)
		}
	}
}

func startServer(port int, server *Mux) {
	Log.Info("Starting webserver on port=%d with ssl=disabled", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), server)
	if err != nil {
		Log.Error("Error starting webserver on port=%d: %s", port, err)
	}
}

func startServerTLS(port int, server *Mux) {
	Log.Info("Starting webserver on port=%d with ssl=enabled with certPath=%s keyPath=%s", port, server.sslCertPath, server.sslKeyPath)
	err := http.ListenAndServeTLS(fmt.Sprintf(":%d", port), server.sslCertPath, server.sslKeyPath, server)
	if err != nil {
		Log.Error("Error starting webserver on port=%d: %s", port, err)
	}
}
