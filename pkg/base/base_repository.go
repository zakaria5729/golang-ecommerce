package base

import (
	"context"
	"errors"
	"reflect"
	"strings"

	e "github.com/easy-comerce/backend/pkg/app_error"
	cfg "github.com/easy-comerce/backend/pkg/config"
	c "github.com/easy-comerce/backend/pkg/constants"
	cu "github.com/easy-comerce/backend/pkg/contextutil"
	l "github.com/easy-comerce/backend/pkg/logger"
	op "github.com/easy-comerce/backend/pkg/option"
	tu "github.com/easy-comerce/backend/pkg/timeutil"
	u "github.com/easy-comerce/backend/pkg/utils"
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
	AppConfig  cfg.AppConfig
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
		AppConfig:  cfg.GetConfig().AppConfig,
	}
}

func (r *baseRepository[T]) Create(ctx context.Context, entity *T) (err error) {
	isWriteAlways := isWriteLogWhenAlways(r)
	userID, _ := cu.GetUserIDFromContext(ctx)

	findAndSetActionValue(r, ctx, entity, c.FieldCreatedBy, userID)
	query := r.db.WithContext(ctx).Model(&entity).Omit(c.FieldUpdatedAt)
	sql := u.GetRawSqlCreate(query, &r.entity)

	if err = query.Create(&entity).Error; err != nil {
		if !isWriteAlways {
			l.Error("❌ Failed to create entity", "entityName", r.entityName, "error", err, "userID", userID, "entity", entity, "method", "Create", "sql", sql)
		}
		return e.WrapServerError("Failed to create record", err)
	}

	if isWriteAlways {
		l.Info("✅ Create entity", "entityName", r.entityName, "userID", userID, "entity", entity, "method", "Create", "sql", sql)
	}
	return nil
}

func (r *baseRepository[T]) CreateInBatch(ctx context.Context, entities *[]T) (err error) {
	isWriteAlways := isWriteLogWhenAlways(r)
	userID, _ := cu.GetUserIDFromContext(ctx)

	findAndSetActionValuesInBatch(r, ctx, entities, c.FieldCreatedBy, userID)
	query := r.db.WithContext(ctx).Model(&entities).Omit(c.FieldUpdatedAt)
	sql := u.GetRawSqlCreate(query, &r.entity)

	if err = query.Create(&entities).Error; err != nil {
		if !isWriteAlways {
			l.Error("❌ Failed to create entity in batch", "entityName", r.entityName, "error", err, "userID", userID, "entities", entities, "method", "CreateInBatch", "sql", sql)
		}
		return e.WrapServerError("Failed to create record in batch", err)
	}

	if isWriteAlways {
		l.Info("✅ Create entity in batch", "entityName", r.entityName, "userID", userID, "entities", entities, "method", "CreateInBatch", "sql", sql)
	}
	return nil
}

func (r *baseRepository[T]) Update(ctx context.Context, id uint, entity *T) (err error) {
	isWriteAlways := isWriteLogWhenAlways(r)
	userID, _ := cu.GetUserIDFromContext(ctx)

	findAndSetActionValue(r, ctx, entity, c.FieldUpdatedBy, userID)
	query := r.db.WithContext(ctx).Model(&entity).Where(c.FieldID+" = ?", id)
	sql := u.GetRawSqlUpdate(query, &r.entity)

	if err = query.Updates(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("No record found with this id")
		}

		if !isWriteAlways {
			l.Error("❌ Failed to update entity", "entityName", r.entityName, "error", err, "id", id, "userID", userID, "entity", entity, "method", "Update", "sql", sql)
		}
		return e.WrapServerError("Failed to update record", err)
	}

	if isWriteAlways {
		l.Info("✅ Update entity", "entityName", r.entityName, "id", id, "userID", userID, "entity", entity, "method", "Update", "sql", sql)
	}
	return nil
}

func (r *baseRepository[T]) HardDelete(ctx context.Context, id uint) (err error) {
	isWriteAlways := isWriteLogWhenAlways(r)
	userID, _ := cu.GetUserIDFromContext(ctx)

	showDeleted := true
	singleEntity, err := r.GetSingleByID(ctx, id, &op.QueryOptions{ShowDeleted: &showDeleted})
	if err != nil || singleEntity == nil {
		return err
	}

	query := r.db.WithContext(ctx).Unscoped().Model(&r.entity).Where(c.FieldID+" = ?", id)
	sql := u.GetRawSqlDelete(query, &r.entity)

	if err = query.Delete(&r.entity).Error; err != nil {
		if isNotFoundError(err) {
			return errors.New("No record found with this id")
		}

		if !isWriteAlways {
			l.Error("❌ Failed to hard delete entity", "entityName", r.entityName, "error", err, "id", id, "userID", userID, "entity", r.entity, "method", "HardDelete", "sql", sql)
		}
		return e.WrapServerError("Failed to delete record", err)
	}

	if isWriteAlways {
		l.Info("✅ Hard delete entity", "entityName", r.entityName, "id", id, "userID", userID, "entity", r.entity, "method", "HardDelete", "sql", sql)
	}
	return nil
}

