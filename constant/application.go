package constant

import (
	"time"
)

const (
	ApplicationName         = "swt-pro"
	PrivateKeyPath          = "./key/jwtRS256.key"
	PublicKeyPath           = "./key/jwtRS256.key.pub"
	LoginExpirationDuration = time.Duration(24) * time.Hour
)
