<template>
  <div>
    <nav>
      <router-link to="/">Accueil</router-link> |
      <router-link to="/me" v-if="isLoggedIn">Mon profil</router-link> |
      <router-link to="/admin" v-if="isAdmin">Admin</router-link> |
      <router-link to="/login" v-if="!isLoggedIn">Connexion</router-link>
      <button v-if="isLoggedIn" @click="logout" style="margin-left:12px;">Déconnexion</button>
    </nav>
    <router-view />
  </div>
</template>
<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
const router = useRouter()
const isLoggedIn = computed(() => !!localStorage.getItem('token'))
const isAdmin = computed(() => {
  try {
    return JSON.parse(localStorage.getItem('roles') || '[]').includes('ROLE_ADMIN')
  } catch {
    return false
  }
})
function logout() {
  localStorage.clear()
  router.push('/login')
}
</script>