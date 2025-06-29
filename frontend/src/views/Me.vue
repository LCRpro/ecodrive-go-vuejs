<template>
  <div>
    <h2>Mon profil</h2>
    <div v-if="user">
      <div>Email: {{ user.email }}</div>
      <div>Nom: {{ user.name }}</div>
      <div>Naissance: <input v-model="user.birthdate" /></div>
      <div>Genre: <input v-model="user.gender" /></div>
      <div>Solde : {{ user.balance?.toFixed(2) ?? '0.00' }} €</div>
      <TransactionHistory />
      <button @click="updateProfile">Mettre à jour</button>

      <div v-if="driverRequest">
        <div v-if="driverRequest.status === 'pending'" style="margin-top:10px;">
          <b>Demande driver en attente de validation</b><br>
          Voiture: {{ driverRequest.car }}<br>
          Plaque: {{ driverRequest.plate }}
        </div>
        <div v-else-if="driverRequest.status === 'accepted'" style="margin-top:10px;">
          <b>Vous êtes driver !</b><br>
          Voiture: {{ driverRequest.car }}<br>
          Plaque: {{ driverRequest.plate }}
        </div>
        <div v-else-if="driverRequest.status === 'refused'" style="margin-top:10px;color:#a00;">
          <b>Votre demande driver a été refusée.</b>
        </div>
      </div>
      <button
        v-if="!driverRequest"
        @click="openDriverModal"
        style="margin-left:12px;margin-top:10px;"
      >Devenir driver</button>
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
import TransactionHistory from '../components/TransactionHistory.vue'
const user = ref(null)
const driverRequest = ref(null)
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

  const drres = await fetch(`http://localhost:8002/driver-requests`)
  if (drres.ok) {
    const list = await drres.json()
    driverRequest.value = list.find(
      d => d.google_id === googleId && (d.status === "pending" || d.status === "accepted")
    ) || null
  }
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
    })
  })
  if (res.ok) alert('Mise à jour OK')
}

function openDriverModal() {
  showModal.value = true
  modalCar.value = ''
  modalPlate.value = ''
}
function closeModal() {
  showModal.value = false
}

async function requestDriver() {
  const token = localStorage.getItem('token')
  const googleId = getGoogleIdFromToken(token)
  if (!googleId) {
    alert('Erreur : google_id introuvable')
    return
  }
  const payload = { google_id: googleId, car: modalCar.value, plate: modalPlate.value }
  const res = await fetch(`http://localhost:8002/driver-requests`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload)
  })
  if (res.ok) {
    alert('Demande envoyée à l’admin !')
    const drres = await fetch(`http://localhost:8002/driver-requests`)
    if (drres.ok) {
      const list = await drres.json()
      driverRequest.value = list.find(
        d => d.google_id === googleId && (d.status === "pending" || d.status === "accepted")
      ) || null
    }
    closeModal()
  } else {
    const err = await res.json().catch(()=>null)
    alert('Erreur lors de la demande : ' + (err?.error || ''))
  }
}
</script>

<style>
.modal-bg {
  position: fixed; left: 0; top: 0; width: 100vw; height: 100vh;
  background: rgba(0,0,0,0.4); display: flex; align-items: center; justify-content: center;
}
.modal { background: #fff; padding: 24px; border-radius: 8px; min-width: 250px; }
</style>