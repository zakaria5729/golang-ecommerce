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
		entityName: name,
	}
}

func (r *baseRepository[T]) Create(ctx context.Context, entity *T) (err error) {
	userID, _ := cu.GetUserIDFromContext(ctx)
	findAndSetActionValue(r, ctx, entity, c.FieldCreatedBy, userID)
	err = r.db.WithContext(ctx).Omit(c.FieldUpdatedAt).Create(entity).Error

	if err != nil {
		l.Error("❌ Failed to create entity", "entityName", r.entityName, "error", err, "userID", userID, "entity", entity, "method", "Create")
		return e.WrapServerError("Failed to create record", err)
	}

	return nil
}

func (r *baseRepository[T]) Update(ctx context.Context, id uint, entity *T) (err error) {
	userID, _ := cu.GetUserIDFromContext(ctx)
	findAndSetActionValue(r, ctx, entity, c.FieldUpdatedBy, userID)
	err = r.db.WithContext(ctx).Model(&entity).Where(c.FieldID+" = ?", id).Updates(entity).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("No record found with this id")
		} else {
			l.Error("❌ Failed to update entity", "entityName", r.entityName, "error", err, "id", id, "userID", userID, "entity", entity, "method", "Update")
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

	var entity T
	userID, _ := cu.GetUserIDFromContext(ctx)
	err = r.db.WithContext(ctx).Unscoped().Where(c.FieldID+" = ?", id).Delete(entity).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("No record found with this id")
		} else {
			l.Error("❌ Failed to hard delete entity", "entityName", r.entityName, "error", err, "id", id, "userID", userID, "entity", singleEntity, "method", "HardDelete")
			err = e.WrapServerError("Failed to delete record", err)
		}

		return err
	}

	return nil
}

func (r *baseRepository[T]) SoftDelete(ctx context.Context, id uint) (err error) {
	var entity T
	userID, _ := cu.GetUserIDFromContext(ctx)
	findAndSetActionValue(r, ctx, &entity, c.FieldDeletedBy, userID)
	findAndSetActionValue(r, ctx, &entity, c.FieldDeletedAt, timeutil.GormNowUTC())

	err = r.db.WithContext(ctx).Unscoped().Omit(c.FieldUpdatedAt).Model(&entity).Where(c.FieldID+" = ?", id).Updates(entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("No record found with this id")
		} else {
			l.Error("❌ Failed to soft delete entity", "entityName", r.entityName, "error", err, "id", id, "userID", userID, "entity", entity, "method", "SoftDelete")
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

	var entity T
	deletedAtExists := findAndSetActionValue(r, ctx, &entity, c.FieldDeletedAt, nil)
	deletedByExists := findAndSetActionValue(r, ctx, &entity, c.FieldDeletedBy, nil)
	query := r.db.WithContext(ctx).Unscoped().Model(&entity)

	if deletedAtExists && deletedByExists {
		query = query.Select(c.FieldDeletedAt, c.FieldDeletedBy)
	} else if deletedAtExists {
		query = query.Select(c.FieldDeletedAt)
	} else if deletedByExists {
		query = query.Select(c.FieldDeletedBy)
	}

	err = query.Omit(c.FieldUpdatedAt).Where(c.FieldID+" = ?", id).Updates(entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("No record found with this id")
		} else {
			l.Error("❌ Failed to undo soft delete entity", "entityName", r.entityName, "error", err, "id", id, "userID", userID, "entity", singleEntity, "method", "UndoSoftDelete")
			err = e.WrapServerError("Failed to undo deleted record", err)
		}

		return err
	}

	return nil
}

func (r *baseRepository[T]) ExistsByID(ctx context.Context, id uint, options *op.QueryOptions) (exists bool, err error) {
	entity := new(T)
	db := r.db.WithContext(ctx).Model(&entity)
	db = unscopedQueryOptions(db, options)
	db = whereQueryOptions(db, options)

	if err = db.Select(c.FieldID).Where(c.FieldID+" = ?", id).Take(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("No record found with this id")
		} else {
			l.Error("❌ Failed to check if entity exists", "entityName", r.entityName, "error", err, "id", id, "options", options, "entity", entity, "method", "ExistsByID")
			err = e.WrapServerError("Failed to check if record exists", err)
		}

		return false, err
	}

	return true, nil
}

