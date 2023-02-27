package utils

import (
	"BCHChat/database/models"
	"regexp"
)

var MetadataTypes = map[string]map[string]int{
	"identity": {
		"name":      1,
		"family":    2,
		"nickname":  3,
		"age":       4,
		"gender":    5,
		"birthdate": 6,
		"picture":   8,
		"national":  9,
	},
	"position": {
		"country":      1,
		"state":        2,
		"city":         3,
		"streetname":   4,
		"streetnumber": 5,
		"residence":    6,
		"coordinates":  9,
	},
	"contact": {
		"email":   1,
		"instant": 2,
		"social":  3,
		"phone":   4,
		"postal":  5,
	},
}

var regexpPatterns = map[string]string{
	"request":    `(?P<scheme>cashid:)(?:[\/]{2})?(?P<domain>[^\/]+)(?P<path>\/[^\?]+)(?P<parameters>\?.+)`,
	"parameters": `(?:(?:[\?\&]a=)(?P<action>[^\&]+))?(?:(?:[\?\&]d=)(?P<data>[^\&]+))?(?:(?:[\?\&]r=)(?P<required>[^\&]+))?(?:(?:[\?\&]o=)(?P<optional>[^\&]+))?(?:(?:[\?\&]x=)(?P<nonce>[^\&]+))?`,
	"metadata":   `(i(?P<identification>[^1-9]*)(?P<name>1)?(?P<family>2)?(?P<nickname>3)?(?P<age>4)?(?P<gender>5)?(?P<birthdate>6)?(?P<picture>8)?(?P<national>9)?)?(p(?P<position>[^1-9]*)(?P<country>1)?(?P<state>2)?(?P<city>3)?(?P<streetname>4)?(?P<streetnumber>5)?(?P<residence>6)?(?P<coordinate>9)?)?(c(?P<contact>[^1-9]*)(?P<email>1)?(?P<instant>2)?(?P<social>3)?(?P<mobilephone>4)?(?P<homephone>5)?(?P<workphone>6)?(?P<postlabel>9)?)?`,
}

// ParseRequest parses a CashID request query string into a CashIDRequest struct.
func ParseRequest(requestUri string) map[string]interface{} {
	// Initialize empty map
	requestRegex := regexp.MustCompile(regexpPatterns["request"])
	paramsRegex := regexp.MustCompile(regexpPatterns["parameters"])
	metadataRegex := regexp.MustCompile(regexpPatterns["metadata"])
	merged := make(map[string]interface{})
	requestParts := MakeMap(requestRegex, requestUri)
	parametersMap := MakeMap(paramsRegex, requestParts["parameters"])
	requiredMap := MakeMap(metadataRegex, parametersMap["required"])
	optionalMap := MakeMap(metadataRegex, parametersMap["optional"])
	merged = MergeMaps(requestParts, parametersMap)
	merged["required"] = requiredMap
	merged["optional"] = optionalMap
	return merged
}

func InvalidateRequest(statusCode int, statusMessage string) models.StatusConfirmation {
	return models.StatusConfirmation{
		Status:  statusCode,
		Message: statusMessage,
	}
}
