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

import (
	"net/http"
	"path"

	"github.com/fabric8io/kubist/api"
	"github.com/fabric8io/kubist/fields"
	"github.com/fabric8io/kubist/labels"
)

type pods struct {
	hc   *http.Client
	base string
}

func newPods(hc *http.Client, namespace string) *pods {
	base := "/api/v1"
	if len(namespace) > 0 {
		base = path.Join(base, "namespaces", namespace)
	}
	base = path.Join(base, "pods")
	return &pods{
		hc:   hc,
		base: base,
	}
}

func (p *pods) Get(name string) (*api.Pod, error) {
	u := path.Join(p.base, name)
	pod := &api.Pod{}
	err := doGet(p.hc, u, pod)
	return pod, err
}

func (p *pods) List(labels labels.Selector, fields fields.Selector) (*api.PodList, error) {
	pl := &api.PodList{}
	err := doGet(p.hc, p.base, pl)
	return pl, err
}

func (p *pods) Replace(pod api.Pod) (*api.Pod, error) {
	return &api.Pod{}, nil
}

func (p *pods) Delete(name string) error {
	u := path.Join(p.base, name)
	err := doDelete(p.hc, u)
	return err
}

func (p *pods) DeleteList(labels labels.Selector, fields fields.Selector) error {
	pl, err := p.List(labels, fields)
	if err != nil {
		return err
	}
	for _, pod := range pl.Items {
		err := p.Delete(pod.ObjectMeta.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *pods) Create(pod *api.Pod) (*api.Pod, error) {
	resp := &api.Pod{}
	err := doPost(p.hc, p.base, pod, resp)
	return resp, err
}
