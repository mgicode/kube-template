// Copyright © 2015 Victor Antonovich <victor@antonovich.me>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"sync"

	"k8s.io/kubernetes/pkg/api"
)

type DependencyManager struct {
	sync.RWMutex
	// Kubernetes client
	client *Client
	// Cached dependencies
	cachedDeps map[string]interface{}
}

func newDependencyManager(client *Client) *DependencyManager {
	return &DependencyManager{
		client:     client,
		cachedDeps: make(map[string]interface{}),
	}
}

func (dm *DependencyManager) cachedDependency(key string) (interface{}, bool) {
	dm.RLock()
	defer dm.RUnlock()
	value, found := dm.cachedDeps[key]
	return value, found
}

func (dm *DependencyManager) cacheDependency(key string, dep interface{}) {
	dm.Lock()
	defer dm.Unlock()

	dm.cachedDeps[key] = dep
}

// Get Kubernetes pods with caching
func (dm *DependencyManager) Pods(namespace, selector string) ([]api.Pod, error) {
	key := fmt.Sprintf("pods(%s,%s)", namespace, selector)
	if value, found := dm.cachedDependency(key); found {
		return value.([]api.Pod), nil
	}
	pods, err := dm.client.Pods(namespace, selector)
	if err != nil {
		return nil, err
	}
	dm.cacheDependency(key, pods)
	return pods, nil
}

// Get Kubernetes services with caching
func (dm *DependencyManager) Services(namespace, selector string) ([]api.Service, error) {
	key := fmt.Sprintf("services(%s,%s)", namespace, selector)
	if value, found := dm.cachedDependency(key); found {
		return value.([]api.Service), nil
	}
	services, err := dm.client.Services(namespace, selector)
	if err != nil {
		return nil, err
	}
	dm.cacheDependency(key, services)
	return services, nil
}