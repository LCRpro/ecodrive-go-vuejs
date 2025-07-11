<template>
  <div style="margin:60px auto;text-align:center;">
    <h2>Connexion en cours...</h2>
    <div v-if="error" style="color:red">{{ error }}</div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ref } from 'vue'

const router = useRouter()
const route = useRoute()
const error = ref('')

onMounted(async () => {
  const token = route.query.token
  if (!token) {
    error.value = "Token manquant, veuillez vous reconnecter."
    router.replace('/login')
    return
  }

  localStorage.clear()
  localStorage.setItem('token', token)

  let payload = {}
  try {
    payload = JSON.parse(atob(token.split('.')[1]))
  } catch {
    error.value = "Token invalide, veuillez vous reconnecter."
    router.replace('/login')
    return
  }
  if (!payload.sub) {
    error.value = "Identifiant Google manquant dans le token."
    router.replace('/login')
    return
  }

  localStorage.setItem('google_id', payload.sub)
  localStorage.setItem('email', payload.email)
  localStorage.setItem('roles', JSON.stringify(payload.roles || []))

  try {
    const res = await fetch('https://user-ecodrive.liamcariou.fr/users/' + payload.sub)
    if (res.ok) {
      const user = await res.json()
      localStorage.setItem('name', user.name || '')
    }
  } catch (e) {
  }

router.replace('/')
  window.location.replace('/')
})
</script>