# Dokumentasi API - Modal POS Server

Dokumentasi ini berisi daftar lengkap endpoint yang tersedia di server Modal POS, termasuk cara penggunaan, struktur permintaan (request), dan respon (response).

## 📌 Informasi Umum

- **Base URL**: `http://<server-address>:<port>/api/v1`
- **Format Data**: Semua data dikirim dan diterima dalam format `application/json` (kecuali upload gambar).
- **Autentikasi**: Menggunakan JWT Bearer Token. Sertakan header `Authorization: Bearer <token>` pada endpoint yang terproteksi.

---

## 🔐 Autentikasi (Public)

### 1. Registrasi Bisnis Baru

Digunakan untuk mendaftarkan akun OWNER sekaligus membuat entitas Bisnis baru.

- **Endpoint**: `POST /auth/register`
- **Payload**:
  ```json
  {
    "email": "owner@example.com",
    "password": "securepassword",
    "business_name": "Toko Maju Bersama"
  }
  ```
- **Respon (201 Created)**:
  ```json
  {
    "message": "registrasi berhasil"
  }
  ```

### 2. Login

Mendapatkan Access Token dan Refresh Token.

- **Endpoint**: `POST /auth/login`
- **Payload**:
  ```json
  {
    "email": "owner@example.com",
    "password": "securepassword"
  }
  ```
- **Respon (200 OK)**:
  ```json
  {
    "access_token": "eyJhbG...",
    "refresh_token": "eyJhbG..."
  }
  ```

### 3. Refresh Token

Mendapatkan Access Token baru menggunakan Refresh Token yang masih valid.

- **Endpoint**: `POST /auth/refresh`
- **Payload**:
  ```json
  {
    "refresh_token": "eyJhbG..."
  }
  ```

### 4. Lupa Password (Forgot Password)

Meminta instruksi reset password. Untuk Owner akan dikirimkan OTP, untuk Staff akan dikirimkan notifikasi ke Owner.

- **Endpoint**: `POST /auth/forgot-password`
- **Payload**:
  ```json
  {
    "email": "user@example.com"
  }
  ```
- **Respon (200 OK)**:
  ```json
  {
    "message": "Instruksi reset password telah dikirim"
  }
  ```

### 5. Reset Password

Mengganti password menggunakan kode OTP (Hanya untuk Owner).

- **Endpoint**: `POST /auth/reset-password`
- **Payload**:
  ```json
  {
    "email": "owner@example.com",
    "code": "123456",
    "new_password": "newsecurepassword"
  }
  ```
- **Respon (200 OK)**: Sama dengan respon Login (Access & Refresh Token).

---

## 👤 Profil & Manajemen Bisnis (Protected)

### 1. Get My Profile

Mendapatkan detail profil user yang sedang login.

- **Endpoint**: `GET /me`
- **Auth**: Terproteksi

### 2. Update Informasi Bisnis (Owner Only)

- **Endpoint**: `PUT /business`
- **Payload**:
  ```json
  {
    "name": "Toko Maju Baru",
    "type": "RETAIL",
    "address": "Jl. Merdeka No. 123",
    "phone": "08123456789",
    "logo_url": "https://storage.com/logo.png"
  }
  ```

---

## 👥 Manajemen Staff (Protected - Owner Only)

### 1. Tambah Staff Baru

- **Endpoint**: `POST /staff`
- **Payload**:
  ```json
  {
    "email": "staff@example.com",
    "password": "staffpassword",
    "role": "STAFF"
  }
  ```

### 2. Ambil Daftar Staff

- **Endpoint**: `GET /staff`

---

## 📦 Manajemen Produk (Protected)

### 1. Upload Gambar Produk

Mengunggah file gambar dan mendapatkan URL publiknya.

- **Endpoint**: `POST /products/upload`
- **Content-Type**: `multipart/form-data`
- **Body**: `file` (Binary Image File)
- **Respon (200 OK)**:
  ```json
  {
    "url": "http://localhost:8080/uploads/unique-image.jpg"
  }
  ```

### 2. Tambah Produk Baru

- **Endpoint**: `POST /products`
- **Payload**:
  ```json
  {
    "name": "Kopi Latte",
    "description": "Kopi susu dengan foam lembut",
    "category": "Minuman",
    "image_url": "http://lo...",
    "variants": [
      {
        "name": "Regular",
        "price": 25000,
        "stock": 50,
        "sku": "KOPI-LATTE-REG"
      },
      {
        "name": "Large",
        "price": 32000,
        "stock": 30,
        "sku": "KOPI-LATTE-LRG"
      }
    ]
  }
  ```

### 3. Ambil Semua Produk

- **Endpoint**: `GET /products`

### 4. Detail Produk Berdasarkan ID

- **Endpoint**: `GET /products/:id`

### 5. Update Produk

- **Endpoint**: `PUT /products/:id`

### 6. Hapus Produk

- **Endpoint**: `DELETE /products/:id`

### 7. Restock Varian

- **Endpoint**: `PATCH /products/variants/:variantId/restock`
- **Payload**:
  ```json
  {
    "quantity": 10
  }
  ```

---

## 🧾 Transaksi / Point of Sale (Protected)

### 1. Checkout (Proses Transaksi)

Membuat transaksi baru dan mengurangi stok produk secara otomatis.

- **Endpoint**: `POST /transactions`
- **Payload**:
  ```json
  {
    "payment_method": "CASH",
    "amount_paid": 100000,
    "items": [
      {
        "product_id": "uuid-product-1",
        "variant_id": "uuid-variant-1",
        "quantity": 2
      }
    ]
  }
  ```

### 2. Ambil Riwayat Transaksi

- **Endpoint**: `GET /transactions`

---

## ⚠️ Status Codes Umum

- `200 OK`: Permintaan berhasil.
- `201 Created`: Data baru berhasil dibuat.
- `400 Bad Request`: Format data tidak valid atau kesalahan input.
- `401 Unauthorized`: Token tidak ada atau tidak valid.
- `403 Forbidden`: User tidak memiliki hak akses (contoh: Staff mencoba akses menu Owner).
- `404 Not Found`: Data tidak ditemukan.
- `409 Conflict`: Konflik data (contoh: Email atau SKU sudah terdaftar).
- `500 Internal Server Error`: Kesalahan pada server.
