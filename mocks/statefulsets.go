// Copyright 2016 Mirantis
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mocks

import (
	"fmt"
	"strings"

	"k8s.io/client-go/kubernetes/typed/apps/v1beta1"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	appsbeta1 "k8s.io/client-go/pkg/apis/apps/v1beta1"
	"k8s.io/client-go/pkg/watch"
	"k8s.io/client-go/rest"
)

type statefulSetClient struct {
}

func pointer(i int32) *int32 {
	return &i
}

// MakeStatefulSet returns a new K8s StatefulSet object for the client to return. If it's name is "fail" it will have labels that will cause it's underlying mock Pods to fail.
func MakeStatefulSet(name string) *appsbeta1.StatefulSet {
	statefulSet := &appsbeta1.StatefulSet{}
	statefulSet.Name = name
	statefulSet.Spec.Replicas = pointer(int32(3))
	statefulSet.Spec.Template.ObjectMeta.Labels = make(map[string]string)
	if name == "fail" {
		statefulSet.Spec.Template.ObjectMeta.Labels["failedpod"] = "yes"
		statefulSet.Status.Replicas = int32(2)
	} else {
		statefulSet.Status.Replicas = int32(3)
	}

	return statefulSet
}

func (r *statefulSetClient) List(opts v1.ListOptions) (*appsbeta1.StatefulSetList, error) {
	var statefulSets []appsbeta1.StatefulSet
	for i := 0; i < 3; i++ {
		statefulSets = append(statefulSets, *MakeStatefulSet(fmt.Sprintf("ready-trolo%d", i)))
	}

	// use ListOptions.LabelSelector to check if there should be any failed statefulSets
	if strings.Index(opts.LabelSelector, "failedstatefulSet=yes") >= 0 {
		statefulSets = append(statefulSets, *MakeStatefulSet("fail"))
	}

	return &appsbeta1.StatefulSetList{Items: statefulSets}, nil
}

func (r *statefulSetClient) Get(name string) (*appsbeta1.StatefulSet, error) {
	status := strings.Split(name, "-")[0]
	if status == "error" {
		return nil, fmt.Errorf("mock service %s returned error", name)
	}

	return MakeStatefulSet(name), nil
}

func (r *statefulSetClient) Create(rs *appsbeta1.StatefulSet) (*appsbeta1.StatefulSet, error) {
	return MakeStatefulSet(rs.Name), nil
}

func (r *statefulSetClient) Update(rs *appsbeta1.StatefulSet) (*appsbeta1.StatefulSet, error) {
	panic("not implemented")
}

func (r *statefulSetClient) UpdateStatus(rs *appsbeta1.StatefulSet) (*appsbeta1.StatefulSet, error) {
	panic("not implemented")
}

func (r *statefulSetClient) Delete(name string, options *v1.DeleteOptions) error {
	panic("not implemented")
}

func (r *statefulSetClient) Watch(opts v1.ListOptions) (watch.Interface, error) {
	panic("not implemented")
}

func (r *statefulSetClient) ProxyGet(scheme string, name string, port string, path string, params map[string]string) rest.ResponseWrapper {
	panic("not implemented")
}

func (r *statefulSetClient) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	panic("not implemented")
}

func (r *statefulSetClient) Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *appsbeta1.StatefulSet, err error) {
	panic("not implemented")
}

// NewStatefulSetClient is a client constructor
func NewStatefulSetClient() v1beta1.StatefulSetInterface {
	return &statefulSetClient{}
}