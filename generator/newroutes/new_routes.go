package newroutes

import (
	"errors"
	"regexp"
	"strings"
)

const (
	beginning = 0

	funcCharWasFoundF  = 1
	funcCharWasFoundU  = 2
	funcCharWasFoundN  = 3
	funcCharWasFoundC  = 4
	funcStringWasFound = 4

	routesCharWasFoundR = 5
	routesCharWasFoundO = 6
	routesCharWasFoundU = 7
	routesCharWasFoundT = 8
	routesCharWasFoundE = 9
	routesCharWasFoundS = 10

	commentSingleLine = 12
	commentMultiLine  = 13

	returnCharWasFoundR  = 15
	returnCharWasFoundE  = 16
	returnCharWasFoundT  = 17
	returnCharWasFoundU  = 18
	returnCharWasFoundR2 = 19
	returnCharWasFoundN  = 20

	returnRoutesStringWasFound = 20
)

var lastCloseBracket int
var currentState int
var currentChar string
var returnRouteAt int
var regexpEmptyChar = regexp.MustCompile(`[\s\t\n\r\f]`)
var funcRoutesWasFound bool
var stack []string
var lastCharWasBackSlash bool
var lastCharWasSlash bool
var lastCharWasStar bool
var stateBeforeComment int

func insideFuncRoutes(i int) {
	if currentChar == "{" && (len(stack) == 0 || (len(stack) > 0 && stack[len(stack)-1] == "{")) {
		stack = append(stack, currentChar)
	} else if currentChar == "}" && stack[len(stack)-1] == "{" {
		stack = stack[:len(stack)-1]
	} else if currentChar == "`" && stack[len(stack)-1] == "`" && !lastCharWasBackSlash {
		stack = stack[:len(stack)-1]
	} else if currentChar == "`" && stack[len(stack)-1] == "{" {
		stack = append(stack, "`")
	} else if currentChar == "\"" && stack[len(stack)-1] == "\"" && !lastCharWasBackSlash {
		stack = stack[:len(stack)-1]
	} else if currentChar == "\"" && stack[len(stack)-1] == "{" {
		stack = append(stack, "\"")
	} else if currentChar == "\"" && stack[len(stack)-1] == "{" {
		stack = append(stack, "\"")
	} else if currentChar == "/" && lastCharWasSlash && (len(stack) == 0 || stack[len(stack)-1] == "{") {
		currentState = commentSingleLine
	} else if currentChar == "*" && lastCharWasSlash && (len(stack) == 0 || stack[len(stack)-1] == "{") {
		currentState = commentMultiLine
	} else if currentChar == "r" && currentState == routesCharWasFoundS {
		currentState = returnCharWasFoundR
		returnRouteAt = i
	} else if currentChar == "e" && currentState == returnCharWasFoundR {
		currentState = returnCharWasFoundE
	} else if currentChar == "t" && currentState == returnCharWasFoundE {
		currentState = returnCharWasFoundT
	} else if currentChar == "u" && currentState == returnCharWasFoundT {
		currentState = returnCharWasFoundU
	} else if currentChar == "r" && currentState == returnCharWasFoundU {
		currentState = returnCharWasFoundR2
	} else if currentChar == "n" && currentState == returnCharWasFoundR2 {
		currentState = returnRoutesStringWasFound
	} else {
		currentState = routesCharWasFoundS
	}

	lastCharWasBackSlash = (currentChar == "\\")
	lastCharWasSlash = (currentChar == "/")
}

func AppendNewCode(newCode string, code string) (string, error) {
	var err error
	var i int
	var outputStr string

	if newCodeAlreadyExists(newCode, code) {
		return "", nil
	}

	funcRoutesWasFound = false
	currentState = beginning
	stack = nil

	for i = 0; i < len(code); i++ {
		currentChar = code[i : i+1]
		if regexpEmptyChar.MatchString(currentChar) && currentState != commentSingleLine {
			continue
		} else if currentChar == "f" && currentState == beginning {
			currentState = funcCharWasFoundF
		} else if currentChar == "u" && currentState == funcCharWasFoundF {
			currentState = funcCharWasFoundU
		} else if currentChar == "n" && currentState == funcCharWasFoundU {
			currentState = funcCharWasFoundN
		} else if currentChar == "c" && currentState == funcCharWasFoundN {
			currentState = funcCharWasFoundC
		} else if currentChar == "R" && currentState == funcCharWasFoundC {
			currentState = routesCharWasFoundR
		} else if currentChar == "o" && currentState == routesCharWasFoundR {
			currentState = routesCharWasFoundO
		} else if currentChar == "u" && currentState == routesCharWasFoundO {
			currentState = routesCharWasFoundU
		} else if currentChar == "t" && currentState == routesCharWasFoundU {
			currentState = routesCharWasFoundT
		} else if currentChar == "e" && currentState == routesCharWasFoundT {
			currentState = routesCharWasFoundE
		} else if currentChar == "s" && currentState == routesCharWasFoundE {
			funcRoutesWasFound = true
			currentState = routesCharWasFoundS
		} else if currentChar == "/" && lastCharWasSlash && currentState != commentMultiLine {
			stateBeforeComment = currentState
			currentState = commentSingleLine
		} else if currentChar == "*" && lastCharWasSlash && currentState != commentSingleLine {
			stateBeforeComment = currentState
			currentState = commentMultiLine
		} else if funcRoutesWasFound && currentState != returnRoutesStringWasFound {
			insideFuncRoutes(i)
		} else if funcRoutesWasFound && currentState == returnRoutesStringWasFound {
			break
		} else {
			currentState = beginning
		}
	}

	if currentState == returnRoutesStringWasFound && len(stack) > 0 {
		outputStr = code[0:returnRouteAt-1] + newCode + "\n\n" + code[returnRouteAt-1:len(code)]
	} else {
		err = errors.New("Could not parse Routes file.")
	}

	return outputStr, err
}

func newCodeAlreadyExists(newCode string, code string) bool {
	var found bool
	var trimRegexp = regexp.MustCompile(`(^[\s\t\n\r\f]+|[\s\t\n\r\f]+$)`)

	found = true
	lines := strings.Split(newCode, "\n")

	for _, line := range lines {
		line = trimRegexp.ReplaceAllString(line, ``)
		if strings.Index(code, line) < 0 {
			found = false
			break
		}
	}

	return found
}
