package command

import (
	"fmt"
	"go/token"
	"strings"
	"unicode"
)

func validatePackageName(language, packageName string) error {
	switch language {
	case "go":
		return validateGoPackageName(packageName)
	case "python":
		return validatePythonPackageName(packageName)
	default:
		return fmt.Errorf("unsupported language: %s (supported: go, python)", language)
	}
}

func validateGoPackageName(packageName string) error {
	if packageName == "" {
		return fmt.Errorf("invalid Go package name: name cannot be empty")
	}
	if token.IsKeyword(packageName) {
		return fmt.Errorf("invalid Go package name %q: it is a Go keyword", packageName)
	}
	if !isGoIdentifier(packageName) {
		return fmt.Errorf("invalid Go package name %q: must start with a letter or underscore and contain only letters, digits, or underscores", packageName)
	}
	return nil
}

func validatePythonPackageName(packageName string) error {
	if packageName == "" {
		return fmt.Errorf("invalid Python package name: name cannot be empty")
	}
	if pythonKeywords[strings.ToLower(packageName)] {
		return fmt.Errorf("invalid Python package name %q: it is a Python keyword", packageName)
	}
	if !isPythonIdentifier(packageName) {
		return fmt.Errorf("invalid Python package name %q: must be a valid Python identifier", packageName)
	}
	return nil
}

func isGoIdentifier(value string) bool {
	for index, r := range value {
		if index == 0 {
			if r != '_' && !unicode.IsLetter(r) {
				return false
			}
			continue
		}
		if r != '_' && !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func isPythonIdentifier(value string) bool {
	for index, r := range value {
		if index == 0 {
			if r != '_' && !unicode.IsLetter(r) {
				return false
			}
			continue
		}
		if r != '_' && !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

var pythonKeywords = map[string]bool{
	"and":      true,
	"as":       true,
	"assert":   true,
	"async":    true,
	"await":    true,
	"break":    true,
	"case":     true,
	"class":    true,
	"continue": true,
	"def":      true,
	"del":      true,
	"elif":     true,
	"else":     true,
	"except":   true,
	"false":    true,
	"finally":  true,
	"for":      true,
	"from":     true,
	"global":   true,
	"if":       true,
	"import":   true,
	"in":       true,
	"is":       true,
	"lambda":   true,
	"match":    true,
	"none":     true,
	"nonlocal": true,
	"not":      true,
	"or":       true,
	"pass":     true,
	"raise":    true,
	"return":   true,
	"true":     true,
	"try":      true,
	"type":     true,
	"while":    true,
	"with":     true,
	"yield":    true,
}
