package dal

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

type Book struct {
	ID     uint   `json:"id" bun:",pk,autoincrement"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type Product struct {
	ID        uint `bun:",pk,autoincrement"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Code      string
	Price     uint
}

type Bonus struct {
	Ename string `json:"ename"`
	Job   string `json:"job"`
	Sal   int    `json:"sal"`
	Comm  int    `json:"comm"`
}

type Dept struct {
	Deptno    int    `json:"deptno"`
	Dname     string `json:"dname"`
	Loc       string `json:"loc"`
}

var DB *bun.DB

// InitDatabase initializes the database and creates the application schema.
func InitDatabase() error {
	dbname := "gonic2.db"
	sqldb, err := sql.Open(sqliteshim.ShimName, dbname)
	if err != nil {
		return fmt.Errorf("open database %q: %w", dbname, err)
	}
	db := bun.NewDB(sqldb, sqlitedialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	log.Printf("[INFO] using db: %s", dbname)
	if _, err := db.NewCreateTable().Model((*Book)(nil)).IfNotExists().Exec(context.Background()); err != nil {
		db.Close()
		return fmt.Errorf("create book table: %w", err)
	}

	/*
		// Create
		db.NewInsert().Model(&Product{Code: "D42", Price: 100}).Exec(context.Background())
		//db.NewInsert().Model(&Book{Title: "Dune", Author: "John Adams"}).Exec(context.Background())

		// Read
		var product Product
		db.NewSelect().Model(&product).Where("id = ?", 1).Scan(context.Background())
	*/
	DB = db
	return nil
}
