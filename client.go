// Copyright (C) 2015 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubist

import "net/http"

func NewClient(config Config) (Client, error) {
	hc, err := newHTTPClient(config)
	if err != nil {
		return nil, err
	}
	return &client{hc: hc}, nil
}

type client struct {
	hc *http.Client
}

func (c *client) Namespace(namespace string) NamespacedClient {
	return &namespacedClient{c, namespace}
}

func (c *client) Pods() PodReader {
	return newPods(c.hc, "")
}

type namespacedClient struct {
	*client
	namespace string
}

func (c *namespacedClient) Pods() PodReadWriter {
	return newPods(c.hc, c.namespace)
}
