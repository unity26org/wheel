package model

var PaginationPath = []string{"commons", "app", "model", "pagination.go"}

var PaginationContent = `package model

import (
	"strconv"
)

func (q *Query) Pagination(page interface{}, perPage interface{}) (int, int, int) {
	type counter struct {
		Entries int
	}

	var currentPage, totalPages, entriesPerPage int
	var result counter

	currentPage = handleCurrentPage(page)
	entriesPerPage = handleEntriesPerPage(perPage)

	q.Db.Table(TableName(q.Table)).Order("", true).Select("COUNT(*) AS entries").Scan(&result)

	totalPages = result.Entries / entriesPerPage
	if (result.Entries % entriesPerPage) > 0 {
		totalPages++
	}

	offset := (currentPage - 1) * entriesPerPage
	q.Db = q.Db.Offset(offset).Limit(entriesPerPage)

	return currentPage, totalPages, result.Entries
}

func handleCurrentPage(page interface{}) int {
	var currentPage int
	var err error

	switch auxPage := page.(type) {
	case int:
		currentPage = auxPage
	case string:
		currentPage, err = strconv.Atoi(auxPage)
		if err != nil {
			currentPage = 1
		}
	default:
		currentPage = 1
	}

	return currentPage
}

func handleEntriesPerPage(perPage interface{}) int {
	var entriesPerPage int
	var err error

	switch auxPerPage := perPage.(type) {
	case int:
		entriesPerPage = auxPerPage
	case string:
		entriesPerPage, err = strconv.Atoi(auxPerPage)
		if err != nil {
			entriesPerPage = 20
		}
	default:
		entriesPerPage = 20
	}

	return entriesPerPage
}`
