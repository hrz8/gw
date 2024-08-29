
```bash
# generate private key
openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048
# generate public key
openssl rsa -pubout -in private_key.pem -out public_key.pem
# generate modulus from public key
openssl rsa -in public_key.pem -pubin -modulus -noout
# generate n value
echo -n "MODULUS" | xxd -r -p | openssl base64 -A | tr '+/' '-_' | tr -d '='
```
