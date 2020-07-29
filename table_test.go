package partoo_test

import (
	"github.com/byrnedo/partoo"
	"github.com/lib/pq"
	"reflect"
	"testing"
)

type baseModel struct {
	ID      string   `sql:"id"`
	Foo     string   `sql:"foo"`
	PQArray []string `sql:"pq_array"`
}

func (t baseModel) TableName() string {
	return "test"
}

func (t *baseModel) Columns() partoo.Cols {
	return partoo.Cols{
		&t.ID,
		&t.Foo,
		pq.Array(&t.PQArray),
	}
}

type manualIDModel struct {
	baseModel
}

func (m manualIDModel) AutoID() bool {
	return false
}

func TestColNames_Prefix(t *testing.T) {

	m := &baseModel{}
	p := partoo.New(partoo.Postgres)
	cols := p.NamedFields(m)
	aliased := cols.Names().Prefix("alias")
	if reflect.DeepEqual(aliased, partoo.ColNames{"alias.id", "alias.foo", "alias.pq_array"}) == false {
		t.Fatal("wrong aliases", aliased)
	}
}

func TestInsert(t *testing.T) {
	m := &baseModel{}
	p := partoo.New(partoo.Postgres)

	sqlStr, args := p.Insert(m)
	if sqlStr != "INSERT INTO test (foo,pq_array) VALUES ($1,$2)" {
		t.Fatal(sqlStr)
	}
	if len(args) != 2 {
		t.Fatal(len(args))
	}

	manualModel := &manualIDModel{}
	sqlStr, args = p.Insert(manualModel)
	if sqlStr != "INSERT INTO test (id,foo,pq_array) VALUES ($1,$2,$3)" {
		t.Fatal(sqlStr)
	}
	if len(args) != 3 {
		t.Fatal(len(args))
	}
}

func TestUpdate(t *testing.T) {
	m := &baseModel{}
	p := partoo.New(partoo.Postgres)

	sqlStr, args := p.Update(m)
	if sqlStr != "UPDATE test SET foo = $1,pq_array = $2" {
		t.Fatal(sqlStr)
	}
	if len(args) != 2 {
		t.Fatal(len(args))
	}
	t.Log(sqlStr)
}

func TestUpdateOne(t *testing.T) {
	m := &baseModel{}
	p := partoo.New(partoo.Postgres)

	sqlStr, args := p.UpdateOne(m)
	if sqlStr != "UPDATE test SET foo = $1,pq_array = $2 WHERE id = $3" {
		t.Fatal(sqlStr)
	}
	if len(args) != 3 {
		t.Fatal(len(args))
	}
	t.Log(sqlStr)
}

func TestPartoo_UpsertOne(t *testing.T) {
	m := &baseModel{}
	p := partoo.New(partoo.Postgres)
	sqlStr, args := p.UpsertOne(m)
	if sqlStr != "INSERT INTO test (foo,pq_array) VALUES ($1,$2) ON CONFLICT (id) DO UPDATE SET foo = $3,pq_array = $4" {
		t.Fatal(sqlStr)
	}
	if len(args) != 4 {
		t.Fatal(len(args))
	}
}
