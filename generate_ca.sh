#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

# CREATE THE PRIVATE KEY FOR OUR CUSTOM CA
openssl genrsa -out certs/ca.key 2048

# GENERATE A CA CERT WITH THE PRIVATE KEY
openssl req -new -x509 -key certs/ca.key -out certs/ca.crt -config certs/ca_config.txt

# CREATE THE PRIVATE KEY FOR OUR gandalf SERVER
openssl genrsa -out certs/gandalf-key.pem 2048

# CREATE A CSR FROM THE CONFIGURATION FILE AND OUR PRIVATE KEY
openssl req -new -key certs/gandalf-key.pem -subj "/CN=gandalf.default.svc" -out gandalf.csr -config certs/gandalf_config.txt

# CREATE THE CERT SIGNING THE CSR WITH THE CA CREATED BEFORE
openssl x509 -req -in gandalf.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/gandalf-crt.pem

# INJECT CA IN THE WEBHOOK CONFIGURATION
export CA_BUNDLE=$(cat certs/ca.crt | base64 | tr -d '\n')
cat _manifest_.yaml | envsubst > manifest.yaml
