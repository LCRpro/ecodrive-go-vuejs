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

      <div v-if="passengerCourses.length" style="margin-top: 24px;">
        <h3>Mes courses passager</h3>
        <ul>
          <li v-for="c in passengerCourses" :key="c.id">
            <b>Départ :</b> {{ c.start_address }}<br>
            <b>Arrivée :</b> {{ c.end_address }}<br>
            <b>Distance :</b> {{ c.distance_km.toFixed(2) }} km<br>
            <b>Prix :</b> {{ c.amount?.toFixed(2) ?? '?' }} €<br>
            <b>CO₂ économisé :</b> {{ (c.co2/1000).toFixed(2) }} kg<br>
            <b>Driver :</b> {{ driverNames[c.driver_id] ?? c.driver_id ?? '—' }}<br>
            <b>Statut :</b> {{ c.status }}
          </li>
        </ul>
      </div>

      <div v-if="user.roles?.includes('ROLE_DRIVER') && driverCourses.length" style="margin-top: 24px;">
        <h3>Mes courses en tant que driver</h3>
        <ul>
          <li v-for="c in driverCourses" :key="c.id">
            <b>Départ :</b> {{ c.start_address }}<br>
            <b>Arrivée :</b> {{ c.end_address }}<br>
            <b>Distance :</b> {{ c.distance_km.toFixed(2) }} km<br>
            <b>Prix :</b> {{ c.amount?.toFixed(2) ?? '?' }} €<br>
            <b>CO₂ économisé :</b> {{ (c.co2/1000).toFixed(2) }} kg<br>
            <b>Passager :</b> {{ passengerNames[c.passenger_id] ?? c.passenger_name ?? c.passenger_id }}<br>
            <b>Statut :</b> {{ c.status }}
          </li>
        </ul>
      </div>
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

const passengerCourses = ref([])
const driverCourses = ref([])

const driverNames = ref({})
const passengerNames = ref({})

const addressCache = new Map()

function getGoogleIdFromToken(token) {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.sub
  } catch {
    return null
  }
}

// Reverse geocoding (OpenStreetMap/Nominatim)
const GOOGLE_API_KEY = import.meta.env.VITE_GOOGLE_MAPS_API_KEY

async function reverseGeocode(lat, lng) {
  const key = `${lat},${lng}`
  if (addressCache.has(key)) return addressCache.get(key)
  const url = `https://maps.googleapis.com/maps/api/geocode/json?latlng=${lat},${lng}&key=${GOOGLE_API_KEY}&language=fr`
  try {
    const res = await fetch(url)
    if (res.ok) {
      const data = await res.json()
      if (data.status === "OK" && data.results.length > 0) {
        const address = data.results[0].formatted_address
        addressCache.set(key, address)
        return address
      }
    }
  } catch (e) {}
  addressCache.set(key, `${lat}, ${lng}`)
  return `${lat}, ${lng}`
}

async function enrichCoursesWithAddresses(courses) {
  for (const c of courses) {
    c.start_address = await reverseGeocode(c.start_lat, c.start_lng)
    c.end_address = await reverseGeocode(c.end_lat, c.end_lng)
  }
}

async function resolveDriverNames(courses) {
  const ids = [...new Set(courses.map(c => c.driver_id).filter(Boolean))]
  for (const id of ids) {
    if (!driverNames.value[id]) {
      const res = await fetch(`http://localhost:8002/users/${id}`)
      if (res.ok) {
        const user = await res.json()
        driverNames.value[id] = user.name
      } else {
        driverNames.value[id] = id 
      }
    }
  }
}
async function resolvePassengerNames(courses) {
  const ids = [...new Set(courses.map(c => c.passenger_id).filter(Boolean))]
  for (const id of ids) {
    if (!passengerNames.value[id]) {
      const res = await fetch(`http://localhost:8002/users/${id}`)
      if (res.ok) {
        const user = await res.json()
        passengerNames.value[id] = user.name
      } else {
        passengerNames.value[id] = id
      }
    }
  }
}

async function fetchMyCourses(googleId) {
  const resPass = await fetch(`http://localhost:8006/courses?role=passenger&id=${googleId}`)
  if (resPass.ok) {
    passengerCourses.value = await resPass.json()
    await resolveDriverNames(passengerCourses.value)
    await enrichCoursesWithAddresses(passengerCourses.value)
  } else {
    passengerCourses.value = []
  }
  if (user.value && user.value.roles?.includes("ROLE_DRIVER")) {
    const resDrv = await fetch(`http://localhost:8006/courses?role=driver&id=${googleId}`)
    if (resDrv.ok) {
      driverCourses.value = await resDrv.json()
      await resolvePassengerNames(driverCourses.value)
      await enrichCoursesWithAddresses(driverCourses.value)
    } else {
      driverCourses.value = []
    }
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
  await fetchMyCourses(googleId)
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