func (r *baseRepository[T]) SoftDelete(ctx context.Context, id uint) (err error) {
	isWriteAlways := isWriteLogWhenAlways(r)
	userID, _ := cu.GetUserIDFromContext(ctx)

	findAndSetActionValue(r, ctx, &r.entity, c.FieldDeletedBy, userID)
	findAndSetActionValue(r, ctx, &r.entity, c.FieldDeletedAt, tu.GormNowUTC())
	query := r.db.WithContext(ctx).Unscoped().Omit(c.FieldUpdatedAt).Model(&r.entity).Where(c.FieldID+" = ?", id)
	sql := u.GetRawSqlUpdate(query, &r.entity)

	if err = query.Updates(&r.entity).Error; err != nil {
		if isNotFoundError(err) {
			return errors.New("No record found with this id")
		}

		if !isWriteAlways {
			l.Error("❌ Failed to soft delete entity", "entityName", r.entityName, "error", err, "id", id, "userID", userID, "entity", r.entity, "method", "SoftDelete", "sql", sql)
		}
		return e.WrapServerError("Failed to soft delete record", err)
	}

	if isWriteAlways {
		l.Info("✅ Soft delete entity", "entityName", r.entityName, "id", id, "userID", userID, "entity", r.entity, "method", "SoftDelete", "sql", sql)
	}
	return nil
}

func (r *baseRepository[T]) UndoSoftDelete(ctx context.Context, id uint) (err error) {
	isWriteAlways := isWriteLogWhenAlways(r)
	userID, _ := cu.GetUserIDFromContext(ctx)

	showDeleted := true
	singleEntity, err := r.GetSingleByID(ctx, id, &op.QueryOptions{ShowDeleted: &showDeleted})
	if err != nil || singleEntity == nil {
		return err
	}

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
	sql := u.GetRawSqlUpdate(query, &r.entity)

	if err = query.Updates(&r.entity).Error; err != nil {
		if isNotFoundError(err) {
			return errors.New("No record found with this id")
		}

		if !isWriteAlways {
			l.Error("❌ Failed to undo soft delete entity", "entityName", r.entityName, "error", err, "id", id, "userID", userID, "entity", r.entity, "method", "UndoSoftDelete", "sql", sql)
		}
		return e.WrapServerError("Failed to undo deleted record", err)
	}

	if isWriteAlways {
		l.Info("✅ Undo soft delete entity", "entityName", r.entityName, "id", id, "userID", userID, "entity", r.entity, "method", "UndoSoftDelete", "sql", sql)
	}
	return nil
}

func (r *baseRepository[T]) ExistsByID(ctx context.Context, id uint, options *op.QueryOptions) (exists bool, err error) {
	isWriteAlways := isWriteLogWhenAlways(r)
	userID, _ := cu.GetUserIDFromContext(ctx)

	query := r.db.WithContext(ctx).Model(&r.entity)
	query = unscopedQueryOptions(query, options)
	query = whereQueryOptions(query, options)
	query = query.Select("1").Where(c.FieldID+" = ?", id).Limit(1).Scan(&exists)
	sql := u.GetRawSqlExists(query, &exists)

	if err = query.Error; err != nil {
		if isNotFoundError(err) {
			return false, errors.New("No record found with this id")
		}

		if !isWriteAlways {
			l.Error("❌ Failed to check if entity exists", "entityName", r.entityName, "error", err, "id", id, "userID", userID, "options", options, "method", "ExistsByID", "sql", sql)
		}
		return false, e.WrapServerError("Failed to check if record exists", err)
	}

	if isWriteAlways {
		l.Info("✅ Check if entity exists", "entityName", r.entityName, "id", id, "userID", userID, "exists", exists, "options", options, "method", "ExistsByID", "sql", sql)
	}
	return exists, nil
}

