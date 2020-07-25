package model

var PaginationPath = []string{"commons", "app", "model", "pagination.go"}

var PaginationContent = `package model

import (
	"{{ .AppRepository }}/commons/config"
  "{{ .AppRepository }}/commons/log"
	"strconv"
)

func (q *Query) Pagination(page, perPage string) (int, int, int) {
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

func handleCurrentPage(page string) int {
	var currentPage int
	var err error

	currentPage, err = strconv.Atoi(page)
	if err != nil {
		currentPage = 1
	}

	return currentPage
}

func handleEntriesPerPage(perPage string) int {
	var entriesPerPage int
	var err error

	entriesPerPage, err = strconv.Atoi(perPage)
	if err != nil {
		entriesPerPage = config.App.Pagination.Default
	}

	if entriesPerPage > config.App.Pagination.Maximum {
		entriesPerPage = config.App.Pagination.Maximum
		log.Warn.Printf("Maximum value for entries per page can't be greater than %d", config.App.Pagination.Maximum)
	}

	return entriesPerPage
}`
