package routes

var MiddlewarePath = []string{"routes", "middleware.go"}

var MiddlewareContent = `package routes

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"{{ .AppRepository }}/app/handlers"
	"{{ .AppRepository }}/app/user"
	"{{ .AppRepository }}/commons/app/handler"
	"{{ .AppRepository }}/commons/log"
	"{{ .AppRepository }}/db/entities"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info.Println(r.Method + ": " + filterUrlValues(r.URL.Path, r.URL.Query()) + " for " + r.RemoteAddr)

		if strings.Trim(r.Header.Get("Content-Type"), " \n") == "application/json" {
			body, err := ioutil.ReadAll(r.Body)
			if err == nil {
				log.Info.Println("Body JSON:", filterJsonValues(string(body)))
				// put the body content back
				r.Body = ioutil.NopCloser(strings.NewReader(string(body)))
			} else {
				log.Error.Println("loggingMiddlware: ", err)
			}
		} else {
			r.ParseMultipartForm(100 * 1024)
			log.Info.Println("Form-data: " + filterFormValues(r.Form))
		}

		next.ServeHTTP(w, r)
	})
}

func authorizeMiddleware(next http.Handler) http.Handler {
	var userID uint64
	var err error
	var userRole string
	var signedInUser entities.User

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID = 0
		err = nil
		userRole = "public"

		userID, err = checkToken(r.Header.Get("Authorization"))
		if err == nil {
			if signedInUser, err = checkSignedInUser(userID); err == nil {
				userRole = "signed_in"
			}

			if userRole == "signed_in" && checkAdminUser(user.Current) {
				userRole = "admin"
			}
		}

		if GrantPermission(r.RequestURI, r.Method, userRole) {
			next.ServeHTTP(w, r)
		} else if userRole == "public" {
			handler.Error401(w, r)
		} else {
			handler.Error403(w, r)
		}
	})
}

func checkAdminUser(signedInUser entities.User) bool {
	if signedInUser.Admin {
		log.Info.Println("admin access granted to user")
		return true
	} else {
		log.Info.Println("admin access denied to user")
		return false
	}
}

func checkSignedInUser(userID uint64) (entities.User, error) {
	log.Info.Printf("checking user id: %d...\n", userID)

  if userID == 0 {
		log.Info.Println("user is not available")
		return entities.User{}, errors.New("user is not available")
  } else if signedInUser, err := user.Find(userID); err == nil {
		user.SetCurrent(userID)
		log.Info.Println("user was found")
		return signedInUser, nil
	} else {
		log.Info.Println("user was not found")
		return signedInUser, errors.New("user was not found")
	}
}

func checkToken(token string) (uint64, error) {
	log.Info.Println("checking token...")
  
  if token == "" {
		log.Info.Println("token was not sent")
		return 0, nil
  } else {
    return validateToken(token)
  }
}

func validateToken(token string) (uint64, error) {
	log.Info.Println("validating token...")
  
  userID, err := handlers.SessionCheck(token)
  if err == nil {
	  log.Info.Println("token is valid")
	  return userID, nil
  } else {
	  log.Info.Println("invalid token")
	  return 0, errors.New("invalid token")
  }
}

func filterParamsValues(queries map[string][]string) map[string][]string {
	var filter = regexp.MustCompile(` + "`" + `(?i)(password)|(token)` + "`" + `)
	queries_filtered := make(map[string][]string)

	for key := range queries {
		if filter.MatchString(key) {
			queries_filtered[key] = []string{"[FILTERED]"}
		} else {
			queries_filtered[key] = []string{}
			for _, element := range queries[key] {
				queries_filtered[key] = append(queries_filtered[key], element)
			}
		}

	}

	return queries_filtered
}

func filterUrlValues(path string, queries map[string][]string) string {
	var firstParam = true
	queries_filtered := filterParamsValues(queries)

	for key := range queries_filtered {
		if firstParam {
			path = path + "?"
			firstParam = false
		} else {
			path = path + "&"
		}

		path = path + key + "=" + strings.Join(queries_filtered[key], " ")
	}

	return path
}

func filterFormValues(queries map[string][]string) string {
	var buffer bytes.Buffer
	var index int
	queries_filtered := filterParamsValues(queries)

	index = 0
	buffer.WriteString("{ ")

	for key := range queries_filtered {
		buffer.WriteString("\"")
		buffer.WriteString(key)
		buffer.WriteString("\": \"")

		buffer.WriteString(strings.Join(queries_filtered[key], " "))
		buffer.WriteString("\"")

		if (index + 1) != len(queries_filtered) {
			buffer.WriteString(", ")
		}

		index++
	}

	buffer.WriteString(" }")

	return buffer.String()
}

func filterJsonValues(inputJson string) string {
	type Point struct {
		StartAt int
		EndAt   int
	}

	var filter = regexp.MustCompile(` + "`" + `(?i)(password)|(token)` + "`" + `)
	var regexpWhiteSpace = regexp.MustCompile(` + "`" + `[\s\t\n]{1}` + "`" + `)
	var stack []string
	var key, currentChar string
	var valueStartAt, valueEndAt, keyStartAt, keyEndAt int
	var points []Point

	substring := []rune(inputJson)
	isCharBeforeEscape := false
	isKey := false
	isValue := false

	for i := 0; i < len(inputJson); i++ {
		currentChar = string(inputJson[i])

		if regexpWhiteSpace.MatchString(currentChar) && (len(stack) == 0 || stack[len(stack)-1] != ` + "`" + `"` + "`" + ` || stack[len(stack)-1] == ` + "`" + `:` + "`" + `) {
			currentChar = currentChar
		} else if currentChar == "{" && (len(stack) == 0 || stack[len(stack)-1] == "{") {
			stack = append(stack, "{")
		} else if currentChar == ":" && (len(stack) == 0 || stack[len(stack)-1] == "{") {
			stack = append(stack, ":")
		} else if currentChar == "}" && (len(stack) > 0 && stack[len(stack)-1] == "{") {
			stack = stack[:len(stack)-1]
		} else if currentChar == ` + "`" + `"` + "`" + ` && (len(stack) > 0 && stack[len(stack)-1] == ":") && !isCharBeforeEscape {
			stack = stack[:len(stack)-1]
			stack = append(stack, ` + "`" + `"` + "`" + `)
			isValue = true
			isKey = false
			valueStartAt = i + 1
		} else if currentChar == ` + "`" + `"` + "`" + ` && (len(stack) > 0 && stack[len(stack)-1] == ` + "`" + `"` + "`" + `) && isValue && !isCharBeforeEscape {
			stack = stack[:len(stack)-1]
			isValue = false
			isKey = false
			valueEndAt = i

			if filter.MatchString(key) {
				points = append(points, Point{StartAt: valueStartAt, EndAt: valueEndAt})
			}
		} else if currentChar == ` + "`" + `"` + "`" + ` && (len(stack) > 0 && stack[len(stack)-1] == "{") && !isCharBeforeEscape {
			isKey = true
			isValue = false
			stack = append(stack, ` + "`" + `"` + "`" + `)
			keyStartAt = i + 1
		} else if currentChar == ` + "`" + `"` + "`" + ` && (len(stack) > 0 && stack[len(stack)-1] == ` + "`" + `"` + "`" + `) && isKey && !isCharBeforeEscape {
			isKey = false
			isValue = false
			stack = stack[:len(stack)-1]
			keyEndAt = i
			key = string(substring[keyStartAt:keyEndAt])
		}

		if currentChar == ` + "`" + `\` + "`" + ` && (len(stack) > 0 && stack[len(stack)-1] == ` + "`" + `"` + "`" + `) {
			isCharBeforeEscape = true
		} else if isCharBeforeEscape {
			isCharBeforeEscape = false
		}
	}

	if len(points) > 0 {
		for i := len(points) - 1; i >= 0; i-- {
			inputJson = string(substring[0:points[i].StartAt]) + "[FILTERED]" + string(substring[points[i].EndAt:len(inputJson)])
			substring = []rune(inputJson)
		}
	}

	return inputJson
}`
