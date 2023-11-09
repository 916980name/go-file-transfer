package common

type contextKey int

const (
	VIPER_RSA_PRI    = "rsa.privateKey"
	VIPER_RSA_PUB    = "rsa.publicKey"
	VIPER_AES_KEY    = "aes.key"
	VIPER_AES_IV     = "aes.iv"
	VIPER_HOST_URL   = "host-url"
	LOGIN_SHARE_PATH = "ls"

	CTX_USER_KEY contextKey = iota
	// ...
)
