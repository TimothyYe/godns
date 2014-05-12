package main

type status struct {
	code       string
	message    string
	created_at string
}

type domain struct {
	id      int
	name    string
	status  string
	records string
	owner   string
}

type domain_list struct {
	ret_status status
	domains    []domain
}
