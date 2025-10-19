package model

type CreateReviewRequest struct {
	ProductID uint   `json:"product_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
}

type UpdateReviewRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}
