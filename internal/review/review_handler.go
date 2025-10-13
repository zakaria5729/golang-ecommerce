package review

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	c "github.com/easy-comerce/backend/pkg/constants"
	cu "github.com/easy-comerce/backend/pkg/contextutil"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/easy-comerce/backend/pkg/validator"
)

type ReviewHandler struct {
	service *ReviewService
}

func NewReviewHandler(service *ReviewService) *ReviewHandler {
	return &ReviewHandler{
		service: service,
	}
}

func (h *ReviewHandler) GetAllReviewsPaginatedPublic(w http.ResponseWriter, r *http.Request) {
	err, paginatedResponse := getAllReviewsPaginatedData(r, nil, h.service)
	response.SendResponse(w, paginatedResponse, err, http.StatusInternalServerError)
}

func (h *ReviewHandler) GetReviewByIdPublic(w http.ResponseWriter, r *http.Request) {
	err, review := getReviewDataByID(r, nil, h.service)
	response.SendResponse(w, review, err, http.StatusInternalServerError)
}

func (h *ReviewHandler) GetAllReviewsPaginated(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, paginatedResponse := getAllReviewsPaginatedData(r, showDeleted, h.service)
	response.SendResponse(w, paginatedResponse, err, http.StatusInternalServerError)
}

func (h *ReviewHandler) GetReviewByID(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, review := getReviewDataByID(r, showDeleted, h.service)
	response.SendResponse(w, review, err, http.StatusInternalServerError)
}

func (h *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	userID, err := cu.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req CreateReviewRequest
	if !utils.DecodeJSON(w, r, &req, "CreateReview") {
		return
	}

	if validationErrors := validateReviewRequest(req.Rating, req.Comment); validationErrors.HasErrors() {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	review, err := h.service.CreateReview(*userID, req.ProductID, req.Rating, req.Comment)
	response.SendResponse(w, review, err, http.StatusInternalServerError)
}

func (h *ReviewHandler) UpdateReviewByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := cu.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	err, validationErrors := updateReviewData(r, userID, h.service)
	if validationErrors != nil {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	response.SendResponse(w, "Review updated successfully", err, http.StatusInternalServerError)
}

func (h *ReviewHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	err, validationErrors := updateReviewData(r, nil, h.service)
	if validationErrors != nil {
		response.SendValidationErrorJSON(w, "Validation failed", validationErrors)
		return
	}

	response.SendResponse(w, "Review updated successfully", err, http.StatusInternalServerError)
}

func (h *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue(c.FieldID)
	err := h.service.DeleteReview(r.Context(), idStr, nil)
	response.SendResponse(w, "Review deleted successfully", err, http.StatusInternalServerError)
}

func (h *ReviewHandler) DeleteReviewByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := cu.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendErrorJSON(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	idStr := r.PathValue(c.FieldID)
	err = h.service.DeleteReview(r.Context(), idStr, userID)
	response.SendResponse(w, "Review deleted successfully", err, http.StatusInternalServerError)
}

func (h *ReviewHandler) UndoDeletedReview(w http.ResponseWriter, r *http.Request) {
	err := h.service.UndoDeletedReview(r.Context(), r.PathValue(c.FieldID))
	response.SendResponse(w, "Undo review deleted successfully", err, http.StatusInternalServerError)
}

func (h *ReviewHandler) GetReviewsByProductPublic(w http.ResponseWriter, r *http.Request) {
	err, reviews := getReviewsByProductData(r, nil, h.service)
	response.SendResponse(w, reviews, err, http.StatusInternalServerError)
}

func (h *ReviewHandler) GetReviewsByProduct(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, reviews := getReviewsByProductData(r, showDeleted, h.service)
	response.SendResponse(w, reviews, err, http.StatusInternalServerError)
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

	paginatedResponse, err := h.service.GetReviewsByUser(*userID, showDeleted, productIDFilter, ratingFilter, pageStr, pageSizeStr, sortBy, sortOrder)
	response.SendResponse(w, paginatedResponse, err, http.StatusInternalServerError)
}

func (h *ReviewHandler) GetProductRatingStatsPublic(w http.ResponseWriter, r *http.Request) {
	err, stats := getProductRatingStatsData(r, nil, h.service)
	response.SendResponse(w, stats, err, http.StatusInternalServerError)
}

func (h *ReviewHandler) GetProductRatingStats(w http.ResponseWriter, r *http.Request) {
	showDeleted := utils.ParseBoolPtr(r.URL.Query().Get(c.ShowDeleted))
	err, stats := getProductRatingStatsData(r, showDeleted, h.service)
	response.SendResponse(w, stats, err, http.StatusInternalServerError)
}

func getReviewDataByID(r *http.Request, showDeleted *bool, service *ReviewService) (error, *Review) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil || *id == 0 {
		return errors.New("invalid review ID"), nil
	}

	review, err := service.GetReviewByID(*id, showDeleted)
	if err != nil {
		return errors.New("review not found"), nil
	}

	return nil, review
}

func getAllReviewsPaginatedData(r *http.Request, showDeleted *bool, service *ReviewService) (error, *models.PaginatedResponse) {
	q := r.URL.Query()

	pageStr := q.Get(c.Page)
	pageSizeStr := q.Get(c.PageSize)
	productIDFilter := q.Get(c.ReviewProductID)
	userIDFilter := q.Get(c.ReviewUserID)
	ratingFromFilter := q.Get(c.From)
	ratingToFilter := q.Get(c.To)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)

	paginatedResponse, err := service.GetAllReviewsPaginated(showDeleted, productIDFilter, userIDFilter, ratingFromFilter, ratingToFilter, pageStr, pageSizeStr, sortBy, sortOrder)
	if err != nil {
		return errors.New("failed to fetch reviews"), nil
	}

	return nil, paginatedResponse
}

func getReviewsByProductData(r *http.Request, showDeleted *bool, service *ReviewService) (error, *models.PaginatedResponse) {
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

	reviews, err := service.GetReviewsByProduct(*productID, showDeleted, ratingFilter, pageStr, pageSizeStr, sortBy, sortOrder)
	if err != nil {
		return errors.New("failed to fetch reviews"), nil
	}

	return nil, reviews
}

func getProductRatingStatsData(r *http.Request, showDeleted *bool, service *ReviewService) (error, *ProductRatingStatsResponse) {
	productID, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || productID == nil || *productID == 0 {
		return errors.New("invalid product ID"), nil
	}

	stats, err := service.GetProductRatingStats(*productID, showDeleted)
	if err != nil {
		return errors.New("failed to fetch rating statistics"), nil
	}

	return nil, stats
}

func updateReviewData(r *http.Request, userID *uint, service *ReviewService) (error, validator.ValidationErrors) {
	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil || id == nil {
		return errors.New("Invalid review ID"), nil
	}

	var req UpdateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.New("Invalid request body"), nil
	}

	validationErrors := validateReviewRequest(req.Rating, req.Comment)
	if validationErrors.HasErrors() {
		return errors.New("Invalid review ID"), validationErrors
	}

	err = service.UpdateReview(r.Context(), *id, userID, req.Rating, req.Comment)
	if err != nil {
		return err, nil
	}

	return nil, nil
}

func validateReviewRequest(rating int, comment string) validator.ValidationErrors {
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
