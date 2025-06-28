<template>
  <div>
    <h2>Mon profil</h2>
    <div v-if="user">
      <div>Email: {{ user.email }}</div>
      <div>Nom: {{ user.name }}</div>
      <div>Naissance: <input v-model="user.birthdate" /></div>
      <div>Genre: <input v-model="user.gender" /></div>
      <div>Voiture: <input v-model="user.car" /></div>
      <div>Plaque: <input v-model="user.plate" /></div>
      <button @click="updateProfile">Mettre à jour</button>
      <button @click="openDriverModal" style="margin-left:12px;">Devenir driver</button>
    </div>
    <div v-else>Chargement...</div>

    <div v-if="showModal" class="modal-bg">
      <div class="modal">
        <h3>Demander à devenir driver</h3>
        <input v-model="modalCar" placeholder="Votre voiture" />
        <input v-model="modalPlate" placeholder="Votre plaque" />
        <button @click="requestDriver">Envoyer</button>
        <button @click="closeModal" style="margin-left:12px;">Annuler</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
const user = ref(null)
const router = useRouter()

const showModal = ref(false)
const modalCar = ref('')
const modalPlate = ref('')

function getGoogleIdFromToken(token) {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.sub
  } catch {
    return null
  }
}

onMounted(async () => {
  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
    return
  }
  const googleId = getGoogleIdFromToken(token)
  if (!googleId) {
    router.push('/login')
    return
  }
  const res = await fetch(`http://localhost:8002/users/${googleId}`)
  if (res.ok) user.value = await res.json()
  else router.push('/login')
})

async function updateProfile() {
  const token = localStorage.getItem('token')
  const googleId = getGoogleIdFromToken(token)
  const res = await fetch(`http://localhost:8002/users/${googleId}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      Birthdate: user.value.birthdate,
      Gender: user.value.gender,
      Car: user.value.car,
      Plate: user.value.plate,
    })
  })
  if (res.ok) alert('Mise à jour OK')
}

// Modal devenir driver
function openDriverModal() {
  showModal.value = true
  modalCar.value = user.value.car || ''
  modalPlate.value = user.value.plate || ''
}
function closeModal() {
  showModal.value = false
}

async function requestDriver() {
  const token = localStorage.getItem('token')
  const googleId = getGoogleIdFromToken(token)
  await fetch(`http://localhost:8003/admin/driver-requests`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Authorization: 'Bearer ' + token },
    body: JSON.stringify({ google_id: googleId, car: modalCar.value, plate: modalPlate.value })
  })
  alert('Demande envoyée à l’admin !')
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