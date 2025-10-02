package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/easy-comerce/backend/internal/feature/review"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/models"
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

func (h *ReviewHandler) GetAllReviewsPaginatedPublic(w http.ResponseWriter, r *http.Request) {
	err, paginatedResponse := h.getAllReviewsPaginatedData(r, nil)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, paginatedResponse)
}

func (h *ReviewHandler) GetReviewByIdPublic(w http.ResponseWriter, r *http.Request) {
	err, review := h.getReviewDataByID(r, nil)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, review)
}

func (h *ReviewHandler) GetAllReviewsPaginated(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, paginatedResponse := h.getAllReviewsPaginatedData(r, showDeleted)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, paginatedResponse)
}

func (h *ReviewHandler) GetReviewByID(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, review := h.getReviewDataByID(r, showDeleted)

	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, review)
}

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

func (h *ReviewHandler) UpdateReviewByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	err, validationErrors := h.updateReviewData(r, userID)
	if validationErrors != nil {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendCommonResponseJSON(w, "Review updated successfully")
}

func (h *ReviewHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	err, validationErrors := h.updateReviewData(r, nil)
	if validationErrors != nil {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendCommonResponseJSON(w, "Review updated successfully")
}

func (h *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue(c.FieldID)
	err := h.useCase.DeleteReview(idStr, nil)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendDeleteJSON(w, "Review deleted successfully")
}

func (h *ReviewHandler) DeleteReviewByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	idStr := r.PathValue(c.FieldID)
	err = h.useCase.DeleteReview(idStr, userID)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendDeleteJSON(w, "Review deleted successfully")
}

func (h *ReviewHandler) UndoDeletedReview(w http.ResponseWriter, r *http.Request) {
	err := h.useCase.UndoDeletedReview(r.PathValue(c.FieldID))
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendDeleteJSON(w, "Undo review deleted successfully")
}

func (h *ReviewHandler) GetReviewsByProductPublic(w http.ResponseWriter, r *http.Request) {
	err, reviews := h.getReviewsByProductData(r, nil)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, reviews)
}

func (h *ReviewHandler) GetReviewsByProduct(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, reviews := h.getReviewsByProductData(r, showDeleted)

	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, reviews)
}

func (h *ReviewHandler) GetReviewsByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	q := r.URL.Query()
	ratingFilter := q.Get(c.ReviewRating)
	productIDFilter := q.Get(c.ReviewProductID)
	pageStr := q.Get(c.Page)
	pageSizeStr := q.Get(c.PageSize)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)
	showDeleted := utils.ParseBoolPtr(q.Get(c.ShowDeleted))

	paginatedResponse, err := h.useCase.GetReviewsByUser(*userID, showDeleted, productIDFilter, ratingFilter, pageStr, pageSizeStr, sortBy, sortOrder)
	if err != nil {
		response.SendErrorJSON(w, "Failed to fetch reviews", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, paginatedResponse)
}

func (h *ReviewHandler) GetProductRatingStatsPublic(w http.ResponseWriter, r *http.Request) {
	err, stats := h.getProductRatingStatsData(r, nil)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, stats)
}

func (h *ReviewHandler) GetProductRatingStats(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, stats := h.getProductRatingStatsData(r, showDeleted)

	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, stats)
}

func (h *ReviewHandler) getReviewDataByID(r *http.Request, showDeleted *bool) (error, *review.Review) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		return errors.New("invalid review ID"), nil
	}

	review, err := h.useCase.GetReviewByID(*id, showDeleted)
	if err != nil {
		logger.Logger.Error("Failed to fetch review by ID", "method", "GetReviewByID", "error", err, "id", id)
		return errors.New("review not found"), nil
	}

	return nil, review
}

func (h *ReviewHandler) getAllReviewsPaginatedData(r *http.Request, showDeleted *bool) (error, *models.PaginatedResponse) {
	q := r.URL.Query()

	pageStr := q.Get(c.Page)
	pageSizeStr := q.Get(c.PageSize)
	productIDFilter := q.Get(c.ReviewProductID)
	userIDFilter := q.Get(c.ReviewUserID)
	ratingFromFilter := q.Get(c.From)
	ratingToFilter := q.Get(c.To)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	paginatedResponse, err := h.useCase.GetAllReviewsPaginated(showDeleted, productIDFilter, userIDFilter, ratingFromFilter, ratingToFilter, pageStr, pageSizeStr, sortBy, sortOrder)
	if err != nil {
		return errors.New("failed to fetch reviews"), nil
	}

	return nil, paginatedResponse
}

func (h *ReviewHandler) getReviewsByProductData(r *http.Request, showDeleted *bool) (error, *models.PaginatedResponse) {
	productID, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || productID == nil || *productID == 0 {
		return errors.New("invalid product ID"), nil
	}

	q := r.URL.Query()
	ratingFilter := q.Get(c.ReviewRating)
	pageStr := q.Get(c.Page)
	pageSizeStr := q.Get(c.PageSize)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	reviews, err := h.useCase.GetReviewsByProduct(*productID, showDeleted, ratingFilter, pageStr, pageSizeStr, sortBy, sortOrder)
	if err != nil {
		return errors.New("failed to fetch reviews"), nil
	}

	return nil, reviews
}

func (h *ReviewHandler) getProductRatingStatsData(r *http.Request, showDeleted *bool) (error, *review.ProductRatingStatsResponse) {
	productID, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || productID == nil || *productID == 0 {
		return errors.New("invalid product ID"), nil
	}

	stats, err := h.useCase.GetProductRatingStats(*productID, showDeleted)
	if err != nil {
		return errors.New("failed to fetch rating statistics"), nil
	}

	return nil, stats
}

func (h *ReviewHandler) updateReviewData(r *http.Request, userID *uint) (error, validator.ValidationErrors) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil {
		return errors.New("Invalid review ID"), nil
	}

	var req review.UpdateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.New("Invalid request body"), nil
	}

	validationErrors := h.validateReviewRequest(req.Rating, req.Comment)
	if validationErrors.HasErrors() {
		return errors.New("Invalid review ID"), validationErrors
	}

	err = h.useCase.UpdateReview(*id, userID, req.Rating, req.Comment)
	if err != nil {
		return err, nil
	}

	return nil, nil
}

func (h *ReviewHandler) validateReviewRequest(rating int, comment string) validator.ValidationErrors {
	var errors validator.ValidationErrors

	if rating < 1 || rating > 5 {
		errors.AddError("rating", "Rating must be between 1 and 5")
	}

	if utils.Trim(comment) == "" {
		errors.AddError("comment", "Comment is required")
	} else if len(comment) > c.MaxReviewCommentLength {
		errors.AddError("comment", "Comment must not exceed"+strconv.Itoa(c.MaxReviewCommentLength)+" characters")
	}

	return errors
}
