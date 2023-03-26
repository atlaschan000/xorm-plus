package xplus

import (
	"github.com/go-xorm/xorm"
	"github.com/jinzhu/copier"
)

type SessionFunc func(session *xorm.Session) error

type Page[T any] struct {
	Current   int
	Size      int
	Total     int64
	TotalPage int64
	Records   []*T
}

func Insert[T any](obj *T) (int64, error) {
	return dbEngine.Insert(obj)
}

func InsertBatch[T any](objs ...*T) (int64, error) {
	values := make([]interface{}, len(objs))
	for i, obj := range objs {
		values[i] = obj
	}
	return dbEngine.Insert(values...)
}

func Update[T any](obj *T, cols ...string) (int64, error) {
	return dbEngine.Cols(cols...).Update(obj)
}

func UpdateById[T any](id any, entity *T) (int64, error) {
	return dbEngine.Id(id).Update(entity)
}

func UpdateFields[T any](q *Query[T], entity *T) (int64, error) {
	return q.session.Update(entity)
}

func DeleteById[T any](id any) (int64, error) {
	var entity T
	return dbEngine.ID(id).Delete(&entity)
}

func DeleteByIds[T any](ids any) error {
	query := NewQuery[T]()
	query.In("id", ids)
	var table T
	_, err := query.session.Delete(table)
	return err
}

func Delete[T any](q *Query[T]) error {
	var table T
	_, err := q.session.Delete(table)
	return err
}

func Transaction(f SessionFunc) error {
	session := dbEngine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	if err := f(session); err != nil {
		_ = session.Rollback()
		return err
	}
	if err := session.Commit(); err != nil {
		return err
	}
	return nil
}

func GetById[T any](id any) (*T, bool, error) {
	var entity T
	has, err := dbEngine.ID(id).Get(&entity)
	return &entity, has, err
}

func SelectById[T any](id any) (*T, error) {
	query := NewQuery[T]()
	query.Eq("id", id)
	var entity T
	_, err := query.Get(&entity)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func SelectOne[T any](q *Query[T]) (*T, error) {
	var entity T
	_, err := q.Limit(1).Get(&entity)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func SelectOneModel[T any, R any](q *Query[T]) (*R, error) {
	var entity T
	if has, err := q.Limit(1).Get(&entity); err != nil {
		return nil, err
	} else if !has {
		return nil, nil
	}
	r := new(R)
	if err := copier.Copy(r, entity); err != nil {
		return nil, err
	}
	return r, nil
}

func SelectByIds[T any](ids any) ([]*T, error) {
	query := NewQuery[T]()
	query.In("id", ids)
	var entities []*T
	err := query.Find(entities)
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func SelectList[T any](q *Query[T]) ([]*T, error) {
	var entities []*T
	err := q.Find(entities)
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func SelectListModel[T any, R any](q *Query[T]) ([]*R, error) {
	var entity T
	rows, err := q.session.Rows(entity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]*R, 0)
	for rows.Next() {
		t := new(T)
		if err := rows.Scan(t); err != nil {
			return nil, err
		}
		r := new(R)
		_ = copier.Copy(r, t)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}

func SelectPage[T any](page *Page[T], q *Query[T]) (*Page[T], error) {

	countQuery := *q
	var entity T
	total, err := countQuery.Count(entity)
	if err != nil {
		return nil, err
	}
	page.Total = total
	page.TotalPage = (total + int64(page.Size) - 1) / int64(page.Size)
	var results []*T
	err = q.Paginate(page.Current, page.Size).Find(results)
	if err != nil {
		return nil, err
	}
	page.Records = results
	return page, nil
}

func SelectPageModel[T any, R any](page *Page[R], q *Query[T]) (*Page[R], error) {

	countQuery := *q
	var entity T
	total, err := countQuery.Count(entity)
	if err != nil {
		return nil, err
	}
	page.Total = total
	page.TotalPage = (total + int64(page.Size) - 1) / int64(page.Size)

	var results []*R
	q.Paginate(page.Current, page.Size)
	rows, err := q.session.Rows(entity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]*R, 0)
	for rows.Next() {
		t := new(T)
		if err := rows.Scan(t); err != nil {
			return nil, err
		}
		r := new(R)
		_ = copier.Copy(r, t)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	page.Records = results
	return page, nil
}

func SelectExist[T any](q *Query[T]) (bool, error) {
	return q.session.Exist()
}

func SelectSubExistOne[T any](q *Query[T], f func(sub *Query[T])) (*T, error) {
	var entity T
	_, err := q.Exists(f).Get(&entity)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func SelectSubExistList[T any](q *Query[T], f func(sub *Query[T])) ([]*T, error) {
	var entities []*T
	err := q.Exists(f).Find(entities)
	if err != nil {
		return nil, err
	}
	return entities, nil
}
