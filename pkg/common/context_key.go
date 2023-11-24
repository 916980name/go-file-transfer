package common

import "fmt"

type contextKey int
type ShareKey int

const (
	VIPER_RSA_PRI    = "rsa.privateKey"
	VIPER_RSA_PUB    = "rsa.publicKey"
	VIPER_AES_KEY    = "aes.key"
	VIPER_AES_IV     = "aes.iv"
	VIPER_HOST_URL   = "host-url"
	LOGIN_SHARE_PATH = "ls"

	CTX_USER_KEY contextKey = iota
	// ...

	SHARE_TYPE_LOGIN ShareKey = iota
	SHARE_TYPE_MESSAGE
	SHARE_TYPE_FILE
)

func init() {
	if FLAG_DEBUG {
		fmt.Println(CTX_USER_KEY)
		fmt.Println(SHARE_TYPE_LOGIN)
		fmt.Println(SHARE_TYPE_MESSAGE)
		fmt.Println(SHARE_TYPE_FILE)
	}
}
