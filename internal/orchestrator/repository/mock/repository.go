package mock

import (
	"cmp"
	"github.com/xLeSHka/calc/internal/models"
	"github.com/xLeSHka/calc/internal/orchestrator/repository"
	"github.com/xLeSHka/calc/internal/pkg/counter"
	"slices"
	"sync"
)

type Repository struct {
	Expressions map[int64]*models.Expression
	Counter     *counter.Counter
	mu          *sync.RWMutex
}

func New(
	counter *counter.Counter,
) *Repository {
	return &Repository{
		Expressions: make(map[int64]*models.Expression),
		Counter:     counter,
		mu:          &sync.RWMutex{},
	}
}

func (r *Repository) CreateExpression(expression string) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := r.Counter.ExprInc()
	expr := &models.Expression{
		ID:         id,
		Expression: expression,
		Status:     "Waiting",
	}
	r.Expressions[id] = expr
	return id, nil
}
func (r *Repository) GetExpression(id int64) (*models.Expression, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	expr, ok := r.Expressions[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return expr, nil
}
func (r *Repository) GetExpressions(size, page int) ([]*models.Expression, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	keys := make([]int64, 0, len(r.Expressions))
	for k := range r.Expressions {
		keys = append(keys, k)
	}
	slices.SortFunc(keys, func(a, b int64) int {
		return -cmp.Compare(a, b)
	})
	exprs := make([]*models.Expression, 0, size)
	for i := page * size; (i-page*size) < size && i < len(r.Expressions); i++ {
		exprs = append(exprs, r.Expressions[keys[i]])
	}
	if len(exprs) == 0 {
		return nil, 0, repository.ErrNotFound
	}
	return exprs, int64(len(r.Expressions)), nil
}
func (r *Repository) UpdateExpression(expression *models.Expression) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	expr, ok := r.Expressions[expression.ID]
	if !ok {
		return repository.ErrNotFound
	}
	expression.Expression = expr.Expression
	r.Expressions[expression.ID] = expression
	return nil
}
func (r *Repository) SetResult(id int64, result float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	expr, ok := r.Expressions[id]
	if !ok {
		return repository.ErrNotFound
	}
	expr.Result = &result
	return nil
}
