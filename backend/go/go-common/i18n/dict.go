package i18n

import "maps"

import "fmt"

type Dict struct {
	*Struct[map[Code]string]
}

// NewDict 创建新的 Dict 实例
func NewDict(defaultLanguage Language) *Dict {
	return &Dict{
		Struct: NewStruct[map[Code]string](defaultLanguage),
	}
}

func (d *Dict) SetCode(language Language, code Code, value string) {
	l, ok := d.languages[language]
	if !ok {
		l = make(map[Code]string)
		d.languages[language] = l
	}
	l[code] = value
}

func (d *Dict) DeleteCode(language Language, code Code) {
	m, b := d.LoadLanguage(language)
	if b {
		delete(m, code)
	}
}

func (d *Dict) LoadOrDefault(language Language, code Code, defaultValue string, args ...any) string {
	m, b := d.LoadLanguage(language)
	if b {
		v, ok := m[code]
		if ok {
			return fmt.Sprintf(v, args...)
		}
	}
	return defaultValue
}

func (d *Dict) TryLoad(language Language, code Code, args ...any) (string, bool) {
	m, b := d.LoadLanguage(language)
	if b {
		return fmt.Sprintf(m[code], args...), true
	}
	return "", false
}

func (d *Dict) Load(language Language, code Code, args ...any) string {
	v, _ := d.TryLoad(language, code, args...)
	return v
}

func (d *Dict) MustLoad(language Language, code Code, args ...any) string {
	v, b := d.TryLoad(language, code, args...)
	if !b {
		panic("code not found")
	}
	return v
}

func (d *Dict) HasCode(language Language, code Code) bool {
	if m, ok := d.LoadLanguage(language); ok {
		_, exists := m[code]
		return exists
	}
	return false
}

func (d *Dict) GetAllCodes(language Language) []Code {
	if m, ok := d.LoadLanguage(language); ok {
		codes := make([]Code, 0, len(m))
		for code := range m {
			codes = append(codes, code)
		}
		return codes
	}
	return nil
}

func (d *Dict) BatchSetCodes(language Language, codes map[Code]string) {
	l, ok := d.languages[language]
	if !ok {
		l = make(map[Code]string, len(codes))
		d.languages[language] = l
	}
	maps.Copy(l, codes)
}
