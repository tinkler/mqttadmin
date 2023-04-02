package route

type Model[T any, S any] struct {
	Data T `json:",omitempty"`
	Args S `json:",omitempty"`
}

type Res[T any] struct {
	Data T `json:",omitempty"`
}
