<template>
  <div
    class="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 via-gray-950 to-violet-900 px-2 py-10">
    <div class="w-full max-w-6xl">
        <div class="flex flex-col items-center mb-12 mt-6">
        <span class="mb-3 px-4 py-1 rounded-full bg-gradient-to-r from-violet-700 to-emerald-500 text-white font-bold text-xs shadow uppercase tracking-widest">
          Mon profil
        </span>
        <h2
          class="text-4xl md:text-5xl font-extrabold bg-gradient-to-r from-emerald-400 via-violet-400 to-emerald-400 bg-clip-text text-transparent drop-shadow-lg text-center"
        >
          Mon espace personnel
        </h2>
        <div class="text-gray-400 text-sm mt-2 font-medium">Gérez votre profil, votre activité, vos trajets et vos demandes driver.</div>
      </div>
      <UserStatsCards :balance="user?.balance?.toFixed(2) ?? '0.00'" :co2-passenger="co2Passenger"
        :co2-driver="co2Driver" :km-passenger="kmPassenger" :km-driver="kmDriver" class="mb-8" />

      <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
        <UserInfoCard v-if="user" :user="user" @updateProfile="updateProfile" />
        <DriverRequestCard :driver-request="driverRequest" :user="user" @openModal="openDriverModal" />
        <div class="md:col-span-2">
          <TransactionHistory />
        </div>
      </div>

      <CoursesCard v-if="passengerCourses.length" :courses="passengerCourses" :names="driverNames" mode="passenger"
        class="mb-8" />
      <CoursesCard v-if="user?.roles?.includes('ROLE_DRIVER') && driverCourses.length" :courses="driverCourses"
        :names="passengerNames" mode="driver" class="mb-8" />

      <div v-if="showModal" class="fixed inset-0 flex items-center justify-center z-50 bg-black/60">
        <div
          class="bg-gradient-to-br from-gray-900/95 via-gray-900/90 to-violet-950/80 border border-gray-800 p-8 rounded-2xl shadow-2xl min-w-[320px] w-full max-w-md relative animate-fade-in">
          <div class="flex items-center gap-3 mb-5">
            <svg class="w-8 h-8 text-violet-400" fill="none" viewBox="0 0 24 24">
              <rect x="2" y="10" width="20" height="6" rx="2" fill="#7c3aed" />
              <rect x="6" y="7" width="12" height="5" rx="2" fill="#a78bfa" />
              <circle cx="7" cy="18" r="2" fill="#d1d5db" />
              <circle cx="17" cy="18" r="2" fill="#d1d5db" />
            </svg>
            <h3 class="text-xl font-bold text-white">Demander à devenir driver</h3>
          </div>
          <form @submit.prevent="requestDriver" class="space-y-5">
            <div>
              <label class="block text-gray-400 mb-1 ml-1 text-sm">Votre voiture</label>
              <input v-model="modalCar" placeholder="Modèle, marque…"
                class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow transition-all duration-200 font-mono text-lg"
                autocomplete="off" required />
            </div>
            <div>
              <label class="block text-gray-400 mb-1 ml-1 text-sm">Votre plaque</label>
              <input v-model="modalPlate" placeholder="AA-123-BB"
                class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow transition-all duration-200 font-mono text-lg"
                autocomplete="off" required />
            </div>
            <div class="flex gap-4 mt-6">
              <button type="submit"
                class="flex-1 py-3 rounded-xl bg-emerald-600 hover:bg-emerald-700 text-white text-lg font-bold shadow transition-all">Envoyer</button>
              <button type="button" @click="closeModal"
                class="flex-1 py-3 rounded-xl bg-gray-700 hover:bg-gray-800 text-white text-lg font-bold shadow transition-all">Annuler</button>
            </div>
          </form>
          <button @click="closeModal" class="absolute top-3 right-4 text-gray-400 hover:text-white text-2xl"
            title="Fermer">&times;</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import UserInfoCard from '../components/UserInfoCard.vue'
import DriverRequestCard from '../components/DriverRequestCard.vue'
import TransactionHistory from '../components/TransactionHistory.vue'
import CoursesCard from '../components/CoursesCard.vue'
import UserStatsCards from '../components/UserStatsCards.vue'
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
  } catch (e) { }
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
      const res = await fetch(`https://user-ecodrive.liamcariou.fr/users/${id}`)
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
      const res = await fetch(`https://user-ecodrive.liamcariou.fr/users/${id}`)
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
  const resPass = await fetch(`https://driver-ecodrive.liamcariou.fr/courses?role=passenger&id=${googleId}`)
  if (resPass.ok) {
    passengerCourses.value = await resPass.json()
    await resolveDriverNames(passengerCourses.value)
    await enrichCoursesWithAddresses(passengerCourses.value)
  } else {
    passengerCourses.value = []
  }
  if (user.value && user.value.roles?.includes("ROLE_DRIVER")) {
    const resDrv = await fetch(`https://driver-ecodrive.liamcariou.fr/courses?role=driver&id=${googleId}`)
    if (resDrv.ok) {
      driverCourses.value = await resDrv.json()
      await resolvePassengerNames(driverCourses.value)
      await enrichCoursesWithAddresses(driverCourses.value)
    } else {
      driverCourses.value = []
    }
  }
}

const co2Passenger = computed(() =>
  (passengerCourses.value.reduce((sum, c) => sum + (c.co2 || 0), 0) / 1000).toFixed(2)
)
const co2Driver = computed(() =>
  (driverCourses.value.reduce((sum, c) => sum + (c.co2 || 0), 0) / 1000).toFixed(2)
)
const kmPassenger = computed(() =>
  passengerCourses.value.reduce((sum, c) => sum + (c.distance_km || 0), 0).toFixed(1)
)
const kmDriver = computed(() =>
  driverCourses.value.reduce((sum, c) => sum + (c.distance_km || 0), 0).toFixed(1)
)

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
  const res = await fetch(`https://user-ecodrive.liamcariou.fr/users/${googleId}`)
  if (res.ok) user.value = await res.json()
  else router.push('/login')

  const drres = await fetch(`https://user-ecodrive.liamcariou.fr/driver-requests`)
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
  const res = await fetch(`https://user-ecodrive.liamcariou.fr/users/${googleId}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      Birthdate: user.value.birthdate,
      Gender: user.value.gender,
      Phone: user.value.phone,
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
  const res = await fetch(`https://user-ecodrive.liamcariou.fr/driver-requests`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload)
  })
  if (res.ok) {
    alert('Demande envoyée à l’admin !')
    const drres = await fetch(`https://user-ecodrive.liamcariou.fr/driver-requests`)
    if (drres.ok) {
      const list = await drres.json()
      driverRequest.value = list.find(
        d => d.google_id === googleId && (d.status === "pending" || d.status === "accepted")
      ) || null
    }
    closeModal()
  } else {
    const err = await res.json().catch(() => null)
    alert('Erreur lors de la demande : ' + (err?.error || ''))
  }
}
</script>

<style>
.modal-bg {
  position: fixed;
  left: 0;
  top: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal {
  background: #fff;
  padding: 24px;
  border-radius: 8px;
  min-width: 250px;
}
</style>