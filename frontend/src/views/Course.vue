<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-900 via-gray-950 to-violet-900 flex flex-col md:flex-row items-center justify-center px-2 py-10 gap-8">
    <div class="w-full max-w-sm p-7 rounded-2xl shadow-2xl bg-gray-900 bg-opacity-90 border border-gray-800 mb-8 md:mb-0">
      <template v-if="activeCourse">
        <h2 class="text-white text-2xl font-extrabold bg-gradient-to-r from-emerald-400 via-violet-400 to-emerald-400 bg-clip-text text-transparent drop-shadow-lg text-center mb-6">Course</h2>
       <div class="bg-gradient-to-r from-gray-800/90 via-gray-900/95 to-violet-950/70 rounded-xl px-6 py-5 mb-3 shadow-inner border border-gray-800 space-y-4">
  <div class="flex items-center justify-between">
    <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Départ</span>
    <span class="text-emerald-400 font-mono text-base font-semibold text-right">
      {{ startAddress || 'Chargement...' }}
    </span>
  </div>
  <div class="flex items-center justify-between">
    <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Arrivée</span>
    <span class="text-violet-400 font-mono text-base font-semibold text-right">
      {{ endAddress || 'Chargement...' }}
    </span>
  </div>
  <div class="flex items-center justify-between">
    <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Distance</span>
    <span class="font-mono text-blue-200 text-base font-semibold text-right">
      <span class="inline-block bg-blue-900/70 px-2 py-1 rounded-md">{{ activeCourse.distance_km }} km</span>
    </span>
  </div>
  <div class="flex items-center justify-between">
    <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Durée</span>
    <span class="font-mono text-cyan-200 text-base font-semibold text-right">
      <span class="inline-block bg-cyan-900/60 px-2 py-1 rounded-md">{{ duration || '...' }}</span>
    </span>
  </div>
  <div class="flex items-center justify-between">
    <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Prix</span>
    <span class="font-mono text-amber-300 text-lg font-bold text-right">
      <span class="inline-block bg-amber-800/80 px-3 py-1 rounded-lg">{{ activeCourse.amount }} €</span>
    </span>
  </div>
  <div class="flex items-center justify-between">
    <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Statut</span>
    <span>
      <span
        v-if="activeCourse.status === 'requested'"
        class="inline-block px-3 py-1 rounded-full bg-amber-600 text-white text-xs font-semibold"
      >En attente</span>
      <span
        v-else-if="activeCourse.status === 'accepted'"
        class="inline-block px-3 py-1 rounded-full bg-emerald-600 text-white text-xs font-semibold"
      >Prise en charge</span>
      <span
        v-else-if="activeCourse.status === 'in_progress'"
        class="inline-block px-3 py-1 rounded-full bg-blue-600 text-white text-xs font-semibold"
      >En cours</span>
      <span
        v-else-if="activeCourse.status === 'completed'"
        class="inline-block px-3 py-1 rounded-full bg-violet-700 text-white text-xs font-semibold"
      >Terminée</span>
      <span
        v-else-if="activeCourse.status === 'cancelled'"
        class="inline-block px-3 py-1 rounded-full bg-rose-700 text-white text-xs font-semibold"
      >Annulée</span>
      <span
        v-else
        class="inline-block px-3 py-1 rounded-full bg-gray-700 text-gray-200 text-xs font-semibold"
      >{{ displayStatus(activeCourse.status) }}</span>
    </span>
  </div>
