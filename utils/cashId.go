package utils

import (
	"BCHat/database/models"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

var regexpPatterns = map[string]string{
	"request":    `(?P<scheme>cashid:)(?:[\/]{2})?(?P<domain>[^\/]+)(?P<path>\/[^\?]+)(?P<parameters>\?.+)`,
	"parameters": `(?:(?:[\?\&]a=)(?P<action>[^\&]+))?(?:(?:[\?\&]d=)(?P<data>[^\&]+))?(?:(?:[\?\&]r=)(?P<required>[^\&]+))?(?:(?:[\?\&]o=)(?P<optional>[^\&]+))?(?:(?:[\?\&]x=)(?P<nonce>[^\&]+))?`,
	"metadata":   `(i(?P<identification>[^1-9]*)(?P<name>1)?(?P<family>2)?(?P<nickname>3)?(?P<age>4)?(?P<gender>5)?(?P<birthdate>6)?(?P<picture>8)?(?P<national>9)?)?(p(?P<position>[^1-9]*)(?P<country>1)?(?P<state>2)?(?P<city>3)?(?P<streetname>4)?(?P<streetnumber>5)?(?P<residence>6)?(?P<coordinate>9)?)?(c(?P<contact>[^1-9]*)(?P<email>1)?(?P<instant>2)?(?P<social>3)?(?P<mobilephone>4)?(?P<homephone>5)?(?P<workphone>6)?(?P<postlabel>9)?)?`,
}

type metadataField struct {
	Name string
	Code int
}

var MetadataTypes = map[string][]metadataField{
	"identity": {
		{Name: "name", Code: 1},
		{Name: "family", Code: 2},
		{Name: "nickname", Code: 3},
		{Name: "age", Code: 4},
		{Name: "gender", Code: 5},
		{Name: "birthdate", Code: 6},
		{Name: "picture", Code: 8},
		{Name: "national", Code: 9},
	},
	"position": {
		{Name: "country", Code: 1},
		{Name: "state", Code: 2},
		{Name: "city", Code: 3},
		{Name: "streetname", Code: 4},
		{Name: "streetnumber", Code: 5},
		{Name: "residence", Code: 6},
		{Name: "coordinates", Code: 9},
	},
	"contact": {
		{Name: "email", Code: 1},
		{Name: "instant", Code: 2},
		{Name: "social", Code: 3},
		{Name: "phone", Code: 4},
		{Name: "postal", Code: 5},
	},
}

func CreateRequest(action, data string, metadata map[string][]string) string {
	// generate a random nonce.
	nonce := rand.Intn(900000000) + 100000000

	// Check if the nonce is already used, and regenerate until it does not exist.
	for ApcuExists(fmt.Sprintf("cashid_nonce_%d", nonce)) {
		// generate a random nonce.
		nonce = rand.Intn(900000000) + 100000000
	}

	// Initialize an empty parameter list.
	parameters := make(map[string]string)

	// If a specific action was requested, add it to the parameter list.
	if action != "" {
		parameters["a"] = fmt.Sprintf("a=%s", action)
	}

	// If specific data was requested, add it to the parameter list.
	if data != "" {
		parameters["d"] = fmt.Sprintf("d=%s", data)
	}

	encodeMetadata(metadata, parameters, "required")
	encodeMetadata(metadata, parameters, "optional")

	// Append the nonce to the parameter list.
	parameters["x"] = fmt.Sprintf("x=%d", nonce)

	// Form the request URI from the configured values.
	var requestURI string
	requestURI = "cashid:bchat.xyz?/cashid?"
	for _, v := range parameters {
		requestURI += fmt.Sprintf("%s&", v)
	}
	requestURI = requestURI[:len(requestURI)-1]

	// Store the request and nonce in local cache.
	entry := &CacheEntry{available: true, expirationTime: time.Now().Add(60 * 15)}
	cache.Store(fmt.Sprintf("cashid_request_%d", nonce), entry)
	// Return the request URI to indicate success.
	return requestURI
}

func encodeMetadata(metadata map[string][]string, parameters map[string]string, metadataType string) {
	typeChar := metadataType[0:1]
	if metadataOpt, ok := metadata[metadataType]; ok {
		parameters[typeChar] = fmt.Sprintf("%s=", typeChar)
		// Convert the metadata from []string to map[string]map[string]string.
		for _, field := range metadataOpt {
			for mdType, metadataFields := range MetadataTypes {
				typeCode := mdType[0:1]
				for _, mdField := range metadataFields {
					if field == mdField.Name {
						if !(strings.Contains(parameters[typeChar], typeCode)) {
							parameters[typeChar] += fmt.Sprintf("%s", typeCode)
						}
						parameters[typeChar] += fmt.Sprintf("%d", mdField.Code)
					}
				}
			}
		}
	}
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
