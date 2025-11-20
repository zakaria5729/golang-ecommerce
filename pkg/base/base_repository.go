package base

import (
	"context"
	"errors"
	"reflect"

	e "github.com/easy-comerce/backend/pkg/app_error"
	c "github.com/easy-comerce/backend/pkg/constants"
	cu "github.com/easy-comerce/backend/pkg/contextutil"
	l "github.com/easy-comerce/backend/pkg/logger"
	op "github.com/easy-comerce/backend/pkg/option"
	"github.com/easy-comerce/backend/pkg/timeutil"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type BaseRepository[T any] interface {
	Create(ctx context.Context, entity *T) (err error)
	CreateInBatch(ctx context.Context, entities *[]T) (err error)
	Update(ctx context.Context, id uint, entity *T) (err error)
	SoftDelete(ctx context.Context, id uint) (err error)
	UndoSoftDelete(ctx context.Context, id uint) (err error)
	HardDelete(ctx context.Context, id uint) (err error)
	Count(ctx context.Context, options *op.QueryOptions) (count int64, err error)
	ExistsByID(ctx context.Context, id uint, options *op.QueryOptions) (exists bool, err error)
	GetSingleByID(ctx context.Context, id uint, options *op.QueryOptions, selectFields ...string) (entity *T, err error)
	GetSingleBy(ctx context.Context, options *op.QueryOptions, selectFields ...string) (entity *T, err error)
	GetAll(ctx context.Context, options *op.QueryOptions, selectFields ...string) (entities []T, err error)
	GetAllPaginated(ctx context.Context, pageNo int, pageSize int, options *op.QueryOptions, selectFields ...string) (entities []T, count int64, err error)
	GetDB() (db *gorm.DB)
	WithTx(tx *gorm.DB) (repo BaseRepository[T])
	Transaction(ctx context.Context, fn func(repo BaseRepository[T]) error) (err error)
}

type baseRepository[T any] struct {
	db         *gorm.DB
	entity     T
	entityName string
}

func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
	var entity T
	name := reflect.TypeOf(entity).Name()
	if name == "" {
		name = reflect.TypeOf(entity).Elem().Name()
	}

	return &baseRepository[T]{
		db:         db,
		entity:     entity,
		entityName: name,
	}
}

func (r *baseRepository[T]) Create(ctx context.Context, entity *T) (err error) {
	userID, _ := cu.GetUserIDFromContext(ctx)
	findAndSetActionValue(r, ctx, entity, c.FieldCreatedBy, userID)
	query := r.db.WithContext(ctx).Model(&entity).Omit(c.FieldUpdatedAt)
	sql := utils.GetRawSqlCreate(query, &r.entity)

	if err = query.Create(&entity).Error; err != nil {
		l.Error("❌ Failed to create entity", "entityName", r.entityName, "error", err, "userID", userID, "entity", entity, "method", "Create", "sql", sql)
		return e.WrapServerError("Failed to create record", err)
	}

	return nil
}

func (r *baseRepository[T]) CreateInBatch(ctx context.Context, entities *[]T) (err error) {
	userID, _ := cu.GetUserIDFromContext(ctx)
	findAndSetActionValuesInBatch(r, ctx, entities, c.FieldCreatedBy, userID)
	query := r.db.WithContext(ctx).Model(&entities).Omit(c.FieldUpdatedAt)
	sql := utils.GetRawSqlCreate(query, &r.entity)

	if err = query.Create(&entities).Error; err != nil {
		l.Error("❌ Failed to create entity in batch", "entityName", r.entityName, "error", err, "userID", userID, "entities", entities, "method", "CreateBatch", "sql", sql)
		return e.WrapServerError("Failed to create record in batch", err)
	}

	return nil
}

func (r *baseRepository[T]) Update(ctx context.Context, id uint, entity *T) (err error) {
	userID, _ := cu.GetUserIDFromContext(ctx)
	findAndSetActionValue(r, ctx, entity, c.FieldUpdatedBy, userID)
	query := r.db.WithContext(ctx).Model(&entity).Where(c.FieldID+" = ?", id)
	sql := utils.GetRawSqlUpdate(query, &r.entity)

	if err = query.Updates(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("No record found with this id")
		} else {
			l.Error("❌ Failed to update entity", "entityName", r.entityName, "error", err, "id", id, "userID", userID, "entity", entity, "method", "Update", "sql", sql)
			err = e.WrapServerError("Failed to update record", err)
		}

		return err
	}

	return nil
}

