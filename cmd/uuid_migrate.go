package main

import (
	"log"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/internal/address"
	"github.com/easy-comerce/backend/internal/attribute_option"
	"github.com/easy-comerce/backend/internal/attribute_type"
	"github.com/easy-comerce/backend/internal/brand"
	"github.com/easy-comerce/backend/internal/category"
	"github.com/easy-comerce/backend/internal/cms/entity"
	"github.com/easy-comerce/backend/internal/color"
	"github.com/easy-comerce/backend/internal/contact_info"
	"github.com/easy-comerce/backend/internal/permission"
	"github.com/easy-comerce/backend/internal/product_stats"
	"github.com/easy-comerce/backend/internal/review"
	"github.com/easy-comerce/backend/internal/role"
	"github.com/easy-comerce/backend/internal/size_category"
	"github.com/easy-comerce/backend/internal/size_option"
	"github.com/easy-comerce/backend/internal/user"
	"github.com/easy-comerce/backend/internal/wishlist"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/tokenutil"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

func InitMigrateUUID() {
	// initUuidValues[notification.Notification]()

	initUuidValues[address.AddressEntity]()
	initUuidValues[attribute_option.AttributeOptionEntity]()
	initUuidValues[attribute_type.AttributeTypeEntity]()
	initUuidValues[brand.BrandEntity]()
	initUuidValues[category.CategoryEntity]()
	initUuidValues[color.ColorEntity]()
	initUuidValues[product_stats.ProductStatsEntity]()
	initUuidValues[review.ReviewEntity]()
	initUuidValues[size_category.SizeCategoryEntity]()
	initUuidValues[size_option.SizeOptionEntity]()
	initUuidValues[wishlist.WishlistEntity]()

	initUuidValues[entity.CmsPageEntity]()
	initUuidValues[entity.CmsSectionEntity]()
	initUuidValues[contact_info.ContactInfoEntity]()
	initUuidValues[permission.PermissionEntity]()
	initUuidValues[role.RoleEntity]()
	initUuidValues[user.UserEntity]()

	log.Printf("\n")
}

func initUuidValues[T any]() {
	log.Printf("\n")
	err := deleteUuidIfExists[T]()
	if err != nil {
		log.Printf("Execution stop at deleteUuidIfExists")
		return
	}

	err = createUuid[T]()
	if err != nil {
		log.Printf("Execution stop at createUuid")
		return
	}

	err = updateUuidValues[T]()
	if err != nil {
		log.Printf("Execution stop at updateUuidValues")
		return
	}

	err = setUuidUnique[T]()
	if err != nil {
		log.Printf("Execution stop at setUuidUnique")
		return
	}
}

func deleteUuidIfExists[T any]() error {
	db := db.GetDB()
	tableName := getTableName[T]()

	err := db.Exec(`
		ALTER TABLE ` + tableName + ` 
		DROP COLUMN IF EXISTS uuid;
	`).Error

	if err != nil {
		return err
	}

	err = db.Exec(`
    	ALTER TABLE ` + tableName + ` 
    	DROP CONSTRAINT IF EXISTS ` + tableName + `_uuid_unique;
	`).Error

	if err != nil {
		return err
	}

	log.Printf("-------> %v :: Dropped UUID column and CONSTRAINT from if it existed", tableName)
	return nil
}

func createUuid[T any]() error {
	db := db.GetDB()
	tableName := getTableName[T]()

	err := db.Exec(`
		ALTER TABLE ` + tableName + ` 
		ADD COLUMN IF NOT EXISTS uuid UUID;
	`).Error

	if err != nil {
		log.Printf("-------> %v :: Failed to setUuid %v", err, tableName)
		return err
	}

	err = db.Exec(`
		UPDATE ` + tableName + ` 
		SET uuid = NULL;
	`).Error

	if err != nil {
		log.Printf("-------> %v :: Failed to setUuid %v", err, tableName)
		return err
	}

	log.Printf("-------> %v :: Create first UUID column", tableName)
	return nil
}

func updateUuidValues[T any]() error {
	db := db.GetDB()
	tableName := getTableName[T]()

	var ids []uint
	if err := db.Model(new(T)).Unscoped().Where("uuid IS NULL").Pluck("id", &ids).Error; err != nil {
		log.Printf("-------> %v :: " + tableName + "Failed to get entities with NULL UUID: " + err.Error())
		return err
	}

	log.Printf("-------> %v :: Found %v rows for uuid", tableName, len(ids))
	count := 0

	for _, id := range ids {
		uuid, err := utils.RetryWithDelay(tokenutil.GenerateNewUUID, constants.RetryLimit, 10)
		if err != nil {
			log.Printf("-------> %v :: Failed to generate UUID: %v", err, tableName)
			return err
		}

		err = db.Model(new(T)).Unscoped().Where("id = ?", id).Update("uuid", uuid).Error
		if err != nil {
			log.Printf("-------> %v :: Failed to update entity ID %d: %v", id, err, tableName)
			return err
		}

		count++
	}

	log.Printf("-------> %v :: Update %v rows for uuid", tableName, count)
	return nil
}

func setUuidUnique[T any]() error {
	db := db.GetDB()
	tableName := getTableName[T]()

	err := db.Exec(`
		ALTER TABLE ` + tableName + ` 
		ALTER COLUMN uuid SET NOT NULL;
	`).Error

	if err != nil {
		log.Printf("-------> %v :: Failed to setUuidUnique %v", err, tableName)
		panic(err)
	}

	err = db.Exec(`
		ALTER TABLE ` + tableName + ` 
		ADD CONSTRAINT ` + tableName + `_uuid_unique UNIQUE (uuid);
	`).Error

	if err != nil {
		log.Printf("-------> %v :: Failed to setUuidUnique %v", err, tableName)
		panic(err)
	}

	log.Printf("-------> %v :: Done UUID column update and create constraint", tableName)
	return nil
}

func getTableName[T any]() string {
	db := db.GetDB()
	var entity T

	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(&entity); err != nil {
		panic(err)
	}

	return stmt.Schema.Table
}
