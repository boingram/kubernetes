/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package internalversion

import (
	rest "k8s.io/client-go/rest"
	"k8s.io/kube-apiextensions-server/pkg/client/clientset/internalclientset/scheme"
)

type ApiextensionsInterface interface {
	RESTClient() rest.Interface
	CustomResourcesGetter
}

// ApiextensionsClient is used to interact with features provided by the apiextensions.k8s.io group.
type ApiextensionsClient struct {
	restClient rest.Interface
}

func (c *ApiextensionsClient) CustomResources() CustomResourceInterface {
	return newCustomResources(c)
}

// NewForConfig creates a new ApiextensionsClient for the given config.
func NewForConfig(c *rest.Config) (*ApiextensionsClient, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &ApiextensionsClient{client}, nil
}

// NewForConfigOrDie creates a new ApiextensionsClient for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *ApiextensionsClient {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new ApiextensionsClient for the given RESTClient.
func New(c rest.Interface) *ApiextensionsClient {
	return &ApiextensionsClient{c}
}

func setConfigDefaults(config *rest.Config) error {
	g, err := scheme.Registry.Group("apiextensions.k8s.io")
	if err != nil {
		return err
	}

	config.APIPath = "/apis"
	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}
	if config.GroupVersion == nil || config.GroupVersion.Group != g.GroupVersion.Group {
		gv := g.GroupVersion
		config.GroupVersion = &gv
	}
	config.NegotiatedSerializer = scheme.Codecs

	if config.QPS == 0 {
		config.QPS = 5
	}
	if config.Burst == 0 {
		config.Burst = 10
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *ApiextensionsClient) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