func (r *baseRepository[T]) HardDelete(ctx context.Context, id uint) (err error) {
	showDeleted := true
	singleEntity, err := r.GetSingleByID(ctx, id, &op.QueryOptions{ShowDeleted: &showDeleted})
	if err != nil || singleEntity == nil {
		return err
	}

	userID, _ := cu.GetUserIDFromContext(ctx)
	query := r.db.WithContext(ctx).Unscoped().Model(&r.entity).Where(c.FieldID+" = ?", id)
	sql := utils.GetRawSqlDelete(query, &r.entity)

	if err = query.Delete(&r.entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("No record found with this id")
		} else {
			l.Error("❌ Failed to hard delete entity", "entityName", r.entityName, "error", err, "id", id, "userID", userID, "entity", r.entity, "method", "HardDelete", "sql", sql)
			err = e.WrapServerError("Failed to delete record", err)
		}

		return err
	}

	return nil
}

func (r *baseRepository[T]) SoftDelete(ctx context.Context, id uint) (err error) {
	userID, _ := cu.GetUserIDFromContext(ctx)
	findAndSetActionValue(r, ctx, &r.entity, c.FieldDeletedBy, userID)
	findAndSetActionValue(r, ctx, &r.entity, c.FieldDeletedAt, timeutil.GormNowUTC())
	query := r.db.WithContext(ctx).Unscoped().Omit(c.FieldUpdatedAt).Model(&r.entity).Where(c.FieldID+" = ?", id)
	sql := utils.GetRawSqlUpdate(query, &r.entity)

	if err = query.Updates(&r.entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("No record found with this id")
		} else {
			l.Error("❌ Failed to soft delete entity", "entityName", r.entityName, "error", err, "id", id, "userID", userID, "entity", r.entity, "method", "SoftDelete", "sql", sql)
			err = e.WrapServerError("Failed to soft delete record", err)
		}

		return err
	}

	return nil
}

func (r *baseRepository[T]) UndoSoftDelete(ctx context.Context, id uint) (err error) {
	showDeleted := true
	singleEntity, err := r.GetSingleByID(ctx, id, &op.QueryOptions{ShowDeleted: &showDeleted})
	if err != nil || singleEntity == nil {
		return err
	}

	userID, _ := cu.GetUserIDFromContext(ctx)
	deletedAtExists := findAndSetActionValue(r, ctx, &r.entity, c.FieldDeletedAt, nil)
	deletedByExists := findAndSetActionValue(r, ctx, &r.entity, c.FieldDeletedBy, nil)
	query := r.db.WithContext(ctx).Unscoped().Model(&r.entity).Omit(c.FieldUpdatedAt).Where(c.FieldID+" = ?", id)

	if deletedAtExists && deletedByExists {
		query = query.Select(c.FieldDeletedAt, c.FieldDeletedBy)
	} else if deletedAtExists {
		query = query.Select(c.FieldDeletedAt)
	} else if deletedByExists {
		query = query.Select(c.FieldDeletedBy)
	}
	sql := utils.GetRawSqlUpdate(query, &r.entity)

	if err = query.Updates(&r.entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("No record found with this id")
		} else {
			l.Error("❌ Failed to undo soft delete entity", "entityName", r.entityName, "error", err, "id", id, "userID", userID, "entity", r.entity, "method", "UndoSoftDelete", "sql", sql)
			err = e.WrapServerError("Failed to undo deleted record", err)
		}

		return err
	}

	return nil
}

func (r *baseRepository[T]) ExistsByID(ctx context.Context, id uint, options *op.QueryOptions) (exists bool, err error) {
	userID, _ := cu.GetUserIDFromContext(ctx)
	query := r.db.WithContext(ctx).Model(&r.entity)
	query = unscopedQueryOptions(query, options)
	query = whereQueryOptions(query, options)
	query = query.Select("1").Where(c.FieldID+" = ?", id).Limit(1).Scan(&exists)
	sql := utils.GetRawSqlExists(query, &exists)

	if err = query.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("No record found with this id")
		} else {
			l.Error("❌ Failed to check if entity exists", "entityName", r.entityName, "error", err, "id", id, "userID", userID, "options", options, "method", "ExistsByID", "sql", sql)
			err = e.WrapServerError("Failed to check if record exists", err)
		}

		return false, err
	}

	return exists, nil
}

