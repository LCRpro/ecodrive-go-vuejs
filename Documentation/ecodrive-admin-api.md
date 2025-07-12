# 📘 API – Service `admin` (Ecodrive)

Ce service est réservé à la gestion des actions administratives : validation des demandes de chauffeur, consultation des utilisateurs, gestion des rôles.

---

## 🌐 Base URL

```
http://localhost:8003
https://admin-ecodrive.liamcariou.fr
```

---

## 📄 Endpoints disponibles

### 📄 GET `/admin/users`

Retourne la liste de tous les utilisateurs depuis le service user.

**Headers requis :**
- Authorization: Bearer `<JWT_TOKEN>` avec `ROLE_ADMIN`

**Réponse :**
```json
[
  {
    "id": 1,
    "name": "John",
    "roles": "[\"ROLE_PASSAGER\"]"
  }
]
```

---

### ✏️ PATCH `/admin/users/:id/accept-driver`

Attribue le rôle `ROLE_DRIVER` à un utilisateur et met à jour son véhicule.

**Body JSON :**
```json
{
  "car": "Tesla",
  "plate": "AB-123-CD"
}
```

**Réponse :**
```json
{
  "id": 1,
  "roles": "[\"ROLE_PASSAGER\", \"ROLE_DRIVER\"]",
  "car": "Tesla",
  "plate": "AB-123-CD"
}
```

---

### 📄 GET `/admin/driver-requests`

Récupère toutes les demandes de conducteurs depuis le service `user`.

**Réponse :**
```json
[
  {
    "user_id": "abc123",
    "car": "Toyota",
    "plate": "AA-999-AA",
    "status": "pending"
  }
]
```

---

### ✏️ PATCH `/admin/driver-requests/:id`

Répond à une demande conducteur avec `approve` ou `reject`.

**Body JSON :**
```json
{
  "action": "approve"
}
```

**Réponse :**
```json
{
  "id": "1",
  "status": "approved"
}
```

---

### ➕ POST `/admin/driver-requests`

Crée une nouvelle demande de conducteur (normalement utilisée par le frontend, mais exposée côté admin).

**Body JSON :**
```json
{
  "user_id": "abc123",
  "car": "Toyota",
  "plate": "AA-999-AA"
}
```

**Réponse :**
```json
{
  "user_id": "abc123",
  "car": "Toyota",
  "plate": "AA-999-AA",
  "status": "pending"
}
```

---

## 🧠 Middleware

Toutes les routes `/admin/**` sont protégées par :

- `AuthRequired()` → valide le token JWT
- `AdminRequired()` → vérifie le rôle `ROLE_ADMIN` dans les claims
