<template>
  <div>
    <nav class="bg-gray-900 text-white px-4 py-3 flex items-center justify-between shadow-lg">
<router-link to="/" class="flex items-center gap-3 group select-none">
  <img
    src="/logo-ecodrive.png"
    alt="Logo EcoDrive"
    class="w-12 h-12 rounded-full shadow-lg border-2 border-emerald-400 group-hover:scale-110 group-hover:shadow-emerald-600 transition-transform duration-200 bg-white object-cover"
  />
  <span class="text-3xl font-extrabold font-mono bg-gradient-to-r from-emerald-400 via-violet-400 to-emerald-300 bg-clip-text text-transparent drop-shadow group-hover:brightness-110 transition-colors">
    Eco<span class="text-violet-400">Drive</span>
  </span>
</router-link>
      <div class="hidden md:flex items-center">
        <router-link to="/" class="nav-link" exact-active-class="text-violet-400">Accueil</router-link>
        <span class="separator" />
        <router-link v-if="isLoggedIn" to="/me" class="nav-link">Mon profil</router-link>
        <span v-if="isLoggedIn" class="separator" />
        <router-link v-if="isAdmin" to="/admin" class="nav-link">Admin</router-link>
        <span v-if="isAdmin" class="separator" />
        <router-link v-if="isLoggedIn" to="/course" class="nav-link">Commander une course</router-link>
        <span v-if="isLoggedIn" class="separator" />
        <router-link v-if="isLoggedIn && isDriver" to="/driver" class="nav-link">Courses à prendre</router-link>
        <span v-if="isLoggedIn && isDriver" class="separator" />
        <router-link v-if="isLoggedIn" to="/support" class="nav-link">Support</router-link>
        <span v-if="isLoggedIn && isDriver" class="separator" />
        <router-link
          v-if="isLoggedIn"
          to="/deposit"
          class="px-3 py-2 rounded-lg font-semibold bg-emerald-600 hover:bg-emerald-700 transition-colors shadow mx-1"
        >Déposer</router-link>
        <span v-if="isLoggedIn" class="separator" />
        <router-link
          v-if="isLoggedIn"
          to="/withdraw"
          class="px-3 py-2 rounded-lg font-semibold bg-violet-600 hover:bg-violet-700 transition-colors shadow mx-1"
        >Retirer</router-link>
        <span v-if="isLoggedIn" class="separator" />
        <router-link v-if="!isLoggedIn" to="/login" class="nav-link">Connexion</router-link>
        <BalanceDisplay v-if="isLoggedIn" class="ml-2 mr-2" />
        <button
          v-if="isLoggedIn"
          @click="logout"
          class="ml-2 px-4 py-1 rounded-lg bg-violet-600 hover:bg-violet-700 transition-colors font-semibold"
        >Déconnexion</button>
      </div>

      <button @click="isOpen = !isOpen" class="md:hidden focus:outline-none flex items-center" aria-label="Menu">
        <svg :class="isOpen ? 'rotate-90' : ''" class="w-7 h-7 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path v-if="!isOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"/>
          <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
        </svg>
      </button>
    </nav>

    <transition name="slide-fade">
      <div
        v-if="isOpen"
        class="md:hidden bg-gray-900 text-white px-4 py-2 flex flex-col gap-2 shadow-lg"
      >
        <router-link to="/" class="nav-link-mobile" @click="isOpen = false">Accueil</router-link>
        <router-link to="/me" v-if="isLoggedIn" class="nav-link-mobile" @click="isOpen = false">Mon profil</router-link>
        <router-link to="/admin" v-if="isAdmin" class="nav-link-mobile" @click="isOpen = false">Admin</router-link>
        <router-link to="/course" v-if="isLoggedIn" class="nav-link-mobile" @click="isOpen = false">Commander une course</router-link>
        <router-link to="/driver" v-if="isLoggedIn && isDriver" class="nav-link-mobile" @click="isOpen = false">Courses à prendre</router-link>
        <router-link to="/deposit" v-if="isLoggedIn" class="w-full text-center py-2 rounded-lg bg-emerald-600 hover:bg-emerald-700 font-semibold transition-colors" @click="isOpen = false">Déposer</router-link>
        <router-link to="/withdraw" v-if="isLoggedIn" class="w-full text-center py-2 rounded-lg bg-violet-600 hover:bg-violet-700 font-semibold transition-colors" @click="isOpen = false">Retirer</router-link>
        <router-link to="/login" v-if="!isLoggedIn" class="nav-link-mobile" @click="isOpen = false">Connexion</router-link>
        <BalanceDisplay v-if="isLoggedIn" class="my-2" />
        <button v-if="isLoggedIn" @click="logout" class="w-full text-left px-2 py-2 rounded-lg bg-violet-600 hover:bg-violet-700 font-semibold transition-colors">Déconnexion</button>
      </div>
    </transition>

    <router-view />


<footer class="w-full bg-gray-900 text-gray-400 pt-6 pb-8 px-6 flex flex-col md:flex-row items-center justify-between border-t border-gray-800 gap-3">
  <div class="text-sm text-center md:text-left">
    &copy; {{ new Date().getFullYear() }} Ecodrive. Tous droits réservés.
  </div>
  <div class="flex flex-col md:flex-row gap-2 md:gap-6 text-sm">
    <router-link to="/mentions" class="hover:text-violet-400 underline transition">Mentions légales</router-link>
    <router-link to="/conditions" class="hover:text-violet-400 underline transition">Conditions générales de vente</router-link>
  </div>
</footer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import BalanceDisplay from './components/BalanceDisplay.vue'

const isOpen = ref(false)
const token = ref(localStorage.getItem('token'))
const roles = ref([])

const refreshAuth = () => {
  token.value = localStorage.getItem('token')
  try {
    roles.value = JSON.parse(localStorage.getItem('roles') || '[]')
  } catch { roles.value = [] }
}
onMounted(() => {
  refreshAuth()
  window.addEventListener('storage', refreshAuth)
})
onUnmounted(() => {
  window.removeEventListener('storage', refreshAuth)
})
const isLoggedIn = computed(() => !!token.value)
const isAdmin = computed(() => roles.value.includes('ROLE_ADMIN'))
const isDriver = computed(() => roles.value.includes('ROLE_DRIVER'))

const router = useRouter()
function logout() {
  localStorage.clear()
  refreshAuth()
  isOpen.value = false
  router.push('/login')
}
</script>

<style scoped>
.nav-link {
  @apply px-3 py-2 rounded-lg hover:bg-violet-700 transition-colors font-medium;
}
.separator {
  display: inline-block;
  height: 24px;
  width: 1.5px;
  background: #444;
  margin: 0 0.5rem;
  vertical-align: middle;
}
.nav-link-mobile {
  @apply px-2 py-2 rounded-md hover:bg-violet-700 transition-colors;
}
.slide-fade-enter-active {
  transition: all 0.2s cubic-bezier(.57,.21,.69,1.25);
}
.slide-fade-leave-active {
  transition: all 0.15s cubic-bezier(.57,.21,.69,1.25);
}
.slide-fade-enter-from,
.slide-fade-leave-to {
  opacity: 0;
  transform: translateY(-12px);
}
</style>