# Usa una immagine di base con Go preinstallato
FROM golang:latest

# Imposta la directory di lavoro nel percorso del codice Go
WORKDIR /go/src/client

# Copia il codice sorgente del servizio Go RPC nella directory di lavoro del container
COPY . .

RUN go build -o client ./client