</div>
        <button 
          v-if="canCancel" 
          @click="cancelCourse" 
          class="w-full py-3 rounded-lg bg-red-600 hover:bg-red-700 text-white text-lg font-bold shadow transition mt-6"
        >
          Annuler la course
        </button>
        <button 
          @click="fetchActiveCourse" 
          class="w-full py-3 rounded-lg bg-violet-600 hover:bg-violet-700 text-white text-lg font-bold shadow transition mt-4"
        >
          Rafraîchir
        </button>
        <div v-if="cancelError" class="mt-4 rounded-lg bg-red-700 bg-opacity-80 text-white p-3 text-center font-semibold">
          {{ cancelError }}
        </div>
      </template>
      <template v-else>
        <h2 class="text-2xl font-extrabold bg-gradient-to-r from-emerald-400 via-violet-400 to-emerald-400 bg-clip-text text-transparent drop-shadow-lg text-center font-bold mb-6">Commander une course</h2>
        <div class="mb-4">
          <label class="text-gray-300 mb-1 block">Destination :</label>
          <input
            ref="inputRef"
            v-model="search"
            @input="searchAddr"
            placeholder="Ex: 6 rue de Paris..."
            autocomplete="off"
            class="bg-gray-800 border border-gray-700 text-white placeholder-gray-400 rounded-lg px-4 py-3 mb-3 w-full focus:outline-none focus:ring-2 focus:ring-violet-500"
          />
          <ul v-if="suggestions.length" class="bg-gray-800 border border-gray-700 rounded-lg shadow-lg max-h-60 overflow-auto absolute z-50 w-full text-white">
            <li 
              v-for="addr in suggestions" 
              :key="addr.place_id" 
              @click="pickSuggestion(addr)" 
              class="cursor-pointer px-4 py-3 hover:bg-violet-700 transition"
            >
              {{ addr.description }}
            </li>
          </ul>
        </div>
        <button 
          @click="requestCourse" 
          :disabled="!destination" 
          class="w-full py-3 rounded-lg bg-emerald-600 hover:bg-emerald-700 text-white text-lg font-bold shadow transition disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Valider
        </button>
        <div v-if="error" class="mt-4 rounded-lg bg-red-700 bg-opacity-80 text-white p-3 text-center font-semibold">
          {{ error }}
        </div>
        <div v-if="price" class="mt-6 text-gray-300 space-y-1">
          <div class="font-bold text-white">Prix estimé : <span class="font-mono">{{ price }} €</span></div>
          <div v-if="km">Distance : <span class="font-mono">{{ km }} km</span></div>
          <div v-if="duration">Durée estimée : <span class="font-mono">{{ duration }}</span></div>
        </div>
      </template>
    </div>
    <div class="w-full max-w-4xl bg-gray-900 bg-opacity-80 border border-gray-800 rounded-2xl shadow-2xl p-3 h-full flex items-center justify-center">
      <GoogleMap
        :from="mapFrom"
        :to="mapTo"
        :follow="activeCourse && activeCourse.status === 'in_progress'"
        :zoomToCurrent="zoomOnMe"
        :currentPosition="myPosition"
        class="w-full h-[500px] rounded-2xl"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import GoogleMap from '../components/GoogleMap.vue'

const myPosition = ref(null)
const zoomOnMe = ref(false)
const activeCourse = ref(null)
let geoWatchId = null

function startGeoloc() {
  if (navigator.geolocation) {
    geoWatchId = navigator.geolocation.watchPosition(
      pos => {
        myPosition.value = {
          lat: pos.coords.latitude,
          lng: pos.coords.longitude
        }
      },
      err => {},
      { enableHighAccuracy: true }
    )
  }
}
function stopGeoloc() {
  if (geoWatchId) {
    navigator.geolocation.clearWatch(geoWatchId)
    geoWatchId = null
  }
}
watch(() => activeCourse.value?.status, (status) => {
  if (status === 'in_progress') {
    zoomOnMe.value = true
    startGeoloc()
  } else {
    zoomOnMe.value = false
    stopGeoloc()
  }
})

const search = ref('')
const suggestions = ref([])
const destination = ref(null)
const origin = ref(null)
const price = ref('')
const km = ref('')
const error = ref('')
const autocomplete = ref(null)
const inputRef = ref(null)
const startAddress = ref('')
const endAddress = ref('')
const duration = ref('')
const cancelError = ref('')

const mapFrom = computed(() => {
  if (activeCourse.value) return { lat: activeCourse.value.start_lat, lng: activeCourse.value.start_lng }
  return origin.value
})
const mapTo = computed(() => {
  if (activeCourse.value) return { lat: activeCourse.value.end_lat, lng: activeCourse.value.end_lng }
  return destination.value
})

const canCancel = computed(() =>
  activeCourse.value && (
    activeCourse.value.status === 'requested' ||
    activeCourse.value.status === 'accepted'
  )
)

function displayStatus(status) {
  switch (status) {
    case 'requested': return "En attente d'un driver"
    case 'accepted': return "Prise en charge de votre course"
    case 'in_progress': return "Course en cours"
    case 'completed': return "Terminée"
    case 'cancelled': return "Annulée"
    default: return status
  }
}

onMounted(async () => {
  getLocation()
  setTimeout(() => { initAutocomplete() }, 800)
  await fetchActiveCourse()
})

watch(activeCourse, async (val) => {
  if (val) {
    startAddress.value = await getAddressFromCoords({ lat: val.start_lat, lng: val.start_lng })
    endAddress.value = await getAddressFromCoords({ lat: val.end_lat, lng: val.end_lng })
    await fetchCourseDuration()
  } else {
    startAddress.value = ''
    endAddress.value = ''
    duration.value = ''
  }
})

