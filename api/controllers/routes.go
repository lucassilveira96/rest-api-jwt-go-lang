package controllers

import "github.com/lucassilveira96/rest-api-jwt-go-lang/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//Products routes
	s.Router.HandleFunc("/products", middlewares.SetMiddlewareJSON(s.CreateProduct)).Methods("POST")
	s.Router.HandleFunc("/products", middlewares.SetMiddlewareJSON(s.GetAllProducts)).Methods("GET")
	s.Router.HandleFunc("/products/{id}", middlewares.SetMiddlewareJSON(s.GetProduct)).Methods("GET")
	s.Router.HandleFunc("/products/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateProduct))).Methods("PUT")
	s.Router.HandleFunc("/products/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteProduct)).Methods("DELETE")

	//Product Categories routes
	s.Router.HandleFunc("/product/categories", middlewares.SetMiddlewareJSON(s.CreateProductCategory)).Methods("POST")
	s.Router.HandleFunc("/product/categories", middlewares.SetMiddlewareJSON(s.GetAllProductCategories)).Methods("GET")
	s.Router.HandleFunc("/product/categories/{id}", middlewares.SetMiddlewareJSON(s.GetProductCategory)).Methods("GET")
	s.Router.HandleFunc("/product/categories/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateProductCategory))).Methods("PUT")
	s.Router.HandleFunc("/product/categories/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteProductCategory)).Methods("DELETE")
}
