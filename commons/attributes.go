package commons

type Attributes interface {
	SetAttr(name string, value interface{}) (replace bool)
	GetAttr(name string) (value interface{}, has bool)
}

//给对象添加属性
type attributes map[string]interface{}

func (this attributes) SetAttr(name string, value interface{}) (replace bool) {
	_, replace = this[name]
	(this)[name] = value
	return
}

func (this attributes) GetAttr(name string) (value interface{}, has bool) {
	value, has = this[name]
	return
}

func NewAttributes() Attributes {
	return attributes(map[string]interface{}{})
}
