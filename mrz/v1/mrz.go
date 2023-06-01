package mrz_v1

import (
	reflect "reflect"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

type protoreflectMessage interface {
	ProtoReflect() protoreflect.Message
	Reset()
}

type typedModel[T protoreflectMessage, S protoreflectMessage] struct {
	Data T
	Args S
}

func NewTypedModel[T protoreflectMessage, S protoreflectMessage]() *typedModel[T, S] {
	return &typedModel[T, S]{}
}

func NewDataTypedModel[T protoreflectMessage]() *typedModel[T, *anypb.Any] {
	return NewTypedModel[T, *anypb.Any]()
}

func ToTypedModel[T protoreflectMessage, S protoreflectMessage](data *anypb.Any) *typedModel[T, S] {
	m := new(Model)
	if err := data.UnmarshalTo(m); err != nil {
		panic(err)
	}
	tm := typedModel[T, S]{}
	if f := m.GetData(); f != nil {
		var t T
		t = reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)
		if err := f.UnmarshalTo(t); err != nil {
			panic(err)
		}
		tm.Data = t
	}
	if f := m.GetArgs(); f != nil {
		var s S
		s = reflect.New(reflect.TypeOf(s).Elem()).Interface().(S)
		if err := f.UnmarshalTo(s); err != nil {
			panic(err)
		}
		tm.Args = s
	}
	return &tm
}

func ToDataTypedModel[T protoreflectMessage](data *anypb.Any) *typedModel[T, *anypb.Any] {
	return ToTypedModel[T, *anypb.Any](data)
}

func (m *typedModel[T, S]) ToAny() *anypb.Any {
	model := &Model{}
	var err error
	model.Data, err = anypb.New(m.Data)
	if err != nil {
		panic(err)
	}
	model.Args, err = anypb.New(m.Args)
	if err != nil {
		panic(err)
	}
	any, err := anypb.New(model)
	if err != nil {
		panic(err)
	}
	return any
}

type typedRes[T protoreflectMessage, S protoreflectMessage] struct {
	Data T
	Resp S
}

func NewTypedRes[T protoreflectMessage, S protoreflectMessage]() *typedRes[T, S] {
	return &typedRes[T, S]{}
}

func NewDataTypedRes[T protoreflectMessage]() *typedRes[T, *anypb.Any] {
	return NewTypedRes[T, *anypb.Any]()
}

func ToTypedRes[T protoreflectMessage, S protoreflectMessage](any *anypb.Any) *typedRes[T, S] {
	m := new(Res)
	err := any.UnmarshalTo(
		m,
	)
	if err != nil {
		panic(err)
	}
	tm := typedRes[T, S]{}
	if f := m.GetData(); f != nil {
		var t T
		if rt := reflect.TypeOf(t); rt.Kind() == reflect.Ptr {
			t = reflect.New(rt.Elem()).Interface().(T)
		}
		if err := f.UnmarshalTo(t); err != nil {
			panic(err)
		}
		tm.Data = t
	}
	if f := m.GetResp(); f != nil {
		var s S
		if rs := reflect.TypeOf(s); rs.Kind() == reflect.Ptr {
			s = reflect.New(rs.Elem()).Interface().(S)
		}
		if err := f.UnmarshalTo(s); err != nil {
			panic(err)
		}
		tm.Resp = s
	}
	return &tm
}

func ToDataTypedRes[T protoreflectMessage](any *anypb.Any) *typedRes[T, *anypb.Any] {
	return ToTypedRes[T, *anypb.Any](any)
}

func (m *typedRes[T, S]) ToAny() *anypb.Any {
	res := &Res{}
	var err error
	res.Data, err = anypb.New(m.Data)
	if err != nil {
		panic(err)
	}
	res.Resp, err = anypb.New(m.Resp)
	if err != nil {
		panic(err)
	}
	any, err := anypb.New(res)
	if err != nil {
		panic(err)
	}
	return any
}
