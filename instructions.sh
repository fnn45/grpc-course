#!/usr/bin/env bash

# Output files
# ca.key Certificate Authority private key file (this shouldnt be shared in real life)
# ca.crt   Certificate Authority trust certificate (this should be shared in real life)
# server.key: Server private key, password protected (this shouldnt be shared)
# server.csr: Server Certificate signing request (this should be shared with the ca owner)
# server.crt Server Certificate signed by the CA (thid would be sent back by the CA owner)
# server.pem Conversion of server.key into a format GRPC likes (this shouldnt be shared)

# Summary
# Private files: ca.key, server.key, server.pem, server crt
# Share files: ca.crt (needed by client), server.csr (needed by the CA)

#Changes these CN's to match hosts in tour enviroment if needed.
SERVER_CN=localhost

#Step 1 Generate Certificate authority + Trust Certificate (ca.crt)
openssl genrsa -passout pass:1111 -des3 -out ca.key 4096
openssl req -passin pass:1111 -new -x509 -days 365 -key ca.key -out ca.crt -subj "/CN=${SERVER_CN}"

#Step 2: Generate the Server private Key (server.key)
openssl genrsa -passout pass:1111 -des3 -out server.key 4096

#Step 3 Get a certificate signing request from the CA (server.csr)
openssl req -passin pass:1111 -new -key server.key -out server.csr -subj "/CN=${SERVER_CN}"

# Step 4: Sign the certificate with the CA we created ( its called self signing) - server.crt
openssl x509 -req -passin pass:1111 -days 365 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt

# Step 5 Convert the server certificate to .pem format (server.pem) - usable by gRPC
openssl pkcs8 -topk8 -nocrypt -passin pass:1111 -in server.key -out server.pem
