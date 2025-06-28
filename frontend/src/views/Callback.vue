<template>
  <div style="margin:60px auto;text-align:center;">
    <h2>Connexion en cours...</h2>
    <div v-if="error" style="color:red">{{ error }}</div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

onMounted(async () => {
  const token = route.query.token
  if (!token) {
    router.replace('/login')
    return
  }
  localStorage.setItem('token', token)

  // Décoder le GoogleID depuis le JWT (payload)
  function getGoogleIdFromToken(token) {
    try {
      return JSON.parse(atob(token.split('.')[1])).sub
    } catch { return null }
  }
  const googleId = getGoogleIdFromToken(token)
  if (!googleId) {
    router.replace('/login')
    return
  }

  // Appel au service-user pour récupérer les rôles
  const res = await fetch('http://localhost:8002/users/' + googleId)
  if (res.ok) {
    const user = await res.json()
    localStorage.setItem('roles', user.roles) // c'est bien un string JSON
    localStorage.setItem('email', user.email)
    localStorage.setItem('name', user.name)
    localStorage.setItem('google_id', user.google_id)
    // etc. selon ce que tu veux afficher
  }
  router.replace('/')
})
</script>