func (r *baseRepository[T]) Count(ctx context.Context, options *op.QueryOptions) (count int64, err error) {
	userID, _ := cu.GetUserIDFromContext(ctx)
	query := r.db.WithContext(ctx).Model(&r.entity)
	query = unscopedQueryOptions(query, options)
	query = whereQueryOptions(query, options)
	sql := utils.GetRawSqlCount(query, &count)

	if err = query.Count(&count).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("No record found")
		} else {
			l.Error("❌ Failed to count entities", "entityName", r.entityName, "error", err, "userID", userID, "options", options, "method", "Count", "sql", sql)
			err = e.WrapServerError("Failed to count records", err)
		}

		return 0, err
	}

	return count, nil
}

func (r *baseRepository[T]) GetSingleBy(ctx context.Context, options *op.QueryOptions, selectFields ...string) (entity *T, err error) {
	userID, _ := cu.GetUserIDFromContext(ctx)
	query := r.db.WithContext(ctx).Model(&r.entity)
	if len(selectFields) > 0 {
		query = query.Select(selectFields)
	}

	query = unscopedQueryOptions(query, options)
	query = whereQueryOptions(query, options)
	query = orderByQueryOptions(query, options)
	query = preloadQueryOptions(query, options)
	sql := utils.GetRawSqlFirst(query, &r.entity)

	if err = query.First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("No record found with this id")
		} else {
			l.Error("❌ Failed to get entity", "entityName", r.entityName, "error", err, "userID", userID, "options", options, "entity", entity, "method", "GetSingleBy", "sql", sql)
			err = e.WrapServerError("Failed to get record", err)
		}

		return nil, err
	}

	return entity, nil
}

func (r *baseRepository[T]) GetSingleByID(ctx context.Context, id uint, options *op.QueryOptions, selectFields ...string) (entity *T, err error) {
	userID, _ := cu.GetUserIDFromContext(ctx)
	query := r.db.WithContext(ctx).Model(&r.entity).Where(c.FieldID+" = ?", id)
	if len(selectFields) > 0 {
		query = query.Select(selectFields)
	}

	query = unscopedQueryOptions(query, options)
	query = whereQueryOptions(query, options)
	query = orderByQueryOptions(query, options)
	query = preloadQueryOptions(query, options)
	sql := utils.GetRawSqlFirst(query, &r.entity)

	if err = query.First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("No record found with this id")
		} else {
			l.Error("❌ Failed to get entity", "entityName", r.entityName, "error", err, "id", id, "userID", userID, "options", options, "entity", entity, "method", "GetSingleByID", "sql", sql)
			err = e.WrapServerError("Failed to get record", err)
		}

		return nil, err
	}

	return entity, nil
}

func (r *baseRepository[T]) GetAll(ctx context.Context, options *op.QueryOptions, selectFields ...string) (entities []T, err error) {
	userID, _ := cu.GetUserIDFromContext(ctx)
	query := r.db.WithContext(ctx).Model(&entities)
	if len(selectFields) > 0 {
		query = query.Select(selectFields)
	}

	query = unscopedQueryOptions(query, options)
	query = whereQueryOptions(query, options)
	query = orderByQueryOptions(query, options)
	query = preloadQueryOptions(query, options)
	sql := utils.GetRawSqlFind(query, &r.entity)

	if err = query.Find(&entities).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("No record found")
		} else {
			l.Error("❌ Failed to get entities", "entityName", r.entityName, "error", err, "userID", userID, "options", options, "method", "GetAll", "sql", sql)
			err = e.WrapServerError("Failed to get records", err)
		}

		return nil, err
	}

	return entities, nil
}

