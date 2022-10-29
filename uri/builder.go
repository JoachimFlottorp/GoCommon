package uri

import "strings"

type URIBuild interface {
	SetURL(string) URIBuild
	AddParam(string, string) URIBuild
	Build() string
}

type builder struct {
	url string

	params map[string]string
}

func NewURLBuilder() URIBuild {
	return &builder{
		url:    "",
		params: make(map[string]string),
	}
}

func (u *builder) SetURL(url string) URIBuild {
	u.url = url

	return u
}

func (u *builder) AddParam(key, value string) URIBuild {
	u.params[key] = value

	return u
}

func (u *builder) Build() string {
	var final = u.url

	final += "?"

	for key, value := range u.params {
		final += key + "=" + value + "&"
	}

	final = strings.TrimSuffix(final, "&")

	return final
}