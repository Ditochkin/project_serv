package API

import (
	"db_lab7/config"
	"db_lab7/db"
	"db_lab7/types"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type API struct {
	config *config.Config
	router *mux.Router
	store  *db.Store
}

func InitApi() (*API, error) {
	res := new(API)
	var err error
	res.config, err = config.GetConfig()

	if err != nil {
		return nil, err
	}

	res.router = mux.NewRouter()
	return res, nil
}

func (a *API) Start() error {
	a.configureRouter()
	a.configureDB()
	fmt.Println(a.store.Open())

	return http.ListenAndServe(a.config.Port, a.router)
}

func (a *API) Stop() {
	a.store.Close()
}

func (a *API) configureDB() {
	a.store = db.New(a.config)
}

func (a *API) configureRouter() {
	a.router.HandleFunc("/create_user", a.handleCreateUser())
	a.router.HandleFunc("/sign_in", a.handleSignIn())
	a.router.HandleFunc("/sign_out", a.handleSignOut())

	a.router.HandleFunc("/get_products", a.handleGetAllProducts())
	a.router.HandleFunc("/get_categories", a.handleGetAllCategories())

	a.router.HandleFunc("/add_category", a.handleAddCategory())
	a.router.HandleFunc("/add_product", a.handleAddProduct())
	a.router.HandleFunc("/add_product_category", a.handleAddProductCategory())
	a.router.HandleFunc("/add_order", a.handleAddOrder())

	a.router.HandleFunc("/delete_category", a.handleDeleteCategory())
	a.router.HandleFunc("/delete_product", a.handleDeleteProduct())
	a.router.HandleFunc("/delete_product_category", a.handleDeleteProductCategory())
	a.router.HandleFunc("/delete_order", a.handleDeleteOrder())

	a.router.HandleFunc("/change_category_name", a.handleUpdateCategoryName())
	a.router.HandleFunc("/change_category_description", a.handleUpdateCategoryDescription())

	a.router.HandleFunc("/change_product_name", a.handleUpdateProductName())
	a.router.HandleFunc("/change_product_description", a.handleUpdateProductDescription())
	a.router.HandleFunc("/change_product_price", a.handleUpdateProductPrice())
	a.router.HandleFunc("/change_product_quantity", a.handleUpdateProductQuantity())
}

func (a *API) handleSignOut() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		c := &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		}
		http.SetCookie(writer, c)

		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) handleSignIn() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, _, err := a.GetIDAndRoleFromToken(writer, request)
		if err == nil {
			writer.WriteHeader(http.StatusOK)
			return
		}

		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "can't read body", http.StatusBadRequest)
			return
		}
		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "can't close body", http.StatusInternalServerError)
			return
		}
		var usr types.User
		err = json.Unmarshal(body, &usr)
		if err != nil {
			http.Error(writer, "can't close body", http.StatusInternalServerError)
			return
		}
		token, err := a.generateTokensByCred(usr.Username, usr.Password)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		setTokenCookies(writer, token)
		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) handleCreateUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, role, err := a.GetIDAndRoleFromToken(writer, request)
		if err != nil {
			http.Error(writer, "You are not logged in. Sign In please", http.StatusBadRequest)
			return
		}
		if role != "admin" {
			http.Error(writer, "You are not admin and you have no right for this act.", http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "error in reading body", http.StatusBadRequest)
			return
		}
		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "can't close body", http.StatusInternalServerError)
			return
		}
		var usr types.User
		err = json.Unmarshal(body, &usr)
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}
		_, err = a.store.Exec(db.CreateUserQuery, usr.Name, usr.Username, generatePasswordHash(usr.Password), usr.Role)
		if err != nil {
			if err.Error() == "UNIQUE constraint failed: users.Username" {
				http.Error(writer, "Username is already in use. Try to use another one.", http.StatusBadGateway)
				return
			}
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) GetProductQuantity(name string) (int, error) {
	fmt.Println(name)

	rows, err := a.store.Query(db.GetProductQuantityQuery, name)
	if err != nil {
		return 0, err
	}

	defer rows.Close()
	var quantity int

	for rows.Next() {
		err := rows.Scan(&quantity)
		fmt.Println(quantity)
		if err == nil {
			fmt.Println(err)
			return quantity, nil
		}
	}

	return 0, err
}

