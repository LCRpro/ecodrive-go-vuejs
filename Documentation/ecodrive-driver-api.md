# 📘 API – Service `driver` (Ecodrive)

Ce service gère les courses créées par les passagers, ainsi que les actions des conducteurs : acceptation, démarrage, annulation et complétion de courses.

---

## 🌐 Base URL

```
http://localhost:8006
https://driver-ecodrive.liamcariou.fr
```

---

## 📄 Endpoints disponibles

### ➕ POST `/courses`

Crée une nouvelle course.

**Body JSON :**
```json
{
  "passenger_id": "abc123",
  "passenger_name": "Jean Dupont",
  "passenger_email": "jean@example.com",
  "start_lat": 48.85,
  "start_lng": 2.35,
  "end_lat": 48.86,
  "end_lng": 2.36,
  "distance_km": 12.4
}
```

**Réponse :**
```json
{
  "id": 1,
  "passenger_id": "abc123",
  "amount": 13.92,
  "status": "requested",
  ...
}
```

---

### 📄 GET `/courses`

Liste les courses. Paramètres possibles :

- `role=passenger&id=abc123` → courses d’un passager
- `role=driver&id=driver456` → courses d’un conducteur
- `driverView=1&driver_id=driver456` → vue chauffeur (priorité aux siennes puis aux libres)
- `status=requested` → filtrage par statut

---

### 📄 GET `/courses/:id`

Retourne une course spécifique.

---

### ✅ PATCH `/courses/:id/accept`

Le conducteur accepte une course.

**Body JSON :**
```json
{
  "driver_id": "driver456"
}
```

---

### ▶️ PATCH `/courses/:id/start`

Le conducteur démarre la course. Cette action :

- Change le statut à `in_progress`
- Débite le passager
- Crédite le conducteur (80%) et l’app (20%)

---

### 🏁 PATCH `/courses/:id/complete`

Marque la course comme complétée (`status = completed`).

---

### ❌ PATCH `/courses/:id/cancel`

Annule une course (si elle n’est pas déjà terminée).

---

### ✏️ PATCH `/courses/:id`

Modifie dynamiquement certains champs d’une course.

**Body JSON exemple :**
```json
{
  "status": "cancelled",
  "driver_id": "driver456"
}
```

---

## 🧠 Modèle `Course`

```go
type Course struct {
  ID             uint64     `gorm:"primaryKey" json:"id"`
  PassengerID    string     `json:"passenger_id"`
  PassengerName  string     `json:"passenger_name"`
  PassengerEmail string     `json:"passenger_email"`
  StartLat       float64    `json:"start_lat"`
  StartLng       float64    `json:"start_lng"`
  EndLat         float64    `json:"end_lat"`
  EndLng         float64    `json:"end_lng"`
  DistanceKm     float64    `json:"distance_km"`
  CO2            float64    `json:"co2"`
  Amount         float64    `json:"amount"`
  Status         string     `json:"status"`
  DriverID       string     `json:"driver_id"`
  AcceptedAt     *time.Time `json:"accepted_at"`
  CreatedAt      time.Time  `json:"created_at"`
}
```
