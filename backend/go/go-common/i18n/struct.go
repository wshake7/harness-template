package i18n

type Language string
type Code string

type Struct[T any] struct {
	languages       map[Language]T
	defaultLanguage Language
}

// NewStruct 创建新的 Struct 实例
func NewStruct[T any](defaultLanguage Language) *Struct[T] {
	return &Struct[T]{
		languages:       make(map[Language]T),
		defaultLanguage: defaultLanguage,
	}
}

func (s *Struct[T]) GetDefaultLanguage() Language {
	return s.defaultLanguage
}

func (s *Struct[T]) SetDefaultLanguage(language Language) {
	s.defaultLanguage = language
}

func (s *Struct[T]) SetLanguage(language Language, data T) {
	s.languages[language] = data
}

func (s *Struct[T]) DeleteLanguage(language Language) {
	delete(s.languages, language)
}

func (s *Struct[T]) LoadLanguage(language Language) (T, bool) {
	l, ok := s.languages[language]
	return l, ok
}

func (s *Struct[T]) LoadLanguageOrDefault(language Language) (T, bool) {
	l, ok := s.languages[language]
	if ok {
		return l, ok
	}
	l, ok = s.languages[s.defaultLanguage]
	return l, ok
}

func (s *Struct[T]) HasLanguage(language Language) bool {
	_, ok := s.languages[language]
	return ok
}

func (s *Struct[T]) GetAllLanguages() []Language {
	languages := make([]Language, 0, len(s.languages))
	for lang := range s.languages {
		languages = append(languages, lang)
	}
	return languages
}
