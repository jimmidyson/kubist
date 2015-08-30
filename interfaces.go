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
	"github.com/fabric8io/kubist/api"
	"github.com/fabric8io/kubist/fields"
	"github.com/fabric8io/kubist/labels"
)

type Client interface {
	Namespace(namespace string) NamespacedClient
	Pods() PodReader
}

type NamespacedClient interface {
	Pods() PodReadWriter
}

type ResourceReadWriter interface {
	Pods() PodReadWriter
}

type PodReader interface {
	List(labels.Selector, fields.Selector) (*api.PodList, error)
}

type PodReadWriter interface {
	PodReader
	Get(name string) (*api.Pod, error)
	Create(pod *api.Pod) (*api.Pod, error)
	Replace(pod api.Pod) (*api.Pod, error)
	Delete(name string) error
	DeleteList(labels.Selector, fields.Selector) error
}
