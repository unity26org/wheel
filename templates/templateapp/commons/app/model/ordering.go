package model

var OrderingPath = []string{"commons", "app", "model", "ordering.go"}

var OrderingContent = `package model

import ()

func (q *Query) Ordering(order string) {
	if order != "" {
		q.Db = q.Db.Order(order)
	}
}`
