<template>
  <div style="margin:60px auto;text-align:center;">
    <h2>Connexion en cours...</h2>
    <div v-if="error" style="color:red">{{ error }}</div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()
const error = ref('')

onMounted(() => {
  const { token, roles, email, name, error: err } = route.query
  if (token) {
    localStorage.setItem('token', token)
    if (roles) localStorage.setItem('roles', roles)
    if (email) localStorage.setItem('email', email)
    if (name) localStorage.setItem('name', name)
    router.replace('/')
  } else if (err) {
    error.value = err
  } else {
    error.value = "Aucun token reçu"
  }
})
</script>