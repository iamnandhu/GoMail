package utils

import (
	"errors"
	"reflect"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	// Regular expression to detect NoSQL injection patterns
	noSQLInjectPattern = regexp.MustCompile(`(?i)\$|\{|\}|;|\.\.|\\|\/\*|\*\/|\\'|\\"|==|!=|>=|<=|=>|=<|&&|\|\||>>|<<|\$(?:[a-z_]+)`)

	// Common MongoDB operators that are allowed after validation
	allowedOperators = map[string]bool{
		"$eq":  true,
		"$gt":  true,
		"$gte": true,
		"$in":  true,
		"$lt":  true,
		"$lte": true,
		"$ne":  true,
		"$nin": true,
	}

	// MongoDB query errors
	ErrInvalidFilter      = errors.New("invalid filter format")
	ErrPotentialInjection = errors.New("potential NoSQL injection detected")
	ErrInvalidOperator    = errors.New("invalid MongoDB operator")
)

// SanitizeMongoFilter sanitizes a MongoDB filter to prevent NoSQL injection
// It accepts interface{} and returns a sanitized bson.M or an error
func SanitizeMongoFilter(filter interface{}) (bson.M, error) {
	// If nil, return empty filter
	if filter == nil {
		return bson.M{}, nil
	}

	// Handle common filter types
	switch v := filter.(type) {
	case primitive.ObjectID:
		// ObjectID is safe
		return bson.M{"_id": v}, nil
	case string:
		// Check if it's a valid ObjectID
		if len(v) == 24 && noSQLInjectPattern.FindString(v) == "" {
			objID, err := primitive.ObjectIDFromHex(v)
			if err == nil {
				return bson.M{"_id": objID}, nil
			}
		}
		// Check for potential injection in string
		if noSQLInjectPattern.FindString(v) != "" {
			return nil, ErrPotentialInjection
		}
		// Otherwise, treat as a simple string filter (not ideal, better to specify field)
		return bson.M{}, nil
	case bson.M:
		// Recursively check each field in the bson.M
		return sanitizeBsonM(v)
	case map[string]interface{}:
		// Convert to bson.M and sanitize
		return sanitizeMap(v)
	case map[string]string:
		// Convert string map to interface map and sanitize
		interfaceMap := make(map[string]interface{})
		for k, v := range v {
			interfaceMap[k] = v
		}
		return sanitizeMap(interfaceMap)
	default:
		// Handle unexpected filter types
		val := reflect.ValueOf(filter)
		if val.Kind() == reflect.Struct {
			// Convert struct to map and sanitize
			data, err := bson.Marshal(filter)
			if err != nil {
				return nil, ErrInvalidFilter
			}
			var result bson.M
			if err := bson.Unmarshal(data, &result); err != nil {
				return nil, ErrInvalidFilter
			}
			return sanitizeBsonM(result)
		}
		// Unhandled filter type
		return nil, ErrInvalidFilter
	}
}

// sanitizeBsonM sanitizes a bson.M document
func sanitizeBsonM(doc bson.M) (bson.M, error) {
	result := bson.M{}

	for key, value := range doc {
		// Check for potential injection in keys
		if noSQLInjectPattern.FindString(key) != "" && !isAllowedOperator(key) {
			return nil, ErrPotentialInjection
		}

		// Check value based on type
		switch v := value.(type) {
		case primitive.ObjectID:
			// ObjectID is safe
			result[key] = v
		case string:
			// Check for injection patterns in string values
			if noSQLInjectPattern.FindString(v) != "" {
				return nil, ErrPotentialInjection
			}
			result[key] = v
		case bson.M:
			// Recursively sanitize nested document
			sanitized, err := sanitizeBsonM(v)
			if err != nil {
				return nil, err
			}
			result[key] = sanitized
		case map[string]interface{}:
			// Convert to bson.M and sanitize
			sanitized, err := sanitizeMap(v)
			if err != nil {
				return nil, err
			}
			result[key] = sanitized
		case []interface{}:
			// Sanitize array values
			sanitized, err := sanitizeArray(v)
			if err != nil {
				return nil, err
			}
			result[key] = sanitized
		default:
			// Primitive types (numbers, booleans, etc.) are generally safe
			result[key] = v
		}
	}

	return result, nil
}

// sanitizeMap converts and sanitizes a map[string]interface{}
func sanitizeMap(m map[string]interface{}) (bson.M, error) {
	bsonM := bson.M{}
	for k, v := range m {
		bsonM[k] = v
	}
	return sanitizeBsonM(bsonM)
}

// sanitizeArray sanitizes array values
func sanitizeArray(arr []interface{}) ([]interface{}, error) {
	result := make([]interface{}, 0, len(arr))
	
	for _, item := range arr {
		switch v := item.(type) {
		case string:
			// Check for injection patterns in string values
			if noSQLInjectPattern.FindString(v) != "" {
				return nil, ErrPotentialInjection
			}
			result = append(result, v)
		case bson.M:
			// Recursively sanitize nested document
			sanitized, err := sanitizeBsonM(v)
			if err != nil {
				return nil, err
			}
			result = append(result, sanitized)
		case map[string]interface{}:
			// Convert to bson.M and sanitize
			sanitized, err := sanitizeMap(v)
			if err != nil {
				return nil, err
			}
			result = append(result, sanitized)
		case []interface{}:
			// Recursively sanitize nested array
			sanitized, err := sanitizeArray(v)
			if err != nil {
				return nil, err
			}
			result = append(result, sanitized)
		default:
			// Primitive types are safe
			result = append(result, v)
		}
	}
	
	return result, nil
}

// isAllowedOperator checks if a string is an allowed MongoDB operator
func isAllowedOperator(key string) bool {
	if strings.HasPrefix(key, "$") {
		return allowedOperators[key]
	}
	return false
} 