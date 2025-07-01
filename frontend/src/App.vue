<template>
  <div>
    <nav>
      <router-link to="/">Accueil</router-link> |
      <router-link to="/me" v-if="isLoggedIn">Mon profil</router-link> |
      <router-link to="/admin" v-if="isAdmin">Admin</router-link>
    <router-link to="/course" v-if="isLoggedIn">Commander une course</router-link>
<router-link to="/driver" v-if="isLoggedIn && roles.includes('ROLE_DRIVER')">Courses à prendre</router-link>
      <BalanceDisplay v-if="isLoggedIn"/>
      <router-link to="/deposit" v-if="isLoggedIn" style="margin-left:20px;">Déposer</router-link>
      <router-link to="/withdraw" v-if="isLoggedIn" style="margin-left:10px;">Retirer</router-link>
      <router-link to="/login" v-if="!isLoggedIn">Connexion</router-link>
      <button v-if="isLoggedIn" @click="logout" style="margin-left:12px;">Déconnexion</button>
    </nav>
    <router-view />
  </div>
</template>
<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import BalanceDisplay from './components/BalanceDisplay.vue'
const router = useRouter()
const roles = computed(() => {
  try {
    return JSON.parse(localStorage.getItem('roles') || '[]')
  } catch { return [] }
})
const isLoggedIn = computed(() => !!localStorage.getItem('token'))
const isAdmin = computed(() => roles.value.includes('ROLE_ADMIN'))
const isPassager = computed(() => roles.value.includes('ROLE_PASSAGER'))
const isDriver = computed(() => roles.value.includes('ROLE_DRIVER'))

function logout() {
  localStorage.clear()
  router.push('/login')
}
</script>