func (r *baseRepository[T]) Count(ctx context.Context, options *op.QueryOptions) (count int64, err error) {
	isWriteAlways := isWriteLogWhenAlways(r)
	userID, _ := cu.GetUserIDFromContext(ctx)

	query := r.db.WithContext(ctx).Model(&r.entity)
	query = unscopedQueryOptions(query, options)
	query = whereQueryOptions(query, options)
	sql := u.GetRawSqlCount(query, &count)

	if err = query.Count(&count).Error; err != nil {
		if isNotFoundError(err) {
			return 0, errors.New("No record found")
		}

		if !isWriteAlways {
			l.Error("❌ Failed to count entities", "entityName", r.entityName, "error", err, "userID", userID, "options", options, "method", "Count", "sql", sql)
		}
		return 0, e.WrapServerError("Failed to count records", err)
	}

	if isWriteAlways {
		l.Info("✅ Count entities", "entityName", r.entityName, "userID", userID, "count", count, "options", options, "method", "Count", "sql", sql)
	}
	return count, nil
}

func (r *baseRepository[T]) GetSingleBy(ctx context.Context, options *op.QueryOptions, selectFields ...string) (entity *T, err error) {
	isWriteAlways := isWriteLogWhenAlways(r)
	userID, _ := cu.GetUserIDFromContext(ctx)

	query := r.db.WithContext(ctx).Model(&r.entity)
	if len(selectFields) > 0 {
		query = query.Select(selectFields)
	}

	query = unscopedQueryOptions(query, options)
	query = whereQueryOptions(query, options)
	query = orderByQueryOptions(query, options)
	query = preloadQueryOptions(query, options)
	sql := u.GetRawSqlFirst(query, &r.entity)

	if err = query.First(&entity).Error; err != nil {
		if isNotFoundError(err) {
			return nil, errors.New("No record found with this id")
		}

		if !isWriteAlways {
			l.Error("❌ Failed to get entity", "entityName", r.entityName, "error", err, "userID", userID, "options", options, "entity", entity, "method", "GetSingleBy", "sql", sql)
		}
		return nil, e.WrapServerError("Failed to get record", err)
	}

	if isWriteAlways {
		l.Info("✅ Get entity", "entityName", r.entityName, "userID", userID, "options", options, "entity", entity, "method", "GetSingleBy", "sql", sql)
	}
	return entity, nil
}

func (r *baseRepository[T]) GetSingleByID(ctx context.Context, id uint, options *op.QueryOptions, selectFields ...string) (entity *T, err error) {
	isWriteAlways := isWriteLogWhenAlways(r)
	userID, _ := cu.GetUserIDFromContext(ctx)

	query := r.db.WithContext(ctx).Model(&r.entity).Where(c.FieldID+" = ?", id)
	if len(selectFields) > 0 {
		query = query.Select(selectFields)
	}

	query = unscopedQueryOptions(query, options)
	query = whereQueryOptions(query, options)
	query = orderByQueryOptions(query, options)
	query = preloadQueryOptions(query, options)
	sql := u.GetRawSqlFirst(query, &r.entity)

	if err = query.First(&entity).Error; err != nil {
		if isNotFoundError(err) {
			return nil, errors.New("No record found with this id")
		}

		if !isWriteAlways {
			l.Error("❌ Failed to get entity", "entityName", r.entityName, "error", err, "id", id, "userID", userID, "options", options, "entity", entity, "method", "GetSingleByID", "sql", sql)
		}
		return nil, e.WrapServerError("Failed to get record", err)
	}

	if isWriteAlways {
		l.Info("✅ Get entity", "entityName", r.entityName, "id", id, "userID", userID, "options", options, "entity", entity, "method", "GetSingleByID", "sql", sql)
	}
	return entity, nil
}

func (r *baseRepository[T]) GetAll(ctx context.Context, options *op.QueryOptions, selectFields ...string) (entities []T, err error) {
	isWriteAlways := isWriteLogWhenAlways(r)
	userID, _ := cu.GetUserIDFromContext(ctx)

	query := r.db.WithContext(ctx).Model(&entities)
	if len(selectFields) > 0 {
		query = query.Select(selectFields)
	}

	query = unscopedQueryOptions(query, options)
	query = whereQueryOptions(query, options)
	query = orderByQueryOptions(query, options)
	query = preloadQueryOptions(query, options)
	sql := u.GetRawSqlFind(query, &r.entity)

	if err = query.Find(&entities).Error; err != nil {
		if isNotFoundError(err) {
			return nil, errors.New("No record found")
		}

		if !isWriteAlways {
			l.Error("❌ Failed to get entities", "entityName", r.entityName, "error", err, "userID", userID, "options", options, "method", "GetAll", "sql", getRawSqlWithAnalyze(r, sql))
		}
		return nil, e.WrapServerError("Failed to get records", err)
	}

	if isWriteAlways {
		l.Info("✅ Get entities", "entityName", r.entityName, "userID", userID, "options", options, "method", "GetAll", "sql", getRawSqlWithAnalyze(r, sql))
	}
	return entities, nil
}

