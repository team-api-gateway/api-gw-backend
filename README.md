# API Gateway-Backend
Dies ist das Backend für das API-Gateway 

## Setup
Starten der Datenbank: ```docker-compose up```

Trage die Umgebungsvariablen in eine Datei ".env" ein (siehe .env.example)

Starten des Backends: ```go run cmd/main.go```


Erneuern der Spec: ```goas --module-path . --output swagger.json --debug --main-file-path ./cmd/main.go --handler-path .```

