package sjson

import (
	jsoniter "github.com/json-iterator/go"
)

type snakedNamedExtension struct {
	jsoniter.DummyExtension
}

func (ex *snakedNamedExtension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {
	for _, field := range structDescriptor.Fields {
		names := make([]string, len(field.ToNames))
		for i, old := range field.ToNames {
			names[i] = ToSnackedName(old)
		}
		field.ToNames = names
		field.FromNames = names
	}
}
