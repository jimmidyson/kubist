package kubist

import "net/url"

type Config struct {
	// The Kubernetes master URL.
	Master *url.URL
	// The HTTP basic authentication credentials for the targets.
	BasicAuth *BasicAuth
	// The bearer token for the targets.
	BearerToken string
	// The bearer token file for the targets.
	BearerTokenFile string
	// The ca cert to use for the targets.
	CAFile string
	// The client cert authentication credentials for the targets.
	ClientCert *ClientCert
	// Disable validation of server certificate
	Insecure bool
}

// ClientCert contains client cert credentials.
type ClientCert struct {
	CertFile string
	KeyFile  string
}

// BasicAuth contains basic HTTP authentication credentials.
type BasicAuth struct {
	Username string
	Password string
}

func InClusterConfig() Config {
	m, _ := url.Parse("https://kubernetes.default.svc/")
	return Config{
		Master:          m,
		CAFile:          "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
		BearerTokenFile: "/var/run/secrets/kubernetes.io/serviceaccount/token",
	}
}
