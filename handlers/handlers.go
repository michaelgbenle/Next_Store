package handlers

import (
	"fmt"
	"github.com/decadevs/next_store/database"
	"github.com/decadevs/next_store/middleware"
	"github.com/decadevs/next_store/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//external database

//HANDLE WELCOME PAGE
func WelcomepageHandler(c *gin.Context) {

	//call the HTML Method of the context to render the template
	c.HTML(
		//setup the status of the template
		http.StatusOK,

		//use which template or deploy which template
		"index.html",
		gin.H{
			"title": "Next-Store"})

}

//HANDLER MARKETPLACE
func MarketPlaceHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "marketplace.html", gin.H{"message": "Market Place"})
}

var counter int

//HANDLER TO POST ON THE DATABASE
func AdminPostProductHandler(c *gin.Context) {
	check := middleware.Authentication(c)
	if check != false {
		//gets the data from form(front end)
		name := strings.TrimSpace(c.PostForm("product_Name"))
		quantity := strings.TrimSpace(c.PostForm("product_quantity_left"))
		price := strings.TrimSpace(c.PostForm("product_price"))
		cat := strings.TrimSpace(c.PostForm("category"))

		//converting the price & quantity to integer due to their format in the db
		p, _ := strconv.Atoi(price)
		q, _ := strconv.Atoi(quantity)
		prodImg := strings.TrimSpace(c.PostForm("product_img"))

		//ensure fields are filled before a product can be posted
		if name == "" || price == "" || quantity == "" || prodImg == "" {
			c.Redirect(http.StatusMovedPermanently, "/seller/addproductspage")
			return
		}

		//populating the model struct with the (model) values
		model := gorm.Model{
			ID:        0,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: nil,
		}

		//populating the product struct with the (product) values
		product := &models.Product{
			Model:                model,
			Name:                 name,
			Price:                p,
			Quantity:             q,
			Productcategory:      cat,
			Productimg:           prodImg,
			TotalProductLaunched: counter,
		}
		var keepRecordOfProduct []models.Product
		keepRecordOfProduct = append(keepRecordOfProduct, *product)

		for i := 0; i < len(keepRecordOfProduct); i++ {
			if keepRecordOfProduct != nil {
				counter++
			}
		}

		//initialize database to db
		db := c.MustGet("db").(*gorm.DB)
		//create a product table in database

		fmt.Println("this is the counter", counter)
		db.Create(&product)

		//redirect back to the seller addproduct page
		c.Redirect(http.StatusMovedPermanently, "/seller/addproductspage")
	} else {
		c.HTML(http.StatusOK, "seller_login.html", gin.H{
			"error": "Please Login",
		})
	}

}

//HANDLER TO GET THE SELLER PAGE
func AdminGetProductHandler(c *gin.Context) {
	check := middleware.Authentication(c)
	log.Println(check)
	log.Println("out")
	if check != false {
		log.Println("in")
		var products []models.Product
		//initialize the db
		db := c.MustGet("db").(*gorm.DB)
		//loops through the database & adds each instance to product slice
		db.Find(&products)

		//renders the products slice to the seller page
		c.HTML(http.StatusOK, "seller_page.html", gin.H{"data": products})

	} else {
		c.HTML(http.StatusOK, "seller_login.html", gin.H{
			"error": "Please Login",
		})
	}
}

//FUNCTION FOR THE SELLER TO DELETE PRODUCT
func AdminDeleteProductHandler(c *gin.Context) {
	check := middleware.Authentication(c)
	if check != false {
		// initialise the database
		db := c.MustGet("db").(*gorm.DB)

		err := db.Delete(&models.Product{}, c.Param("id")).Error
		if err != nil {
			return
		}

		c.Redirect(302, "/seller/addproductspage")
	} else {
		c.HTML(http.StatusOK, "seller_login.html", gin.H{
			"error": "Please Login",
		})
	}

}

//LAUNCH PRODUCT IN MARKET
func AdminPostInMarketHandler(c *gin.Context) {
	check := middleware.Authentication(c)
	if check != false {
		c.HTML(http.StatusOK, "seller_login.html", gin.H{
			"error": "Please Login",
		})
	}
}

func AdminLaunchProductHandler(c *gin.Context) {

	var products []models.Product
	//Get the page
	db := c.MustGet("db").(*gorm.DB)
	//finds the record(product data) from the database
	db.Find(&products)
	//sends the record(product data) to the market place page
	//c.HTML(http.StatusOK, "seller_page.html", gin.H{"data": products})
	c.HTML(http.StatusOK, "marketplace.html", gin.H{"data": products})

}

//SELLER EDIT PRODUCT
func SellerEditProductHandler(c *gin.Context) {
	check := middleware.Authentication(c)
	if check != false {
		db := c.MustGet("db").(*gorm.DB)
		var product models.Product
		if err := db.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
			return
		}

		c.HTML(http.StatusOK, "seller_editproduct.html", gin.H{"data": product})
	} else {
		c.HTML(http.StatusOK, "seller_login.html", gin.H{
			"error": "Please Login",
		})
	}

}