func (a *API) GetOrderQuantity(name string) (int, error) {
	fmt.Println("name: ", name)

	rows, err := a.store.Query(db.GetOrderQuantityQuery, name)
	if err != nil {
		return 0, err
	}

	defer rows.Close()
	var quantity int

	for rows.Next() {
		err := rows.Scan(&quantity)
		fmt.Println("Order q:  ", quantity)
		if err == nil {
			fmt.Println(err)
			return quantity, nil
		}
	}

	return 0, err
}

func (a *API) GetAllProducts() ([]types.Product, error) {
	rows, err := a.store.Query(db.SelectAllProducts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var id int
	var name string
	var description string
	var price float32
	var quantity int

	var products []types.Product

	for rows.Next() {
		err := rows.Scan(&id, &name, &description, &price, &quantity)
		if err != nil {
			fmt.Println(err)
			continue
		}

		product := types.Product{
			ProductName:        name,
			ProductDescription: description,
			ProductPrice:       price,
			ProductQuantity:    quantity,
		}

		products = append(products, product)

		fmt.Println(id, name)
	}

	return products, nil
}

func (a *API) handleGetAllProducts() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		products, err := a.GetAllProducts()
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		response := types.ProductsResponse{
			Products: products,
		}

		responseData, err := json.Marshal(response)
		if err != nil {
			http.Error(writer, "failed to marshal response", http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(responseData)
	}
}

func (a *API) GetAllCategories() ([]types.Category, error) {
	rows, err := a.store.Query(db.SelectAllCategories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var id int
	var name string
	var description string

	var categories []types.Category

	for rows.Next() {
		err := rows.Scan(&id, &name, &description)
		if err != nil {
			fmt.Println(err)
			continue
		}

		category := types.Category{
			CategoryName:        name,
			CategoryDescription: description,
		}

		categories = append(categories, category)

		fmt.Println(id, name)
	}

	return categories, nil
}

func (a *API) handleGetAllCategories() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		categories, err := a.GetAllCategories()
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		response := types.CategoriesResponse{
			Categories: categories,
		}

		responseData, err := json.Marshal(response)
		if err != nil {
			http.Error(writer, "failed to marshal response", http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(responseData)
	}
}

func (a *API) handleAddCategory() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// _, _, err := a.GetIDAndRoleFromToken(writer, request)
		// if err != nil {
		// 	http.Error(writer, "You are not logged in. Sign In please", http.StatusBadRequest)
		// 	return
		// }

		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "error in reading request", http.StatusBadRequest)
			return
		}

		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		var category types.Category
		err = json.Unmarshal(body, &category)
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		if category.CategoryName == "" {
			http.Error(writer, "CategoryName is empty", http.StatusInternalServerError)
			return
		}

		if category.CategoryDescription == "" {
			http.Error(writer, "CategoryDescription is empty", http.StatusInternalServerError)
			return
		}

		_, err = a.store.Exec(db.AddCategoryQuery, category.CategoryName, category.CategoryDescription)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		// Отправка JSON-файла в ответ
		response := types.Response{
			Message: "Success",
		}
		responseData, err := json.Marshal(response)
		if err != nil {
			http.Error(writer, "failed to marshal response", http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(responseData)
	}
}

