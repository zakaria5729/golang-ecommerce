package review

type ReviewPaginationParams struct {
	ShowDeleted *bool  `form:"show_deleted"`
	ProductID   *uint  `form:"product_id"`
	UserID      *uint  `form:"user_id"`
	RatingFrom  *int   `form:"rating_from"`
	RatingTo    *int   `form:"rating_to"`
	Page        int    `form:"page"`
	PageSize    int    `form:"page_size"`
	SortBy      string `form:"sort_by"`
	SortOrder   string `form:"sort_order"`
}
