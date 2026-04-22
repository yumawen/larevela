package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	SolanaRpcURL string `json:"solanaRpcURL,optional"`
}
