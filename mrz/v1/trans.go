package mrz_v1

import (
	"reflect"

	"github.com/tinkler/mqttadmin/pkg/jsonz/sjson"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

func GetData[T any, S protoreflectMessage, J protoreflectMessage](m *typedModel[S, J]) T {
	var t T
	if m == nil {
		return t
	}
	jd, err := sjson.Marshal(m.Data)
	if err != nil {
		panic(err)
	}
	isPtr := false
	if rs := reflect.TypeOf(t); rs.Kind() == reflect.Ptr {
		isPtr = true
		t = reflect.New(rs.Elem()).Interface().(T)
	}
	if isPtr {
		if err := sjson.Unmarshal(jd, t); err != nil {
			panic(err)
		}
	} else {
		if err := sjson.Unmarshal(jd, &t); err != nil {
			panic(err)
		}
	}
	return t
}

func GetArgs[T any, S protoreflectMessage](m *typedModel[S, *structpb.Struct], name string) T {
	var t T
	if m == nil {
		return t
	}
	if m.Args == nil {
		return t
	}
	jd, err := sjson.Marshal(m.Args.Fields[name].AsInterface())
	if err != nil {
		panic(err)
	}
	isPtr := false
	if rs := reflect.TypeOf(t); rs.Kind() == reflect.Ptr {
		isPtr = true
		t = reflect.New(rs.Elem()).Interface().(T)
	}
	if isPtr {
		if err := sjson.Unmarshal(jd, t); err != nil {
			panic(err)
		}
	} else {
		if err := sjson.Unmarshal(jd, &t); err != nil {
			panic(err)
		}
	}
	return t
}

func SetResData[T any, S protoreflectMessage, J protoreflectMessage](m *typedRes[S, J], t T) {
	if m == nil {
		return
	}
	jd, err := sjson.Marshal(t)
	if err != nil {
		panic(err)
	}
	if rv := reflect.ValueOf(m.Data); rv.IsNil() {
		m.Data = reflect.New(rv.Type().Elem()).Interface().(S)
	}
	if err := sjson.Unmarshal(jd, m.Data); err != nil {
		panic(err)
	}
}

func SetResResp[T any, S protoreflectMessage](m *typedRes[S, *structpb.Value], t T) {
	if m == nil {
		return
	}
	if m.Resp == nil {
		m.Resp, _ = structpb.NewValue(t)
	}
}

func SetResRespList[T any, S protoreflectMessage](m *typedRes[S, *structpb.Value], t []T) {
	if m == nil {
		return
	}
	ts := make([]interface{}, len(t))
	for i, v := range t {
		ts[i] = v
	}
	if m.Resp == nil {
		m.Resp, _ = structpb.NewValue(ts)
	}

}

func SetResRespNil[S protoreflectMessage](m *typedRes[S, *structpb.Value]) {
	if m == nil {
		return
	}
	if m.Resp == nil {
		m.Resp, _ = structpb.NewValue(nil)
	}
}
