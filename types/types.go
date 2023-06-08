package types

type Country struct {
	Id          int
	CountryName string
}

type University struct {
	Id             int
	CountryName    string
	UniversityName string
}

type RankingCriteria struct {
	Id           int
	SystemID     int
	CriteriaName string
}

type ChangeStudentStaffRatio struct {
	UniversityName string
	Year           int
	NewStaffRatio  int
}

type AddUniversityRankingYear struct {
	UniversityName string
	CriteriaName   string
	Year           int
	Score          int
}

type Publisher struct {
	Id            int
	PublisherName string
}

type ChangePublisher struct {
	NewPublisherName string
	PublisherName    string
}

type GamePublisher struct {
	GameName      string
	PublisherName string
}

type PlatformYear struct {
	Year int
}

type Category struct {
	CategoryName        string
	CategoryDescription string
}

type NewCategory struct {
	CategoryName           string
	NewCategoryName        string
	NewCategoryDescription string
}

type Product struct {
	ProductName        string
	ProductDescription string
	ProductPrice       float32
	ProductQuantity    int
}

type ProductsResponse struct {
	Products []Product `json:"products"`
}

type CategoriesResponse struct {
	Categories []Category `json:"products"`
}

type NewProduct struct {
	ProductName           string
	NewProductName        string
	NewProductDescription string
	NewProductPrice       int
	NewProductQuantity    int
}

type ProductCategory struct {
	ProductName  string
	CategoryName string
}

type Order struct {
	ProductName     string
	ProductQuantity int
}

type User struct {
	Id       int64
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Response struct {
	Message string `json:"message"`
}
