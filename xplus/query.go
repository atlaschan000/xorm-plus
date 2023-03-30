package xplus

import (
	"xorm.io/builder"
	"xorm.io/xorm"
)

type Query[T any] struct {
	engine  *xorm.Engine
	session *xorm.Session
}

func NewQuery[T any]() *Query[T] {
	return &Query[T]{
		engine:  dbEngine,
		session: dbEngine.NewSession(),
	}
}

func (q *Query[T]) Where(condition string, args ...interface{}) *Query[T] {
	q.session = q.session.Where(condition, args...)
	return q
}

func (q *Query[T]) Select(columns string) *Query[T] {
	q.session = q.session.Select(columns)
	return q
}

func (q *Query[T]) OrderBy(order string) *Query[T] {
	q.session = q.session.OrderBy(order)
	return q
}

func (q *Query[T]) Skip(count int) *Query[T] {
	q.session = q.session.Limit(-1, count)
	return q
}

func (q *Query[T]) Limit(count int) *Query[T] {
	q.session = q.session.Limit(count)
	return q
}

func (q *Query[T]) GroupBy(keys string) *Query[T] {
	q.session = q.session.GroupBy(keys)
	return q
}

func (q *Query[T]) Having(condition string) *Query[T] {
	q.session = q.session.Having(condition)
	return q
}

func (q *Query[T]) Join(joinType string, tableName interface{}, condition string, args ...interface{}) *Query[T] {
	q.session = q.session.Join(joinType, tableName, condition, args...)
	return q
}

func (q *Query[T]) Alias(alias string) *Query[T] {
	q.session = q.session.Alias(alias)
	return q
}

func (q *Query[T]) Eq(column string, value interface{}) *Query[T] {
	q.session.And(column+" = ?", value)
	return q
}

func (q *Query[T]) Ne(column string, value interface{}) *Query[T] {
	q.session.And(column+" != ?", value)
	return q
}

func (q *Query[T]) AllEq(m map[string]interface{}, isAnd ...bool) *Query[T] {
	var op string
	if len(isAnd) > 0 && !isAnd[0] {
		op = " or "
	} else {
		op = " and "
	}

	for k, v := range m {
		if v == nil {
			q.session = q.session.Where(k + " is null" + op)
		} else {
			q.session = q.session.Where(k+"=?", v).And(op)
		}
	}

	return q
}

func (q *Query[T]) Gt(column string, value interface{}) *Query[T] {
	q.session.And(column+" > ?", value)
	return q
}

func (q *Query[T]) Ge(column string, value interface{}) *Query[T] {
	q.session.And(column+" >= ?", value)
	return q
}

func (q *Query[T]) Lt(column string, value interface{}) *Query[T] {
	q.session.And(column+" < ?", value)
	return q
}

func (q *Query[T]) Le(column string, value interface{}) *Query[T] {
	q.session.And(column+" <= ?", value)
	return q
}

func (q *Query[T]) Like(column string, value interface{}) *Query[T] {
	q.session.And(column+" LIKE ?", value)
	return q
}

func (q *Query[T]) NotLike(column string, value interface{}) *Query[T] {
	q.session.And(column+" NOT LIKE ?", value)
	return q
}

func (q *Query[T]) LikeLeft(column string, value interface{}) *Query[T] {
	q.session.And(column+" LIKE ?", "%"+value.(string))
	return q
}

func (q *Query[T]) NotLikeLeft(column string, pattern string) *Query[T] {
	q.session = q.session.Where(column+" not like ?", "%"+pattern)
	return q
}

func (q *Query[T]) LikeRight(column string, value interface{}) *Query[T] {
	q.session.And(column+" LIKE ?", value.(string)+"%")
	return q
}

func (q *Query[T]) NotLikeRight(column string, pattern string) *Query[T] {
	q.session = q.session.Where(column+" not like ?", pattern+"%")
	return q
}

func (q *Query[T]) IsNull(column string) *Query[T] {
	q.session.And(column + " IS NULL")
	return q
}

func (q *Query[T]) IsNotNull(column string) *Query[T] {
	q.session.And(column + " IS NOT NULL")
	return q
}

func (q *Query[T]) In(column string, values interface{}) *Query[T] {
	q.session.In(column, values)
	return q
}

func (q *Query[T]) InBuilder(column string, builder *builder.Builder) *Query[T] {
	q.session = q.session.In(column, builder)
	return q
}

func (q *Query[T]) NotIn(column string, values interface{}) *Query[T] {
	q.session = q.session.NotIn(column, values)
	return q
}

func (q *Query[T]) NotInBuilder(column string, builder *builder.Builder) *Query[T] {
	q.session = q.session.NotIn(column, builder)
	return q
}

func (q *Query[T]) Between(column string, start, end interface{}) *Query[T] {
	q.session = q.session.Where(column+" between ? and ?", start, end)
	return q
}

func (q *Query[T]) NotBetween(column string, start, end interface{}) *Query[T] {
	q.session = q.session.Where(column+" not between ? and ?", start, end)
	return q
}

func (q *Query[T]) And(f func(q *Query[T])) *Query[T] {
	newQuery := NewQuery[T]()
	f(newQuery)
	q.session.And(newQuery.session.Conds())
	return q
}

func (q *Query[T]) Or(f func(q *Query[T])) *Query[T] {
	newQuery := NewQuery[T]()
	f(newQuery)
	q.session.Or(newQuery.session.Conds())
	return q
}

func (q *Query[T]) Func(f func(q *Query[T]) *Query[T]) *Query[T] {
	return f(q)
}

func (q *Query[T]) Paginate(page, pageSize int) *Query[T] {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	start := (page - 1) * pageSize
	q.session = q.session.Limit(pageSize, start)
	return q
}

func (q *Query[T]) Cols(columns ...string) *Query[T] {
	q.session = q.session.Cols(columns...)
	return q
}

func (q *Query[T]) SetExpr(column string, expression string) *Query[T] {
	q.session = q.session.SetExpr(column, expression)
	return q
}

func (q *Query[T]) Find(beans []*T) ([]*T, error) {
	err := q.session.Find(&beans)
	return beans, err
}

func (q *Query[T]) Get(bean *T) (bool, error) {
	return q.session.Get(bean)
}

func (q *Query[T]) Count(bean T) (int64, error) {
	return q.session.Count(bean)
}

func (q *Query[T]) Exists(f func(q *Query[T])) *Query[T] {
	subQuery := q.session.Select("1")
	f(&Query[T]{session: subQuery})
	q.session = q.session.Where("EXISTS (?)", subQuery)
	return q
}

func (q *Query[T]) sumInt(bean T,fields string) (int64,error) {
	return q.session.SumInt(bean,fields)
}

func (q *Query[T]) sum(bean T,fields string) (float64,error) {
	return q.session.Sum(bean,fields)
}
