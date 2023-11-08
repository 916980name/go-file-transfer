package common

type contextKey int

const (
	VIPER_HOST_URL   = "host-url"
	LOGIN_SHARE_PATH = "ls"

	CTX_USER_KEY contextKey = iota
	// ...
)