function initAutocomplete() {
  if (window.google && window.google.maps && window.google.maps.places) {
    autocomplete.value = new window.google.maps.places.AutocompleteService()
  }
}
async function searchAddr() {
  if (!autocomplete.value) { initAutocomplete(); return }
  if (search.value.length < 3) { suggestions.value = []; return }
  autocomplete.value.getPlacePredictions(
    { input: search.value, types: ['geocode'], componentRestrictions: { country: 'fr' } },
    (preds, status) => {
      suggestions.value = status === "OK" ? preds : []
    }
  )
}
function pickSuggestion(addr) {
  const service = new window.google.maps.places.PlacesService(inputRef.value)
  service.getDetails({ placeId: addr.place_id }, (result, status) => {
    if (status === "OK" && result.geometry) {
      destination.value = {
        lat: result.geometry.location.lat(),
        lng: result.geometry.location.lng()
      }
      suggestions.value = []
      calcRoute()
      search.value = addr.description
    } else {
      error.value = "Impossible de récupérer l'adresse"
    }
  })
}

function getLocation() {
  navigator.geolocation.getCurrentPosition(
    pos => {
      origin.value = { lat: pos.coords.latitude, lng: pos.coords.longitude }
    },
    err => {
      error.value = "Impossible de récupérer votre position : vérifiez les permissions navigateur ou réessayez."
    }
  )
}

async function fetchActiveCourse() {
  const token = localStorage.getItem('token')
  if (!token) return
  const payload = JSON.parse(atob(token.split('.')[1]))
  const res = await fetch('http://localhost:8006/courses?role=passenger&id=' + payload.sub, {
    headers: { Authorization: 'Bearer ' + token }
  })
  if (res.ok) {
    const data = await res.json()
    const ongoing = (data || []).find(c =>
      c.status === 'requested' ||
      c.status === 'accepted' ||
      c.status === 'in_progress'
    )
    activeCourse.value = ongoing || null
  }
}

async function calcRoute() {
  if (!origin.value || !destination.value) return
  if (!window.google || !window.google.maps) return
  const service = new window.google.maps.DirectionsService()
  service.route({
    origin: {
      lat: Number(origin.value.lat),
      lng: Number(origin.value.lng)
    },
    destination: {
      lat: Number(destination.value.lat),
      lng: Number(destination.value.lng)
    },
    travelMode: 'DRIVING'
  }, (result, status) => {
    if (status === 'OK') {
      const meters = result.routes[0].legs[0].distance.value
      km.value = (meters / 1000).toFixed(2)
      price.value = (4 + 0.8 * km.value).toFixed(2)
      duration.value = result.routes[0].legs[0].duration.text
    }
  })
}

async function fetchCourseDuration() {
  if (!activeCourse.value) return
  if (!window.google || !window.google.maps) return
  const service = new window.google.maps.DirectionsService()
  service.route({
    origin: {
      lat: activeCourse.value.start_lat,
      lng: activeCourse.value.start_lng
    },
    destination: {
      lat: activeCourse.value.end_lat,
      lng: activeCourse.value.end_lng
    },
    travelMode: 'DRIVING'
  }, (result, status) => {
    if (status === 'OK') {
      duration.value = result.routes[0].legs[0].duration.text
    }
  })
}

async function requestCourse() {
  if (!origin.value || !destination.value) return
  const token = localStorage.getItem('token')
  const payload = JSON.parse(atob(token.split('.')[1]))
  const passengerName = localStorage.getItem('name') || payload.name || "Nom inconnu"
  const res = await fetch('http://localhost:8006/courses', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Authorization: 'Bearer ' + token },
    body: JSON.stringify({
      passenger_id: payload.sub,
      passenger_name: passengerName,
      passenger_email: payload.email,
      start_lat: origin.value.lat,
      start_lng: origin.value.lng,
      end_lat: destination.value.lat,
      end_lng: destination.value.lng,
      distance_km: parseFloat(km.value),
      co2: 0
    })
  })
  if (res.ok) {
    alert("Course créée, en attente d'un driver !")
    fetchActiveCourse()
  } else {
    const err = await res.json().catch(()=>null)
    error.value = err?.error || "Erreur lors de la création de la course"
  }
}

async function cancelCourse() {
  cancelError.value = ''
  if (!activeCourse.value) return
  const token = localStorage.getItem('token')
  const res = await fetch(`http://localhost:8006/courses/${activeCourse.value.id}/cancel`, {
    method: 'PATCH',
    headers: { Authorization: 'Bearer ' + token }
  })
  if (res.ok) {
    await fetchActiveCourse()
  } else {
    cancelError.value = 'Erreur lors de l\'annulation de la course.'
  }
}

async function getAddressFromCoords({lat, lng}) {
  const GOOGLE_MAPS_API_KEY = import.meta.env.VITE_GOOGLE_MAPS_API_KEY
  const api = `https://maps.googleapis.com/maps/api/geocode/json?latlng=${lat},${lng}&key=${GOOGLE_MAPS_API_KEY}`
  const res = await fetch(api)
  const data = await res.json()
  return data.results?.[0]?.formatted_address || ''
}
</script>