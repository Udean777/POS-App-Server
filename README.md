# Modal POS - Server (Backend)

Sisi server dari aplikasi Modal POS, dibangun menggunakan bahasa pemrograman Go dengan fokus pada performa tinggi, konkurensi, dan skalabilitas.

## 🚀 Tech Stack

Server ini menggunakan library dan framework pilihan untuk memastikan keamanan dan stabilitas data:

- **Language**: [Go](https://go.dev/) (v1.26.1+)
- **Web Framework**: [Gin Gonic](https://gin-gonic.com/)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **Autentikasi**: [JWT](https://github.com/golang-jwt/jwt) (Access & Refresh Tokens)
- **Storage**: [Cloudflare R2](https://www.cloudflare.com/products/r2/) / [AWS S3](https://aws.amazon.com/s3/) (dengan fallback ke Local Storage)
- **Live Reload**: [Air](https://github.com/air-verse/air)

## ✨ Fitur Utama

1.  **Manajemen Bisnis**: Pendaftaran owner dan pengaturan informasi toko secara dinamis.
2.  **Manajemen Staff**: Sistem role-based access control (RBAC) untuk Owner dan Staff.
3.  **Katalog Produk**: Manajemen produk dengan dukungan banyak varian (size, warna, dll), SKU unik, dan kategori.
4.  **Sistem Point of Sale (POS)**: Pemrosesan transaksi checkout dengan manajemen stok otomatis (atomik).
5.  **Analytics & History**: Pelacakan riwayat transaksi per bisnis.
6.  **Penyimpanan Gambar**: Integrasi dengan Cloud Storage untuk penyimpanan aset produk yang efisien.
7.  **Sistem Reset Password**: Fitur lupa password dengan sistem OTP untuk Owner dan sistem ajuan untuk Staff.

## 🏗️ Arsitektur Proyek

Backend ini menerapkan **Clean Architecture** untuk memisahkan tanggung jawab antara layer:

- `cmd/api/`: Titik masuk aplikasi (main.go).
- `internal/delivery/`: Layer komunikasi (HTTP Handler & Middleware).
- `internal/usecase/`: Layer logika bisnis utama.
- `internal/domain/`: Layer definisi entitas, struct, dan interface.
- `internal/repository/`: Layer akses data (GORM PostgreSQL).
- `pkg/`: Library pendukung yang bersifat general.

## 🛠️ Cara Menjalankan

### Persiapan

1. Pastikan Go (v1.26+) telah terinstal.
2. Pastikan database PostgreSQL sudah berjalan.
3. Buat file `.env` di direktori ini berdasarkan contoh berikut:

```env
PORT=8080
DB_URL=postgres://user:password@localhost:5432/pos_db?sslmode=disable
JWT_SECRET=rahasia-anda
APP_URL=http://localhost:8080

# Cloud Storage (Opsional)
R2_ACCOUNT_ID=...
R2_ACCESS_KEY_ID=...
R2_SECRET_ACCESS_KEY=...
R2_BUCKET_NAME=...
R2_PUBLIC_DOMAIN=...
```

### Jalankan untuk Development

Gunakan **Air** untuk fitur live reload otomatis saat ada perubahan kode:

```bash
air
```

Atau jalankan secara manual:

```bash
go run cmd/api/main.go
```

## 📖 Dokumentasi API

Untuk detail mengenai endpoint, payload, dan respon, silakan baca:
👉 **[DOKUMENTASI_API.md](./DOKUMENTASI_API.md)**

---

> [!IMPORTANT]
> Migrasi database dilakukan secara otomatis saat aplikasi pertama kali dijalankan melalui fungsi `db.AutoMigrate()`.
