<template>
  <div>
    <nav style="margin-bottom:24px;">
      <router-link v-if="isLoggedIn" to="/">Accueil</router-link> |
      <router-link v-if="isLoggedIn" to="/me">Mon profil</router-link> |
      <router-link v-if="!isLoggedIn" to="/login">Login</router-link>
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
function logout() {
  localStorage.clear()
  router.push('/login')
}
</script>