package db

import (
	"database/sql"
	"db_lab7/config"

	_ "github.com/mattn/go-sqlite3"
)

const (
	SelectAllProducts = `SELECT * FROM Products;`

	SelectAllCategories = `SELECT * FROM Categories;`

	SelectAllOrders = `SELECT Products.Name, Orders.Quantity
					   FROM Orders
					   JOIN Products ON Orders.ProductId = Products.Id;`

	GetProductQuantityQuery = `SELECT Quantity FROM Products WHERE Name = (?)`

	GetOrderQuantityQuery = `SELECT Quantity FROM Orders WHERE ProductId IN 
							(SELECT id FROM Products WHERE Name = (?))`

	AddCategoryQuery = `INSERT INTO Categories (Name, Description)
						VALUES ((?), (?));`

	AddProductQuery = `INSERT INTO Products (Name, Description, Price, Quantity)
					   VALUES ((?), (?), (?), (?));`

	AddProductCategoryQuery = `INSERT INTO Product_category (ProductId, CategoryId)
					   		   VALUES ((SELECT id FROM Products WHERE Name = (?)), 
								  	   (SELECT id FROM Categories WHERE Name = (?)));`

	AddOrderQuery = `INSERT INTO Orders (UserId, ProductId, Quantity) 
					 VALUES ((1), 
					 (SELECT id FROM Products WHERE Name = (?)), 
					 ((?)));`

	DeleteCategoryQuery = `DELETE FROM Product_category WHERE CategoryId IN 
						   (SELECT id FROM Categories WHERE Name = (?));
						   DELETE FROM Categories WHERE Name = (?);`

	DeleteProductQuery = `DELETE FROM Product_category WHERE ProductId IN 
						  (SELECT id FROM Products WHERE Name = (?));
						  DELETE FROM Products WHERE Name = (?);`

	DeleteProductCategoryQuery = `DELETE FROM Product_category WHERE 
								  ProductId IN (SELECT id FROM Products WHERE Name = (?)) AND
								  CategoryId IN (SELECT id FROM Categories WHERE Name = (?));`

	DeleteOrderQuery = `DELETE FROM Orders WHERE
						ProductId IN (SELECT id FROM Products WHERE Name = (?));`

	ChangeCategoryNameQuery        = `UPDATE Categories SET Name = (?) WHERE Name = (?);`
	ChangeCategoryDescriptionQuery = `UPDATE Categories SET Description = (?) WHERE Name = (?);`

	ChangeProductNameQuery        = `UPDATE Products SET Name = (?) WHERE Name = (?);`
	ChangeProductDescriptionQuery = `UPDATE Products SET Description = (?) WHERE Name = (?);`
	ChangeProductPriceQuery       = `UPDATE Products SET Price = (?) WHERE Name = (?);`
	ChangeProductQuantityQuery    = `UPDATE Products SET Quantity = (?) WHERE Name = (?);`

	CreateUserQuery = `INSERT INTO Users (UserName, Password, Email, Role)
						VALUES ((?),(?),(?),(?))`

	GetUserQuery = `SELECT * FROM Users WHERE UserName = (?) AND Password = (?);`
)

type Store struct {
	db  *sql.DB
	dsn string
}

func New(config *config.Config) *Store {
	return &Store{
		dsn: config.DSN,
	}
}

func (s *Store) Open() error {
	db, err := sql.Open("sqlite3", s.dsn)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *Store) Query(querySTR string, args ...any) (*sql.Rows, error) {
	return s.db.Query(querySTR, args...)
}

func (s *Store) Exec(querySTR string, args ...any) (sql.Result, error) {
	return s.db.Exec(querySTR, args...)
}

func (s *Store) Close() {
	s.db.Close()
}
