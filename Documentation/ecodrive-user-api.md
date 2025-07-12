# 📘 API – Service user (Ecodrive)

Ce service gère les utilisateurs, les rôles (passager/admin/conducteur), les demandes de devenir chauffeur, et les soldes de l’application.

---


## 🌐 Base URL

```
http://localhost:8002  
https://user-ecodrive.liamcariou.fr
```


## 🔐 Authentification

Certaines routes exigent un JWT (Authorization: Bearer <token>) généré à la connexion ou via /update-token/:id.


## 📄 Endpoints disponibles

### ➕ POST /users/find-or-create

Crée un utilisateur ou met à jour un utilisateur existant selon son email ou son Google ID.

**Body JSON :**
```json
{
  "google_id": "abc123",
  "email": "user@example.com",
  "name": "John Doe"
}
```

Réponse :
```json
{
  "id": 1,
  "google_id": "abc123",
  "email": "user@example.com",
  "name": "John Doe",
  "roles": "[\"ROLE_PASSAGER\"]"
}
```


### 📄 GET /users

Retourne la liste de tous les utilisateurs.


### 📄 GET /users/:id

Retourne un utilisateur par google_id.


### ✏️ PATCH /users/:id

Met à jour les données d’un utilisateur par google_id.

```json
{
  "balance": 42.5
}
```


## 💰 Gestion du compte global (AppAccount)

### 📄 GET /app-balance

Retourne le solde global de l’application.


### 💸 PATCH /app-account/credit

Crédite le solde global.

```json
{
  "amount": 100.0
}
```


## 🚗 Demandes pour devenir conducteur

### 📝 POST /driver-requests

Soumet une demande pour devenir chauffeur.

```json
{
  "google_id": "abc123",
  "car": "Renault Clio",
  "plate": "AB-123-CD"
}
```


### 📄 GET /driver-requests

Liste toutes les demandes en attente ou traitées.


### ✅ PATCH /driver-requests/:id

Répond à une demande (admin uniquement).

```json
{
  "action": "accept" 
}
```


### ✏️ PATCH /driver-requests/:id/edit

Modifie les champs de la demande (admin).


## 🛡️ Gestion des rôles & tokens

### 🧙‍♂️ PATCH /become-admin

Ajoute le rôle ROLE_ADMIN à l'utilisateur identifié par le JWT fourni dans le header.


### 🔁 POST /update-token/:id

Génère un nouveau JWT pour un utilisateur donné (par google_id).

Réponse :
```json
{
  "token": "eyJhbGciOiJIUzI1...",
  "user": "..."
{ ... }
}
```


## 📦 Modèle User
```go
type User struct 
{
  ID        uint
  GoogleID  string
  Email     string
  Name      string
  Roles     string  
  Birthdate string
  Gender    string
  Phone     string
  Car       string
  Plate     string
  Balance   float64
}
```
