package shared

import "encoding/json"

type Nullable[T any] struct {
	Value  T
	IsNull bool
}

func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if n.IsNull {
		return []byte("null"), nil
	}
	return json.Marshal(n.Value)
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.IsNull = true
		var v T
		n.Value = v
		return nil
	}

	err := json.Unmarshal(data, &n.Value)
	if err != nil {
		return err
	}
	n.IsNull = false
	return nil
}
