package viperc

import (
	"fmt"
	"github.com/creasty/defaults"
	"github.com/fsnotify/fsnotify"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
	"path"
	"path/filepath"
	"reflect"
	"strings"
)

func watchConfig(obj any, reloads ...func()) {
	viper.WatchConfig()

	// Note: OnConfigChange is called twice on Windows
	viper.OnConfigChange(func(e fsnotify.Event) {
		t := reflect.TypeOf(obj).Elem()
		v := reflect.New(t)
		reflect.ValueOf(obj).Elem().Set(v.Elem()) // reset object

		err := viper.Unmarshal(obj)
		if err != nil {
			fmt.Println("viper.Unmarshal error: ", err)
		} else {
			for _, reload := range reloads {
				reload()
			}
		}
	})
}

func ParseFile(filePath string, obj any, reloads ...func()) (*viper.Viper, error) {
	viper0 := viper.New()
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return nil, fmt.Errorf("obj must be a non-nil pointer")
	}

	confFileAbs, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}

	filePathStr, filename := filepath.Split(confFileAbs)
	ext := strings.TrimLeft(path.Ext(filename), ".")
	if ext == "" {
		return nil, fmt.Errorf("file %s suffix is empty", ext)
	}
	filename = strings.ReplaceAll(filename, "."+ext, "")

	viper0.AddConfigPath(filePathStr) // path
	viper0.SetConfigName(filename)    // file name
	viper0.SetConfigType(ext)         // get the configuration type from the file name
	if err = viper0.ReadInConfig(); err != nil {
		return nil, err
	}
	if err = defaults.Set(obj); err != nil {
		return nil, err
	}

	err = viper0.UnmarshalExact(obj, func(dc *mapstructure.DecoderConfig) {
		dc.ZeroFields = false // 不清零已有字段（保留默认值）
	})
	if err != nil {
		return nil, err
	}
	if len(reloads) > 0 {
		watchConfig(obj, reloads...)
	}
	return viper0, nil
}

func ParseGlobalFile(filePath string, obj any, reloads ...func()) error {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return fmt.Errorf("obj must be a non-nil pointer")
	}

	confFileAbs, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	filePathStr, filename := filepath.Split(confFileAbs)
	ext := strings.TrimLeft(path.Ext(filename), ".")
	if ext == "" {
		return fmt.Errorf("file %s suffix is empty", ext)
	}
	filename = strings.ReplaceAll(filename, "."+ext, "")

	viper.AddConfigPath(filePathStr) // path
	viper.SetConfigName(filename)    // file name
	viper.SetConfigType(ext)         // get the configuration type from the file name
	if err = viper.ReadInConfig(); err != nil {
		return err
	}
	if err = defaults.Set(obj); err != nil {
		return err
	}

	err = viper.UnmarshalExact(obj, func(dc *mapstructure.DecoderConfig) {
		dc.ZeroFields = false // 不清零已有字段（保留默认值）
	})
	if len(reloads) > 0 {
		watchConfig(obj, reloads...)
	}
	return nil
}
