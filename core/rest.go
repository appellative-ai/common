package core

import (
	"fmt"
	"net/http"
	"reflect"
)

// Micro-REST

// Exchange - http exchange
type Exchange func(r *http.Request) (*http.Response, error)

// ExchangeLink - interface to link http Exchanges. Used in the collective exchange
type ExchangeLink func(next Exchange) Exchange

// Chainable - interface to create a link
type Chainable[T any] interface {
	Link(t T) T
}

// BuildNetwork - build network
func BuildNetwork(operatives []any) Exchange {
	return buildNetwork[Exchange, Chainable[Exchange]](operatives)
}

// buildNetwork - build a chain of links - panic on nil or invalid type links
func buildNetwork[T any, U Chainable[T]](operatives []any) (head T) {
	if len(operatives) == 0 {
		panic("operatives list is nil")
	}
	for i := len(operatives) - 1; i >= 0; i-- {
		if operatives[i] == nil {
			panic(fmt.Sprintf("operative is nil at index: %v", i))
		}
		// Check for a next function
		if fn, ok := operatives[i].(func(next T) T); ok {
			head = fn(head)
			continue
		}
		// Check for a Chainable interface
		if c, ok := operatives[i].(U); ok {
			head = c.Link(head)
			continue
		}
		panic(fmt.Sprintf("invalid operative type: %v", reflect.TypeOf(operatives[i])))
	}
	return head
}
