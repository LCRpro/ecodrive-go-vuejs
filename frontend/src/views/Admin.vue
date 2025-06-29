<template>
  <div>
    <h2>Panel Admin</h2>
    <div v-if="appBalance !== null" style="margin-bottom: 16px;">
  <b>Solde de la plateforme :</b> {{ appBalance.toFixed(2) }} €
</div>
    <div v-if="loading">Chargement...</div>
    <div v-else>
      <table border="1" cellpadding="5">
        <thead>
          <tr>
            <th>ID</th>
            <th>Email</th>
            <th>Nom</th>
            <th>Rôles</th>
            <th>Voiture</th>
            <th>Plaque</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="u in users" :key="u.id">
            <td>{{ u.id }}</td>
            <td>{{ u.email }}</td>
            <td>{{ u.name }}</td>
            <td>{{ JSON.parse(u.roles || "[]").join(", ") }}</td>
            <td>{{ u.car || '-' }}</td>
            <td>{{ u.plate || '-' }}</td>
          </tr>
        </tbody>
      </table>

      <h3 style="margin-top:32px;">Demandes driver</h3>
      <table border="1" cellpadding="5">
        <thead>
          <tr>
            <th>ID</th>
            <th>Nom</th>
            <th>Email</th>
            <th>Voiture</th>
            <th>Plaque</th>
            <th>Status</th>
            <th>Action</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="req in driverRequests" :key="req.ID">
            <td>{{ req.ID }}</td>
<td>{{ getUserByGoogleId(req.google_id || req.GoogleID)?.name || '-' }}</td>
<td>{{ getUserByGoogleId(req.google_id || req.GoogleID)?.email || '-' }}</td>
            <td>{{ req.car }}</td>
            <td>{{ req.plate }}</td>
            <td>{{ req.status }}</td>
            <td>
              <button v-if="req.status === 'pending'" @click="handleRequest(req.ID, 'accept')">Accepter</button>
              <button v-if="req.status === 'pending'" @click="handleRequest(req.ID, 'refuse')" style="margin-left:8px;">Refuser</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>

  <h3 style="margin-top:32px;">Transactions</h3>
<table border="1" cellpadding="5" v-if="transactions.length">
  <thead>
    <tr>
      <th>Date</th>
      <th>Utilisateur</th>
      <th>Type</th>
      <th>Montant (€)</th>
      <th>Détails</th>
    </tr>
  </thead>
  <tbody>
    <tr v-for="tx in transactions" :key="tx.id">
      <td>{{ new Date(tx.created_at).toLocaleString() }}</td>
      <td>
        {{ getUserByGoogleId(tx.google_id)?.email || tx.google_id }}
      </td>
      <td>{{ tx.type === 'deposit' ? 'Dépôt' : 'Retrait' }}</td>
      <td>{{ tx.amount.toFixed(2) }}</td>
      <td v-if="tx.type === 'deposit'">Carte</td>
      <td v-else>IBAN : {{ tx.iban_mask }}</td>
    </tr>
  </tbody>
</table>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
const router = useRouter()
const appBalance = ref(null)
const users = ref([])
const driverRequests = ref([])
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
})

async function fetchUsers() {
  loading.value = true
  const token = localStorage.getItem('token')
  const res = await fetch('http://localhost:8003/admin/users', {
    headers: { Authorization: 'Bearer ' + token }
  })
  users.value = res.ok ? await res.json() : []
  loading.value = false
}

async function fetchDriverRequests() {
  const token = localStorage.getItem('token')
  const res = await fetch('http://localhost:8003/admin/driver-requests', {
    headers: { Authorization: 'Bearer ' + token }
  })
  driverRequests.value = res.ok ? await res.json() : []
}

function getUserByGoogleId(google_id) {
  return users.value.find(u => u.google_id === google_id)
}

async function handleRequest(id, action) {
  const token = localStorage.getItem('token')
  await fetch('http://localhost:8003/admin/driver-requests/' + id, {
    method: 'PATCH',
    headers: { Authorization: 'Bearer ' + token, 'Content-Type': 'application/json' },
    body: JSON.stringify({ action })
  })
  await fetchDriverRequests()
  await fetchUsers()
}

async function fetchAppBalance() {
  const res = await fetch('http://localhost:8002/app-balance')
  if (res.ok) {
    const data = await res.json()
    appBalance.value = data.balance
  }
}

const transactions = ref([])
onMounted(async () => {
  const token = localStorage.getItem('token')
  const res = await fetch('http://localhost:8004/transactions', {
    headers: { Authorization: 'Bearer ' + token }
  })
  transactions.value = res.ok ? await res.json() : []
})
</script>

<style>
.modal-bg {
  position: fixed; left: 0; top: 0; width: 100vw; height: 100vh;
  background: rgba(0,0,0,0.4); display: flex; align-items: center; justify-content: center;
}
.modal { background: #fff; padding: 24px; border-radius: 8px; min-width: 250px; }
</style>