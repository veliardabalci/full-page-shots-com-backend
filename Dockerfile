# Resmi Golang imajını kullan
FROM golang:1.23.4

# Çalışma dizinini ayarla
WORKDIR /app

# Gerekli dosyaları kopyala (özellikle go.mod ve go.sum dosyalarını)
COPY go.mod go.sum ./

# Go modüllerini indir
RUN go mod download

# Uygulama dosyalarını konteynıra kopyala
COPY . .

# Uygulamayı derle
RUN go build -o main .

# Portu aç
EXPOSE 8000

# Uygulamayı çalıştır
CMD ["./main"]
