# generate key pair
ssh-keygen -m PEM -b 2048 -t rsa
# original public key format could not parsed by golang, format it
ssh-keygen -e -m PEM -f rsa.pem > rsa.pem.pub