//SELLER UPDATE PRODUCT
func SellerUpdateProductHandler(c *gin.Context) {

	check := middleware.Authentication(c)
	if check != false {
		id := c.Param("id")
		name := strings.TrimSpace(c.PostForm("product_Name"))
		quantity := strings.TrimSpace(c.PostForm("product_quantity_left"))
		price := strings.TrimSpace(c.PostForm("product_price"))
		cat := strings.TrimSpace(c.PostForm("category"))
		p, _ := strconv.Atoi(price)
		q, _ := strconv.Atoi(quantity)
		prodImg := strings.TrimSpace(c.PostForm("product_img"))

		if name == "" || price == "" || quantity == "" || prodImg == "" {
			c.Redirect(301, "/sellers/editPost/{{.ID}}")
			return
		}

		product := &models.Product{
			Name:            name,
			Price:           p,
			Quantity:        q,
			Productcategory: cat,
			Productimg:      prodImg,
		}

		db := c.MustGet("db").(*gorm.DB)
		db.Model(models.Product{}).Where("id = ?", id).Updates(product, true)
		//db.Model(&product).Select("*").Update(models.Product{Name: name, Price: p, Quantity: q, Productcategory: cat, Productimg: prodImg})
		c.Redirect(302, "/seller/addproductspage")
	} else {
		c.HTML(http.StatusOK, "seller_login.html", gin.H{
			"error": "Please Login",
		})
	}

}

//BUYER PAGE
func BuyerPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "buyer_page.html", gin.H{
		"Message": "Product removed from cart Successfully",
	})
}

var storeCost []int
var TtCost int

//HANDLER ORDER PRODUCT
func AddToCartHandler(c *gin.Context) {
	check := middleware.Authentication(c)
	if check != false {
		//get product and store in database
		var products models.Product

		user := &models.User{ID: 5}
		fmt.Println("This is the user==>", user)

		userID := user.ID

		db := c.MustGet("db").(*gorm.DB)

		db.Find(&products).Where("id = ?", "products_id")

		storeCost = append(storeCost, products.Price)

		for _, val := range storeCost {
			TtCost += val
		}

		fmt.Println("this is the total cost", TtCost)

		carts := &models.Cart{
			Id:         0,
			Name:       products.Name,
			Quantity:   products.Quantity,
			Price:      products.Price,
			Image:      products.Productimg,
			ProductID:  products.ID,
			BuyerID:    userID,
			Buyer:      models.Buyer{},
			TotalPrice: TtCost,
		}

		err := db.Create(carts).Error

		if err != nil {
			log.Println(err)
		}

		fmt.Println("this is the cart total ", carts.TotalPrice)

		//display product in cart
		var cart []models.Cart

		id := carts.ProductID

		db.Find(&cart).Where("id = ?", id)

		c.HTML(http.StatusOK, "buyer_page.html", gin.H{
			"order": cart,
		})

	} else {
		c.HTML(http.StatusOK, "buyer_signup.html", gin.H{
			"error": "Please Login",
		})
	}

}

//REMOVE PRODUCT FROM CART
func RemoveProductFromCartHandler(c *gin.Context) {
	check := middleware.Authentication(c)
	if check != false {
		//initialise the database
		db := c.MustGet("db").(*gorm.DB)

		err := db.Delete(&models.Cart{}, c.Param("product_id")).Error
		if err != nil {
			return
		}
		c.Redirect(302, "/buyerpage")
	} else {
		c.HTML(http.StatusOK, "buyer_signup.html", gin.H{
			"error": "Please Login",
		})
	}

}

//HANDLER FOR PAYMENT DETAILS
func PaymentHandler(c *gin.Context) {
	check := middleware.Authentication(c)
	if check != false {
		var cart models.Cart

		db := c.MustGet("db").(*gorm.DB)

		db.Find(&cart).Where("id = ?", "product_id")

		c.HTML(http.StatusOK, "buyerpayment.html", gin.H{
			"order": cart,
		})
	} else {
		c.HTML(http.StatusOK, "buyer_signup.html", gin.H{
			"error": "Please Login or create an account",
		})
	}

}

//HANDLER SELLER SIGN UP PAGE
func SellerLoginPageHandler(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"seller_login.html",
		gin.H{
			"message": "Seller Login"})
}

//HANDLER BUYER SIGN UP PAGE
func BuyerSignUpPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "buyer_signup.html", gin.H{"products": "Buyer Sign In"})
}

//BUYER SIGN-UP HANDLER
func BuyerSignUpHandler(c *gin.Context) {
	// get the user from the form and populate the user struct
	user := &models.User{}
	user.Name = c.PostForm("name")
	user.Password = c.PostForm("password")
	user.Address = c.PostForm("address")
	user.Username = c.PostForm("username")
	user.Email = c.PostForm("email")
	user.Phonenumber = c.PostForm("phonenumber")

	// check the database if email exists
	_, err := database.FindUserByEmail(user.Email)
	if err == nil {
		log.Println("finding email", err)
		c.HTML(http.StatusOK, "seller_login.html", gin.H{
			"error": "user with this email already exits, please login",
		})
		return
	}
	// hashing the password
	user.PasswordHash = user.PasswordHasher()
	// add the user to the db if email does not exist and password is hashed
	err = database.CreateNewUser(user)
	if err != nil {
		log.Println("creating user", err)
		c.HTML(http.StatusOK, "seller_login.html", gin.H{
			"error": "internal server error",
		})
		return
	}

	// encoding the cookie string
	cookiePart := rand.Int()
	cookie := []byte(user.Email)

	// setting a cookie for the user
	c.SetCookie("seasalt", fmt.Sprintf("%v awesome : %v : %v", user.ID, cookie, cookiePart), 3600*24, "", "", true, true)

	// user stored to database and user redirected to the homepage
	c.Redirect(http.StatusMovedPermanently, "/sellers/launchproduct")
	return
}

