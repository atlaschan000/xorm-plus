package tests

import (
	"fmt"
	"github.com/atlaschan000/xorm-plus/xplus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

var db *xorm.Engine

func init() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=true"
	if database, err := xorm.NewEngine("mysql", dsn); err != nil {
		panic(err)
	} else {
		db = database
	}
	xplus.Init(db)
}

type Test struct {
	Id       int64
	ParentId int64
}

type T2 struct {
	ParentId int64
}

func TestSelectOne(t *testing.T) {

	list, err := xplus.SelectPageModel[Test, T2](&xplus.Page[T2]{Current: 1, Size: 2}, xplus.NewQuery[Test]())
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Equal(t, list.TotalPage, int64(4))
}

func TestInsert(t *testing.T) {

	item := &Test{ParentId: 6}
	insert, err := xplus.Insert[Test](item)
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Equal(t, insert, int64(1))
}

func TestDeleteByIds(t *testing.T) {

	ids := []int64{1, 2, 3}
	err := xplus.DeleteByIds[Test](ids)
	assert.Equal(t, err, nil)
}

func TestUpdateField(t *testing.T) {

	q := xplus.NewQuery[Test]()
	q.Eq("id", 8).SetExpr("value", "value+1").Cols("parent_id")
	entity := &Test{ParentId: 6}
	rows, err := xplus.UpdateFields[Test](q, entity)
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Equal(t, rows, int64(1))

}

func TestUpdateEntity(t *testing.T) {

	entity, exist, err := xplus.GetById[Test](6)
	if err != nil {
		fmt.Println(err.Error())
	}
	if !exist {
		fmt.Println("record not exist")
	}
	entity.ParentId = 4
	rows, err := xplus.UpdateById(entity.Id, entity)
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Equal(t, rows, int64(0))
}

func TestTransaction(t *testing.T) {

	err := xplus.Transaction(func(session *xorm.Session) error {
		entity, exist, err := xplus.GetById[Test](6)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		if !exist {
			fmt.Println("record not exist")
		}
		entity.ParentId = 11
		_, err = xplus.UpdateById(entity.Id, entity)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		return nil
	})

	assert.Equal(t, err, nil)
}
