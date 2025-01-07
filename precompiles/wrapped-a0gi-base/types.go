package wrappeda0gibase

import "math/big"

type Supply = struct {
	Cap   *big.Int "json:\"cap\""
	Total *big.Int "json:\"total\""
}
