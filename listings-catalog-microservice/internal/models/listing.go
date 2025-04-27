package models

type Listing struct {
	ID          int64
	Title       string
	Description string
	Quantity    int64
	Category    string
	Closed      bool
	Price       int64
	Creator     int64
}
