# Stock Dashboard API - Product Endpoints

## Product API Documentation

### Base URL

```
http://localhost:8080/api
```

### Authentication

Semua endpoint produk memerlukan JWT token dalam header:

```
Authorization: Bearer <your_jwt_token>
```

## Product Endpoints

### 1. Get All Products dengan Filter

**Endpoint:** `GET /api/products`

**Query Parameters (Optional):**

- `name` - Filter berdasarkan nama produk (pencarian partial/ILIKE)
- `category` - Filter berdasarkan kategori (pencarian partial/ILIKE)
- `min_price` - Harga minimum
- `max_price` - Harga maksimum
- `min_stock` - Stok minimum
- `max_stock` - Stok maksimum
- `limit` - Batasan jumlah hasil (pagination)
- `offset` - Offset untuk pagination

**Contoh Request:**

```bash
# Semua produk
GET /api/products

# Filter berdasarkan nama
GET /api/products?name=laptop

# Filter berdasarkan kategori
GET /api/products?category=elektronik

# Filter berdasarkan range harga
GET /api/products?min_price=100000&max_price=500000

# Filter berdasarkan stok
GET /api/products?min_stock=10

# Kombinasi filter dengan pagination
GET /api/products?category=elektronik&min_price=100000&limit=10&offset=0
```

**Response:**

```json
{
  "message": "Products fetched successfully",
  "products": [
    {
      "id": 1,
      "name": "Laptop Gaming",
      "price": 15000000,
      "stock": 5,
      "category": "Elektronik",
      "created_at": "2025-07-09T06:50:27Z",
      "updated_at": "2025-07-09T06:50:27Z"
    }
  ],
  "count": 1
}
```

### 2. Get Product by ID

**Endpoint:** `GET /api/products/:id`

**Contoh Request:**

```bash
GET /api/products/1
```

### 3. Search Products

**Endpoint:** `GET /api/products/search`

**Query Parameters:**

- `q` - Term pencarian (required) - mencari di nama dan kategori

**Contoh Request:**

```bash
GET /api/products/search?q=laptop
```

### 4. Create Product

**Endpoint:** `POST /api/products`

**Request Body:**

```json
{
  "name": "Laptop Gaming",
  "price": 15000000,
  "stock": 5,
  "category": "Elektronik"
}
```

### 5. Update Product

**Endpoint:** `PUT /api/products/:id`

**Request Body:**

```json
{
  "name": "Laptop Gaming Updated",
  "price": 16000000,
  "stock": 3,
  "category": "Elektronik"
}
```

### 6. Delete Product

**Endpoint:** `DELETE /api/products/:id`

## Contoh Penggunaan Filter

### 1. Filter Produk Elektronik dengan Harga di atas 1 juta

```bash
GET /api/products?category=elektronik&min_price=1000000
```

### 2. Filter Produk dengan Stok Rendah (kurang dari 10)

```bash
GET /api/products?max_stock=9
```

### 3. Filter Produk dengan Nama yang Mengandung "Laptop"

```bash
GET /api/products?name=laptop
```

### 4. Pagination - Halaman 2 dengan 5 item per halaman

```bash
GET /api/products?limit=5&offset=5
```

### 5. Filter Kompleks

```bash
GET /api/products?category=elektronik&min_price=500000&max_price=2000000&min_stock=1&limit=10
```

## Response Format

### Success Response

```json
{
  "message": "Success message",
  "products": [...],
  "count": 10
}
```

### Error Response

```json
{
  "message": "Error message",
  "error": "Detailed error description"
}
```

## Database Schema

### Products Table

```sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL CHECK (price > 0),
    stock INTEGER NOT NULL CHECK (stock >= 0),
    category VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Validasi

- `name`: Required, string
- `price`: Required, harus lebih besar dari 0
- `stock`: Required, harus lebih besar atau sama dengan 0
- `category`: Required, string
