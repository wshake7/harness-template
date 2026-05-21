package file

import (
	"github.com/bytedance/sonic"
	"go-common/utils/catch"
	"testing"
)

func TestReadDir(t *testing.T) {
	dir, err := ScanDir("./", &ReadDirOption{
		ResourceType: AllResource,
		FileRegExp:   "",
		DirRegExp:    "",
	})
	if err != nil {
		t.Error(err)
	}
	t.Logf("fileInfos: %s", catch.Try1(sonic.MarshalString(dir)))
}