func (r *baseRepository[T]) GetAllPaginated(ctx context.Context, pageNo int, pageSize int, options *op.QueryOptions, selectFields ...string) (entities []T, count int64, err error) {
	isWriteAlways := isWriteLogWhenAlways(r)
	userID, _ := cu.GetUserIDFromContext(ctx)

	query := r.db.WithContext(ctx).Model(&entities)
	if len(selectFields) > 0 {
		query = query.Select(selectFields)
	}

	query = unscopedQueryOptions(query, options)
	query = whereQueryOptions(query, options)
	countSql := u.GetRawSqlCount(query, &count)

	if err := query.Count(&count).Error; err != nil {
		if isNotFoundError(err) {
			return nil, 0, errors.New("No record found")
		}

		if !isWriteAlways {
			l.Error("❌ Failed to count entities", "entityName", r.entityName, "error", err, "userID", userID, "options", options, "method", "GetAllPaginated", "countSql", countSql)
		}
		return nil, 0, e.WrapServerError("Failed to count records", err)
	}

	if isWriteAlways {
		l.Info("✅ Count entities", "entityName", r.entityName, "userID", userID, "options", options, "method", "GetAllPaginated", "countSql", countSql)
	}

	query = orderByQueryOptions(query, options)
	query = preloadQueryOptions(query, options)
	query = query.Offset(u.GetOffset(pageNo, pageSize)).Limit(pageSize)
	rawSql := u.GetRawSqlFind(query, &r.entity)

	if err := query.Find(&entities).Error; err != nil {
		if isNotFoundError(err) {
			return nil, 0, errors.New("No record found")
		}

		if !isWriteAlways {
			l.Error("❌ Failed to get entities", "entityName", r.entityName, "error", err, "userID", userID, "pageNo", pageNo, "pageSize", pageSize, "options", options, "method", "GetAllPaginated", "sql", getRawSqlWithAnalyze(r, rawSql))
		}
		return nil, 0, e.WrapServerError("Failed to get records", err)
	}

	if isWriteAlways {
		l.Info("✅ Get entities", "entityName", r.entityName, "userID", userID, "pageNo", pageNo, "pageSize", pageSize, "options", options, "method", "GetAllPaginated", "sql", getRawSqlWithAnalyze(r, rawSql))
	}
	return entities, count, nil
}

func (r *baseRepository[T]) GetDB() (db *gorm.DB) {
	return r.db
}

func (r *baseRepository[T]) WithTx(tx *gorm.DB) (repo BaseRepository[T]) {
	return &baseRepository[T]{
		db:         tx,
		entity:     r.entity,
		entityName: r.entityName,
		AppConfig:  r.AppConfig,
	}
}

func (r *baseRepository[T]) Transaction(ctx context.Context, fn func(r BaseRepository[T]) error) (err error) {
	isWriteAlways := isWriteLogWhenAlways(r)
	userID, _ := cu.GetUserIDFromContext(ctx)

	err = r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := r.WithTx(tx)
		return fn(txRepo)
	})

	if err != nil {
		if !isWriteAlways {
			l.Error("❌ Failed to execute transaction", "entityName", r.entityName, "error", err, "userID", userID, "method", "Transaction")
		}
		return err
	}

	if isWriteAlways {
		l.Info("✅ Execute transaction", "entityName", r.entityName, "userID", userID, "method", "Transaction")
	}
	return nil
}

func getRawSqlWithAnalyze[T any](r *baseRepository[T], rawSql string) any {
	if rawSql == "" {
		return nil
	}

	if rawSql != "" && strings.EqualFold(r.AppConfig.DBStatsLogType, c.DBStatsLogSqlAnalyze) {
		return &map[string]string{
			"rawSql":     rawSql,
			"analyzeSql": u.GetSqlExplainAnalyze(r.db, rawSql),
		}
	}

	return rawSql
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
		if orderClause := u.BuildSortingOrders(options.SortOptions, &options.SortableFields); orderClause != "" {
			db = db.Order(orderClause)
		}
	} else {
		if orderClause := u.BuildSortingOrder(options.SortBy, options.SortOrder, &options.SortableFields); orderClause != "" {
			db = db.Order(orderClause)
		}
	}
	return db
}

func isWriteLogWhenAlways[T any](r *baseRepository[T]) bool {
	return strings.EqualFold(r.AppConfig.WriteLogWhen, c.WriteLogWhenAlways)
}

func isNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
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
