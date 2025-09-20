package handler

import (
	"net/http"
	"strconv"

	"github.com/easy-comerce/backend/internal/feature/review"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/middleware"
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
		logger.Logger.Error("Failed to fetch review by ID", "method", "GetReviewByID", "error", err, "id", id, "include", includeStr)
		response.SendErrorJSON(w, "Review not found", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, review)
}

// **REQUIRED
func (h *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req review.CreateReviewRequest
	if !utils.DecodeJSON(w, r, &req, "CreateReview") {
		return
	}

	if validationErrors := h.validateReviewRequest(req.Rating, req.Comment); validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	review, err := h.useCase.CreateReview(*userID, req.ProductID, req.Rating, req.Comment)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, review, http.StatusCreated)
}

// **REQUIRED
func (h *ReviewHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil {
		response.SendErrorJSON(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	var req review.UpdateReviewRequest
	if !utils.DecodeJSON(w, r, &req, "CreateReview") {
		return
	}

	if validationErrors := h.validateReviewRequest(req.Rating, req.Comment); validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	err = h.useCase.UpdateReview(*id, *userID, req.Rating, req.Comment)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendCommonResponseJSON(w, "Review updated successfully")
}

// **REQUIRED
func (h *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	idStr := r.PathValue(constants.FieldID)
	err = h.useCase.DeleteReview(idStr, *userID)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendDeleteJSON(w, "Review deleted successfully")
}

// **REQUIRED
func (h *ReviewHandler) GetReviewsByProduct(w http.ResponseWriter, r *http.Request) {
	productID, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || productID == nil || *productID == 0 {
		response.SendErrorJSON(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	q := r.URL.Query()
	includeStr := q.Get(constants.Include)
	ratingFilter := q.Get(constants.ReviewRating)
	pageStr := q.Get(constants.Page)
	pageSizeStr := q.Get(constants.PageSize)
	sortBy := q.Get(constants.SortBy)
	sortOrder := q.Get(constants.SortOrder)

	reviews, err := h.useCase.GetReviewsByProduct(*productID, includeStr, ratingFilter, pageStr, pageSizeStr, sortBy, sortOrder)
	if err != nil {
		response.SendErrorJSON(w, "Failed to fetch reviews", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, reviews)
}

// **REQUIRED
func (h *ReviewHandler) GetReviewsByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	q := r.URL.Query()
	includeStr := q.Get(constants.Include)
	ratingFilter := q.Get(constants.ReviewRating)
	productIDFilter := q.Get(constants.ReviewProductID)
	pageStr := q.Get(constants.Page)
	pageSizeStr := q.Get(constants.PageSize)
	sortBy := q.Get(constants.SortBy)
	sortOrder := q.Get(constants.SortOrder)

	paginatedResponse, err := h.useCase.GetReviewsByUser(*userID, productIDFilter, includeStr, ratingFilter, pageStr, pageSizeStr, sortBy, sortOrder)
	if err != nil {
		response.SendErrorJSON(w, "Failed to fetch reviews", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, paginatedResponse)
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

// **REQUIRED
func (h *ReviewHandler) validateReviewRequest(rating, comment string) validator.ValidationErrors {
	var errors validator.ValidationErrors

	ratingInt, err := utils.ParseInt(rating)
	if err != nil || rating == "" || ratingInt == nil || *ratingInt < 1 || *ratingInt > 5 {
		errors.AddError("rating", "Rating must be between 1 and 5")
	}

	if utils.Trim(comment) != "" {
		errors.AddError("comment", "Comment is required")
	} else if len(comment) > constants.MaxReviewCommentLength {
		errors.AddError("comment", "Comment must not exceed"+strconv.Itoa(constants.MaxReviewCommentLength)+" characters")
	}

	return errors
}
