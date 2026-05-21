package file

import (
	"os"
	"path/filepath"
	"regexp"
)

type ResourceType uint8

const (
	AllResource = ResourceType(iota)
	FileResource
	DirResource
)

type RegExp string

const (
	AllRegExp    = RegExp("")
	ConfRegExp   = RegExp("^.*\\.conf$")
	JsonRegExp   = RegExp("^.*\\.json$")
	YamlRegExp   = RegExp("^.*\\.(yaml|yml)$")
	LogRegExp    = RegExp("^.*\\.log$")
	HiddenRegExp = RegExp("^\\..*") // 隐藏文件/目录
)

type ReadDirOption struct {
	ResourceType ResourceType
	FileRegExp   RegExp
	DirRegExp    RegExp
}

func DefaultReadDirOption() *ReadDirOption {
	return &ReadDirOption{
		ResourceType: AllResource,
		FileRegExp:   AllRegExp,
		DirRegExp:    AllRegExp,
	}
}

type FileDetail struct {
	Name    string
	Path    string
	IsDir   bool
	Size    int64
	ModTime int64
	Ext     string
}

func ScanDir(path string, option *ReadDirOption) ([]FileDetail, error) {
	// 如果 option 为 nil，使用默认配置
	if option == nil {
		option = DefaultReadDirOption()
	}

	// 转换为绝对路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	// 读取目录内容
	entries, err := os.ReadDir(absPath)
	if err != nil {
		return nil, err
	}

	// 编译文件正则表达式
	var fileRegex *regexp.Regexp
	if option.FileRegExp != AllRegExp {
		fileRegex, err = regexp.Compile(string(option.FileRegExp))
		if err != nil {
			return nil, err
		}
	}

	// 编译目录正则表达式
	var dirRegex *regexp.Regexp
	if option.DirRegExp != AllRegExp {
		dirRegex, err = regexp.Compile(string(option.DirRegExp))
		if err != nil {
			return nil, err
		}
	}

	var result []FileDetail
	for _, entry := range entries {
		isDir := entry.IsDir()
		name := entry.Name()

		// 根据 ResourceType 过滤
		switch option.ResourceType {
		case FileResource:
			if isDir {
				continue
			}
			// 文件正则匹配
			if fileRegex != nil && !fileRegex.MatchString(name) {
				continue
			}
		case DirResource:
			if !isDir {
				continue
			}
			// 目录正则匹配
			if dirRegex != nil && !dirRegex.MatchString(name) {
				continue
			}
		case AllResource:
			// 分别应用对应的正则
			if isDir {
				if dirRegex != nil && !dirRegex.MatchString(name) {
					continue
				}
			} else {
				if fileRegex != nil && !fileRegex.MatchString(name) {
					continue
				}
			}
		}

		// 获取详细信息
		info, e := entry.Info()
		if e != nil {
			continue // 跳过无法获取信息的项
		}

		result = append(result, FileDetail{
			Name:    name,
			Path:    filepath.Join(absPath, name),
			IsDir:   isDir,
			Size:    info.Size(),
			ModTime: info.ModTime().Unix(),
		})
	}

	return result, nil
}

func ExistFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func ExistDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func Exist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func Load(filename string) (string, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
