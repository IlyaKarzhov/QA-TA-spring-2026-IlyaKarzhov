package models

type Statistics struct {
	Likes     int `json:"likes"`
	ViewCount int `json:"viewCount"`
	Contacts  int `json:"contacts"`
}

type CreateItemRequest struct {
	Name       string     `json:"name"`
	Price      int        `json:"price"`
	SellerID   int        `json:"sellerID"`
	Statistics Statistics `json:"statistics"`
}

// CreateItemRequestWithoutPrice is used to test missing price field validation.
type CreateItemRequestWithoutPrice struct {
	Name       string     `json:"name"`
	SellerID   int        `json:"sellerID"`
	Statistics Statistics `json:"statistics"`
}
