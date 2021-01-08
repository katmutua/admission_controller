
## Creatinga webhook

source: https://docs.openshift.com/container-platform/3.11/architecture/additional_concepts/dynamic_admission_controllers.html


some common uses of admission control
- mutating resources to inject sidecar containers into pods
- restricting projects to block some resources from a project
- custom resource validation to perform complex validation on dependant fields


Types of admission webhook

- mutating admission webhook
  invoked during the mutation phase of the admission process which allows modification of the resource
  content before it is persisted.

  An example of this is on openshift container platform there is the pods node selector that adds annotation on a namespace
  to find a label selector and add it as part of the pod specification

  the mutating webhook definition will look like this
  ```
  apiVersion: admissionregistration.k8s.io/v1beta1
    kind: MutatingWebhookConfiguration
    metadata:
      name: <controller_name> # The name for the admission webhook object.
    webhooks:
    - name: <webhook_name> # The name of the webhook to call.
      clientConfig: # Information about how to connect to, trust, and send data to the webhook server.
        service:
          namespace:  The project where the front-end service is created
          name: The name of the front-end service.
         path: <webhook_url> #The webhook URL used for admission requests.
        caBundle: <cert>  #	A PEM-encoded CA certificate that signs the server certificate used by the webhook server.
      rules: # Rules that define when the API server should use this controller.
      - operations:
        - <operation>
        apiGroups:
        - ""
        apiVersions:
        - "*"
        resources:
        - <resource>
      failurePolicy: <policy>

  ```

  - Validating admission webhook
    These are invoked in the validation phase of admission control.
    This phase allows the enforcement of invariants on a particular API
    resources to ensure that the resource does not change again.


    The Pod Node Selector is also a type of validation admission by ensuring that
    all nodeSelector fields are constrained by the node selector restrictions on the
    project

  Sample Validating Admission Webhook Configuration

  ```
  apiVersion: admissionregistration.k8s.io/v1beta1
  kind: ValidatingWebhookConfiguration
  metadata:
    name: <controller_name>
  webhooks:
  - name: <webhook_name>
    clientConfig:
      service:
        namespace: default  
        name: kubernetes
       path: <webhook_url>
      caBundle: <cert>
    rules:
    - operations: # The operation(s) that triggers the API server to call this controller:
      - <operation>
      apiGroups:
      - ""
      apiVersions:
      - "*"
      resources:
      - <resource>
    failurePolicy: <policy>  	
        # Specifies how the policy should proceed if the webhook admission server is unavailable.
        # Either Ignore (allow/fail open) or Fail (block/fail closed).
  ```

### Creating the Admission Webhook

- First deploy the external webhook server and ensure that it is working properly.
  [ try create a local server with go and run the code with the validate endpoint ]

- Configure a mutating or a validating admission webhook object in a YAML file

- Create a front-end service for the admission webhook
  ```
  apiVersion: v1
  kind: Service
  metadata:
    labels:
      role: webhook #Free-form label to trigger the webhook
    name: <name>
  spec:
    selector:
     role: webhook # 	Free-form label to trigger the webhook.
  ```
- Apply the manifest
   - You can add the components in one manifest file and apply as follows
       ```
       kubectl apply -f manifest.yaml  
       ```
- Add the admission webhook name to pods you want to be controlled by the webhook
  Add the name of the webhook to the pod spec
  ```
  apiVersion: v1
  kind: Pod
  metadata:
    labels:
      role: webhook
    name: <name>
  spec:
    containers:
      - name: <name>
        image: myrepo/myimage:latest
        imagePullPolicy: <policy>
        ports:
         - containerPort: 8000
  ```

##### Admission Webhook Example
```
apiVersion: admissionregistration.k8s.io/v1beta1
kind: ValidatingWebhookConfiguration
metadata:
  name: namespacereservations.admission.online.openshift.io
webhooks:
- name: namespacereservations.admission.online.openshift.io
  clientConfig:
    service:
      namespace: default
      name: webhooks
     path: /apis/admission.online.openshift.io/v1beta1/namespacereservations
    caBundle: KUBE_CA_HERE
  rules:
  - operations:
    - CREATE
    apiGroups:
    - ""
    apiVersions:
    - "b1"
    resources:
    - namespaces
  failurePolicy: Ignore
```
The following is an example pod that will be evaluated by the admission webhook named webhook:
  ```
  apiVersion: v1
  kind: Pod
  metadata:
    labels:
      role: webhook
    name: webhook
  spec:
    containers:
      - name: webhook
        image: myrepo/myimage:latest
        imagePullPolicy: IfNotPresent
        ports:
  - containerPort: 8000
  ```

  The following is the front-end service for the webhook:

  ```
  apiVersion: v1
  kind: Service
  metadata:
   labels:
     role: webhook
   name: webhook
  spec:
   ports:
     - port: 443
       targetPort: 8000
   selector:
  role: webhook
  ```
