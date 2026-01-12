# Mekari Expense Management System

## Fitur Utama

- **Access Control (Role-Based)**: 
  - **Employee**: Hanya dapat melihat dan mengelola data pengeluaran milik sendiri.
  - **Manager**: Memiliki hak akses penuh untuk melihat seluruh data pengeluaran karyawan dan melakukan approval/rejection.
- **Business Logic Approval**:
  - Pengeluaran **< Rp 1.000.000**: Otomatis disetujui (*Auto-Approved*).
  - Pengeluaran **>= Rp 1.000.000**: Membutuhkan persetujuan manual dari Manager (*Pending*).
- **Pagination & Server-Side Filtering**: Navigasi data menggunakan limit dan offset untuk performa optimal di sisi backend.

## Tech Stack

### Backend
- **Language**: Golang 1.21+
- **Framework**: Standard `net/http` dengan GORM
- **Database**: PostgreSQL / MySQL (GORM supported)
- **CORS Handling**: Custom Middleware untuk header kustom (`X-User-ID`, `X-User-Role`)

### Frontend
- **Framework**: Nuxt 3 (Vue.js)
- **Styling**: Tailwind CSS & Nuxt UI
- **State Management**: Vue Refs & Cookies untuk session token

## Prasyarat Jalankan Aplikasi

1. **Go** terinstal di mesin Anda.
2. **Node.js (LTS)** terinstal.
3. Database (PostgreSQL/MySQL) sudah berjalan.

### 1. Backend (Golang)
Masuk ke folder backend, sesuaikan konfigurasi database di file `.env` atau `config.go`, lalu jalankan:
```bash
go mod tidy
go run main.go
```

Pastikan server berjalan di http://localhost:8080.

### 2. Frontend (Nuxt 3)
Masuk ke folder frontend, lalu jalankan perintah berikut:
```bash
pnpm install
pnpm run dev
```

Aplikasi akan tersedia di http://localhost:3000.

### 3. Skenario Pengujian
Untuk menguji fitur Access Control, gunakan kredensial simulasi berikut:

| Role | Email | Password | Hak Akses | 
| ----------- | ----------- | ----------- | ----------- |
| Manager | manager@mekari.com | password123 | Lihat semua data & Approve/Reject |
| Employee | employee@mekari.com | password123 | Hanya lihat data sendiri |


---
Mekari Challenge - Submission oleh [Arifful Fikri]