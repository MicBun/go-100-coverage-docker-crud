package database

import "gorm.io/gorm"

type TestFunc func(*gorm.DB)

func RunTest(testFunc TestFunc) {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	if err := Migrate(db); err != nil {
		panic(err)
	}

	tx := db.Begin()
	defer tx.Rollback()
	testFunc(tx)
}
