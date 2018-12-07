package idl

import (
	"fmt"
	"github.com/lugu/qiloop/meta/signature"
	"strings"
)

// Scope represents the scope of a type in which type declarations are
// resolved.
type Scope interface {
	Add(name string, typ signature.Type) error
	Search(name string) (signature.Type, error)
	Extend(namespace string, other Scope) error
}

type scopeImpl struct {
	local  map[string]signature.Type
	global map[string]Scope
}

// NewScope returns a new scope.
func NewScope() Scope {
	return &scopeImpl{
		local:  make(map[string]signature.Type),
		global: make(map[string]Scope),
	}
}
func (s *scopeImpl) Add(name string, typ signature.Type) error {
	_, ok := s.local[name]
	if !ok {
		s.local[name] = typ
		return nil
	}
	return fmt.Errorf("already present in scope: %s", name)
}

func (s *scopeImpl) searchLocal(name string) (signature.Type, error) {
	t, ok := s.local[name]
	if ok {
		return t, nil
	}
	return nil, fmt.Errorf("not found in scope: %s", name)
}

func (s *scopeImpl) searchGlobal(namespace, name string) (signature.Type, error) {
	scope, ok := s.global[namespace]
	if ok {
		t, err := scope.Search(name)
		if err != nil {
			return nil, fmt.Errorf("can not find %s in %s (%s)",
				name, namespace, err)
		}
		return t, nil
	}
	return nil, fmt.Errorf("unknown namespace: %s", namespace)
}

func (s *scopeImpl) Search(name string) (signature.Type, error) {
	split := strings.SplitN(name, ".", 2)
	if len(split) > 1 {
		return s.searchGlobal(split[0], split[1])
	}
	return s.searchLocal(name)
}

func (s *scopeImpl) Extend(namespace string, other Scope) error {
	_, ok := s.global[namespace]
	if !ok {
		s.global[namespace] = other
		return nil
	}
	return fmt.Errorf("duplicated namespace: %s", namespace)
}
