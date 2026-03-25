package models

type ItemStatistics struct {
	Likes     int `json:"likes"`
	ViewCount int `json:"viewCount"`
	Contacts  int `json:"contacts"`
}

type ItemResponse struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Price      int            `json:"price"`
	SellerID   int            `json:"sellerId"`
	Statistics ItemStatistics `json:"statistics"`
	CreatedAt  string         `json:"createdAt"`
}

// CreateItemResponse represents the actual server response for POST /api/1/item.
// Example: {"status":"Сохранили объявление - bfc8b374-ad8b-46ce-8a33-cde0c23e1039"}
type CreateItemResponse struct {
	Status string `json:"status"`
}

type StatisticResponse struct {
	ViewCount int `json:"viewCount"`
	Contacts  int `json:"contacts"`
	Likes     int `json:"likes"`
}

type ErrorResponse struct {
	Result struct {
		Message  string                 `json:"message"`
		Messages map[string]interface{} `json:"messages"`
	} `json:"result"`
	Status string `json:"status"`
}
