package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/easy-comerce/backend/internal/feature/review"
	"github.com/easy-comerce/backend/pkg/constants"
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
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(constants.ShowDeleted))
	err, paginatedResponse := h.getAllReviewsPaginatedData(r, showDeleted)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, paginatedResponse)
}

func (h *ReviewHandler) GetReviewByID(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(constants.ShowDeleted))
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

func (h *ReviewHandler) UndoDeletedReview(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	idStr := r.PathValue(constants.FieldID)
	err = h.useCase.UndoDeletedReview(idStr, *userID)
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
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(constants.ShowDeleted))
	err, reviews := h.getReviewsByProductData(r, showDeleted)

	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
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
	ratingFilter := q.Get(constants.ReviewRating)
	productIDFilter := q.Get(constants.ReviewProductID)
	pageStr := q.Get(constants.Page)
	pageSizeStr := q.Get(constants.PageSize)
	sortBy := q.Get(constants.SortBy)
	sortOrder := q.Get(constants.SortOrder)
	showDeleted := utils.ParseBoolPtr(q.Get(constants.ShowDeleted))

	paginatedResponse, err := h.useCase.GetReviewsByUser(*userID, showDeleted, productIDFilter, includeStr, ratingFilter, pageStr, pageSizeStr, sortBy, sortOrder)
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
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(constants.ShowDeleted))
	err, stats := h.getProductRatingStatsData(r, showDeleted)

	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, stats)
}

func (h *ReviewHandler) getReviewDataByID(r *http.Request, showDeleted *bool) (error, *review.Review) {
	id, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || id == nil || *id == 0 {
		return errors.New("invalid review ID"), nil
	}

	q := r.URL.Query()
	includeStr := q.Get(constants.Include)
	review, err := h.useCase.GetReviewByID(*id, showDeleted, includeStr)
	if err != nil {
		logger.Logger.Error("Failed to fetch review by ID", "method", "GetReviewByID", "error", err, "id", id, "include", includeStr)
		return errors.New("review not found"), nil
	}

	return nil, review
}

func (h *ReviewHandler) getAllReviewsPaginatedData(r *http.Request, showDeleted *bool) (error, *models.PaginatedResponse) {
	q := r.URL.Query()

	includeStr := q.Get(constants.Include)
	pageStr := q.Get(constants.Page)
	pageSizeStr := q.Get(constants.PageSize)
	productIDFilter := q.Get(constants.ReviewProductID)
	userIDFilter := q.Get(constants.ReviewUserID)
	ratingFilter := q.Get(constants.ReviewRating)
	sortBy := q.Get(constants.SortBy)
	sortOrder := q.Get(constants.SortOrder)

	paginatedResponse, err := h.useCase.GetAllReviewsPaginated(showDeleted, includeStr, productIDFilter, userIDFilter, ratingFilter, pageStr, pageSizeStr, sortBy, sortOrder)
	if err != nil {
		return errors.New("failed to fetch reviews"), nil
	}

	return nil, paginatedResponse
}

func (h *ReviewHandler) getReviewsByProductData(r *http.Request, showDeleted *bool) (error, *models.PaginatedResponse) {
	productID, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || productID == nil || *productID == 0 {
		return errors.New("invalid product ID"), nil
	}

	q := r.URL.Query()
	includeStr := q.Get(constants.Include)
	ratingFilter := q.Get(constants.ReviewRating)
	pageStr := q.Get(constants.Page)
	pageSizeStr := q.Get(constants.PageSize)
	sortBy := q.Get(constants.SortBy)
	sortOrder := q.Get(constants.SortOrder)

	reviews, err := h.useCase.GetReviewsByProduct(*productID, showDeleted, includeStr, ratingFilter, pageStr, pageSizeStr, sortBy, sortOrder)
	if err != nil {
		return errors.New("failed to fetch reviews"), nil
	}

	return nil, reviews
}

func (h *ReviewHandler) getProductRatingStatsData(r *http.Request, showDeleted *bool) (error, *review.ProductRatingStatsResponse) {
	productID, err := utils.ParseUint(r.PathValue(constants.FieldID))
	if err != nil || productID == nil || *productID == 0 {
		return errors.New("invalid product ID"), nil
	}

	stats, err := h.useCase.GetProductRatingStats(*productID, showDeleted)
	if err != nil {
		return errors.New("failed to fetch rating statistics"), nil
	}

	return nil, stats
}

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
