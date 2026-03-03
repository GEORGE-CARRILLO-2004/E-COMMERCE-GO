package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"golan/application/usecases"
	adapterHTTP "golan/infraestructure/in/http"
	"golan/infraestructure/out/persistence/mysql"
	"golan/infraestructure/out/security/jwt"
)

func main() {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root:manuelin2004@tcp(localhost:3306)/ecommerce_db?parseTime=true"
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Printf("Warning: failed to ping db %v\n", err)
	}

	customerRepo := mysql.NewCustomerRepository(db)
	productRepo := mysql.NewProductRepository(db)
	cartRepo := mysql.NewCartRepository(db)
	orderRepo := mysql.NewOrderRepository(db)
	paymentRepo := mysql.NewPaymentRepository(db)
	reviewRepo := mysql.NewReviewRepository(db)

	tokenProvider := jwt.NewJWTProvider("supersecretkey")

	customerUC := usecases.NewCustomerService(customerRepo, tokenProvider)
	productUC := usecases.NewProductService(productRepo, customerRepo)
	cartUC := usecases.NewCartService(cartRepo, productRepo, customerRepo)
	orderUC := usecases.NewOrderService(orderRepo, cartRepo, productRepo, paymentRepo, customerRepo)
	reviewUC := usecases.NewReviewService(reviewRepo, customerRepo, productRepo, orderRepo)

	customerHandler := adapterHTTP.NewCustomerHandler(customerUC, tokenProvider)
	productHandler := adapterHTTP.NewProductHandler(productUC, tokenProvider)
	cartOrderHandler := adapterHTTP.NewCartOrderHandler(cartUC, orderUC, tokenProvider)
	reviewHandler := adapterHTTP.NewReviewHandler(reviewUC, tokenProvider)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/register", customerHandler.Register)
	mux.HandleFunc("/api/login", customerHandler.Login)
	mux.HandleFunc("/api/profile", customerHandler.UpdateProfile)
	mux.HandleFunc("/api/password", customerHandler.ChangePassword)

	mux.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			productHandler.ListProducts(w, r)
		case http.MethodPost:
			productHandler.CreateProduct(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/cart", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			cartOrderHandler.GetCart(w, r)
		case http.MethodPost:
			cartOrderHandler.AddToCart(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/cart/remove", cartOrderHandler.RemoveFromCart)
	mux.HandleFunc("/api/cart/clear", cartOrderHandler.ClearCart)
	mux.HandleFunc("/api/cart/checkout", cartOrderHandler.CheckoutCart)

	mux.HandleFunc("/api/reviews", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			reviewHandler.GetProductReviews(w, r)
		case http.MethodPost:
			reviewHandler.CreateReview(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			cartOrderHandler.GetMyOrders(w, r)
		} else {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static"))))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/templates/index.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("GoShop E-Commerce (DDD/Hexagonal) => http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
