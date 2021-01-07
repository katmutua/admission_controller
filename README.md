# admission_controller

- How to create an admission controller in Kubernetes
    - we can dynamically configure what resources are subject to what admission webhooks via
      `ValidatingWebhookConfiguration` or `MutatingWebhookConfiguration`.

- Kubernetes API is based on REST model; gives possibility of managing all workloads using HTTP requests

Two types of admission controllers in Kubernetes
  - Validating admission controller:
    proxies the requests to the subscribed web hooks
    Kubernetes API registers the webhooks based on the resource type and the request method

    In case the validation webhook rejects the request, the Kubernetes API  returns a failed HTTP
    response to the user  otherwise continues to the next admission.

  - Mutating admission controller
    Modifies the resource submitted by the user so that you can create defaults or validate the schema

    GOAL: Create a simple validation controller which enabled us to influence the pod creation
    Controller name "gandalf" and will reject all new pods with a name different than "shire-pod"

#### Needed:

   The Kubernetes server needs to know when to send an incoming request to our admissions controller
    Kubernetes philosophy advocates for using a declarative strategy.

    1. Define a ValidationWebhookConfiguration that gives the information needed to the API

     `apiVersion: admissionregistration.k8s.io/v1beta1
      kind: ValidatingWebhookConfiguration  # The name for the webhook admission object.
      metadata:
        name: <controller_name>  # The name for the webhook admission object.
      webhooks:
      - name: <webhook_name> # The name of the webhook to call.
        clientConfig: #  Information about how to connect to, trust, and send data to the webhook server.
          service:
            namespace: default  #  The project where the front-end service is created.
            name: kubernetes #  The name of the front-end service.
           path: <webhook_url> #  The webhook URL used for admission requests.
          caBundle: <cert> #  A PEM-encoded CA certificate that signs the server certificate used by the webhook server.
        rules: #  Rules that define when the API server should use this controller.
        - operations:
          - <operation>
          apiGroups:
          - ""
          apiVersions:
          - "*"
          resources:
          - <resource>
        failurePolicy: <policy>  # Specifies how the policy should proceed if the webhook admission server is unavailable. Either Ignore (allow/fail open) or Fail (block/fail closed).
        `
    clientConfig,: defines where our service can be found (it can be an external URL)
                  and the path which our validation server will listen on
                  Since security is always important, adding the cert authority will tell the Kubernetes API
                  to use HTTPS and validate our server using the passed asset.

   The second part specifies which rules the API will follow to decide if a request should be forwarded for validation
   or not

   Here it is configured that only requests with method equal to CREATE and resource type pod will be
   forwarded.

#### Generating the certificates and the CA
  - run the generate_Ca.sh script to generate your certificates
  Besides creating the certificates and the CA, the script later injects it into the manifest used to deploy our server.
      cat manifest.yaml | grep caBundle

  - We need to create a secret to place the certificates. After we apply the manifest, the pod will be able to store the
  secret files into a directory.
```
  kubectl create secret generic gandalf -n default --from-file=key.pem=certs/gandalf-key.pem --from-file=cert.pem=certs/gandalf-crt.pem
```
#### Deploying the Controller
  - we will use a deployment with a single replica which mounts the certs generated to expose a secure REST endpoint where the pod request will be submitted
      `< see the deployment config in the manifest >`

apply the manifest

``` kubectl apply -f manifest.yaml
```
Now the server should be running and ready to validate the creation of new pods.

#### Verify that the validation controller works
Let's try creating a pod with a non matching name
 `< see invalid config in invalid pod config >`

```
$ kubectl apply -f non-shire-app.yaml
```
We should get an error

```
  Error from server: error when creating "non-shire-app.yaml": admission webhook "gandalf-webhook" denied the request: Keep calm and don't add more crap to the cluster!
```

#### Building local image

`docker build -t  <image_name> .
  docker tag
`
