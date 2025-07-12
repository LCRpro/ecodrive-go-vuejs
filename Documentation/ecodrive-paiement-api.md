# 📘 API – Service `paiement` (Ecodrive)

Ce service gère les dépôts, les retraits, et les historiques de transactions des utilisateurs.

---

## 🌐 Base URL

```
http://localhost:8004
https://paiement-ecodrive.liamcariou.fr
```

---

## 📄 Endpoints disponibles

### ➕ POST `/deposit`

Effectue un dépôt sur le compte d’un utilisateur.

**Body JSON :**
```json
{
  "google_id": "abc123",
  "amount": 50,
  "card_number": "1234567812345678",
  "expiry": "12/29",
  "cvv": "123"
}
```

**Réponse :**
```json
{
  "id": 1,
  "google_id": "abc123",
  "type": "deposit",
  "amount": 50,
  "created_at": "2025-07-12T12:34:56Z"
}
```

---

### ➖ POST `/withdraw`

Effectue un retrait vers un IBAN si le solde est suffisant.

**Body JSON :**
```json
{
  "google_id": "abc123",
  "amount": 30,
  "iban": "FR7630006000011234567890189"
}
```

**Réponse :**
```json
{
  "id": 2,
  "google_id": "abc123",
  "type": "withdrawal",
  "amount": 30,
  "iban_mask": "FR**********0189",
  "created_at": "2025-07-12T13:00:00Z"
}
```

---

### 📄 GET `/transactions/:google_id`

Retourne toutes les transactions de l'utilisateur spécifié.

**Exemple :**
```
GET /transactions/abc123
```

**Réponse :**
```json
[
  {
    "id": 1,
    "google_id": "abc123",
    "type": "deposit",
    "amount": 50,
    "created_at": "..."
  },
  {
    "id": 2,
    "google_id": "abc123",
    "type": "withdrawal",
    "amount": 30,
    "iban_mask": "FR**********0189",
    "created_at": "..."
  }
]
```

---

### 📄 GET `/transactions`

Retourne toutes les transactions enregistrées (admin recommandé).

---

## 🧠 Modèle `Transaction`

```go
type Transaction struct {
  ID        uint      `gorm:"primaryKey" json:"id"`
  GoogleID  string    `json:"google_id"`
  Type      string    `json:"type"` // "deposit" ou "withdrawal"
  Amount    float64   `json:"amount"`
  IBANMask  string    `json:"iban_mask,omitempty"`
  CreatedAt time.Time `json:"created_at"`
}
```
