// Copyright 2013 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubist

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/fabric8io/kubist/api"
)

// newHTTPClient returns a http.Client using the specified config.
func newHTTPClient(config Config) (*http.Client, error) {
	rt := http.DefaultTransport

	tlsConfig := &tls.Config{InsecureSkipVerify: config.Insecure}

	// Load client cert if specified.
	if config.ClientCert != nil {
		cert, err := tls.LoadX509KeyPair(config.ClientCert.CertFile, config.ClientCert.KeyFile)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	caCertPool := x509.NewCertPool()
	if len(config.CAFile) > 0 {
		// Load CA cert.
		caCert, err := ioutil.ReadFile(config.CAFile)
		if err != nil {
			return nil, err
		}
		caCertPool.AppendCertsFromPEM(caCert)
	}
	tlsConfig.RootCAs = caCertPool

	tlsConfig.BuildNameToCertificate()

	tr := rt.(*http.Transport)
	tr.TLSClientConfig = tlsConfig

	if len(config.BearerToken) == 0 && len(config.BearerTokenFile) > 0 {
		bt, err := ioutil.ReadFile(config.BearerTokenFile)
		if err != nil {
			return nil, fmt.Errorf("Unable to read bearer token file(%s): %v", config.BearerTokenFile, err)
		}
		config.BearerToken = string(bt)
	}
	if len(config.BearerToken) > 0 {
		rt = &bearerAuthRoundTripper{
			bearerToken: config.BearerToken,
			rt:          rt,
		}
	} else if config.BasicAuth != nil {
		rt = &basicAuthRoundTripper{
			username: config.BasicAuth.Username,
			password: config.BasicAuth.Password,
			rt:       rt,
		}
	}

	rt = &reqUpdaterRoundTripper{
		master: config.Master,
		rt:     rt,
	}

	return &http.Client{Transport: rt}, nil
}

type reqUpdaterRoundTripper struct {
	master *url.URL
	rt     http.RoundTripper
}

func (rt *reqUpdaterRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req = cloneRequest(req)
	req.URL.Scheme = rt.master.Scheme
	req.URL.Host = rt.master.Host
	req.Header.Set("Accept", "application/json")
	return rt.rt.RoundTrip(req)
}

type bearerAuthRoundTripper struct {
	bearerToken string
	rt          http.RoundTripper
}

func (rt *bearerAuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if len(req.Header.Get("Authorization")) == 0 {
		req = cloneRequest(req)
		req.Header.Set("Authorization", "Bearer "+rt.bearerToken)
	}

	return rt.rt.RoundTrip(req)
}

type basicAuthRoundTripper struct {
	username string
	password string
	rt       http.RoundTripper
}

func (rt *basicAuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if len(req.Header.Get("Authorization")) != 0 {
		return rt.rt.RoundTrip(req)
	}
	req = cloneRequest(req)
	req.SetBasicAuth(rt.username, rt.password)
	return rt.rt.RoundTrip(req)
}

func cloneRequest(r *http.Request) *http.Request {
	// Shallow copy of the struct.
	r2 := new(http.Request)
	*r2 = *r
	// Deep copy of the Header.
	r2.Header = make(http.Header)
	for k, s := range r.Header {
		r2.Header[k] = s
	}
	return r2
}

func handleError(r *http.Response) error {
	defer r.Body.Close()

	e := &KubernetesError{
		StatusCode: r.StatusCode,
		StatusText: r.Status,
	}

	if b, err := ioutil.ReadAll(r.Body); err == nil {
		e.Response = string(b)
	}

	s := &api.Status{}
	if err := json.NewDecoder(r.Body).Decode(s); err == nil {
		e.Status = s
	}

	return e
}

func doGet(hc *http.Client, path string, o interface{}) error {
	r, err := hc.Get(path)
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return handleError(r)
	}

	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(o)
	if err != nil {
		return err
	}
	return nil
}

func doDelete(hc *http.Client, path string) error {
	req, err := http.NewRequest("DELETE", path, nil)
	r, err := hc.Do(req)
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return handleError(r)
	}

	return nil
}

func doPost(hc *http.Client, path string, reqObj interface{}, respObj interface{}) error {
	contentJson, _ := json.Marshal(reqObj)
	contentReader := bytes.NewReader(contentJson)
	r, err := hc.Post(path, "application/json", contentReader)
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return handleError(r)
	}

	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(respObj)
	if err != nil {
		return err
	}
	return nil
}
