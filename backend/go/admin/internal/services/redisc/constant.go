package redisc

const (
	KeyGlobalEncryptPublicKey = "global:encrypt:public:key"
)

type DtoKeyPair struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}
