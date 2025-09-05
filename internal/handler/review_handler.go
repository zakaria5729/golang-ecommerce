package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/easy-comerce/backend/internal/feature/review"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
)

type ReviewHandler struct {
	useCase *review.ReviewUseCase
}

func NewReviewHandler() *ReviewHandler {
	return &ReviewHandler{
		useCase: review.NewReviewUseCase(),
	}
}

func (h *ReviewHandler) GetAllReviews(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	includeStr := q.Get(constants.Include)
	productIDFilter := q.Get(review.ReviewProductID)
	userIDFilter := q.Get(review.ReviewUserID)
	ratingFilter := q.Get(review.ReviewRating)
	sortBy := q.Get(constants.SortBy)
	sortOrder := q.Get(constants.SortOrder)

	reviews, err := h.useCase.GetAllReviews(includeStr, productIDFilter, userIDFilter, ratingFilter, sortBy, sortOrder)
	if err != nil {
		response.SendErrorJSON(w, "Failed to fetch reviews", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, reviews)
}

func (h *ReviewHandler) GetAllReviewsPaginated(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	includeStr := q.Get(constants.Include)
	pageStr := q.Get(constants.Page)
	pageSizeStr := q.Get(constants.PageSize)
	productIDFilter := q.Get(review.ReviewProductID)
	userIDFilter := q.Get(review.ReviewUserID)
	ratingFilter := q.Get(review.ReviewRating)
	sortBy := q.Get(constants.SortBy)
	sortOrder := q.Get(constants.SortOrder)

	paginatedResponse, err := h.useCase.GetAllReviewsPaginated(includeStr, productIDFilter, userIDFilter, ratingFilter, pageStr, pageSizeStr, sortBy, sortOrder)
	if err != nil {
		response.SendErrorJSON(w, "Failed to fetch reviews", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, paginatedResponse)
}

func (h *ReviewHandler) GetReviewByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil {
		response.SendErrorJSON(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	q := r.URL.Query()
	includeStr := q.Get(constants.Include)

	review, err := h.useCase.GetReviewByID(fmt.Sprintf("%d", *id), includeStr)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.SendErrorJSON(w, "Review not found", http.StatusNotFound)
			return
		}
		response.SendErrorJSON(w, "Failed to fetch review", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, review)
}

func (h *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		ProductID string `json:"product_id"`
		Rating    string `json:"rating"`
		Comment   string `json:"comment"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		response.SendErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	validationErrors := h.useCase.ValidateReviewInput(requestData.Rating, requestData.Comment)
	if len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	userID := h.getUserID()
	review, err := h.useCase.CreateReview(userID, requestData.ProductID, requestData.Rating, requestData.Comment)
	if err != nil {
		if strings.Contains(err.Error(), "already reviewed") {
			response.SendErrorJSON(w, "You have already reviewed this product", http.StatusConflict)
			return
		}
		response.SendErrorJSON(w, "Failed to create review", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, review, http.StatusCreated)
}

func (h *ReviewHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil {
		response.SendErrorJSON(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	var requestData struct {
		Rating  string `json:"rating"`
		Comment string `json:"comment"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		response.SendErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	validationErrors := h.useCase.ValidateReviewInput(requestData.Rating, requestData.Comment)
	if len(validationErrors) > 0 {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	userID := h.getUserID()
	review, err := h.useCase.UpdateReview(fmt.Sprintf("%d", *id), userID, requestData.Rating, requestData.Comment)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.SendErrorJSON(w, "Review not found or not owned by user", http.StatusNotFound)
			return
		}
		response.SendErrorJSON(w, "Failed to update review", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, review)
}

func (h *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil {
		response.SendErrorJSON(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	userID := h.getUserID()
	err = h.useCase.DeleteReview(fmt.Sprintf("%d", *id), userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.SendErrorJSON(w, "Review not found or not owned by user", http.StatusNotFound)
			return
		}
		response.SendErrorJSON(w, "Failed to delete review", http.StatusInternalServerError)
		return
	}

	response.SendDeleteJSON(w, "Review deleted successfully")
}

func (h *ReviewHandler) GetReviewsByProduct(w http.ResponseWriter, r *http.Request) {
	productID, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || productID == nil {
		response.SendErrorJSON(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	q := r.URL.Query()
	includeStr := q.Get(constants.Include)
	ratingFilter := q.Get(review.ReviewRating)
	sortBy := q.Get(constants.SortBy)
	sortOrder := q.Get(constants.SortOrder)

	reviews, err := h.useCase.GetReviewsByProduct(fmt.Sprintf("%d", *productID), includeStr, ratingFilter, sortBy, sortOrder)
	if err != nil {
		response.SendErrorJSON(w, "Failed to fetch reviews", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, reviews)
}

func (h *ReviewHandler) GetReviewsByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || userID == nil {
		response.SendErrorJSON(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	q := r.URL.Query()
	includeStr := q.Get(constants.Include)
	ratingFilter := q.Get(review.ReviewRating)
	sortBy := q.Get(constants.SortBy)
	sortOrder := q.Get(constants.SortOrder)

	reviews, err := h.useCase.GetReviewsByUser(fmt.Sprintf("%d", *userID), includeStr, ratingFilter, sortBy, sortOrder)
	if err != nil {
		response.SendErrorJSON(w, "Failed to fetch reviews", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, reviews)
}

func (h *ReviewHandler) GetProductRatingStats(w http.ResponseWriter, r *http.Request) {
	productID, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || productID == nil {
		response.SendErrorJSON(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	stats, err := h.useCase.GetProductRatingStats(fmt.Sprintf("%d", *productID))
	if err != nil {
		response.SendErrorJSON(w, "Failed to fetch rating statistics", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, stats)
}

func (h *ReviewHandler) getUserID() uint {
	return uint(1)
}
