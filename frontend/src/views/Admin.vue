<template>
  <div
    class="min-h-screen flex flex-col items-center bg-gradient-to-br from-gray-900 via-gray-950 to-violet-900 px-2 py-10">
    <div class="w-full max-w-7xl">
      <div class="flex flex-col items-center mb-12 mt-6">
        <span class="mb-3 px-4 py-1 rounded-full bg-gradient-to-r from-violet-700 to-emerald-500 text-white font-bold text-xs shadow uppercase tracking-widest">
          Espace sécurisé
        </span>
        <h2
          class="text-4xl md:text-5xl font-extrabold bg-gradient-to-r from-emerald-400 via-violet-400 to-emerald-400 bg-clip-text text-transparent drop-shadow-lg text-center"
        >
          Panel Admin
        </h2>
        <div class="text-gray-400 text-sm mt-2 font-medium">Gestion des utilisateurs, chauffeurs, courses et finances</div>
      </div>
      <div v-if="appBalance !== null" class="mb-6">
        <div
          class="bg-gradient-to-br from-gray-900/95 via-gray-900/90 to-violet-950/80 rounded-2xl shadow-xl border border-gray-800 px-8 py-6 flex flex-col items-center justify-center mb-8">
          <div class="text-lg text-gray-400 mb-2 font-medium">Chiffre d'affaires</div>
          <div class="text-3xl font-extrabold text-emerald-400 font-mono">{{ appBalance.toFixed(2) }} €</div>
        </div>
      </div>
      <div v-if="loading" class="text-gray-400">Chargement...</div>
      <div v-else>
        <AdminUserTable :users="users" />
        <AdminDriverRequests :driverRequests="driverRequests" :getUserByGoogleId="getUserByGoogleId"
          @handleRequest="handleRequest" />
                  <AdminCoursesTable :courses="courses" :getUserByGoogleId="getUserByGoogleId" />

        <AdminTransactions :transactions="transactions" :getUserByGoogleId="getUserByGoogleId" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import AdminUserTable from '../components/AdminUserTable.vue'
import AdminDriverRequests from '../components/AdminDriverRequests.vue'
import AdminTransactions from '../components/AdminTransactions.vue'
import AdminCoursesTable from '../components/AdminCoursesTable.vue'

const courses = ref([])
const router = useRouter()
const appBalance = ref(null)
const users = ref([])
const driverRequests = ref([])
const transactions = ref([])
const loading = ref(true)

onMounted(async () => {
  const token = localStorage.getItem('token')
  const roles = JSON.parse(localStorage.getItem('roles') || '[]')
  if (!token || !roles.includes('ROLE_ADMIN')) {
    router.push('/')
    return
  }
  await fetchUsers()
  await fetchDriverRequests()
  await fetchAppBalance()
  await fetchTransactions()
  await fetchCourses()
})

async function fetchCourses() {
  const token = localStorage.getItem('token')
  const res = await fetch('https://driver-ecodrive.liamcariou.fr/courses', {
    headers: { Authorization: 'Bearer ' + token }
  })
  courses.value = res.ok ? await res.json() : []
}

async function fetchUsers() {
  loading.value = true
  const token = localStorage.getItem('token')
  const res = await fetch('https://admin-ecodrive.liamcariou.fr/admin/users', {
    headers: { Authorization: 'Bearer ' + token }
  })
  users.value = res.ok ? await res.json() : []
  loading.value = false
}

async function fetchDriverRequests() {
  const token = localStorage.getItem('token')
  const res = await fetch('https://admin-ecodrive.liamcariou.fr/admin/driver-requests', {
    headers: { Authorization: 'Bearer ' + token }
  })
  driverRequests.value = res.ok ? await res.json() : []
}

function getUserByGoogleId(google_id) {
  return users.value.find(u => u.google_id === google_id)
}

async function handleRequest(id, action) {
  const token = localStorage.getItem('token')
  await fetch('https://admin-ecodrive.liamcariou.fr/admin/driver-requests/' + id, {
    method: 'PATCH',
    headers: { Authorization: 'Bearer ' + token, 'Content-Type': 'application/json' },
    body: JSON.stringify({ action })
  })
  await fetchDriverRequests()
  await fetchUsers()
}

async function fetchAppBalance() {
  const res = await fetch('https://user-ecodrive.liamcariou.fr/app-balance')
  if (res.ok) {
    const data = await res.json()
    appBalance.value = data.balance
  }
}

async function fetchTransactions() {
  const token = localStorage.getItem('token')
  const res = await fetch('https://paiement-ecodrive.liamcariou.fr/transactions', {
    headers: { Authorization: 'Bearer ' + token }
  })
  transactions.value = res.ok ? await res.json() : []
}
</script>