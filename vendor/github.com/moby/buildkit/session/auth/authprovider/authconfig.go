package authprovider

type AuthTLSConfig struct {
	RootCAs  []string
	KeyPairs []TLSKeyPair
}

type TLSKeyPair struct {
	Key         string
	Certificate string
}
