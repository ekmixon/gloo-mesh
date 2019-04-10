package utils

import (
	"io/ioutil"
	"log"
)

func DeployBookInfoIstio(namespace string) error {
	return KubectlApply(namespace, IstioBookInfoYaml)
}

func DeployBookInfoAppmesh(namespace string) error {
	return KubectlApply(namespace, AppmeshBookInfoYaml)
}

func DeployBookInfoWithInject(namespace, istioNamespace string) error {
	injected, err := IstioInject(istioNamespace, IstioBookInfoYaml)
	if err != nil {
		return err
	}

	return KubectlApply(namespace, injected)
}

// loads bookinfo from <root>/test/e2e/istio/files/bookinfo.yaml
var IstioBookInfoYaml = func() string {
	bookinfoYamlFile := MustTestFile("bookinfo.yaml")
	b, err := ioutil.ReadFile(bookinfoYamlFile)
	if err != nil {
		log.Fatalf("failed to read bookinfo for test: %v", err)
	}
	return string(b)
}()

const (
	AppmeshBookInfoYaml = BookinfoYaml + appmeshServices
)

const BookinfoYaml = `# source: istio-1.0.3/samples/bookinfo/platform/kube/bookinfo.yaml
# Copyright 2017 Istio Authors
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.

##################################################################################################
# Details service
##################################################################################################
apiVersion: v1
kind: Service
metadata:
  name: details
  labels:
    app: details
spec:
  ports:
  - port: 9080
    name: http
  selector:
    app: details
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: details-v1
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: details
        version: v1
        vn: details-v1
    spec:
      containers:
      - name: details
        image: istio/examples-bookinfo-details-v1:1.8.0
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 9080
---
##################################################################################################
# Ratings service
##################################################################################################
apiVersion: v1
kind: Service
metadata:
  name: ratings
  labels:
    app: ratings
spec:
  ports:
  - port: 9080
    name: http
  selector:
    app: ratings
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: ratings-v1
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: ratings
        version: v1
        vn: ratings-v1
    spec:
      containers:
      - name: ratings
        image: istio/examples-bookinfo-ratings-v1:1.8.0
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 9080
---

apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: reviews-v1
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: reviews
        version: v1
        vn: reviews-v1
    spec:
      containers:
      - name: reviews
        image: istio/examples-bookinfo-reviews-v1:1.8.0
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 9080
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: reviews-v2
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: reviews
        version: v2
        vn: reviews-v2
    spec:
      containers:
      - name: reviews
        image: istio/examples-bookinfo-reviews-v2:1.8.0
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 9080
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: reviews-v3
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: reviews
        version: v3
        vn: reviews-v3
    spec:
      containers:
      - name: reviews
        image: istio/examples-bookinfo-reviews-v3:1.8.0
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 9080
---
##################################################################################################
# Productpage services
##################################################################################################
apiVersion: v1
kind: Service
metadata:
  name: productpage
  labels:
    app: productpage
spec:
  ports:
  - port: 9080
    name: http
  selector:
    app: productpage
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: productpage-v1
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: productpage
        version: v1
        vn: productpage-v1
    spec:
      containers:
      - name: productpage
        image: istio/examples-bookinfo-productpage-v1:1.8.0
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 9080
---
`

const appmeshServices = `
apiVersion: v1
kind: Service
metadata:
  name: reviews
  labels:
    app: reviews
spec:
  ports:
    - port: 9080
      name: http
  selector:
    app: reviews
    version: v1
    vn: reviews-v1
---
apiVersion: v1
kind: Service
metadata:
  name: reviews-v2
  labels:
    app: reviews
spec:
  ports:
    - port: 9080
      name: http
  selector:
    app: reviews
    version: v2
    vn: reviews-v2
---
apiVersion: v1
kind: Service
metadata:
  name: reviews-v3
  labels:
    app: reviews
spec:
  ports:
    - port: 9080
      name: http
  selector:
    app: reviews
    version: v3
    vn: reviews-v3
`
