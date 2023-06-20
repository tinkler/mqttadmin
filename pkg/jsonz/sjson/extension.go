package sjson

import (
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type snakedNamedExtension struct {
	jsoniter.DummyExtension
}

func (ex *snakedNamedExtension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {
	for _, field := range structDescriptor.Fields {
		names := make([]string, len(field.ToNames))
		tag, _, _ := strings.Cut(field.Field.Tag().Get("json"), ",")
		for i, old := range field.ToNames {
			if tag != "" {
				names[i] = tag
			} else {
				names[i] = ToSnackedName(old)
			}
		}
		field.ToNames = names
		field.FromNames = names
	}
}
