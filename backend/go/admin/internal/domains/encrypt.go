package domains

const KeyGlobalEncryptPublicKey = "global:encrypt:public:key"

type EncryptKeyPair struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}
