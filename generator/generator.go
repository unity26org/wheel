package generator

import (
	"github.com/unity26org/wheel/generator/newapp"
	"github.com/unity26org/wheel/generator/newcrud"
)

func NewApp(options map[string]interface{}) error {
	return newapp.Generate(options)
}

func NewCrud(entityName string, columns []string, options map[string]bool) error {
	return newcrud.Generate(entityName, columns, options)
}
