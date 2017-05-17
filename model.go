package main

import (
	"fmt"

	"github.com/Hepri/case-transformer"
)

type Service struct {
	PackageName string
	Name        string
	Methods     []*Method
	Imports     []string
}

type Method struct {
	Name    string
	Params  []*Param
	Results []*Result
}

func (m *Method) CamelCaseName() string {
	return case_transformer.StringToCamelCase(m.Name)
}

func (m *Method) PascalCaseName() string {
	return case_transformer.StringToPascalCase(m.Name)
}

func (m *Method) EndpointRequestName() string {
	return fmt.Sprintf("%sRequest", m.Name)
}

func (m *Method) EndpointResponseName() string {
	return fmt.Sprintf("%sResponse", m.Name)
}

type Param struct {
	Name string
	Type string
}

func (p *Param) CamelCaseName() string {
	return case_transformer.StringToCamelCase(p.Name)
}

func (p *Param) PascalCaseName() string {
	return case_transformer.StringToPascalCase(p.Name)
}

type Result struct {
	Name    string
	Type    string
	GenName string
}

func (r *Result) CamelCaseName() string {
	return case_transformer.StringToCamelCase(r.Name)
}

func (r *Result) PascalCaseName() string {
	return case_transformer.StringToPascalCase(r.Name)
}
