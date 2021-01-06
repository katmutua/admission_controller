# admission_controller

- How to create an admission controller in Kubernetes

Kubernetes API
 - based on REST model; gives possibility of managing all workloads using HTTP requests


Two types of admission controllers in Kubernetes
  - Validating admission controller :
    proxies the requests to the subscribed web hooks
    Kubernetes API registers the webhooks based on the resource type and the request method

    In case the validation webhook rejects the request, the Kubernetes API  returns a failed HTTP
    response to the user  otherwise continues to the next admission.


  - Mutating admission controller
    Modifies the resource submitted by the user so that you can create defaults or validate the schema


    GOAL: Create a simple validation controller which enabld us to influence the pod creation
    Controller name "gandalf" and will reject all new pods with a name different than "shire-pod"


#### Needed:

    The Kubernetes server needs to know when to send an incoming request to our admissions controller
    Kubernetes philosophy advocates deo using a declarative strategy.

    1. Define a ValidationWebhookConfiguration that gives the information needed to the API
```
       apiVersion: admissionregistration.k8s.io/v1beta1
       kind: ValidationWebhookConfiguration
       metadata:
         name: gandalf
      webhooks:
        - name: gandalf
          clientConfig:
            service:
            name: gandalf
            namespace: default
            path "/validate"
          caBundle: "${CA_BUNDLE}"
        rules:
          - operations ["CREATE"]
            apiGroups[""]
            apiVersions: ["v1"]
            resources: ["pods"]
```
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
      kubectl create secret generic gandalf -n default \
        --from-file=key.pem=certs/gandalf-key.pem \
        --from-file=cert.pem=certs/gandalf-crt.pem
```
#### Deploying the Controller
  - we will use a deployment with a single replica which mounts the certs generated to expose a secure REST endpoint where the pod request will be submitted
```
        apiVersion: apps/v1beta1
      kind: Deployment
      metadata:
        name: gandalf
        namespace: default
      spec:
        replicas: 1
        template:
          spec:
            containers:
              - name: webhook
                image: giantswarm/gandalf:1.0.0
                ...
                volumeMounts:
                  - name: webhook-certs
                    mountPath: /etc/certs
              ...
            volumes:
              - name: webhook-certs
                secret:
                  secretName: gandalf
      ---
      apiVersion: v1
      kind: Service
      metadata:
        name: gandalf
        namespace: default
      spec:
        ports:
        - name: webhook
          port: 443
          targetPort: 8080
          ...
```

apply the manifest

``` kubectl apply -f manifest.yaml
```
Now the server should be running and ready to validate the creation of new pods.


#### Verify that the validation controller works
Let's try creating a pod with a non matching name
```
apiVersion: v1
kind: Pod
metadata:
  name: non-shire-app
spec:
  containers:
  - image: busybox
    name: non-shire-app
```

```
$ kubectl apply -f non-shire-app.yaml
```
We should get an error

```
Error from server: error when creating "non-shire-app.yaml": admission webhook "gandalf-webhook" denied the request: Keep calm and don't add more crap to the cluster!
```
