package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/tinkler/mqttadmin/pkg/jsonz/cjson"
	"github.com/tinkler/mqttadmin/pkg/jsonz/sjson"
)

// GenerateRoutes generate routes for package using template.
func GenerateChiRoutes(path string, pkg *Package, dep map[string]*Package) error {
	f, err := os.Create(filepath.Join(path, fmt.Sprintf("%s.go", pkg.Name)))
	if err != nil {
		return err
	}
	defer f.Close()

	debugPathFile, err := os.Create(filepath.Join(path, fmt.Sprintf("%s_path.go", pkg.Name)))
	if err != nil {
		return err
	}
	defer debugPathFile.Close()

	{
		f.WriteString(fmt.Sprintln("// Code generated by github.com/tinkler/mqttadmin; DO NOT EDIT."))
		f.WriteString(fmt.Sprintln("package route"))

		debugPathFile.WriteString(fmt.Sprintln("// Code generated by github.com/tinkler/mqttadmin; DO NOT EDIT."))
		debugPathFile.WriteString(fmt.Sprintln("package route"))
	}

	var (
		usedDepStruct = make(map[string]map[string]bool)
		hasMethods    bool
	)
	for _, s := range pkg.Structs {
		fields := make([]Field, 0)
		for _, m := range s.Methods {
			fields = append(fields, m.Args...)
			fields = append(fields, m.Rets...)
			hasMethods = true
		}
		for _, f := range fields {
			typ := f.Type
		FIND:
			typ = strings.TrimPrefix(typ, "*")
			if _, isBt := tsTypeMap[typ]; !isBt {
				if FindStruct(pkg, typ) == nil {
					if dep != nil {
						if nameSlice := strings.Split(typ, "."); len(nameSlice) > 1 {
							if _, isDep := dep[nameSlice[0]]; isDep {
								if FindStruct(dep[nameSlice[0]], nameSlice[1]) == nil {
									return fmt.Errorf("type %s not found", typ)
								} else {
									if usedDepStruct[nameSlice[0]] == nil {
										usedDepStruct[nameSlice[0]] = make(map[string]bool)
									}
									usedDepStruct[nameSlice[0]][nameSlice[1]] = true
								}
							} else {
								if strings.HasPrefix(typ, "[]") {
									typ = strings.TrimPrefix(typ, "[]")
									goto FIND
								} else {
									return fmt.Errorf("package %s not found", nameSlice[0])
								}
							}
						}

					}
				}
			}
		}
	}

	{
		f.WriteString(fmt.Sprintln("import ("))
		if hasMethods {
			f.WriteString(fmt.Sprintln("\t\"net/http\""))
			f.WriteString(fmt.Sprintln())
		}

		f.WriteString(fmt.Sprintln("\t\"github.com/go-chi/chi/v5\""))
		if hasMethods {
			f.WriteString(fmt.Sprintln("\t\"github.com/tinkler/mqttadmin/pkg/jsonz/sjson\""))
			f.WriteString(fmt.Sprintf("\t\"%s\"\n", pkg.ImportsMap[pkg.Name]))
			f.WriteString(fmt.Sprintln("\t\"github.com/tinkler/mqttadmin/pkg/status\""))
		}
		for _, importName := range pkg.Imports {
			if _, used := usedDepStruct[importName]; used {
				f.WriteString(fmt.Sprintf("\t\"%s\"\n", pkg.ImportsMap[importName]))
			}
		}

		f.WriteString(fmt.Sprintln(")"))
	}

	{
		t := template.Must(template.New("model").Funcs(template.FuncMap{
			"toType": func(typ string) string {
				if _, isBt := tsTypeMap[typ]; isBt {
					return typ
				}
				if sli := strings.Split(typ, "."); len(sli) > 1 {
					return typ
				}
				if strings.HasPrefix(typ, "*") {
					typ = strings.TrimPrefix(typ, "*")
					return fmt.Sprintf("*%s.%s", pkg.Name, typ)
				}
				if strings.HasPrefix(typ, "[]") {
					typ = strings.TrimPrefix(typ, "[]")
					if strings.HasPrefix(typ, "*") {
						typ = strings.TrimPrefix(typ, "*")
						return fmt.Sprintf("[]*%s.%s", pkg.Name, typ)
					}
					return fmt.Sprintf("[]%s.%s", pkg.Name, typ)
				}
				return fmt.Sprintf("%s.%s", pkg.Name, typ)
			},
			"toSnack": func(s string) string {
				return sjson.ToSnackedName(s)
			},
			"toFulle": func(s string) string {
				return cjson.SnakeCaseToFullCamelCase(s)
			},
			"toMinus": func(s string) string {
				return strings.ReplaceAll(sjson.ToSnackedName(s), "_", "-")
			},
		}).Parse(goChiRouteTemplate))
		if err := t.Execute(f, pkg); err != nil {
			return err
		}

		if hasMethods {
			t := template.Must(template.New("model").Funcs(template.FuncMap{
				"toType": func(typ string) string {
					if _, isBt := tsTypeMap[typ]; isBt {
						return typ
					}
					if sli := strings.Split(typ, "."); len(sli) > 1 {
						return typ
					}
					if strings.HasPrefix(typ, "*") {
						typ = strings.TrimPrefix(typ, "*")
						return fmt.Sprintf("*%s.%s", pkg.Name, typ)
					}
					if strings.HasPrefix(typ, "[]") {
						typ = strings.TrimPrefix(typ, "[]")
						if strings.HasPrefix(typ, "*") {
							typ = strings.TrimPrefix(typ, "*")
							return fmt.Sprintf("[]*%s.%s", pkg.Name, typ)
						}
						return fmt.Sprintf("[]%s.%s", pkg.Name, typ)
					}
					return fmt.Sprintf("%s.%s", pkg.Name, typ)
				},
				"toSnack": func(s string) string {
					return sjson.ToSnackedName(s)
				},
				"toFulle": func(s string) string {
					return cjson.SnakeCaseToFullCamelCase(s)
				},
				"toMinus": func(s string) string {
					return strings.ReplaceAll(sjson.ToSnackedName(s), "_", "-")
				},
			}).Parse(goChiRouteDebugPathTemplate))
			if err := t.Execute(debugPathFile, pkg); err != nil {
				return err
			}
		}
	}

	return nil
}

