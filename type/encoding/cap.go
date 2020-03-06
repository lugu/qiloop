package encoding

type Capability interface {
	MetaObjectCache() bool
	ObjectPtrUID() bool
}

type defaultCap struct{}

func (d defaultCap) MetaObjectCache() bool {
	return false
}

func (d defaultCap) ObjectPtrUID() bool {
	return false
}

func DefaultCap() Capability {
	return defaultCap{}
}

