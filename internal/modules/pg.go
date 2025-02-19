package modules

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/sullivtr/k8s_platform/internal/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PGSDK struct {
	db *gorm.DB
}

func InitPGDB(autoMigrate bool, host, username, password, dbName, environment string) PGSDK {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", username, password, host, dbName)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	if autoMigrate {
		migrate(db, environment)
	}

	return PGSDK{db: db}
}

func migrate(db *gorm.DB, environment string) {
	if err := db.AutoMigrate(
		&types.Group{},
		&types.Permission{},
		&types.User{},
		&types.GroupPermissions{},
		&types.GroupUsers{},
		&types.MySQLDBInfo{},
		&types.DynamicAppConfig{}); err != nil {
		log.Fatalln(err)
	}

	allAdminPermissionID, err := uuid.Parse("1b434611-5fe8-4ed0-b0b4-1307f9945b34")
	if err != nil {
		log.Fatalln(err)
	}
	allAdminPermission := types.Permission{
		ID:     &allAdminPermissionID,
		Name:   "AllAdmin",
		AppTag: "*",
	}

	if db.Migrator().HasTable(&types.Permission{}) {
		if err := db.First(&types.Permission{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			//Insert init data
			result := db.Create(&allAdminPermission)
			if result.Error != nil {
				log.Fatalln(err)
			}
		}
	}

	gid, err := uuid.Parse("1b434611-5fe8-4ed0-b0b4-1307f9945b34")
	if err != nil {
		log.Fatalln(err)
	}
	if db.Migrator().HasTable(&types.Group{}) {
		if err := db.First(&types.Group{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			//Insert init data
			groupZero := types.Group{
				ID:   &gid,
				Name: "Admin",
				Permissions: []*types.Permission{
					&allAdminPermission,
				},
			}
			result := db.Create(&groupZero)

			if result.Error != nil {
				log.Fatalln(err)
			}
		}
	}

	// Insert the default dynamic app config
	dynamicAppConfigDefault := types.DynamicAppConfig{
		ID: 1,
		Data: types.DynamicConfigJSONB{
			DefaultReplicaScaleLimit: 100,
			EnableK8sGlobalReadOnly:  true,
			K8sClusterName:           "updateme",
		},
	}

	if db.Migrator().HasTable(&types.DynamicAppConfig{}) {
		if err := db.Where("id = 1").First(&types.DynamicAppConfig{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			//Insert default data row
			result := db.Create(&dynamicAppConfigDefault)
			if result.Error != nil {
				log.Fatalln(err)
			}
		}
	}

	// create the hostdash tester user if the environment is development
	if environment == "Development" {
		uid, err := uuid.Parse("1b434611-5fe8-4ed0-b0b4-1307f9945b34")
		if err != nil {
			log.Fatalln(err)
		}
		if db.Migrator().HasTable(&types.User{}) {
			if err := db.First(&types.User{ID: &uid}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
				//Insert init data
				userZero := types.User{
					ID:    &uid,
					Name:  "hostdashtester",
					Email: "hostdashtester@gmail.com",
				}
				result := db.Create(&userZero)
				if result.Error != nil {
					log.Fatalln(err)
				}
			}
		}

		if db.Migrator().HasTable(&types.GroupUsers{}) {
			if err := db.First(&types.GroupUsers{UserID: uid, GroupID: gid}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
				//Insert init data
				groupUserZero := types.GroupUsers{
					UserID:  uid,
					GroupID: gid,
				}
				result := db.Create(&groupUserZero)
				if result.Error != nil {
					log.Fatalln(err)
				}
			}
		}

	}
}
