# 🏗️ Architecture du projet Ecodrive

Ce projet suit une architecture **microservices** avec une séparation claire entre le frontend Vue.js et les services backend Go (un par domaine fonctionnel).

---

## 📁 Structure globale

```bash
├── frontend
│   ├── Dockerfile
│   ├── index.html
│   ├── jsconfig.json
│   ├── nginx.conf
│   ├── package-lock.json
│   ├── package.json
│   ├── public
│   │   ├── favicon.ico
│   │   └── logo-ecodrive.png
│   ├── README.md
│   ├── src
│   │   ├── App.vue
│   │   ├── assets
│   │   ├── components
│   │   ├── main.js
│   │   ├── router
│   │   └── views
│   └── vite.config.js
├── README.md
├── service-admin
│   ├── admin_test.go
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── handlers.go
│   ├── main.go
│   └── middleware.go
├── service-auth
│   ├── auth_test.go
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── handlers.go
│   └── main.go
├── service-driver
│   ├── Dockerfile
│   ├── driver_test.go
│   ├── go.mod
│   ├── go.sum
│   ├── handlers.go
│   ├── main.go
│   └── models.go
├── service-paiement
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── handlers.go
│   ├── main.go
│   ├── models.go
│   └── paiement_test.go
├── service-support
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── handlers.go
│   ├── main.go
│   ├── models.go
│   └── support_test.go
└── service-user
    ├── Dockerfile
    ├── go.mod
    ├── go.sum
    ├── handlers.go
    ├── main.go
    ├── models.go
    └── user_test.go
```

---

Chaque microservice est indépendant, avec son propre `main.go`, son `Dockerfile`, ses tests, et ses handlers. Le `frontend` utilise **Vite + Vue.js**.

