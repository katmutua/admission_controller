#!/bin/bash
# Create a temp certs directory
mkdir `pwd`/certs || true

# Create the certs with cfssl. Use the configs in
# dockerfiles/cfssl-certs/configs to create a cert specific for the
# webhook server.
# The configured service name is "admissions-webhook"
# The configured namespace is "admissions"
# THe cert is valid for admissions-webhook.admissions.svc:443
CFSSL_IMAGE=`docker build dockerfiles/cfssl-certs/. -q`
docker run -it --rm -v `pwd`/certs:/tmp/ $CFSSL_IMAGE sh -c '
  cfssl gencert -initca /configs/ca.json | cfssljson -bare ca && \
  cfssl gencert -ca=ca.pem -ca-key=ca-key.pem /configs/server.json | cfssljson -bare server
  cp *.pem /tmp
'

# Create the namespace that you'll place the webhook server in
kubectl create namespace admissions || true

# Create the secrets. If run again, replace the exising secrets.
kubectl -n admissions delete secret admissions-webhook-tls  || true
kubectl -n admissions create secret generic admissions-webhook-tls --from-file=certs/server.pem --from-file=certs/server-key.pem --from-file=certs/ca.pem

# Install the admissions webhook server.
helm upgrade --install admissions-webhook charts/admissions-webhook --namespace admissions

# Install the validating webhook configuration
export CA=`base64 -i certs/ca.pem`
envsubst < manifests/ValidatingWebhookConfiguration.yaml | kubectl apply -f -
