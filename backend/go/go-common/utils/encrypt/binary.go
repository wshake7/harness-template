package encrypt

var binary = New([]byte("234sdfn234ksjdf"))

func Encrypt(bytes []byte) {
	j := 0
	for i := range bytes {
		bytes[i] = bytes[i] ^ binary.secret[j]
		j = (j + 1) % binary.length
	}
}

func SetSecret(secret0 []byte) {
	binary.SetSecret(secret0)
}

type Binary struct {
	secret []byte
	length int
}

func New(secret []byte) *Binary {
	return &Binary{
		secret: secret,
		length: len(secret),
	}
}

func (b *Binary) SetSecret(secret []byte) {
	b.secret = secret
	b.length = len(secret)
}