// go route template
const goChiRouteTemplate = `
func Routes{{.Name | toFulle}}(m chi.Router) {
	m.Route("/{{.Name}}", func(r chi.Router) {
		{{range .Structs}}{{$struct := .}}{{range .Methods}}
		r.Post("/{{$struct.Name | toSnack}}/{{.Name | toMinus}}", func(w http.ResponseWriter, r *http.Request) {
			m := Model[*{{$struct.Name | toType}}, {{if .Args}}struct{
				{{range .Args}}{{.Name | toFulle}} {{.Type}} 
				{{end}} } {{else}}any{{end}}]{}
			err := sjson.Bind(r, &m)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			res := Res[*{{$struct.Name | toType}},{{if ge (len .Rets) 1}}{{(index .Rets 0).Type | toType}}{{else}}any{{end}}]{Data:m.Data}
			{{if .Args}}{{if ge (len .Rets) 1}}res.Resp, err = m.Data.{{.Name}}(r.Context(){{if .Args}}, {{$rl := (len .Args)}}{{range $index,$arg := .Args}}m.Args.{{$arg.Name | toFulle}}{{if lt $index $rl}}, {{end}}{{end}}{{end}})
			{{else}}err = m.Data.{{.Name}}(r.Context(){{if .Args}}, {{$rl := (len .Args)}}{{range $index,$arg := .Args}}m.Args.{{$arg.Name | toFulle}}{{if lt $index $rl}}, {{end}}{{end}}{{end}})
			{{end}}{{else}}{{if ge (len .Rets) 1}}res.Resp, err = m.Data.{{.Name}}(r.Context())
			{{else}}err = m.Data.{{.Name}}(r.Context())
			{{end}}{{end}}
			if status.HttpError(w, err) {
				return
			}
			if sjson.HttpWrite(w, res) {
				return
			}

		}){{end}}{{end}}
	})
}
`

// go route debug path template
const goChiRouteDebugPathTemplate = `{{$pkgName := .Name}}
func init() {
	{{range .Structs}}{{$struct := .}}{{range .Methods}}
	routePathMap["/{{$pkgName}}/{{$struct.Name | toSnack}}/{{.Name | toMinus}}"] = "{{.Filename}}:{{.StartLine}}"{{end}}{{end}}
}

`

// GenerateTSCode generate ts code
func GenerateTSCode(path string, pkg *Package, dep map[string]*Package) error {

	f, err := os.Create(filepath.Join(path, fmt.Sprintf("%s.ts", sjson.ToSnackedName(pkg.Name))))
	if err != nil {
		return err
	}
	defer f.Close()

	{
		f.WriteString(fmt.Sprintln("// Code generated by github.com/tinkler/mqttadmin; DO NOT EDIT."))
	}

	var (
		usedDepStruct = make(map[string]map[string]bool)
	)
	for _, s := range pkg.Structs {
		fields := s.Fields
		for _, m := range s.Methods {
			fields = append(fields, m.Args...)
			fields = append(fields, m.Rets...)
		}
		for _, f := range fields {
			typ := f.Type
		FIND:
			typ = strings.TrimPrefix(typ, "*")
			if _, isBt := tsTypeMap[typ]; !isBt {
				if FindStruct(pkg, typ) == nil {
					if dep != nil {
						if nameSlice := strings.Split(typ, "."); len(nameSlice) > 1 {
							if _, isDep := dep[nameSlice[0]]; isDep {
								if FindStruct(dep[nameSlice[0]], nameSlice[1]) == nil {
									return fmt.Errorf("type %s not found", typ)
								} else {
									if usedDepStruct[nameSlice[0]] == nil {
										usedDepStruct[nameSlice[0]] = make(map[string]bool)
									}
									usedDepStruct[nameSlice[0]][nameSlice[1]] = true
								}
							} else {
								if strings.HasPrefix(typ, "[]") {
									typ = strings.TrimPrefix(typ, "[]")
									goto FIND
								} else {
									return fmt.Errorf("package %s not found", nameSlice[0])
								}
							}
						}

					}
				}
			}
		}
	}
	if len(usedDepStruct) > 0 {
		for _, importName := range pkg.Imports {
			if _, ok := usedDepStruct[importName]; ok {
				f.WriteString(fmt.Sprintf("import { %s } from './%s';\n", strings.Join(func() []string {
					var s []string
					for _, ds := range dep[importName].Structs {
						if _, ok := usedDepStruct[importName][ds.Name]; ok {
							s = append(s, ds.Name)
						}
					}
					return s
				}(), ", "), importName))
			}
		}
	}

	if len(pkg.Structs) > 0 {
		t := template.Must(template.New("model").Funcs(template.FuncMap{
			"toType": func(goType string) string {
				suffix := ""
			FIND:
				goType = strings.TrimPrefix(goType, "*")
				if bt, isBt := tsTypeMap[goType]; isBt {
					return bt + suffix
				}
				if s := FindStruct(pkg, goType); s != nil {
					return s.Name + suffix
				}
				if dep != nil {
					if nameSlice := strings.Split(goType, "."); len(nameSlice) > 1 {
						if pkg, isDep := dep[nameSlice[0]]; isDep {
							if s := FindStruct(pkg, nameSlice[1]); s != nil {
								return s.Name + suffix
							}
						}
					}
				}
				if strings.HasPrefix(goType, "[]") {
					suffix += "[]"
					goType = strings.TrimPrefix(goType, "[]")
					goto FIND
				}

				return ""
			},
			"toCamel": func(s string) string {
				return cjson.ToCamel(s)
			},
			"toSnack": func(s string) string {
				return sjson.ToSnackedName(s)
			},
			"toMinus": func(s string) string {
				return strings.ReplaceAll(sjson.ToSnackedName(s), "_", "-")
			},
			"toDefault": func(s string) string {
				if strings.HasSuffix(s, "[]") {
					return "[]"
				}
				if bt, isBt := tsDefaultValue[tsTypeMap[s]]; isBt {
					return bt
				}
				s = strings.TrimPrefix(s, "*")
				if s := FindStruct(pkg, s); s != nil {
					return s.Name + "()"
				}
				if dep != nil {
					if sli := strings.Split(s, "."); len(sli) > 1 {
						if pkg, isDep := dep[sli[0]]; isDep {
							if s := FindStruct(pkg, sli[1]); s != nil {
								return s.Name + "()"
							}
						}
					}

				}
				return ""
			},
		}).Parse(tsModelTemplate))

		err = t.Execute(f, pkg)
		if err != nil {
			return err
		}
	}

	return nil
}

