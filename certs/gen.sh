rm *.key *.crt *.req *.srl
# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca.key -out ca.crt -subj "/C=ID/ST=Jakarta/L=Jakarta/O=BP Tapera/OU=TI/CN=*.tapera.go.id/emailAddress=naray.citra@tapera.go.id"

echo "CA's self-signed certificate"
openssl x509 -in ca.crt -noout -text
# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout server.key -out server.req -subj "/C=ID/ST=Jakarta/L=Jakarta Selatan/O=BP Tapera/OU=TI/CN=*.int.tapera.go.id/emailAddress=naray.citra@tapera.go.id"

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in server.req -days 365 -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -extfile server-ext.cnf

# 4. Generate client's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout client.key -out client.req -subj "/C=ID/ST=Jakarta/L=Jakarta Barat/O=BP Tapera Client/OU=TI/CN=*.client.tapera.go.id/emailAddress=naray.citra@tapera.go.id"

# 5. Use CA's private key to sign client's CSR and get back the signed certificate
openssl x509 -req -in client.req -days 365 -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -extfile client-ext.cnf

openssl verify -CAfile ca.crt server.crt

openssl verify -CAfile ca.crt client.crt
