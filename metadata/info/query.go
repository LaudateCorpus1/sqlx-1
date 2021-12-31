package info

import (
	"fmt"
	"github.com/viant/sqlx/metadata/database"
)

//Query represents dialect metadata queries
type (
	Query struct {
		Kind     Kind
		SQL      string
		Criteria Criteria
		database.Product
	}
	Criterion struct {
		Name   string
		Column string
	}
	Criteria []*Criterion

	Queries []*Query
)

//NewQuery creates a new query
func NewQuery(kind Kind, SQL string, info database.Product, criteria ...*Criterion) *Query {
	return &Query{
		Kind:     kind,
		SQL:      SQL,
		Product:  info,
		Criteria: criteria,
	}
}

func (c Criteria) Supported() int {
	supported := 0
	for _, item := range c {
		if item.Column != "" {
			supported++
		}
	}
	return supported
}

func (c Criteria) Validate(kind Kind) error {
	criteria := kind.Criteria()
	if len(c) != len(criteria) {
		return fmt.Errorf("invalid query '%v': expected %v criteria, but query defined %v", kind, len(criteria), len(c))
	}
	for i, item := range c {
		if item.Name != criteria[i] {
			return fmt.Errorf("invalid query criterion '%v': expected %v, but had %v", kind, item.Name, criteria[i])
		}
	}
	return nil
}

func (q Queries) Len() int {
	return len(q)
}

// Swap is part of sort.Interface.
func (q Queries) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

// Less is part of sort.Interface.
func (q Queries) Less(i, j int) bool {
	return q[i].Product.Major < q[j].Product.Major && q[i].Product.Minor < q[j].Product.Minor
}

//Match matches queries for version, or latest version
func (q Queries) Match(info *database.Product) *Query {
	switch len(q) {
	case 0:
		return nil
	case 1:
		return q[0]
	}
	for _, candidate := range q {
		if candidate.Product.Major >= info.Major {
			if candidate.Product.Minor >= info.Minor {
				return candidate
			}
		}
	}
	//by default return the latest version
	return q[len(q)-1]
}

//NewCriterion creates a new criteria, name refers to kind.Crtiera, column to local vendor column, use '?' for already defined placeholder, %v for substitution
func NewCriterion(name, column string) *Criterion {
	return &Criterion{
		Name:   name,
		Column: column,
	}
}
