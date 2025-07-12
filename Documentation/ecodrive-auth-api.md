# 📘 API – Service `auth` (Ecodrive)

Ce service gère l’authentification des utilisateurs via Google OAuth2 et la génération de jetons JWT valides pour les autres microservices.

---

## 🌐 Base URL

```
http://localhost:8000
https://auth-ecodrive.liamcariou.fr
```

---

## 🔐 Authentification via Google OAuth

### 🔗 GET `/auth/google`

Redirige l’utilisateur vers la page de connexion Google.

**Réponse :**

- Redirection 307 vers `https://accounts.google.com/o/oauth2/auth?...`

---

### 🔗 GET `/auth/callback`

Callback appelé par Google après connexion. Ce endpoint :

- Récupère les infos utilisateur depuis Google.
- Appelle `/users/find-or-create` sur `user-service`.
- Génère un token JWT avec les rôles de l’utilisateur.
- Redirige vers `FRONTEND_CALLBACK_URL` avec le token en query param.

**Réponse :**

- Redirection 307 vers : `http://localhost:5173/callback?token=<jwt>`

---

## 🎯 JWT généré

Le JWT contient les claims suivants :

```json
{
  "sub": "google_id",
  "email": "user@example.com",
  "roles": ["ROLE_PASSAGER", "ROLE_ADMIN"],
  "exp": 1723459200
}
```

---

## 💡 Exemple de logique

1. L'utilisateur clique sur "Connexion avec Google".
2. `/auth/google` le redirige vers la page Google.
3. Google redirige vers `/auth/callback?code=...`.
4. Le service `auth` échange le code contre un token Google.
5. Il appelle `user-service` pour créer ou récupérer l'utilisateur.
6. Il génère un JWT signé avec `JWT_SECRET`.
7. Il redirige vers le frontend avec `?token=...`.

---

## 🧠 Environnement requis (`.env`)

```env
JWT_SECRET=...
GOOGLE_CLIENT_ID=...
GOOGLE_CLIENT_SECRET=...
USER_SERVICE_URL=http://localhost:8002
FRONTEND_CALLBACK_URL=http://localhost:5173/callback
AUTH_REDIRECT_URL=http://localhost:8000/auth/callback
FRONTEND_ORIGIN=http://localhost:5173
```
