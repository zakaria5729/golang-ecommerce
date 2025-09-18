package review

type CreateReviewRequest struct {
	ProductID string `json:"product_id"`
	Rating    string `json:"rating"`
	Comment   string `json:"comment"`
}

type UpdateReviewRequest struct {
	ProductID string `json:"product_id"`
	Rating    string `json:"rating"`
	Comment   string `json:"comment"`
}
