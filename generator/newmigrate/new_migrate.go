package newmigrate

import (
	"errors"
	"regexp"
	"strings"
)

const (
	beginning = 0

	funcCharWasFoundF = 1
	funcCharWasFoundU = 2
	funcCharWasFoundN = 3
	funcCharWasFoundC = 4

	migrateCharWasFoundL  = 5
	migrateCharWasFoundO  = 6
	migrateCharWasFoundA  = 7
	migrateCharWasFoundD  = 8
	migrateCharWasFoundM  = 9
	migrateCharWasFoundI  = 10
	migrateCharWasFoundG  = 11
	migrateCharWasFoundR  = 12
	migrateCharWasFoundA2 = 13
	migrateCharWasFoundT  = 14
	migrateCharWasFoundI2 = 15
	migrateCharWasFoundO2 = 16
	migrateCharWasFoundN  = 17
	migrateCharWasFoundS  = 18

	commentSingleLine = 19
	commentMultiLine  = 20

	lastCloseBracketFound = 21
)

var lastCloseBracket int
var currentState int
var currentChar string
var regexpEmptyChar = regexp.MustCompile(`[\s\t\n\r\f]`)
var stack []string
var lastCharWasBackSlash bool
var lastCharWasSlash bool
var lastCharWasStar bool
var stateBeforeComment int

func insideFuncMigrate(i int) {
	if currentChar == "{" && (len(stack) == 0 || (len(stack) > 0 && stack[len(stack)-1] == "{")) {
		stack = append(stack, currentChar)
	} else if currentChar == "}" && stack[len(stack)-1] == "{" {
		stack = stack[:len(stack)-1]
		lastCloseBracket = i + 1
		if len(stack) == 0 {
			currentState = lastCloseBracketFound
		}
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
	}
}

func insideCommentSingleLine() {
	lastCharWasSlash = false
	if currentChar == "\n" {
		currentState = stateBeforeComment
	}
}

func insideCommentMultiLine() {
	lastCharWasSlash = false
	if currentChar == "/" && lastCharWasStar {
		currentState = stateBeforeComment
	}

	lastCharWasStar = currentChar == "*"
}

func AppendNewCode(newCode string, code string) (string, error) {
	var i int
	var outputStr string

	if newCodeAlreadyExists(newCode, code) {
		return "", nil
	}

	lastCharWasSlash = false
	lastCharWasSlash = false
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
		} else if currentChar == "l" && currentState == funcCharWasFoundC {
			currentState = migrateCharWasFoundL
		} else if currentChar == "o" && currentState == migrateCharWasFoundL {
			currentState = migrateCharWasFoundO
		} else if currentChar == "a" && currentState == migrateCharWasFoundO {
			currentState = migrateCharWasFoundA
		} else if currentChar == "d" && currentState == migrateCharWasFoundA {
			currentState = migrateCharWasFoundD
		} else if currentChar == "M" && currentState == migrateCharWasFoundD {
			currentState = migrateCharWasFoundM
		} else if currentChar == "i" && currentState == migrateCharWasFoundM {
			currentState = migrateCharWasFoundI
		} else if currentChar == "g" && currentState == migrateCharWasFoundI {
			currentState = migrateCharWasFoundG
		} else if currentChar == "r" && currentState == migrateCharWasFoundG {
			currentState = migrateCharWasFoundR
		} else if currentChar == "a" && currentState == migrateCharWasFoundR {
			currentState = migrateCharWasFoundA2
		} else if currentChar == "t" && currentState == migrateCharWasFoundA2 {
			currentState = migrateCharWasFoundT
		} else if currentChar == "i" && currentState == migrateCharWasFoundT {
			currentState = migrateCharWasFoundI2
		} else if currentChar == "o" && currentState == migrateCharWasFoundI2 {
			currentState = migrateCharWasFoundO2
		} else if currentChar == "n" && currentState == migrateCharWasFoundO2 {
			currentState = migrateCharWasFoundN
		} else if currentChar == "s" && currentState == migrateCharWasFoundN {
			currentState = migrateCharWasFoundS
		} else if currentChar == "/" && lastCharWasSlash && currentState != commentMultiLine {
			stateBeforeComment = currentState
			currentState = commentSingleLine
		} else if currentChar == "*" && lastCharWasSlash && currentState != commentSingleLine {
			stateBeforeComment = currentState
			currentState = commentMultiLine
		} else if currentState == migrateCharWasFoundS {
			insideFuncMigrate(i)
		} else if currentState == commentSingleLine {
			insideCommentSingleLine()
		} else if currentState == commentMultiLine {
			insideCommentMultiLine()
		} else {
			currentState = beginning
		}

		lastCharWasBackSlash = (currentChar == "\\")
		lastCharWasSlash = (currentChar == "/")

		if currentState == lastCloseBracketFound {
			outputStr = code[0:lastCloseBracket-1] + "\n    " + newCode + "\n" + code[lastCloseBracket-1:len(code)]
			return outputStr, nil
		}

	}

	return "", errors.New("Could not parse Migrate file.")
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
