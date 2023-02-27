package utils

import (
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func MergeMaps(maps ...map[string]string) map[string]interface{} {
	merged := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			if v != "" {
				merged[k] = v
			}
		}
	}
	return merged
}

func MakeMap(regex *regexp.Regexp, subString string) map[string]string {
	madeMap := make(map[string]string)
	xMatches := regex.FindStringSubmatch(subString)
	for i, name := range regex.SubexpNames() {
		if i != 0 && name != "" && !(xMatches[i] == "") {
			madeMap[name] = xMatches[i]
		}
	}
	return madeMap
}
