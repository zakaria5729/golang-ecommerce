package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/easy-comerce/backend/internal/feature/review"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type ReviewHandler struct {
	useCase *review.ReviewUseCase
}

func NewReviewHandler() *ReviewHandler {
	return &ReviewHandler{
		useCase: review.NewReviewUseCase(),
	}
}

// **REQUIRED
func (h *ReviewHandler) GetAllReviewsPaginated(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	includeStr := q.Get(constants.Include)
	pageStr := q.Get(constants.Page)
	pageSizeStr := q.Get(constants.PageSize)
	productIDFilter := q.Get(constants.ReviewProductID)
	userIDFilter := q.Get(constants.ReviewUserID)
	ratingFilter := q.Get(constants.ReviewRating)
	sortBy := q.Get(constants.SortBy)
	sortOrder := q.Get(constants.SortOrder)

	paginatedResponse, err := h.useCase.GetAllReviewsPaginated(includeStr, productIDFilter, userIDFilter, ratingFilter, pageStr, pageSizeStr, sortBy, sortOrder)
	if err != nil {
		response.SendErrorJSON(w, "Failed to fetch reviews", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, paginatedResponse)
}

// **REQUIRED
func (h *ReviewHandler) GetReviewByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil || *id == 0 {
		response.SendErrorJSON(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	includeStr := r.URL.Query().Get(constants.Include)
	review, err := h.useCase.GetReviewByID(*id, includeStr)
	if err != nil {
		logger.Logger.Error("Failed to fetch review by ID", "method", "GetReviewByID", "error", err, "id", id, "include", include)
		response.SendErrorJSON(w, "Review not found", http.StatusInternalServerError)
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

	if validationErrors := h.validateReviewRequest(requestData.Rating, requestData.Comment); validationErrors.HasErrors() {
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

	if validationErrors := h.validateReviewRequest(requestData.Rating, requestData.Comment); validationErrors.HasErrors() {
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
	if err != nil || id == nil || *id == 0 {
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
	if err != nil || productID == nil || *productID == 0 {
		response.SendErrorJSON(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	q := r.URL.Query()
	includeStr := q.Get(constants.Include)
	ratingFilter := q.Get(constants.ReviewRating)
	sortBy := q.Get(constants.SortBy)
	sortOrder := q.Get(constants.SortOrder)

	reviews, err := h.useCase.GetReviewsByProduct(*productID, includeStr, ratingFilter, sortBy, sortOrder)
	if err != nil {
		response.SendErrorJSON(w, "Failed to fetch reviews", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, reviews)
}

func (h *ReviewHandler) GetReviewsByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || userID == nil || *userID == 0 {
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

// **REQUIRED
func (h *ReviewHandler) GetProductRatingStats(w http.ResponseWriter, r *http.Request) {
	productID, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || productID == nil || *productID == 0 {
		response.SendErrorJSON(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	stats, err := h.useCase.GetProductRatingStats(*productID)
	if err != nil {
		response.SendErrorJSON(w, "Failed to fetch rating statistics", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, stats)
}

func (h *ReviewHandler) getUserID() uint {
	return uint(1)
}

func (h *ReviewHandler) validateReviewRequest(rating, comment string) validator.ValidationErrors {
	var errors validator.ValidationErrors

	if rating == "" {
		errors.AddError("rating", "Rating is required")
	} else {
		if ratingInt, err := strconv.Atoi(rating); err != nil || ratingInt < 1 || ratingInt > 5 {
			errors.AddError("rating", "Rating must be between 1 and 5")
		}
	}

	if comment != "" {
		if len(comment) > 1000 {
			errors.AddError("comment", "Comment must not exceed 1000 characters")
		}
	}

	return errors
}