// typescript template
const tsModelTemplate = `{{$pkgName := .Name}}
{{range .Structs}}
{{range .Comments}}/**
* {{.}}
*/
{{end}}export interface {{.Name}} {
	{{range .Fields}}
	{{range .Comments}}/**
	* {{.}}
	*/
	{{end}}{{.Name | toCamel}}: {{.Type | toType}};
	{{end}}
	{{range .Methods}}
	{{range .Comments}}/**
	* {{.}}
	*/
	{{end}}{{.Name | toCamel}}({{$rl := (len .Args)}}{{range $index,$arg := .Args}}{{$arg.Name | toCamel}}: {{$arg.Type | toType}}{{if lt $index $rl}}, {{end}}{{end}}): {{if ge (len .Rets) 1}}Promise<{{$ret := (index .Rets 0)}}{{$ret.Type | toType}}> {{else}}Promise<void>{{end}};
	{{end}}
}
{{end}}

{{range .Structs}}{{$structName := .Name}}
{{range .Comments}}/**
* {{.}}
*/
{{end}}export function {{.Name}}(): {{.Name}} {
	
	return {
		{{range .Fields}}
		{{.Name | toCamel}}: {{.Type | toDefault}},
		{{end}}
		{{range .Methods}}
		{{.Name | toCamel}}({{range $index,$arg := .Args}}{{$arg.Name | toCamel}}: {{$arg.Type | toType}}, {{end}}): {{if ge (len .Rets) 1}}Promise<{{$ret := (index .Rets 0)}}{{$ret.Type | toType}}> {{else}}Promise<void>{{end}} {
			
			return post{{$structName}}(this, '{{.Name | toMinus}}', { {{range $index,$arg := .Args}}{{$arg.Name | toCamel}}, {{end}} }){{if ge (len .Rets) 1}}.then((res: { data: any }) => res.data as {{$ret := (index .Rets 0)}}{{$ret.Type | toType}}){{end}};
			
		},
		{{end}}
		
	};
	
}

// post data by restful api

function post{{.Name}}({{.Name | toCamel}}: {{.Name}}, method: string, args: {}): Promise<any> {
	const xhr = new XMLHttpRequest();
	xhr.open("POST", ` + "`/{{$pkgName}}//{{$struct.Name | toSnack}}/{{.Name | toSnack}}/${method}`" + `, true);
	xhr.setRequestHeader("Content-Type", "application/json");
	return new Promise((resolve, reject) => {
		xhr.onload = () => {
			if (xhr.status === 200) {
				resolve(xhr.response);
			} else {
				reject(new Error(xhr.statusText));
			}
		};
		xhr.onerror = () => {
			reject(new Error(xhr.statusText));
		};
		xhr.send(JSON.stringify({ data: {{.Name | toCamel}}, args }));
	});
}
{{end}}
`