func SellerLoginHandler(c *gin.Context) {
	// grt the seller from the form and populate the seller struct
	seller := &models.Seller{}
	seller.Password = c.PostForm("password")
	seller.Email = c.PostForm("email")

	userDB, err := database.FindSellerByEmail(seller.Email)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusOK, "seller_login.html", gin.H{
			"error": "invalid email",
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(userDB.PasswordHash), []byte(seller.Password))
	if err != nil {
		log.Printf("error validating password :%v", err)
		c.HTML(http.StatusOK, "seller_login.html", gin.H{
			"error": "invalid password",
		})
		return
	}

	//setting the cookie for the user

	// encoding the cookie string
	cookiePart := rand.Int()
	cookie := []byte(seller.Email)

	// setting a cookie for the user
	c.SetCookie("seasalt", fmt.Sprintf("%v awesome : %v : %v", userDB.ID, cookie, cookiePart), 3600*24, "", "", true, true)

	c.Redirect(302, "/seller/dashboard")
	return
}

func LoginHandler(c *gin.Context) {
	// get the user from the form and populate the user struct
	user := &models.User{}
	user.Password = c.PostForm("password")
	user.Email = c.PostForm("email")

	userDB, err := database.FindUserByEmail(user.Email)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusOK, "buyer_login.html", gin.H{
			"error": "invalid email",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDB.PasswordHash), []byte(user.Password))
	if err != nil {
		log.Printf("error validating password :%v", err)
		c.HTML(http.StatusOK, "buyer_login.html", gin.H{
			"error": "invalid password",
		})
		return
	}

	// encoding the cookie string
	cookiePart := rand.Int()
	cookie := []byte(user.Email)

	// setting a cookie for the user
	c.SetCookie("seasalt", fmt.Sprintf("%v awesome : %v : %v", userDB.ID, cookie, cookiePart), 3600*24, "", "", true, true)

	c.Redirect(302, "/sellers/launchproduct")
	return
}

func LogoutUserHandler(c *gin.Context) {
	user := &models.User{}
	cookiePart := rand.Int()
	cookie := []byte(user.Email)
	c.SetCookie("seasalt", fmt.Sprintf("awesome : %v : %v", cookie, cookiePart), -1, "", "", true, true)
	c.Redirect(http.StatusFound, "/")
	return
}

//HANDLER SELLER PAGE
func SellerPageHandler(c *gin.Context) {

	c.HTML(
		http.StatusOK,
		"seller_page.html",
		gin.H{
			"seller": "Seller Page"})
}

func SearchProduct(c *gin.Context) {
	check := middleware.Authentication(c)
	if check != false {
		//search based on product names
		productName := c.Query("product_Name")

		c.HTML(http.StatusOK, "marketplace.html", gin.H{
			"searchproduct": productName,
		})
	} else {
		c.HTML(http.StatusOK, "buyer_signup.html", gin.H{
			"error": "Please Login or create an account",
		})
	}

}

//SELLER DASHBOARD HANDLER
func AdminDashBoard(c *gin.Context) {
	check := middleware.Authentication(c)
	if check != false {
		db := c.MustGet("db").(*gorm.DB)
		//Seller details
		UserOne := models.User{
			ID:            1,
			Name:          os.Getenv("SELLER_NAME"),
			Email:         os.Getenv("SELLER_EMAIL"),
			Username:      os.Getenv("SELLER_USERNAME"),
			Password:      os.Getenv("SELLER_PASSWORD"),
			Address:       os.Getenv("SELLER_ADDRESS"),
			AccountName:   os.Getenv("SELLER_ACCOUNTNAME"),
			AccountNumber: os.Getenv("SELLER_ACCOUNT_NUMBER"),
			Phonenumber:   os.Getenv("SELLER_PHONENUMBER"),
			BankName:      os.Getenv("SELLER_BANKNAME"),
			PasswordHash:  "",
			TimeCreated:   time.Now().Format("20-02-2021, 23:12"),
		}
		UserOne.PasswordHash = UserOne.PasswordHasher()

		//CREATE SELLER
		var Seller = models.Seller{
			UserOne,
			1,
		}

		db.First(&Seller)

		c.HTML(http.StatusOK, "sellerdashboard.html", gin.H{
			"admin": Seller,
		})
		return
	} else {
		c.HTML(http.StatusOK, "seller_login.html", gin.H{
			"error": "Please Login",
		})
		return
	}
}
