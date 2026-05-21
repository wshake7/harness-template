package encrypt

import (
	"fmt"
	"go-common/utils/str"
	"sort"
	"strings"
	"time"
)

const (
	// RSA算法配置
	RSA_ALGORITHM = "RSA/ECB/OAEPWithSHA-256AndMGF1Padding"
	RSA_KEY_SIZE  = 2048

	// AES算法配置
	AES_ALGORITHM  = "AES/GCM/NoPadding"
	AES_KEY_SIZE   = 32
	GCM_IV_LENGTH  = 12
	GCM_TAG_LENGTH = 16

	// 防重放配置
	REQUEST_EXPIRE_TIME = 5 * time.Minute
	NONCE_EXPIRE_TIME   = REQUEST_EXPIRE_TIME * 2
)

var (
	NONCE_REDIS_KEY_PREFIX = str.ValueF{
		Value: "security:nonce:%s",
	}
)

func UriSort(m map[string]any, filterFn func(key string) bool) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		if filterFn(k) {
			if v := m[k]; v != "" {
				keys = append(keys, k)
			}
		}
	}
	sort.Strings(keys)
	var buf strings.Builder
	for i, k := range keys {
		if i > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(fmt.Sprint(m[k]))
	}
	return buf.String()
}
