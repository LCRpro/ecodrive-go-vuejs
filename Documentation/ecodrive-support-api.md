# 📘 API – Service `support` (Ecodrive)

Ce service gère les tickets d’assistance soumis par les utilisateurs (passagers ou chauffeurs) et permet aux administrateurs d’y répondre.

---

## 🌐 Base URL

```
http://localhost:8007
https://support-ecodrive.liamcariou.fr
```

---

## 📄 Endpoints disponibles

### ➕ POST `/support`

Crée un nouveau ticket de support.

**Body JSON :**
```json
{
  "google_id": "abc123",
  "category": "bug",
  "message": "J’ai un problème avec la course"
}
```

**Réponse :**
```json
{
  "id": 1,
  "google_id": "abc123",
  "category": "bug",
  "message": "J’ai un problème avec la course",
  "status": "ouvert",
  "created_at": "2025-07-12T12:34:56Z"
}
```

---

### 📄 GET `/support/user/:google_id`

Retourne tous les tickets créés par un utilisateur spécifique (via `google_id`).

**Exemple :**
```
GET /support/user/abc123
```

**Réponse :**
```json
[
  {
    "id": 1,
    "google_id": "abc123",
    "category": "bug",
    "message": "J’ai un problème",
    "status": "ouvert",
    "created_at": "...",
    "updated_at": "..."
  }
]
```

---

### 📄 GET `/support/all`

Retourne tous les tickets enregistrés (admin uniquement recommandé).

**Réponse :**
```json
[
  {
    "id": 1,
    "google_id": "abc123",
    "category": "bug",
    "message": "Message",
    "status": "ouvert",
    "created_at": "...",
    "updated_at": "..."
  }
]
```

---

### ✏️ PATCH `/support/reply/:id`

Permet à un administrateur de répondre à un ticket.

**Body JSON :**
```json
{
  "admin_reply": "Merci pour votre retour. Le problème est résolu.",
  "status": "fermé"
}
```

**Réponse :**
```json
{
  "id": 1,
  "google_id": "abc123",
  "category": "bug",
  "message": "J’ai un problème",
  "status": "fermé",
  "admin_reply": "Merci pour votre retour. Le problème est résolu.",
  "updated_at": "..."
}
```

---

## 🧠 Modèle `SupportTicket`

```go
type SupportTicket struct {
  ID         uint      `gorm:"primaryKey" json:"id"`
  GoogleID   string    `json:"google_id"`
  Category   string    `json:"category"`
  Message    string    `json:"message"`
  Status     string    `json:"status"`
  AdminReply string    `json:"admin_reply,omitempty"`
  CreatedAt  time.Time `json:"created_at"`
  UpdatedAt  time.Time `json:"updated_at"`
}
```
