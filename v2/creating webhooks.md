
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
