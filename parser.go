package main

import (
	"fmt"
	"go/ast"
	"strings"
)

func ParseService(file *ast.File) (res *Service, err error) {
	srv := new(Service)
	srv.PackageName = file.Name.Name

	ast.Inspect(file, func(node ast.Node) bool {
		switch t := node.(type) {
		case *ast.TypeSpec:

			switch s := t.Type.(type) {
			case *ast.InterfaceType:
				srv.Name = strings.Replace(t.Name.Name, "Service", "", -1)

				for _, m := range s.Methods.List {
					srv.Methods = append(srv.Methods, parseMethod(m))
				}

				res = srv

				return false
			}
		case *ast.ImportSpec:
			// filter out context
			if t.Path.Value == "\"context\"" {
				return false
			}

			srv.Imports = append(srv.Imports, t.Path.Value)

			return false
		}

		return true
	})

	return
}

func parseMethod(field *ast.Field) *Method {
	m := new(Method)
	m.Name = field.Names[0].Name

	switch t := field.Type.(type) {
	case *ast.FuncType:
		if t.Params != nil {
			for _, p := range t.Params.List {
				m.Params = append(m.Params, parseParam(p, len(m.Params)))
			}
		}

		if t.Results != nil {
			for _, r := range t.Results.List {
				m.Results = append(m.Results, parseResult(r, len(m.Results)))
			}
		}
	}

	return m
}

func parseParam(field *ast.Field, idx int) *Param {
	p := new(Param)
	if len(field.Names) > 0 {
		p.Name = field.Names[0].Name
	} else {
		p.Name = fmt.Sprintf("param%d", idx)
	}
	p.Type = getTypeName(field.Type, "")

	return p
}

func parseResult(field *ast.Field, idx int) *Result {
	r := new(Result)
	r.Type = getTypeName(field.Type, "")
	if len(field.Names) > 0 {
		r.Name = field.Names[0].Name
	} else {
		r.Name = fmt.Sprintf("res%d", idx)
	}
	if idx == 0 {
		r.GenName = "res"
	} else {
		r.GenName = fmt.Sprintf("res%d", idx)
	}

	return r
}

func getTypeName(t ast.Expr, pkgName string) string {
	switch t1 := t.(type) {
	case *ast.StructType:
		return "struct{}"
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.SelectorExpr:
		// we have to override this crap...
		return fmt.Sprintf("%s.%s", t1.X, getTypeName(t1.Sel, ""))
	case *ast.StarExpr:
		return fmt.Sprintf("*%s", getTypeName(t1.X, pkgName))
	case *ast.Ident:
		if pkgName != "" {
			return fmt.Sprintf("%s.%s", pkgName, t1)
		}
		return fmt.Sprintf("%s", t1)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", getTypeName(t1.Key, pkgName), getTypeName(t1.Value, pkgName))
	case *ast.ArrayType:
		l := ""
		if t1.Len != nil {
			// we have an array, not a slice.. pity...
			l = fmt.Sprintf("%s", t1.Len)
		}
		return fmt.Sprintf("[%s]%s", l, getTypeName(t1.Elt, pkgName))
	default:
		return fmt.Sprintf("UKNOWN: +V", t)
	}
}
