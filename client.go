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
	"k8s.io/kubernetes/pkg/api"
	kubeClient "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/labels"
)

type ClientInterface interface {
	Pods(namespace string, selector string) ([]api.Pod, error)
}

type Client struct {
	kubeClient *kubeClient.Client
}

func newClient(cfg *Config) (*Client, error) {
	config := kubeClient.Config{
		Host: cfg.Server,
	}

	c, err := kubeClient.New(&config)
	if err != nil {
		return nil, err
	}

	return &Client{
		kubeClient: c,
	}, nil
}

func resolveNamespace(namespace string) string {
	if len(namespace) > 0 {
		return namespace
	}
	return api.NamespaceDefault
}

func (c *Client) Pods(namespace string, selector string) ([]api.Pod, error) {
	s, err := labels.Parse(selector)
	if err != nil {
		return nil, err
	}
	podList, err := c.kubeClient.Pods(resolveNamespace(namespace)).List(s, nil)
	if err != nil {
		return nil, err
	}
	return podList.Items, nil
}