package models

type Filters struct {
	Name string
}

type Query struct {
	Offset int64
	Limit  int64
	Filters
}
