package option

import (
	"strings"
)

type QueryOptions struct {
	ShowDeleted    *bool
	SortBy         string
	SortOrder      string
	SortableFields []string
	SortOptions    []SortOption
	Preloads       []string
	Filters        map[string]any
}

// USAGE: q.AddFilter("contact_type = ?", contactType)
func (q *QueryOptions) AddFilter(condition string, value any) {
	if q == nil || condition == "" || value == nil {
		return
	}

	if q.Filters == nil {
		q.Filters = make(map[string]any)
	}

	switch v := value.(type) {
	case []any:
		if len(v) > 0 {
			q.Filters[condition] = value
		}
	default:
		q.Filters[condition] = value
	}
}

// USAGE: q.AddInFilter("id", []any{1, 2, 3})
func (q *QueryOptions) AddInFilter(field string, values ...any) {
	addInNotInFilter(q, field, values, false)
}

// USAGE: q.AddNotInFilter("id", []any{1, 2, 3})
func (q *QueryOptions) AddNotInFilter(field string, values ...any) {
	addInNotInFilter(q, field, values, true)
}

// USAGE: q.AddFullLikeOrFilter("zak", "name", "email")
func (q *QueryOptions) AddFullLikeOrFilter(searchTerm string, fields ...string) {
	addLikeFilter(q, searchTerm, "%"+searchTerm+"%", "OR", fields...)
}

func (q *QueryOptions) AddFullLikeAndFilter(searchTerm string, fields ...string) {
	addLikeFilter(q, searchTerm, "%"+searchTerm+"%", "AND", fields...)
}

func (q *QueryOptions) AddPrefixLikeOrFilter(searchTerm string, fields ...string) {
	addLikeFilter(q, searchTerm, searchTerm+"%", "OR", fields...)
}

func (q *QueryOptions) AddLPrefixikeAndFilter(searchTerm string, fields ...string) {
	addLikeFilter(q, searchTerm, searchTerm+"%", "AND", fields...)
}

func (q *QueryOptions) AddSuffixLikeOrFilter(searchTerm string, fields ...string) {
	addLikeFilter(q, searchTerm, "%"+searchTerm, "OR", fields...)
}

func (q *QueryOptions) AddSuffixLikeAndFilter(searchTerm string, fields ...string) {
	addLikeFilter(q, searchTerm, "%"+searchTerm, "AND", fields...)
}

func addLikeFilter(q *QueryOptions, searchTerm string, condition string, orAnd string, fields ...string) {
	if q == nil || searchTerm == "" || len(fields) == 0 {
		return
	}

	if q.Filters == nil {
		q.Filters = make(map[string]any)
	}

	conditions := make([]string, len(fields))
	values := make([]any, len(fields))

	for i, field := range fields {
		conditions[i] = field + " ILIKE ?"
		values[i] = condition
	}

	q.Filters[strings.Join(conditions, " "+orAnd+" ")] = values
}

func addInNotInFilter(q *QueryOptions, field string, values []any, isNotIn bool) {
	if q == nil || field == "" || values == nil || len(values) == 0 {
		return
	}

	if q.Filters == nil {
		q.Filters = make(map[string]any)
	}

	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = "?"
	}

	var condition string
	if isNotIn {
		condition = field + " NOT IN (" + strings.Join(placeholders, ",") + ")"
	} else {
		condition = field + " IN (" + strings.Join(placeholders, ",") + ")"
	}

	q.Filters[condition] = values
}
