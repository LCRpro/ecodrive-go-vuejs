# ecodrive

Ecodrive est une application de mobilité durable inspirée du modèle Uber/Taxi, développée avec une architecture microservices moderne. Elle permet aux passagers de commander une course en temps réel et aux conducteurs de répondre aux trajets proposés.  
Le projet repose sur les technologies suivantes :

## 🔗 Disponible

Le projet est accessible depuis :  
[https://ecodrive.liamcariou.fr/](https://ecodrive.liamcariou.fr/)


## 

- **Frontend** : Vue.js + Vite
- **Backend** : Go (Gin), avec microservices indépendants
- **Base de données** : MySQL
- **Authentification** : JWT, Oauth2
- **Déploiement** : Coolify + (CI/CD)
- **Cartographie** : Google Maps

---

## ⚙️ Lancer le projet en local

Chaque service est isolé et possède son propre environnement `.env`. Voici les étapes pour démarrer chaque service.

---

### 🧩 `service-auth` 

#### Commande pour lancer :
```bash
cd service-auth
go run .
```

#### Contenu du `.env` :
```
DATABASE_URL="mysql://user:mdp@host:port/nom?parseTime=true&loc=Local"
JWT_SECRET="jwt"
GOOGLE_CLIENT_ID="googleclientid"
GOOGLE_CLIENT_SECRET="googleclientsecret"
VITE_GOOGLE_MAPS_API_KEY=vitegoogleapikey
USER_SERVICE_URL=http://localhost:8002
FRONTEND_CALLBACK_URL=http://localhost:5173/callback
AUTH_REDIRECT_URL=http://localhost:8000/auth/callback
FRONTEND_ORIGIN=http://localhost:5173
```

---

### 🧩 `service-admin` 

#### Commande pour lancer :
```bash
cd service-admin
go run .
```

#### Contenu du `.env` :
```
DATABASE_URL="mysql://user:mdp@host:port/nom?parseTime=true&loc=Local"
JWT_SECRET="jwt"
GOOGLE_CLIENT_ID="googleclientid"
GOOGLE_CLIENT_SECRET="googleclientsecret"
VITE_GOOGLE_MAPS_API_KEY=vitegoogleapikey
USER_SERVICE_URL=http://localhost:8002
FRONTEND_ORIGIN=http://localhost:5173
```

---

### 🧩 `service-driver` 

#### Commande pour lancer :
```bash
cd service-driver
go run .
```

#### Contenu du `.env` :
```
DATABASE_URL="mysql://user:mdp@host:port/nom?parseTime=true&loc=Local"
JWT_SECRET="jwt"
GOOGLE_CLIENT_ID="googleclientid"
GOOGLE_CLIENT_SECRET="googleclientsecret"
VITE_GOOGLE_MAPS_API_KEY=vitegoogleapikey
USER_SERVICE_URL=http://localhost:8002
FRONTEND_ORIGIN=http://localhost:5173
```

---

### 🧩 `service-paiement` 

#### Commande pour lancer :
```bash
cd service-paiement
go run .
```

#### Contenu du `.env` :
```
DATABASE_URL="mysql://user:mdp@host:port/nom?parseTime=true&loc=Local"
JWT_SECRET="jwt"
GOOGLE_CLIENT_ID="googleclientid"
GOOGLE_CLIENT_SECRET="googleclientsecret"
VITE_GOOGLE_MAPS_API_KEY=vitegoogleapikey
FRONTEND_ORIGIN=http://localhost:5173
USER_SERVICE_URL=http://localhost:8002
```

---

### 🧩 `service-support` 

#### Commande pour lancer :
```bash
cd service-support
go run .
```

#### Contenu du `.env` :
```
DATABASE_URL="mysql://user:mdp@host:port/nom?parseTime=true&loc=Local"
JWT_SECRET="jwt"
GOOGLE_CLIENT_ID="googleclientid"
GOOGLE_CLIENT_SECRET="googleclientsecret"
VITE_GOOGLE_MAPS_API_KEY=vitegoogleapikey
FRONTEND_ORIGIN=http://localhost:5173
```

---


### 🧩 `service-user` 

#### Commande pour lancer :
```bash
cd service-user
go run .
```

#### Contenu du `.env` :
```
DATABASE_URL="mysql://user:mdp@host:port/nom?parseTime=true&loc=Local"
JWT_SECRET="jwt"
GOOGLE_CLIENT_ID="googleclientid"
GOOGLE_CLIENT_SECRET="googleclientsecret"
VITE_GOOGLE_MAPS_API_KEY=vitegoogleapikey
FRONTEND_ORIGIN=http://localhost:5173
```

---

### 🧩 `frontend` 

#### Commande pour lancer :
```bash
cd frontend
npm install
npm run dev
```

#### Contenu du `.env` :
```
VITE_GOOGLE_MAPS_API_KEY=googlemapsapikey
VITE_USER_SERVICE_URL=http://localhost:8002
VITE_PAIEMENT_SERVICE_URL=http://localhost:8004
VITE_DRIVER_SERVICE_URL=http://localhost:8006
VITE_AUTH_SERVICE_URL=http://localhost:8000
VITE_ADMIN_SERVICE_URL=http://localhost:8003
VITE_SUPPORT_SERVICE_URL=http://localhost:8007
```

---





## Informations

- Si vous avez besoin de plus d'informations sur le repository (Coolify, .env, etc.), n'hésitez pas à me contacter.
