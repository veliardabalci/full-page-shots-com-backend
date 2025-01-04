# Resmi Golang imajını kullan
FROM golang:1.23.4

# Gerekli bağımlılıkları yükle
RUN apt-get update && apt-get install -y \
    chromium \
    libx11-xcb1 \
    libxcomposite1 \
    libxcursor1 \
    libxdamage1 \
    libxi6 \
    libxtst6 \
    libnss3 \
    libcups2 \
    libxrandr2 \
    libasound2 \
    libatk1.0-0 \
    libpangocairo-1.0-0 \
    libxshmfence1 \
    libgbm1 \
    libgtk-3-0 \
    --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*

# Uygulama dosyalarını konteynıra kopyalamıyoruz, host dizinini bağlayacağız
# Çalışma dizinini ayarla (konteyner içinde /app olacak)
WORKDIR /app

# Go modüllerini indir
RUN go mod download

# Uygulamayı derle
RUN go build -o main .

# Portu aç
EXPOSE 8000

# Uygulamayı çalıştır
CMD ["./main"]
