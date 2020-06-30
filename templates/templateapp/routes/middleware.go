package routes

var MiddlewarePath = []string{"routes", "middleware.go"}

var MiddlewareContent = `package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
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
				log.Info.Println("Body JSON:", filterSensitiveDataInJson(string(body)))
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
	queriesFiltered := make(map[string][]string)

	for key := range queries {
		if matchKey(key) {
			queriesFiltered[key] = []string{"[FILTERED]"}
		} else {
			queriesFiltered[key] = []string{}
			for _, element := range queries[key] {
				queriesFiltered[key] = append(queriesFiltered[key], element)
			}
		}

	}

	return queriesFiltered
}

func filterUrlValues(path string, queries map[string][]string) string {
	var firstParam = true
	queriesFiltered := filterParamsValues(queries)

	for key := range queriesFiltered {
		if firstParam {
			path = path + "?"
			firstParam = false
		} else {
			path = path + "&"
		}

		path = path + key + "=" + strings.Join(queriesFiltered[key], " ")
	}

	return path
}

func filterFormValues(queries map[string][]string) string {
	var buffer bytes.Buffer
	var index int
	queriesFiltered := filterParamsValues(queries)

	index = 0
	buffer.WriteString("{ ")

	for key := range queriesFiltered {
		buffer.WriteString("\"")
		buffer.WriteString(key)
		buffer.WriteString("\": \"")

		buffer.WriteString(strings.Join(queriesFiltered[key], " "))
		buffer.WriteString("\"")

		if (index + 1) != len(queriesFiltered) {
			buffer.WriteString(", ")
		}

		index++
	}

	buffer.WriteString(" }")

	return buffer.String()
}

func matchKey(candidate string) bool {
	values := []string{"password", "token", "webhook"}

	for _, value := range values {
		if candidate == value {
			return true
		}
	}

	return false
}

func filterValuesByKeyInJson(v map[string]interface{}) {
	for key := range v {

		newV, ok := v[key].(map[string]interface{})
		if ok {
			filterValuesByKeyInJson(newV)
		}

		if matchKey(key) {
			v[key] = "[FILTERED]"
		}
	}
}

func filterSensitiveDataInJson(inputJson string) string {
	if inputJson == "" {
		return "{}"
	}

	var v map[string]interface{}
	dec := json.NewDecoder(strings.NewReader(inputJson))

	if err := dec.Decode(&v); err != nil {
		log.Warn.Println("Middleware JSON decoding:", err)
		return "(invalid)"
	}

	filterValuesByKeyInJson(v)

	b, err := json.Marshal(&v)
	if err != nil {
		log.Warn.Println("Middleware JSON marshalling:", err)
		return "(invalid)"
	} else {
		return string(b)
	}
}`
