# Dockerfile: Frontend
FROM node:18.17.1-bullseye-slim AS build

WORKDIR /app

# Install pnpm secara global
RUN npm install -g pnpm

# Salin file package manager
COPY package.json pnpm-lock.yaml ./

# Install dependensi tanpa menjalankan script postinstall
RUN pnpm install --frozen-lockfile --ignore-scripts --no-optional
RUN pnpm rebuild

# Salin semua kode aplikasi ke dalam image
COPY . .

# Bangun aplikasi Nuxt
RUN pnpm build

# Ekspose port
EXPOSE 3000

# Mulai aplikasi
CMD ["pnpm", "start"]
