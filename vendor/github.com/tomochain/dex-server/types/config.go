package types

type KeyValue struct {
	Key   string      `json:"key" bson"key"`
	Value interface{} `json:"value" bson "value"`
}
