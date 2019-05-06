/*
Copyright 2019 The Kubernetes Authors.

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

package v2

import (
	"sigs.k8s.io/kubebuilder/pkg/scaffold/input"
)

var _ input.File = &Makefile{}

// Makefile scaffolds the Makefile
type Makefile struct {
	input.Input
	// Image is controller manager image name
	Image string

	// path for controller-tools pkg
	ControllerToolsPath string
}

// GetInput implements input.File
func (c *Makefile) GetInput() (input.Input, error) {
	if c.Path == "" {
		c.Path = "Makefile"
	}
	if c.Image == "" {
		c.Image = "controller:latest"
	}
	if c.ControllerToolsPath == "" {
		c.ControllerToolsPath = "vendor/sigs.k8s.io/controller-tools"
	}
	c.TemplateBody = makefileTemplate
	c.Input.IfExistsAction = input.Error
	return c.Input, nil
}

var makefileTemplate = `
# Image URL to use all building/pushing image targets
IMG ?= {{ .Image }}

all: test manager

# Run tests
test: generate fmt vet manifests
	go test ./api/... ./controllers/... -coverprofile cover.out

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	go run ./main.go

# Install CRDs into a cluster
install: manifests
	kubectl apply -f config/crds/bases

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	kubectl apply -f config/crds/bases
	kustomize build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests:
# TODO(droot): controller-gen will require fix to take new api-path as input, so disabling this for now
# 
#	go run {{ .ControllerToolsPath }}/cmd/controller-gen/main.go all --output-dir=config/crds/bases/

# Run go fmt against code
fmt:
	go fmt ./api/... ./controllers/...

# Run go vet against code
vet:
	go vet ./api/... ./controllers/...

# Generate code
generate:
ifndef GOPATH
	$(error GOPATH not defined, please define GOPATH. Run "go help gopath" to learn more about GOPATH)
endif
	go run ./vendor/k8s.io/code-generator/cmd/deepcopy-gen/main.go -O zz_generated.deepcopy --go-header-file ./hack/boilerplate.go.txt -i ./api/...
	

# Build the docker image
docker-build: test
	docker build . -t ${IMG}
	@echo "updating kustomize image patch file for manager resource"
	sed -i'' -e 's@image: .*@image: '"${IMG}"'@' ./config/default/manager_image_patch.yaml

# Push the docker image
docker-push:
	docker push ${IMG}
`
