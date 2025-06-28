<template>
  <div>
    <h2>Panel Admin</h2>
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
            <th>Action</th>
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
            <td>
              <button 
                v-if="showAcceptDriver(u)"
                @click="openDriverModal(u)"
              >Accepter driver</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="modalUser" class="modal-bg">
      <div class="modal">
        <h3>Accepter {{ modalUser.name }} comme driver</h3>
        <input v-model="car" placeholder="Voiture" />
        <input v-model="plate" placeholder="Plaque" />
        <button @click="acceptDriver">Valider</button>
        <button @click="closeModal" style="margin-left:12px;">Annuler</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
const router = useRouter()
const users = ref([])
const loading = ref(true)
const modalUser = ref(null)
const car = ref('')
const plate = ref('')

onMounted(async () => {
  const token = localStorage.getItem('token')
  const roles = JSON.parse(localStorage.getItem('roles') || '[]')
  if (!token || !roles.includes('ROLE_ADMIN')) {
    router.push('/')
    return
  }
  await fetchUsers()
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

function showAcceptDriver(u) {
  try {
    return !JSON.parse(u.roles || "[]").includes('ROLE_DRIVER')
  } catch {
    return true
  }
}

function openDriverModal(u) {
  modalUser.value = u
  car.value = u.car || ''
  plate.value = u.plate || ''
}
function closeModal() {
  modalUser.value = null
  car.value = ''
  plate.value = ''
}

async function acceptDriver() {
  const token = localStorage.getItem('token')
  await fetch('http://localhost:8003/admin/users/' + modalUser.value.id + '/accept-driver', {
    method: 'PATCH',
    headers: {
      Authorization: 'Bearer ' + token,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ car: car.value, plate: plate.value })
  })
  await fetchUsers()
  closeModal()
}
</script>

<style>
.modal-bg {
  position: fixed; left: 0; top: 0; width: 100vw; height: 100vh;
  background: rgba(0,0,0,0.4); display: flex; align-items: center; justify-content: center;
}
.modal { background: #fff; padding: 24px; border-radius: 8px; min-width: 250px; }
</style>