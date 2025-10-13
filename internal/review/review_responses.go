package review

type ProductRatingStatsResponse struct {
	AverageRating float64     `json:"average_rating"`
	RatingCounts  map[int]int `json:"rating_counts"`
}
