package util

// DataWrapper is a generic struct to wrap data in a response, used for swagger documentation
type DataWrapper[T any] struct {
	Data T `json:"data"`
}

// MessageWrapper wraps a message in a response, used for swagger documentation
type MessageWrapper struct {
	Message string `json:"message"`
}

type ListData[T any] struct {
	List T `json:"list"`
}

type ListDataWrapper[T any] struct {
	Data ListData[T] `json:"data"`
}

func WrapData[T any](data T) DataWrapper[T] {
	return DataWrapper[T]{
		Data: data,
	}
}

func WrapListData[T any](data T) ListDataWrapper[T] {
	return ListDataWrapper[T]{
		Data: ListData[T]{List: data},
	}
}

func WrapMessage(message string) MessageWrapper {
	return MessageWrapper{
		Message: message,
	}
}
