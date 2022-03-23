### Download kubebuilder and Install locally.
```shell
curl -L -o kubebuilder https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)
chmod +x kubebuilder && mv kubebuilder /usr/local/bin/

```

### Create a Project
```shell
mkdir ~/kubebuilder/mysts
cd ~/kubebuilder/mysts
kubebuilder init --domain cncamp.io --repo cncamp.io/mysts
```

### Create API
```shell
kubebuilder create api --group apps --version v1beta1 --kind MyStatefulset
```

### Check Makefile
```shell
Usage:
  make <target>

General
  help             Display this help.

Development
  manifests        Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
  generate         Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
  fmt              Run go fmt against code.
  vet              Run go vet against code.
  test             Run tests.

Build
  build            Build manager binary.
  run              Run a controller from your host.
  docker-build     Build docker image with the manager.
  docker-push      Push docker image with the manager.

Deployment
  install          Install CRDs into the K8s cluster specified in ~/.kube/config.
  uninstall        Uninstall CRDs from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
  deploy           Deploy controller to the K8s cluster specified in ~/.kube/config.
  undeploy         Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
  controller-gen   Download controller-gen locally if necessary.
  kustomize        Download kustomize locally if necessary.
  envtest          Download envtest-setup locally if necessary.
```

### Edit api/v1alpha1/mystatefulset_types.go
```go
// MyStatefulsetSpec defines the desired state of MyStatefulset
type MyStatefulsetSpec struct {
        // INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
        // Important: Run "make" to regenerate code after modifying this file

        // Foo is an example field of MyStatefulset. Edit mystatefulset_types.go to remove/update
        Image string `json:"image,omitempty"`
}

// MyStatefulsetStatus defines the observed state of MyStatefulset
type MyStatefulsetStatus struct {
        // INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
        // Important: Run "make" to regenerate code after modifying this file
        AvaiableReplicas int `json:"avaiableReplicas,omitempty"`
}
```
Define the MyStatefulset spec contains Image and MyStatefulset Status spec contains AvaiableReplicas


### Edit controllers/mystatefulset_controller.go, add permissions to the controller
```go
//+kubebuilder:rbac:groups=apps.cncamp.io,resources=mystatefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.cncamp.io,resources=mystatefulsets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.cncamp.io,resources=mystatefulsets/finalizers,verbs=update
// Add the following
//+kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
```

### Generate crd
```shell
make manifests
```

### Build & install
```shell
make build
make docker-build
make docker-push
make deploy
```

### Enable webhooks

#### Install cert-manager
```shell
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.6.1/cert-manager.yaml
```

#### Create webhooks
```shell
kubebuilder create webhook --group apps --version v1beta1 --kind MyDaemonset --defaulting --programmatic-validation
```

#### Edit code
```go
// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *MyStatefulset) ValidateCreate() error {
        mystatefulsetlog.Info("validate create", "name", r.Name)
        if r.Spec.Replicas > 5 {
            return fmt.Errorof("replicas cannot be larger than 5")
        }

        // TODO(user): fill in your validation logic upon object creation.
        return nil
}
```


#### Enable webhook in config/default/kustomization.yaml
```shell
bases:
...
# [WEBHOOK] To enable webhook, uncomment all the sections with [WEBHOOK] prefix including the one in
# crd/kustomization.yaml
- ../webhook
# [CERTMANAGER] To enable cert-manager, uncomment all sections with 'CERTMANAGER'. 'WEBHOOK' components are required.
- ../certmanager
...
```

#### Redeploy
```shell
make 
make install 
make deploy
```