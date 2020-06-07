/*
Copyright 2020 The Helm Broker Authors.

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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	addonsv1alpha1 "github.com/kyma-project/helm-broker/pkg/apis/addons/v1alpha1"
	versioned "github.com/kyma-project/helm-broker/pkg/client/clientset/versioned"
	internalinterfaces "github.com/kyma-project/helm-broker/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/kyma-project/helm-broker/pkg/client/listers/addons/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// AddonsConfigurationInformer provides access to a shared informer and lister for
// AddonsConfigurations.
type AddonsConfigurationInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.AddonsConfigurationLister
}

type addonsConfigurationInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewAddonsConfigurationInformer constructs a new informer for AddonsConfiguration type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewAddonsConfigurationInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredAddonsConfigurationInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredAddonsConfigurationInformer constructs a new informer for AddonsConfiguration type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredAddonsConfigurationInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AddonsV1alpha1().AddonsConfigurations(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AddonsV1alpha1().AddonsConfigurations(namespace).Watch(context.TODO(), options)
			},
		},
		&addonsv1alpha1.AddonsConfiguration{},
		resyncPeriod,
		indexers,
	)
}

func (f *addonsConfigurationInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredAddonsConfigurationInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *addonsConfigurationInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&addonsv1alpha1.AddonsConfiguration{}, f.defaultInformer)
}

func (f *addonsConfigurationInformer) Lister() v1alpha1.AddonsConfigurationLister {
	return v1alpha1.NewAddonsConfigurationLister(f.Informer().GetIndexer())
}
