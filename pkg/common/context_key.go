package common

import "fmt"

type contextKey int
type ShareKey int
type ShareExpireTypeKey int

const (
	VIPER_RSA_PRI      = "rsa.privateKey"
	VIPER_RSA_PUB      = "rsa.publicKey"
	VIPER_AES_KEY      = "aes.key"
	VIPER_AES_IV       = "aes.iv"
	VIPER_HOST_URL     = "host-url"
	LOGIN_SHARE_PATH   = "ls"
	MESSAGE_SHARE_PATH = "ms"

	CTX_USER_KEY contextKey = iota
	// ...
)
const (
	SHARE_TYPE_LOGIN ShareKey = iota
	SHARE_TYPE_MESSAGE
	SHARE_TYPE_FILE
)

const (
	SHARE_EXPIRE_TYPE_TIMES ShareExpireTypeKey = iota
	SHARE_EXPIRE_TYPE_DURATION
)

func init() {
	if FLAG_DEBUG {
		fmt.Printf("CTX_USER_KEY %d\n", CTX_USER_KEY)
		fmt.Printf("SHARE_TYPE_LOGIN %d\n", SHARE_TYPE_LOGIN)
		fmt.Printf("SHARE_EXPIRE_TYPE_TIMES %d\n", SHARE_EXPIRE_TYPE_TIMES)
	}
}
