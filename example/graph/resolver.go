package graph

import (
	"sync"

	"github.com/danutavadanei/gqlgen-validate/example/graph/model"
)

type Resolver struct {
	mu     sync.Mutex
	users  []*model.User
	nextID int
}