func (r *baseRepository[T]) GetAllPaginated(ctx context.Context, pageNo int, pageSize int, options *op.QueryOptions, selectFields ...string) (entities []T, count int64, err error) {
	userID, _ := cu.GetUserIDFromContext(ctx)
	query := r.db.WithContext(ctx).Model(&entities)
	if len(selectFields) > 0 {
		query = query.Select(selectFields)
	}

	query = unscopedQueryOptions(query, options)
	query = whereQueryOptions(query, options)
	countSql := utils.GetRawSqlCount(query, &count)

	if err := query.Count(&count).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("No record found")
		} else {
			l.Error("❌ Failed to count entities", "entityName", r.entityName, "error", err, "userID", userID, "pageNo", pageNo, "pageSize", pageSize, "options", options, "method", "GetAllPaginated", "countSql", countSql)
			err = e.WrapServerError("Failed to count records", err)
		}

		return nil, 0, err
	}

	query = orderByQueryOptions(query, options)
	query = preloadQueryOptions(query, options)
	query = query.Offset(utils.GetOffset(pageNo, pageSize)).Limit(pageSize)
	sql := utils.GetRawSqlFind(query, &r.entity)

	if err := query.Find(&entities).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("No record found")
		} else {
			l.Error("❌ Failed to get entities", "entityName", r.entityName, "error", err, "userID", userID, "pageNo", pageNo, "pageSize", pageSize, "options", options, "method", "GetAllPaginated", "sql", sql)
			err = e.WrapServerError("Failed to get records", err)
		}

		return nil, 0, err
	}

	return entities, count, nil
}

func (r *baseRepository[T]) GetDB() (db *gorm.DB) {
	return r.db
}

func (r *baseRepository[T]) WithTx(tx *gorm.DB) (repo BaseRepository[T]) {
	return &baseRepository[T]{db: tx}
}

func (r *baseRepository[T]) Transaction(ctx context.Context, fn func(r BaseRepository[T]) error) (err error) {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := r.WithTx(tx)
		return fn(txRepo)
	})
}

func unscopedQueryOptions(db *gorm.DB, options *op.QueryOptions) *gorm.DB {
	if options != nil && options.ShowDeleted != nil && *options.ShowDeleted {
		db = db.Unscoped()
	}
	return db
}

func whereQueryOptions(db *gorm.DB, options *op.QueryOptions) *gorm.DB {
	if options != nil && options.Filters != nil {
		for condition, value := range options.Filters {
			if value != nil {
				switch v := value.(type) {
				case []any:
					db = db.Where(condition, v...)
				default:
					db = db.Where(condition, value)
				}
			}
		}
	}
	return db
}

func preloadQueryOptions(db *gorm.DB, options *op.QueryOptions) *gorm.DB {
	if options != nil {
		for _, preload := range options.Preloads {
			db = db.Preload(preload)
		}
	}
	return db
}

func orderByQueryOptions(db *gorm.DB, options *op.QueryOptions) *gorm.DB {
	if options == nil {
		return db
	}
	if len(options.SortOptions) > 0 {
		if orderClause := utils.BuildSortingOrders(options.SortOptions, &options.SortableFields); orderClause != "" {
			db = db.Order(orderClause)
		}
	} else {
		if orderClause := utils.BuildSortingOrder(options.SortBy, options.SortOrder, &options.SortableFields); orderClause != "" {
			db = db.Order(orderClause)
		}
	}
	return db
}

func findAndSetActionValue[T any](r *baseRepository[T], ctx context.Context, entity *T, fieldName string, value any) bool {
	if entity == nil {
		return false
	}

	slice := []T{*entity}
	fieldExists := findAndSetActionValuesInBatch(r, ctx, &slice, fieldName, value)

	if fieldExists {
		*entity = slice[0]
	}

	return fieldExists
}

func findAndSetActionValuesInBatch[T any](r *baseRepository[T], ctx context.Context, entities *[]T, fieldName string, value any) (fieldExists bool) {
	if entities == nil || len(*entities) <= 0 {
		return false
	}

	var entity T
	stmt := &gorm.Statement{DB: r.db}
	if err := stmt.Parse(&entity); err != nil {
		return false
	}

	field := stmt.Schema.LookUpField(fieldName)
	if field == nil {
		return false
	}

	v := reflect.ValueOf(entities).Elem()
	for i := 0; i < v.Len(); i++ {
		if err := field.Set(ctx, v.Index(i), value); err != nil {
			continue
		}
	}

	return true
}