func (a *API) handleAddProduct() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// _, _, err := a.GetIDAndRoleFromToken(writer, request)
		// if err != nil {
		// 	http.Error(writer, "You are not logged in. Sign In please", http.StatusBadRequest)
		// 	return
		// }

		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "error in reading request", http.StatusBadRequest)
			return
		}

		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		var product types.Product
		err = json.Unmarshal(body, &product)
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		if product.ProductName == "" {
			http.Error(writer, "ProductName is empty", http.StatusInternalServerError)
			return
		}

		if product.ProductDescription == "" {
			http.Error(writer, "ProductDescription is empty", http.StatusInternalServerError)
			return
		}

		if product.ProductPrice < 0 {
			http.Error(writer, "ProductPrice is negative", http.StatusInternalServerError)
			return
		}

		if product.ProductQuantity < 0 {
			http.Error(writer, "ProductQuantity is negative", http.StatusInternalServerError)
			return
		}

		_, err = a.store.Exec(db.AddProductQuery, product.ProductName, product.ProductDescription, product.ProductPrice, product.ProductQuantity)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) handleAddProductCategory() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// _, _, err := a.GetIDAndRoleFromToken(writer, request)
		// if err != nil {
		// 	http.Error(writer, "You are not logged in. Sign In please", http.StatusBadRequest)
		// 	return
		// }

		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "error in reading request", http.StatusBadRequest)
			return
		}

		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		var product_category types.ProductCategory
		err = json.Unmarshal(body, &product_category)
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		if product_category.ProductName == "" {
			http.Error(writer, "ProductName is empty", http.StatusInternalServerError)
			return
		}

		if product_category.CategoryName == "" {
			http.Error(writer, "CategoryName is empty", http.StatusInternalServerError)
			return
		}

		_, err = a.store.Exec(db.AddProductCategoryQuery, product_category.ProductName, product_category.CategoryName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) handleAddOrder() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// _, _, err := a.GetIDAndRoleFromToken(writer, request)
		// if err != nil {
		// 	http.Error(writer, "You are not logged in. Sign In please", http.StatusBadRequest)
		// 	return
		// }

		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "error in reading request", http.StatusBadRequest)
			return
		}

		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		var order types.Order
		err = json.Unmarshal(body, &order)
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		if order.ProductName == "" {
			http.Error(writer, "ProductName is empty", http.StatusInternalServerError)
			return
		}

		if order.ProductQuantity <= 0 {
			http.Error(writer, "ProductQuantity is not positive", http.StatusInternalServerError)
			return
		}

		curQuantity, err := a.GetProductQuantity(order.ProductName)
		fmt.Println(curQuantity)

		if curQuantity < order.ProductQuantity {
			http.Error(writer, "Current product quantity is less than requested", http.StatusInternalServerError)
			return
		}

		_, err = a.store.Exec(db.ChangeProductQuantityQuery, curQuantity-order.ProductQuantity, order.ProductName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = a.store.Exec(db.AddOrderQuery, order.ProductName, order.ProductQuantity)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) handleDeleteCategory() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// _, _, err := a.GetIDAndRoleFromToken(writer, request)
		// if err != nil {
		// 	http.Error(writer, "You are not logged in. Sign In please", http.StatusBadRequest)
		// 	return
		// }

		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "error in reading request", http.StatusBadRequest)
			return
		}

		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		var category types.Category
		err = json.Unmarshal(body, &category)
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		if category.CategoryName == "" {
			http.Error(writer, "CategoryName is empty", http.StatusInternalServerError)
			return
		}

		_, err = a.store.Exec(db.DeleteCategoryQuery, category.CategoryName, category.CategoryName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) handleDeleteProduct() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// _, _, err := a.GetIDAndRoleFromToken(writer, request)
		// if err != nil {
		// 	http.Error(writer, "You are not logged in. Sign In please", http.StatusBadRequest)
		// 	return
		// }

		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "error in reading request", http.StatusBadRequest)
			return
		}

		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		var product types.Product
		err = json.Unmarshal(body, &product)
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		if product.ProductName == "" {
			http.Error(writer, "ProductName is empty", http.StatusInternalServerError)
			return
		}

		_, err = a.store.Exec(db.DeleteProductQuery, product.ProductName, product.ProductName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) handleDeleteProductCategory() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "error in reading request", http.StatusBadRequest)
			return
		}

		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		var product_category types.ProductCategory
		err = json.Unmarshal(body, &product_category)
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		if product_category.ProductName == "" {
			http.Error(writer, "ProductName is empty", http.StatusInternalServerError)
			return
		}

		if product_category.CategoryName == "" {
			http.Error(writer, "CategoryName is empty", http.StatusInternalServerError)
			return
		}

		_, err = a.store.Exec(db.DeleteProductCategoryQuery, product_category.ProductName, product_category.CategoryName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) handleDeleteOrder() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// _, _, err := a.GetIDAndRoleFromToken(writer, request)
		// if err != nil {
		// 	http.Error(writer, "You are not logged in. Sign In please", http.StatusBadRequest)
		// 	return
		// }

		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "error in reading request", http.StatusBadRequest)
			return
		}

		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		var order types.Order
		err = json.Unmarshal(body, &order)
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		if order.ProductName == "" {
			http.Error(writer, "ProductName is empty", http.StatusInternalServerError)
			return
		}

		curQuantity, err := a.GetProductQuantity(order.ProductName)
		orderQuantity, err := a.GetOrderQuantity(order.ProductName)

		_, err = a.store.Exec(db.ChangeProductQuantityQuery, curQuantity+orderQuantity, order.ProductName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = a.store.Exec(db.DeleteOrderQuery, order.ProductName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) handleUpdateCategoryName() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// _, _, err := a.GetIDAndRoleFromToken(writer, request)
		// if err != nil {
		// 	http.Error(writer, "You are not logged in. Sign In please", http.StatusBadRequest)
		// 	return
		// }

		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "error in reading request", http.StatusBadRequest)
			return
		}

		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		var category types.NewCategory
		err = json.Unmarshal(body, &category)
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		if category.CategoryName == "" {
			http.Error(writer, "CategoryName is empty", http.StatusInternalServerError)
			return
		}

		if category.NewCategoryName == "" {
			http.Error(writer, "NewCategoryName is empty", http.StatusInternalServerError)
			return
		}

		_, err = a.store.Exec(db.ChangeCategoryNameQuery, category.NewCategoryName, category.CategoryName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) handleUpdateCategoryDescription() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// _, _, err := a.GetIDAndRoleFromToken(writer, request)
		// if err != nil {
		// 	http.Error(writer, "You are not logged in. Sign In please", http.StatusBadRequest)
		// 	return
		// }

		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "error in reading request", http.StatusBadRequest)
			return
		}

		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		var category types.NewCategory
		err = json.Unmarshal(body, &category)
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		if category.CategoryName == "" {
			http.Error(writer, "CategoryName is empty", http.StatusInternalServerError)
			return
		}

		if category.NewCategoryDescription == "" {
			http.Error(writer, "NewCategoryDescription is empty", http.StatusInternalServerError)
			return
		}

		_, err = a.store.Exec(db.ChangeCategoryDescriptionQuery, category.NewCategoryDescription, category.CategoryName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) handleUpdateProductName() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// _, _, err := a.GetIDAndRoleFromToken(writer, request)
		// if err != nil {
		// 	http.Error(writer, "You are not logged in. Sign In please", http.StatusBadRequest)
		// 	return
		// }

		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "error in reading request", http.StatusBadRequest)
			return
		}

		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		var product types.NewProduct
		err = json.Unmarshal(body, &product)
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		if product.ProductName == "" {
			http.Error(writer, "ProductName is empty", http.StatusInternalServerError)
			return
		}

		if product.NewProductName == "" {
			http.Error(writer, "NewProductName is empty", http.StatusInternalServerError)
			return
		}

		_, err = a.store.Exec(db.ChangeProductNameQuery, product.NewProductName, product.ProductName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) handleUpdateProductPrice() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// _, _, err := a.GetIDAndRoleFromToken(writer, request)
		// if err != nil {
		// 	http.Error(writer, "You are not logged in. Sign In please", http.StatusBadRequest)
		// 	return
		// }

		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "error in reading request", http.StatusBadRequest)
			return
		}

		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		var product types.NewProduct
		err = json.Unmarshal(body, &product)
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		if product.ProductName == "" {
			http.Error(writer, "ProductName is empty", http.StatusInternalServerError)
			return
		}

		if product.NewProductPrice < 0 {
			http.Error(writer, "NewProductPrice is negative", http.StatusInternalServerError)
			return
		}

		_, err = a.store.Exec(db.ChangeProductPriceQuery, product.NewProductPrice, product.ProductName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) handleUpdateProductQuantity() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// _, _, err := a.GetIDAndRoleFromToken(writer, request)
		// if err != nil {
		// 	http.Error(writer, "You are not logged in. Sign In please", http.StatusBadRequest)
		// 	return
		// }

		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "error in reading request", http.StatusBadRequest)
			return
		}

		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		var product types.NewProduct
		err = json.Unmarshal(body, &product)
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		if product.ProductName == "" {
			http.Error(writer, "ProductName is empty", http.StatusInternalServerError)
			return
		}

		if product.NewProductQuantity < 0 {
			http.Error(writer, "NewProductQuantity is negative", http.StatusInternalServerError)
			return
		}

		_, err = a.store.Exec(db.ChangeProductQuantityQuery, product.NewProductQuantity, product.ProductName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) handleUpdateProductDescription() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// _, _, err := a.GetIDAndRoleFromToken(writer, request)
		// if err != nil {
		// 	http.Error(writer, "You are not logged in. Sign In please", http.StatusBadRequest)
		// 	return
		// }

		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "error in reading request", http.StatusBadRequest)
			return
		}

		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		var product types.NewProduct
		err = json.Unmarshal(body, &product)
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		if product.ProductName == "" {
			http.Error(writer, "ProductName is empty", http.StatusInternalServerError)
			return
		}

		if product.NewProductDescription == "" {
			http.Error(writer, "NewProductDescription is empty", http.StatusInternalServerError)
			return
		}

		_, err = a.store.Exec(db.ChangeProductDescriptionQuery, product.NewProductDescription, product.ProductName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}
