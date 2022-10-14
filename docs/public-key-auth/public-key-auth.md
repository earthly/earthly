# Earthly Public Key Authentication

Earthly provides public-key based authentication (in addition to traditional username and password authentication). This guide details how it works.

## What is public-key authentication

Public key authentication provides greater security compared to password authentication; it is achieved by using [asymmetric cryptography](https://en.wikipedia.org/wiki/Public-key_cryptography).

A user generates a pair of private and public keys; the public key is publicly distributed to anyone who wishes to send an encrypted message that only the holder of the private key can decrypt.
It is important that you **never share your private key**, otherwise anyone could use it to access data that is only intended for you.

Similarly, it is possible to sign data using your private key -- any user who has your public key can use it to verify the message was signed by you (or anyone who has access to your private key).

For these reasons, it is crucial that your private key remains private -- as a result, **earthly will never store, or transmit your private key**.

## How does earthly implement public-key authentication

Earthly accounts can be associated with any number of public keys (both `ssh-rsa`, and `ssh-ed25519` public keys are supported). These public keys are stored on the earthly server, in a database
that mimics the `~/.ssh/authorized_keys` file one typically finds on a server.

The client first connects to the earthly server over a https connection; the client responds with a [cryptographically-secure random](https://en.wikipedia.org/wiki/Cryptographically_secure_pseudorandom_number_generator) blob of data.
The client then passes that blob of data to the [ssh-agent](https://en.wikipedia.org/wiki/Ssh-agent) process, which must be running on your local host. This connection occurs by using the local unix-socket as set by the `SSH_AUTH_SOCK` environment variable. The ssh-agent signs the blob of data, and returns the signature -- **earthly will never read your private keys directly**.

This signature is sent to the earthly server; if the signature can be verified using a registered public key, then the server responds with a [JSON Web Token (JWT)](https://en.wikipedia.org/wiki/JSON_Web_Token) which is used
for the duration of your session.
