package routes

import (
	"net/http"
	"stock-dashboard/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	var filter models.ProductFilter

	if name := c.Query("name"); name != "" {
		filter.Name = name
	}
	if category := c.Query("category"); category != "" {
		filter.Category = category
	}
	if minPrice := c.Query("min_price"); minPrice != "" {
		if price, err := strconv.ParseFloat(minPrice, 64); err == nil {
			filter.MinPrice = price
		}
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if price, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filter.MaxPrice = price
		}
	}
	if minStock := c.Query("min_stock"); minStock != "" {
		if stock, err := strconv.Atoi(minStock); err == nil {
			filter.MinStock = stock
		}
	}
	if maxStock := c.Query("max_stock"); maxStock != "" {
		if stock, err := strconv.Atoi(maxStock); err == nil {
			filter.MaxStock = stock
		}
	}
	if sortOrder := c.Query("sort_order"); sortOrder != "" {
		if sortOrder == "asc" || sortOrder == "desc" {
			filter.SortOrder = sortOrder
		}
	}
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			filter.Limit = l
		}
	} else {
		filter.Limit = 10
	}
	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			filter.Offset = p
		}
	} else {
		filter.Offset = 1
	}

	result, err := models.GetProductsWithPagination(filter)
	if err != nil {
		response := models.NewErrorResponse("Failed to fetch products")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	data := gin.H{
		"products":   result.Products,
		"count":      len(result.Products),
		"page":       result.Page,
		"limit":      result.Limit,
		"total":      result.Total,
		"totalPages": result.TotalPages,
	}
	response := models.NewSuccessResponse(data, "Products fetched successfully")
	c.JSON(http.StatusOK, response)
}

func GetProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response := models.NewErrorResponse("Invalid product ID")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var product models.Product
	product.ID = id
	err = product.Get()
	if err != nil {
		response := models.NewErrorResponse("Product not found")
		c.JSON(http.StatusNotFound, response)
		return
	}

	data := gin.H{
		"product": product,
	}
	response := models.NewSuccessResponse(data, "Product fetched successfully")
	c.JSON(http.StatusOK, response)
}

func CreateProduct(c *gin.Context) {
	var product models.Product

	err := c.ShouldBindJSON(&product)
	if err != nil {
		response := models.NewErrorResponse("Invalid request format")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = product.Save()
	if err != nil {
		response := models.NewErrorResponse("Failed to create product")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	data := gin.H{
		"product": product,
	}
	response := models.NewSuccessResponse(data, "Product created successfully")
	c.JSON(http.StatusCreated, response)
}

func UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response := models.NewErrorResponse("Invalid product ID")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var product models.ProductUpdate
	err = c.ShouldBindJSON(&product)
	if err != nil {
		response := models.NewErrorResponse("Invalid request format")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	product.ID = id
	err = product.Update()
	if err != nil {
		response := models.NewErrorResponse("Failed to update product")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	data := gin.H{
		"product": product,
	}
	response := models.NewSuccessResponse(data, "Product updated successfully")
	c.JSON(http.StatusOK, response)
}

func DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response := models.NewErrorResponse("Invalid product ID")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var product models.Product
	product.ID = id
	err = product.Delete()
	if err != nil {
		response := models.NewErrorResponse("Failed to delete product")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	data := gin.H{
		"product_id": id,
	}
	response := models.NewSuccessResponse(data, "Product deleted successfully")
	c.JSON(http.StatusOK, response)
}

func SearchProducts(c *gin.Context) {
	searchTerm := c.Query("q")
	if searchTerm == "" {
		response := models.NewErrorResponse("Search term is required")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	sortOrder := c.Query("sort_order")
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	products, err := models.SearchProducts(searchTerm, sortOrder)
	if err != nil {
		response := models.NewErrorResponse("Failed to search products")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	data := gin.H{
		"products": products,
		"count":    len(products),
		"term":     searchTerm,
	}
	response := models.NewSuccessResponse(data, "Search completed successfully")
	c.JSON(http.StatusOK, response)
}