func GenerateTSAngularDelonCode(path string, pkg *Package, dep map[string]*Package) error {
	var (
		serviceFile     *os.File
		serviceSpecFile *os.File
		modelFile       *os.File
	)
	serviceFile, err := os.Create(filepath.Join(path, fmt.Sprintf("%s.service.ts", sjson.ToSnackedName(pkg.Name))))
	if err != nil {
		return err
	}
	defer serviceFile.Close()

	serviceSpecFile, err = os.Create(filepath.Join(path, fmt.Sprintf("%s.spec.ts", sjson.ToSnackedName(pkg.Name))))
	if err != nil {
		return err
	}
	defer serviceSpecFile.Close()

	modelFile, err = os.Create(filepath.Join(path, fmt.Sprintf("%s.model.ts", sjson.ToSnackedName(pkg.Name))))
	if err != nil {
		return err
	}
	defer modelFile.Close()

	{
		doNotEdit := "// Code generated by github.com/tinkler/mqttadmin; DO NOT EDIT."
		serviceFile.WriteString(fmt.Sprintln(doNotEdit))
		serviceSpecFile.WriteString(fmt.Sprintln(doNotEdit))
		modelFile.WriteString(fmt.Sprintln(doNotEdit))
	}

	var (
		modelUsedDepStruct   = make(map[string]map[string]bool)
		serviceUsedDepStruct = make(map[string]map[string]bool)
	)
	// collect used dep struct of model
	for _, s := range pkg.Structs {
		fields := s.Fields
		for _, m := range s.Methods {
			fields = append(fields, m.Args...)
			fields = append(fields, m.Rets...)
		}
		for _, f := range fields {
			typ := f.Type
		MODELFIND:
			typ = strings.TrimPrefix(typ, "*")
			if _, isBt := tsTypeMap[typ]; !isBt {
				if FindStruct(pkg, typ) == nil {
					if dep != nil {
						if nameSlice := strings.Split(typ, "."); len(nameSlice) > 1 {
							if _, isDep := dep[nameSlice[0]]; isDep {
								if FindStruct(dep[nameSlice[0]], nameSlice[1]) == nil {
									return fmt.Errorf("type %s not found", typ)
								} else {
									if modelUsedDepStruct[nameSlice[0]] == nil {
										modelUsedDepStruct[nameSlice[0]] = make(map[string]bool)
									}
									modelUsedDepStruct[nameSlice[0]][nameSlice[1]] = true
								}
							} else {
								if strings.HasPrefix(typ, "[]") {
									typ = strings.TrimPrefix(typ, "[]")
									goto MODELFIND
								} else {
									return fmt.Errorf("package %s not found", nameSlice[0])
								}
							}
						}

					}
				}
			}
		}
	}
	if len(modelUsedDepStruct) > 0 {
		for _, importName := range pkg.Imports {
			if _, ok := modelUsedDepStruct[importName]; ok {
				modelFile.WriteString(fmt.Sprintf("import { %s } from './%s.model';\n", strings.Join(func() []string {
					var s []string
					for _, ds := range dep[importName].Structs {
						if _, ok := modelUsedDepStruct[importName][ds.Name]; ok {
							s = append(s, ds.Name)
						}
					}
					return s
				}(), ", "), importName))
			}
		}
	}

	// collect used dep struct of service
	for _, s := range pkg.Structs {
		var fields []Field
		for _, m := range s.Methods {
			fields = append(fields, m.Args...)
			fields = append(fields, m.Rets...)
		}
		for _, f := range fields {
			typ := f.Type
		FIND:
			typ = strings.TrimPrefix(typ, "*")
			if _, isBt := tsTypeMap[typ]; !isBt {
				if FindStruct(pkg, typ) == nil {
					if dep != nil {
						if nameSlice := strings.Split(typ, "."); len(nameSlice) > 1 {
							if _, isDep := dep[nameSlice[0]]; isDep {
								if FindStruct(dep[nameSlice[0]], nameSlice[1]) == nil {
									return fmt.Errorf("type %s not found", typ)
								} else {
									if serviceUsedDepStruct[nameSlice[0]] == nil {
										serviceUsedDepStruct[nameSlice[0]] = make(map[string]bool)
									}
									serviceUsedDepStruct[nameSlice[0]][nameSlice[1]] = true
								}
							} else {
								if strings.HasPrefix(typ, "[]") {
									typ = strings.TrimPrefix(typ, "[]")
									goto FIND
								} else {
									return fmt.Errorf("package %s not found", nameSlice[0])
								}
							}
						}

					}
				}
			}
		}
	}
	if len(serviceUsedDepStruct) > 0 {
		for _, importName := range pkg.Imports {
			if _, ok := serviceUsedDepStruct[importName]; ok {
				serviceFile.WriteString(fmt.Sprintf("import { %s } from './%s.model';\n", strings.Join(func() []string {
					var s []string
					for _, ds := range dep[importName].Structs {
						if _, ok := serviceUsedDepStruct[importName][ds.Name]; ok {
							s = append(s, ds.Name)
						}
					}
					return s
				}(), ", "), importName))
			}
		}
	}

	if len(pkg.Structs) > 0 {
		globalFuncMap := template.FuncMap{
			"toType": func(goType string) string {
				suffix := ""
			FIND:
				goType = strings.TrimPrefix(goType, "*")
				if bt, isBt := tsTypeMap[goType]; isBt {
					return bt + suffix
				}
				if s := FindStruct(pkg, goType); s != nil {
					return s.Name + suffix
				}
				if dep != nil {
					if nameSlice := strings.Split(goType, "."); len(nameSlice) > 1 {
						if pkg, isDep := dep[nameSlice[0]]; isDep {
							if s := FindStruct(pkg, nameSlice[1]); s != nil {
								return s.Name + suffix
							}
						}
					}
				}
				if strings.HasPrefix(goType, "[]") {
					suffix += "[]"
					goType = strings.TrimPrefix(goType, "[]")
					goto FIND
				}

				return ""
			},
			"toCamel": func(s string) string {
				return cjson.ToCamel(s)
			},
			"toSnack": func(s string) string {
				return sjson.ToSnackedName(s)
			},
			"toMinus": func(s string) string {
				return strings.ReplaceAll(sjson.ToSnackedName(s), "_", "-")
			},
			"toDefault": func(s string) string {
				if strings.HasPrefix(s, "[]") {
					return "[]"
				}
				if bt, isBt := tsDefaultValue[tsTypeMap[s]]; isBt {
					return bt
				}
				s = strings.TrimPrefix(s, "*")
				if s := FindStruct(pkg, s); s != nil {
					return "new " + s.Name + "(this.http)"
				}
				if dep != nil {
					if sli := strings.Split(s, "."); len(sli) > 1 {
						if pkg, isDep := dep[sli[0]]; isDep {
							if s := FindStruct(pkg, sli[1]); s != nil {
								return "new " + s.Name + "(this.http)"
							}
						}
					}
				}
				return ""
			},
			"firstUpper": func(s string) string {
				return strings.ToUpper(s[:1]) + s[1:]
			},
		}

		{
			customFuncMap := make(template.FuncMap)
			for k, v := range globalFuncMap {
				customFuncMap[k] = v
			}
			customFuncMap["respToType"] = func(goType string) string {
				isArray := strings.HasPrefix(goType, "[]")
				if isArray {
					goType = strings.TrimPrefix(goType, "[]")
				}
				goType = strings.TrimPrefix(goType, "*")
				if bt, isBt := tsTypeMap[goType]; isBt {
					if isArray {
						return bt + "[]"
					}
					return bt
				}
				return "any"
			}
			customFuncMap["respSet"] = func(goType string) string {
				isArray := strings.HasPrefix(goType, "[]")
				if isArray {
					goType = strings.TrimPrefix(goType, "[]")
				}
				if bt, isBt := tsTypeMap[goType]; isBt {
					if isArray {
						return "const resp = res.data.resp as []" + bt
					}
					return "const resp = res.data.resp as " + bt
				}
				goType = strings.TrimPrefix(goType, "*")
				if s := FindStruct(pkg, goType); s != nil {
					if isArray {
						content := "const resp: []" + s.Name + " = [];"
						content += "for(const item of res.data.resp) {"
						content += "const _new = new " + s.Name + "(this.http);"
						for _, f := range s.Fields {
							content += "_new." + cjson.ToCamel(f.Name) + " = item['" + sjson.ToSnackedName(f.Name) + "'];"
						}
						content += "resp.push(_new);"
						content += "}"
						return content
					}
					content := "const resp: " + s.Name + " = new " + s.Name + "(this.http);"
					for _, f := range s.Fields {
						content += "resp." + cjson.ToCamel(f.Name) + " = res.data.resp['" + sjson.ToSnackedName(f.Name) + "'];"
					}
					return content
				}
				return "const resp = res.data.resp"
			}
			t := template.Must(template.New("model").Funcs(customFuncMap).Parse(tsAngularDelonModelTemplate))
			err = t.Execute(modelFile, pkg)
			if err != nil {
				return err
			}
		}

		{

			t := template.Must(template.New("service").Funcs(globalFuncMap).Parse(tsAngularDelonServiceTemplate))
			err = t.Execute(serviceFile, pkg)
			if err != nil {
				return err
			}
		}

		{
			t := template.Must(template.New("serviceSpec").Funcs(globalFuncMap).Parse(tsAngularDelonServiceSpecTemplate))
			err = t.Execute(serviceSpecFile, pkg)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

const tsAngularDelonModelTemplate = `{{$pkgName := .Name}}

import { _HttpClient } from '@delon/theme';
import { modelUrlPrefix  } from './const';
{{range .Structs}}{{$struct := .}}
{{range .Comments}}/**
* {{.}}
*/
{{end}}export interface {{.Name}} {
	{{range .Fields}}
	{{range .Comments}}/**
	* {{.}}
	*/
	{{end}}{{.Name | toCamel}}: {{.Type | toType}};
	{{end}}
	{{range .Methods}}
	{{range .Comments}}/**
	* {{.}}
	*/
	{{end}}{{.Name | toCamel}}({{$rl := (len .Args)}}{{range $index,$arg := .Args}}{{$arg.Name | toCamel}}: {{$arg.Type | toType}}{{if lt $index $rl}}, {{end}}{{end}}): {{if ge (len .Rets) 1}}Promise<{{$ret := (index .Rets 0)}}{{$ret.Type | toType}}> {{else}}Promise<void>{{end}};
	{{end}}
}

export class {{.Name}} {
	{{range .Fields}}
	{{range .Comments}}/**
	* {{.}}
	*/
	{{end}}{{.Name | toCamel}}: {{.Type | toType}} = {{.Type | toDefault}};
	{{end}}

	constructor(
		private http: _HttpClient,
	){}

	{{range .Methods}}
	{{range .Comments}}/**
	* {{.}}
	*/
	{{end}}{{.Name | toCamel}}({{range $index,$arg := .Args}}{{$arg.Name | toCamel}}: {{$arg.Type | toType}}, {{end}}): {{if ge (len .Rets) 1}}Promise<{{$ret := (index .Rets 0)}}{{$ret.Type | toType}}> {{else}}Promise<void>{{end}} {
		return new Promise((resolve, reject) => {
			this.http.post(` + "`" + `${modelUrlPrefix}/{{$pkgName}}/{{$struct.Name | toSnack}}/{{.Name | toMinus}}` + "`" + `, { data: this, args: { {{range $index,$arg := .Args}}{{$arg.Name | toSnack}}: {{$arg.Name | toCamel}}, {{end}} } }).subscribe({
				next: (res: { code: number; data: { data: any, resp: {{if ge (len .Rets) 1}}{{$ret := (index .Rets 0)}}{{$ret.Type | respToType}}{{else}}{}{{end}} }, message: string } ) => {
					if (res.code === 0) {
						{{range $struct.Fields}}this.{{.Name | toCamel}} = res.data.data['{{.Name | toSnack}}'];
						{{end}}
						{{if ge (len .Rets) 1}}{{$ret := (index .Rets 0)}}
						{{$ret.Type | respSet}}
						resolve(resp);
						{{else}}resolve();{{end}}
					} else {
						reject(res.message);
					}
				}, error: (err) => {
					reject(err);
				}
			});
		});
	}
	{{end}}
}

{{end}}


`
const tsAngularDelonServiceTemplate = `{{$pkgName := .Name}}
import { Injectable } from '@angular/core';
import { {{range .Structs}}{{.Name}}, {{end}} } from './{{.Name | toMinus}}.model';
{{range .Imports}}{{if ne . $pkgName}}import { {{. | firstUpper}}Service  } from './{{. | toMinus}}.service';
{{end}}{{end}}
import { _HttpClient } from '@delon/theme';

@Injectable({
	providedIn: 'root'
})
export class {{.Name | firstUpper}}Service {
  
	constructor(
		private http: _HttpClient,
		{{range .Imports}}private {{. | toCamel}}Srv: {{. | firstUpper}}Service,
		{{end}}
		) { }

	{{range .Structs}}{{$structName := .Name}}{{$struct := .}}
	{{range .Comments}}/**
	* {{.}}
	*/
	{{end}}new{{.Name}}(): {{.Name}} {
		return new {{.Name}}(this.http);
	}{{end}}
}
`

const tsAngularDelonServiceSpecTemplate = `{{$pkgName := .Name}}
import { TestBed } from '@angular/core/testing';

import { {{.Name | firstUpper}}Service } from './{{.Name | toMinus}}.service';

describe('{{.Name | firstUpper}}Service', () => {
  let service: {{.Name | firstUpper}}Service;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject({{.Name | firstUpper}}Service);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});

`

// GenerateDartCode generate dart code
func GenerateDartCode(path string, pkg *Package, dep map[string]*Package) error {
	f, err := os.Create(filepath.Join(path, fmt.Sprintf("%s.dart", sjson.ToSnackedName(pkg.Name))))
	if err != nil {
		return err
	}
	defer f.Close()

	{
		f.WriteString(fmt.Sprintln("// Code generated by github.com/tinkler/mqttadmin; DO NOT EDIT."))
		f.WriteString(fmt.Sprintln("import '../../http.dart';"))
		f.WriteString(fmt.Sprintln("import './const.dart';"))
	}
	var (
		usedDepStruct = make(map[string]map[string]bool)
	)
	for _, s := range pkg.Structs {
		fields := s.Fields
		for _, m := range s.Methods {
			fields = append(fields, m.Args...)
			fields = append(fields, m.Rets...)
		}
		for _, f := range fields {
			typ := f.Type
		FIND:
			typ = strings.TrimPrefix(typ, "*")
			if _, isBt := tsTypeMap[typ]; !isBt {
				if FindStruct(pkg, typ) == nil {
					if dep != nil {
						if nameSlice := strings.Split(typ, "."); len(nameSlice) > 1 {
							if _, isDep := dep[nameSlice[0]]; isDep {
								if FindStruct(dep[nameSlice[0]], nameSlice[1]) == nil {
									return fmt.Errorf("type %s not found", typ)
								} else {
									if usedDepStruct[nameSlice[0]] == nil {
										usedDepStruct[nameSlice[0]] = make(map[string]bool)
									}
									usedDepStruct[nameSlice[0]][nameSlice[1]] = true
								}
							} else {
								if strings.HasPrefix(typ, "[]") {
									typ = strings.TrimPrefix(typ, "[]")
									goto FIND
								} else {
									return fmt.Errorf("package %s not found", nameSlice[0])
								}
							}
						}

					}
				}
			}
		}
	}
	if len(usedDepStruct) > 0 {
		for _, importName := range pkg.Imports {
			if _, ok := usedDepStruct[importName]; ok {
				f.WriteString(fmt.Sprintf("import './%[1]s.dart' as $%[1]s show %[2]s;\n", importName, strings.Join(func() []string {
					var s []string
					for _, ds := range dep[importName].Structs {
						if _, ok := usedDepStruct[importName][ds.Name]; ok {
							s = append(s, ds.Name)
						}
					}
					return s
				}(), ", "), importName))
			}
		}
	}

	t := template.Must(template.New("model").Funcs(template.FuncMap{
		"toType": func(goType string) string {
			nullCheckSuffix := ""
			if strings.HasPrefix(goType, "*") {
				nullCheckSuffix = "?"
			}
			goType = strings.ReplaceAll(goType, "*", "")
			if typ := dartTypeMap[goType]; typ != "" {
				return typ + nullCheckSuffix
			}
			if strings.HasPrefix(goType, "[]") {
				var (
					resTypePre string
					resTypeSuf string
					typ        = goType
				)
				for strings.HasPrefix(typ, "[]") {
					resTypePre += "List<"
					resTypeSuf += ">"
					typ = strings.TrimPrefix(typ, "[]")
				}
				if typ := dartTypeMap[typ]; typ != "" {
					return fmt.Sprintf(resTypePre+"%s"+resTypeSuf+nullCheckSuffix, typ)
				}
				if s := FindStruct(pkg, typ); s != nil {
					return fmt.Sprintf(resTypePre+"%s"+resTypeSuf+nullCheckSuffix, typ)
				}
				if dep != nil {
					if nameSplice := strings.Split(typ, "."); len(nameSplice) > 1 {
						if pkg, isDep := dep[nameSplice[0]]; isDep {
							if FindStruct(pkg, nameSplice[1]) != nil {
								return fmt.Sprintf(resTypePre+"%s"+resTypeSuf+nullCheckSuffix, "$"+pkg.Name+"."+nameSplice[1])
							}
						}
					}
				}
				return fmt.Sprintf(resTypePre+"%s"+resTypeSuf+nullCheckSuffix, typ)
			}
			if typ := dartTypeMap[goType]; typ != "" {
				return typ + nullCheckSuffix
			}
			if s := FindStruct(pkg, goType); s != nil {
				return goType + nullCheckSuffix
			}
			if dep != nil {
				if nameSplice := strings.Split(goType, "."); len(nameSplice) > 1 {
					if pkg, isDep := dep[nameSplice[0]]; isDep {
						if FindStruct(pkg, nameSplice[1]) != nil {
							return "$" + pkg.Name + "." + nameSplice[1] + nullCheckSuffix
						}
					}
				}
			}
			return goType + nullCheckSuffix
		},
		"toCamel": func(s string) string {
			return cjson.ToCamel(s)
		},
		"toSnack": func(s string) string {
			return sjson.ToSnackedName(s)
		},
		"toMinus": func(s string) string {
			return strings.ReplaceAll(sjson.ToSnackedName(s), "_", "-")
		},
		"toDefault": func(goType string) string {
			if strings.HasPrefix(goType, "*") {
				return ""
			}
			goType = strings.ReplaceAll(goType, "*", "")
			if typ := dartDefaultValue[dartTypeMap[goType]]; typ != "" {
				return " = " + typ
			}
			if strings.HasPrefix(goType, "[]") {
				return " = " + "[]"
			}
			if dep != nil {
				if nameSplice := strings.Split(goType, "."); len(nameSplice) > 1 {
					if pkg, isDep := dep[nameSplice[0]]; isDep {
						if FindStruct(pkg, nameSplice[1]) != nil {
							return " = $" + pkg.Name + "." + nameSplice[1] + "()"
						}
					}
				}
			}
			return " = " + goType + "()"
		},
		"asToJson": func(goType string, name string) string {
			if typ := dartTypeMap[goType]; typ != "" {
				return ""
			}
			if strings.HasPrefix(goType, "[]") {
				goType = strings.TrimPrefix(goType, "[]")
				if typ := dartTypeMap[goType]; typ != "" {
					return ""
				}
				return ".map((e) => e.toJson()).toList()"
			}
			if strings.HasPrefix(goType, "*") {
				return fmt.Sprintf(" != null ? %[1]s!.toJson() : null", name)
			}
			return name + ".toJson()"
		},
		"fromJson": func(goType string, name string) string {
			goType = strings.ReplaceAll(goType, "*", "")
			if typ := dartTypeMap[goType]; typ != "" {
				return ""
			}
			if strings.HasPrefix(goType, "[]") {
				goType = strings.TrimPrefix(goType, "[]")
				if typ := dartTypeMap[goType]; typ != "" {
					return fmt.Sprintf(" == null ? [] : (json[\"%s\"] as List<dynamic>).map((e) => e as %s).toList()", name, typ)
				}
				if s := FindStruct(pkg, goType); s != nil {
					return fmt.Sprintf(" == null ? [] : (json[\"%s\"] as List<dynamic>).map((e) => %s.fromJson(e)).toList()", name, goType)
				}
				if dep != nil {
					if nameSplice := strings.Split(goType, "."); len(nameSplice) > 1 {
						if pkg, isDep := dep[nameSplice[0]]; isDep {
							if FindStruct(pkg, nameSplice[1]) != nil {
								return fmt.Sprintf(" == null ? [] : (json[\"%s\"] as List<dynamic>).map((e) => $"+"%s."+nameSplice[1]+".fromJson(e)).toList()", name, pkg.Name)
							}
						}
					}
				}
				return fmt.Sprintf(" == null ? [] : (json[\"%s\"] as List<dynamic>).map((e) => %s.fromJson(e)).toList()", name, strings.TrimPrefix(goType, "[]"))
			}
			if s := FindStruct(pkg, goType); s != nil {
				return fmt.Sprintf(" == null ? %[1]s() : %[1]s.fromJson(json[\"%[2]s\"])", goType, name)
			}
			if dep != nil {
				if nameSplice := strings.Split(goType, "."); len(nameSplice) > 1 {
					if pkg, isDep := dep[nameSplice[0]]; isDep {
						if FindStruct(pkg, nameSplice[1]) != nil {
							return fmt.Sprintf(" == null ? %[1]s() : %[1]s.fromJson(json[\"%[2]s\"])", "$"+pkg.Name+"."+nameSplice[1], name)
						}
					}
				}
			}
			return ".fromJson($s)"
		},
		"returnArgs": func(goType string) string {
			goType = strings.ReplaceAll(goType, "*", "")
			if typ := dartTypeMap[goType]; typ != "" {
				return "response.data['data']['resp']"
			}
			if strings.HasPrefix(goType, "[]") {
				goType = strings.TrimPrefix(goType, "[]")
				if typ := dartTypeMap[goType]; typ != "" {
					return "response.data['data']['resp'] as List<" + typ + ">"
				}
				if s := FindStruct(pkg, goType); s != nil {
					return "(response.data['data']['resp'] as List<dynamic>).map((e) => " + goType + ".fromJson(e)).toList()"
				}
				if dep != nil {
					if nameSplice := strings.Split(goType, "."); len(nameSplice) > 1 {
						if pkg, isDep := dep[nameSplice[0]]; isDep {
							if FindStruct(pkg, nameSplice[1]) != nil {
								return "(response.data['data']['resp'] as List<dynamic>).map((e) => $" + pkg.Name + "." + nameSplice[1] + ".fromJson(e)).toList()"
							}
						}
					}
				}
			}
			if s := FindStruct(pkg, goType); s != nil {
				return goType + ".fromJson(response.data['data']['resp'])"
			}
			if dep != nil {
				if nameSplice := strings.Split(goType, "."); len(nameSplice) > 1 {
					if pkg, isDep := dep[nameSplice[0]]; isDep {
						if FindStruct(pkg, nameSplice[1]) != nil {
							return "$" + pkg.Name + "." + nameSplice[1] + ".fromJson(response.data['data']['resp'])"
						}
					}
				}
			}
			return "response.data['data']['resp']"
		},
		"returnDefault": func(goType string) string {
			if strings.HasPrefix(goType, "*") {
				return "return null"
			}
			goType = strings.ReplaceAll(goType, "*", "")
			if typ := dartDefaultValue[dartTypeMap[goType]]; typ != "" {
				return "return " + typ
			}
			if strings.HasPrefix(goType, "[]") {
				return "return " + "[]"
			}
			if dep != nil {
				if nameSplice := strings.Split(goType, "."); len(nameSplice) > 1 {
					if pkg, isDep := dep[nameSplice[0]]; isDep {
						if FindStruct(pkg, nameSplice[1]) != nil {
							return "return  $" + pkg.Name + "." + nameSplice[1] + "()"
						}
					}
				}
			}
			return "return " + goType + "()"
		},
	}).Parse(dartModelTemplate))

	return t.Execute(f, pkg)
}

// dart template
var dartModelTemplate = `{{$name := .Name}}
{{range .Structs}}{{$struct := .}}
class {{.Name}} {
	{{range .Fields}}
	{{.Type | toType}} {{.Name | toCamel}}{{.Type | toDefault}};
	{{end}}
	{{range .Methods}}{{if .Comments}}
	{{range .Comments}}/// {{.}}
	{{end}}{{end}}Future<{{if .Rets}}{{(index .Rets 0).Type | toType}}{{else}}void{{end}}> {{.Name | toCamel}}(
		{{range .Args}}{{.Type | toType}} {{.Name | toCamel}},
		{{end}}
	) async {
		var response = await D.instance.dio.post('$modelUrlPrefix/{{$name}}/{{$struct.Name | toSnack}}/{{.Name | toMinus}}', data: {
			"data": toJson(),
			"args": { {{range .Args}}"{{.Name | toSnack}}": {{.Name | toCamel}},{{end}} }
		});
		if (response.data['code'] == 0) {
			var respModel = {{$struct.Name}}.fromJson(response.data['data']['data']);
			assign(respModel);
			{{if .Rets}}if (response.data['data']['resp'] != null) {
				return {{(index .Rets 0).Type | returnArgs}};
			} else {
				{{(index .Rets 0).Type | returnDefault}};
			}
			{{end}}
		}
		{{if .Rets}}{{(index .Rets 0).Type | returnDefault}};
		{{end}}
	}
	{{end}}
	{{.Name}}();

	assign({{.Name}} other) {
		{{range .Fields}}
		{{.Name | toCamel}} = other.{{.Name | toCamel}};
		{{end}}
	}

	Map<String, dynamic> toJson() {
		return {
			{{range .Fields}}
			"{{.Name | toSnack}}": {{.Name | toCamel}}{{asToJson .Type (.Name | toCamel)}},
			{{end}}
		};
	}
	{{.Name}}.fromJson(Map<String, dynamic> json) {
		{{range .Fields}}
		{{.Name | toCamel}} = json["{{.Name | toSnack}}"]{{fromJson .Type (.Name | toCamel)}};
		{{end}}
	}
}
{{end}}

`

// GenerateSwiftCode generate swift code
func GenerateSwiftCode(path string, pkg *Package) error {
	t := template.Must(template.New("model").Funcs(template.FuncMap{
		"toType": func(goType string) string {
			return swiftTypeMap[goType]
		},
		"toCamel": func(s string) string {
			return cjson.ToCamel(s)
		},
		"toDefault": func(s string) string {
			return swiftDefaultValue[swiftTypeMap[s]]
		},
	}).Parse(swiftModelTemplate))
	f, err := os.Create(filepath.Join(path, cjson.SnakeCaseToFullCamelCase(pkg.Name)+".swift"))
	if err != nil {
		return err
	}
	defer f.Close()
	return t.Execute(f, pkg)
}

// swift template
var swiftModelTemplate = `// Code generated by github.com/tinkler/mqttadmin; DO NOT EDIT.
{{range .Structs}}
struct {{.Name}} {
	{{range .Fields}}
	var {{.Name | toCamel}}: {{.Type | toType}};
	{{end}}
	{{range .Methods}}
	func {{.Name | toCamel}}({{range .Args}}{{.Name | toCamel}}: {{.Type | toType}}, {{end}}) -> {{if .Rets}}{{range .Rets}}{{.Type | toType}}, {{end}} {{else}}Void{{end}} {
		
	}
	{{end}}
}
{{end}}
`