func (r *baseRepository[T]) Count(ctx context.Context, options *op.QueryOptions) (count int64, err error) {
	entity := new(T)
	db := r.db.WithContext(ctx).Model(&entity)
	db = unscopedQueryOptions(db, options)
	db = whereQueryOptions(db, options)

	if err = db.Count(&count).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("No record found")
		} else {
			l.Error("❌ Failed to count entities", "entityName", r.entityName, "error", err, "options", options, "entity", entity, "method", "Count")
			err = e.WrapServerError("Failed to count records", err)
		}

		return 0, err
	}

	return count, nil
}

func (r *baseRepository[T]) GetSingleBy(ctx context.Context, options *op.QueryOptions, selectFields ...string) (entity *T, err error) {
	entity = new(T)
	db := r.db.WithContext(ctx).Model(&entity)
	if len(selectFields) > 0 {
		db = db.Select(selectFields)
	}

	db = unscopedQueryOptions(db, options)
	db = whereQueryOptions(db, options)
	db = orderByQueryOptions(db, options)
	db = preloadQueryOptions(db, options)

	if err = db.First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("No record found with this id")
		} else {
			l.Error("❌ Failed to get entity", "entityName", r.entityName, "error", err, "options", options, "entity", entity, "method", "GetSingleBy")
			err = e.WrapServerError("Failed to get record", err)
		}

		return nil, err
	}

	return entity, nil
}

func (r *baseRepository[T]) GetSingleByID(ctx context.Context, id uint, options *op.QueryOptions, selectFields ...string) (entity *T, err error) {
	entity = new(T)
	db := r.db.WithContext(ctx).Model(&entity).Where(c.FieldID+" = ?", id)
	if len(selectFields) > 0 {
		db = db.Select(selectFields)
	}

	db = unscopedQueryOptions(db, options)
	db = whereQueryOptions(db, options)
	db = orderByQueryOptions(db, options)
	db = preloadQueryOptions(db, options)

	if err = db.First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("No record found with this id")
		} else {
			l.Error("❌ Failed to get entity", "entityName", r.entityName, "error", err, "id", id, "options", options, "entity", entity, "method", "GetSingleByID")
			err = e.WrapServerError("Failed to get record", err)
		}

		return nil, err
	}

	return entity, nil
}

func (r *baseRepository[T]) GetAll(ctx context.Context, options *op.QueryOptions, selectFields ...string) (entities []T, err error) {
	entities = make([]T, 0)
	db := r.db.WithContext(ctx).Model(&entities)
	if len(selectFields) > 0 {
		db = db.Select(selectFields)
	}

	db = unscopedQueryOptions(db, options)
	db = whereQueryOptions(db, options)
	db = orderByQueryOptions(db, options)
	db = preloadQueryOptions(db, options)

	if err = db.Find(&entities).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("No record found")
		} else {
			l.Error("❌ Failed to get entities", "entityName", r.entityName, "error", err, "options", options, "method", "GetAll")
			err = e.WrapServerError("Failed to get records", err)
		}

		return nil, err
	}

	return entities, nil
}

func (r *baseRepository[T]) GetAllPaginated(ctx context.Context, pageNo int, pageSize int, options *op.QueryOptions, selectFields ...string) (entities []T, count int64, err error) {
	db := r.db.WithContext(ctx).Model(new(T))
	if len(selectFields) > 0 {
		db = db.Select(selectFields)
	}

	db = unscopedQueryOptions(db, options)
	db = whereQueryOptions(db, options)

	if err := db.Count(&count).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("No record found")
		} else {
			l.Error("❌ Failed to count entities", "entityName", r.entityName, "error", err, "pageNo", pageNo, "pageSize", pageSize, "options", options, "method", "GetAllPaginated")
			err = e.WrapServerError("Failed to count records", err)
		}

		return nil, 0, err
	}

	db = orderByQueryOptions(db, options)
	db = preloadQueryOptions(db, options)

	entities = make([]T, 0)
	if err := db.Offset(utils.GetOffset(pageNo, pageSize)).Limit(pageSize).Find(&entities).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("No record found")
		} else {
			l.Error("❌ Failed to get entities", "entityName", r.entityName, "error", err, "pageNo", pageNo, "pageSize", pageSize, "options", options, "method", "GetAllPaginated")
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

func findAndSetActionValue[T any](r *baseRepository[T], ctx context.Context, entity *T, field string, value any) (fieldExists bool) {
	stmt := &gorm.Statement{DB: r.db}
	stmt.Parse(entity)

	if field := stmt.Schema.LookUpField(field); field != nil {
		v := reflect.ValueOf(entity).Elem()
		field.Set(ctx, v, value)
		return true
	}

	return false
}
