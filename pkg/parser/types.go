/*
go type
*/
package parser

// tsTypeMap go type map to typescript type
var tsTypeMap = map[string]string{
	"string":  "string",
	"int":     "number",
	"int8":    "number",
	"int16":   "number",
	"int32":   "number",
	"int64":   "number",
	"uint":    "number",
	"uint8":   "number",
	"uint16":  "number",
	"uint32":  "number",
	"uint64":  "number",
	"float32": "number",
	"float64": "number",
	"bool":    "boolean",
	"byte":    "string",
	"rune":    "string",
	"error":   "string",
}

// Typescript default value
var tsDefaultValue = map[string]string{
	"string":  `""`,
	"number":  `0`,
	"boolean": `false`,
}

// dartTypeMap go type map to dart type
var dartTypeMap = map[string]string{
	"string":  "String",
	"int":     "int",
	"int8":    "int",
	"int16":   "int",
	"int32":   "int",
	"int64":   "int",
	"uint":    "int",
	"uint8":   "int",
	"uint16":  "int",
	"uint32":  "int",
	"uint64":  "int",
	"float32": "double",
	"float64": "double",
	"bool":    "bool",
	"byte":    "String",
	"rune":    "String",
	"error":   "String",
}

// Dart default value
var dartDefaultValue = map[string]string{
	"String":  `""`,
	"int":     `0`,
	"double":  `0.0`,
	"bool":    `false`,
	"dynamic": `null`,
}

// swiftTypeMap go type map to swift type
var swiftTypeMap = map[string]string{
	"string":  "String",
	"int":     "Int",
	"int8":    "Int8",
	"int16":   "Int16",
	"int32":   "Int32",
	"int64":   "Int64",
	"uint":    "UInt",
	"uint8":   "UInt8",
	"uint16":  "UInt16",
	"uint32":  "UInt32",
	"uint64":  "UInt64",
	"float32": "Float",
	"float64": "Double",
	"bool":    "Bool",
	"byte":    "String",
	"rune":    "String",
	"error":   "String",
}

// Swift default value
var swiftDefaultValue = map[string]string{
	"String": `""`,
	"Int":    `0`,
	"Int8":   `0`,
	"Int16":  `0`,
	"Int32":  `0`,
	"Int64":  `0`,
	"UInt":   `0`,
	"UInt8":  `0`,
	"UInt16": `0`,
	"UInt32": `0`,
	"UInt64": `0`,
	"Float":  `0.0`,
	"Double": `0.0`,
	"Bool":   `false`,
}
