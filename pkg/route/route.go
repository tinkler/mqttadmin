package route

type Model[T any, S any] struct {
	Data T `json:",omitempty"`
	Args S `json:",omitempty"`
}

type Res[T any, S any] struct {
	Data T `json:",omitempty"`
	Resp S `json:",omitempty"`
}
