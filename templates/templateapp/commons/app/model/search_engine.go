package model

var SearchEnginePath = []string{"commons", "app", "model", "searchengine.go"}

var SearchEngineContent = `package model

import (
	"{{ .AppRepository }}/commons/log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (q *Query) SearchEngine(criteria map[string]string) {
	query, values := Criteria(q.Table, criteria, "AND")

	q.Db = q.Db.Where(query, values...)
}

func Criteria(table interface{}, criteria map[string]string, logic string) (string, []interface{}) {
	var queries []string
	var values []interface{}
	var query string
	var value interface{}
	var err error
	var hasInputValue = regexp.MustCompile(` + "`" + `\?` + "`" + `)

	for criterionKey, criterionValue := range criteria {
		query, value, err = handleCriterion(table, criterionKey, criterionValue)
		if err == nil {
			queries = append(queries, query)

			if hasInputValue.MatchString(query) {
				values = append(values, value)
			}

		} else {
			log.Error.Println("SearchEngine Query", err)
		}
	}

	query = strings.Join(queries, " "+logic+" ")

	return query, values
}

func handleCriterion(table interface{}, key string, value string) (string, interface{}, error) {
	var column, columnType, query string
	var interfaceValue interface{}
	var regexpInclusionQuery = regexp.MustCompile(` + "`" + `IN\s\(\?\)` + "`" + `)
	var regexpIsNullQuery = regexp.MustCompile(` + "`" + `IS\s(NOT\s){0,1}NULL` + "`" + `)
	var err error

	names := strings.Split(key, "_")

	column, query = strings.Join(names[:len(names)-1], "_"), names[len(names)-1]

	columnType, err = GetColumnType(table, column)
	if err != nil {
		log.Error.Println("SearchEngine handleCriterion", err)
		return "", "", err
	}

	query, value = translateQuery(column, query, value)

	if regexpIsNullQuery.MatchString(query) {
		interfaceValue = 1
	} else {
		interfaceValue, err = valueToInterface(columnType, value, regexpInclusionQuery.MatchString(query))
		if err != nil {
			log.Error.Println("SearchEngine handleCriterion", err)
			return "", "", err
		}
	}

	return query, interfaceValue, nil
}

func valueToInterface(columnType string, valueContent string, isQueryInclusion bool) (interface{}, error) {
	var regexpBooleanType = regexp.MustCompile(` + "`" + `bool` + "`" + `)
	var regexpIntType = regexp.MustCompile(` + "`" + `int` + "`" + `)
	var regexpFloatType = regexp.MustCompile(` + "`" + `float|double` + "`" + `)
	var regexpDateTimeType = regexp.MustCompile(` + "`" + `time|date` + "`" + `)
	var returnValue interface{}
	var err error

	if regexpIntType.MatchString(columnType) {
		returnValue, err = convertoToInt(valueContent, isQueryInclusion)
		if err != nil {
			return "", err
		}
	} else if regexpFloatType.MatchString(columnType) {
		returnValue, err = convertToFloat(valueContent, isQueryInclusion)
		if err != nil {
			return "", err
		}
	} else if regexpDateTimeType.MatchString(columnType) {
		returnValue, err = time.Parse(time.RFC3339, valueContent)
		if err != nil {
			return "", err
		}
	} else if regexpBooleanType.MatchString(columnType) {
		returnValue, err = convertToBoolean(valueContent)
		if err != nil {
			return "", err
		}
	} else if isQueryInclusion {
		returnValue = regexpSplit(valueContent, ` + "`" + `\s*,\s*` + "`" + `)
		err = nil
	} else {
		returnValue = valueContent
		err = nil
	}

	return returnValue, err
}

func translateQuery(column string, query string, value string) (string, string) {
	switch query {
	case "cont":
		query = "ILIKE ?"
		value = "%" + value + "%"
	case "notcont":
		query = "NOT ILIKE ?"
		value = "%" + value + "%"
	case "start":
		query = "ILIKE ?"
		value = value + "%"
	case "notstart":
		query = "NOT ILIKE ?"
		value = value + "%"
	case "end":
		query = "ILIKE ?"
		value = "%" + value
	case "notend":
		query = "NOT ILIKE ?"
		value = "%" + value
	case "gt":
		query = "> ?"
	case "lt":
		query = "< ?"
	case "gteq":
		query = ">= ?"
	case "lteq":
		query = "<= ?"
	case "null":
		query = "IS NULL"
	case "isnull":
		query = "IS NULL"
	case "notnull":
		query = "IS NOT NULL"
	case "in":
		query = "IN (?)"
	case "notin":
		query = "NOT IN (?)"
	case "true":
		query = "= 't'"
	case "false":
		query = "= 'f'"
	case "noteq":
		query = "<> ?"
	default:
		query = "= ?"
	}

	return column + " " + query, value
}

func convertoToInt(valueContent string, isQueryInclusion bool) (interface{}, error) {
	var newValues []int64

	newValue, err := strconv.ParseInt(valueContent, 10, 64)

	if err != nil {
		if !isQueryInclusion {
			return 0, err
		}

		newValues, err = splitStringToIntArray(valueContent)
		if err != nil {
			return 0, err
		} else {
			return newValues, nil
		}
	} else {
		return newValue, nil
	}
}

func convertToFloat(valueContent string, isQueryInclusion bool) (interface{}, error) {
	var newValues []float64

	newValue, err := strconv.ParseFloat(valueContent, 64)

	if err != nil {
		if !isQueryInclusion {
			return 0, err
		}

		newValues, err = splitStringToFloatArray(valueContent)
		if err != nil {
			return 0, err
		} else {
			return newValues, nil
		}
	} else {
		return newValue, nil
	}
}

func convertToStringIn(valueContent string) (interface{}, error) {
	stringValues := regexpSplit(valueContent, ` + "`" + `\s*,\s*` + "`" + `)
	return stringValues, nil
}

func convertToBoolean(valueContent string) (interface{}, error) {
	var checkFalse = regexp.MustCompile(` + "`" + `\A(0|f|false|no)\z` + "`" + `)

	valueContent = strings.TrimSpace(valueContent)

	if valueContent != "" {
		if checkFalse.MatchString(valueContent) {
			return false, nil
		} else {
			return true, nil
		}
	} else {
		return false, nil
	}
}

func splitStringToIntArray(value string) ([]int64, error) {
	var intValues []int64
	var intValue int64
	var i int
	var err error

	stringValues := regexpSplit(value, ` + "`" + `\s*,\s*` + "`" + `)

	for i = 0; i < len(stringValues); i++ {
		intValue, err = strconv.ParseInt(stringValues[i], 10, 64)
		if err != nil {
			intValues = nil
			return intValues, err
		}

		intValues = append(intValues, intValue)
	}

	return intValues, nil
}

func splitStringToFloatArray(value string) ([]float64, error) {
	var floatValues []float64
	var floatValue float64
	var i int
	var err error

	stringValues := regexpSplit(value, ` + "`" + `\s*,\s*` + "`" + `)

	for i = 0; i < len(stringValues); i++ {
		floatValue, err = strconv.ParseFloat(stringValues[i], 64)
		if err != nil {
			floatValues = nil
			return floatValues, err
		}

		floatValues = append(floatValues, floatValue)
	}

	return floatValues, nil
}

func regexpSplit(text string, delimeter string) []string {
	reg := regexp.MustCompile(delimeter)
	indexes := reg.FindAllStringIndex(text, -1)
	lastStart := 0
	result := make([]string, len(indexes)+1)
	for i, element := range indexes {
		result[i] = text[lastStart:element[0]]
		lastStart = element[1]
	}
	result[len(indexes)] = text[lastStart:len(text)]
	return result
}`
