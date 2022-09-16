package helper

import "regexp"

var (
	Regexp      regexper
	regPhone    = regexp.MustCompile(`^1[3-9]\d{9}$`)
	regExpEmail = regexp.MustCompile(`^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+`)
)

type regexper struct{}

func (regexper) IsEmail(email string) bool {
	return regExpEmail.MatchString(email)
}

func (regexper) IsPhone(phone string) bool {
	return regPhone.MatchString(phone)
}
