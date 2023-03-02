package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var StatusCodes = map[string]int{
	"SUCCESSFUL": 0,

	"REQUEST_BROKEN":           100,
	"REQUEST_MISSING_SCHEME":   111,
	"REQUEST_MISSING_DOMAIN":   112,
	"REQUEST_MISSING_NONCE":    113,
	"REQUEST_MALFORMED_SCHEME": 121,
	"REQUEST_MALFORMED_DOMAIN": 122,
	"REQUEST_INVALID_DOMAIN":   131,
	"REQUEST_INVALID_NONCE":    132,
	"REQUEST_ALTERED":          141,
	"REQUEST_EXPIRED":          142,
	"REQUEST_CONSUMED":         143,

	"RESPONSE_BROKEN":              200,
	"RESPONSE_MISSING_REQUEST":     211,
	"RESPONSE_MISSING_ADDRESS":     212,
	"RESPONSE_MISSING_SIGNATURE":   213,
	"RESPONSE_MISSING_METADATA":    214,
	"RESPONSE_MALFORMED_ADDRESS":   221,
	"RESPONSE_MALFORMED_SIGNATURE": 222,
	"RESPONSE_MALFORMED_METADATA":  223,
	"RESPONSE_INVALID_METHOD":      231,
	"RESPONSE_INVALID_ADDRESS":     232,
	"RESPONSE_INVALID_SIGNATURE":   233,
	"RESPONSE_INVALID_METADATA":    234,

	"SERVICE_BROKEN":                 300,
	"SERVICE_ADDRESS_DENIED":         311,
	"SERVICE_ADDRESS_REVOKED":        312,
	"SERVICE_ACTION_DENIED":          321,
	"SERVICE_ACTION_UNAVAILABLE":     322,
	"SERVICE_ACTION_NOT_IMPLEMENTED": 323,
	"SERVICE_INTERNAL_ERROR":         331,
}

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

func InvalidateRequest(c *gin.Context, statusCodeName string, statusMessage string) {
	statusCode, exists := StatusCodes[statusCodeName]
	if !exists {
		statusCode = 100 // Default status code for invalid status code names
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": gin.H{
			"code":    statusCode,
			"message": statusMessage,
		},
	})
}
