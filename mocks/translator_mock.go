package mocks

// TranslatorStub satisfies localization.TranslatorInstance.
// Returns the key as message for every Translate call — keeps controller tests
// focused on HTTP status codes and response shape, not i18n strings.
type TranslatorStub struct{}

func NewTranslatorStub() *TranslatorStub { return &TranslatorStub{} }

func (t *TranslatorStub) Translate(key string, params ...string) (string, error) {
	return key, nil
}

func (t *TranslatorStub) Locale() string { return "fa_IR" }
