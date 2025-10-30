package model

type PostNewsLetterRequest struct {
	Email string `json:"email"`
}

type CreateContactInfoRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Message string `json:"message"`
}
