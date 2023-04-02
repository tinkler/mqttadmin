package jsonz

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	// https://github.com/golang/lint/blob/master/lint.go#L770
	CommonInitialisms         = []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "TTL", "UID", "UI", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}
	CommonInitialismsReplacer *strings.Replacer
)

func init() {
	commonInitialismsForReplacer := make([]string, 0, len(CommonInitialisms))
	for _, initialism := range CommonInitialisms {
		commonInitialismsForReplacer = append(commonInitialismsForReplacer, initialism, cases.Title(language.Chinese).String(strings.ToLower(initialism)))
	}

	CommonInitialismsReplacer = strings.NewReplacer(commonInitialismsForReplacer...)
}
