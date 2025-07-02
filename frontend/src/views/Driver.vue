<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-900 via-gray-950 to-violet-900 flex flex-col items-center px-2 py-10 gap-8">
    <template v-if="displayedCourses.length && displayedCourses[0].status === 'requested'">
      <div class="w-full max-w-5xl flex flex-wrap gap-6 justify-center mb-8">
        <div
          v-for="c in displayedCourses"
          :key="c.id"
          class="w-full md:w-[340px] bg-gradient-to-r from-gray-800/90 via-gray-900/90 to-violet-950/60 rounded-xl px-6 py-5 shadow-inner border border-gray-800 space-y-3 flex flex-col"
        >
          <div class="flex justify-between items-center">
            <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Passager</span>
            <span class="font-mono text-emerald-400 text-base font-semibold">{{ c.passenger_name }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Départ</span>
            <span class="font-mono text-violet-400 text-base font-semibold">
              {{ requestedAddresses[c.id]?.start || 'Chargement...' }}
            </span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Arrivée</span>
            <span class="font-mono text-violet-400 text-base font-semibold">
              {{ requestedAddresses[c.id]?.end || 'Chargement...' }}
            </span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Distance</span>
            <span class="font-mono text-blue-200 text-base font-semibold">
              <span class="inline-block bg-blue-900/70 px-2 py-1 rounded-md">{{ c.distance_km }} km</span>
            </span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Prix</span>
            <span class="font-mono text-amber-300 text-lg font-bold">
              <span class="inline-block bg-amber-800/80 px-3 py-1 rounded-lg">{{ c.amount }} €</span>
            </span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Statut</span>
            <span class="inline-block px-3 py-1 rounded-full bg-amber-600 text-white text-xs font-semibold">En attente</span>
          </div>
          <div class="flex gap-2 mt-2">
            <button
              @click="previewCourseId = c.id"
              :class="previewCourseId === c.id ? 'bg-blue-700' : 'bg-blue-600 hover:bg-blue-700'"
              class="flex-1 py-2 rounded-lg text-white font-bold shadow transition"
            >Voir sur la carte</button>
            <button
              @click="accept(c.id)"
              class="flex-1 py-2 rounded-lg bg-emerald-600 hover:bg-emerald-700 text-white font-bold shadow transition"
            >Accepter</button>
          </div>
        </div>
      </div>
      <div class="w-full max-w-6xl bg-gray-900 bg-opacity-80 border border-gray-800 rounded-2xl shadow-2xl p-3 flex items-center justify-center h-[520px]">
        <GoogleMap
          :from="mapFrom"
          :to="mapTo"
          :follow="activeCourse && activeCourse.status === 'in_progress'"
          :zoomToCurrent="zoomOnMe"
          :currentPosition="myPosition"
        />
      </div>
    </template>

   <template v-else-if="isIdle">
  <div class="w-full max-w-lg mx-auto mb-8">
    <div class="p-8 bg-gradient-to-r from-gray-800/90 via-gray-900/95 to-violet-950/80 rounded-2xl shadow-inner border border-gray-800 flex flex-col items-center gap-4">
      <span class="text-4xl text-gray-300 mb-2">🚗</span>
      <span class="text-lg text-white font-semibold text-center">Aucune course à prendre pour le moment.</span>
      <span class="text-gray-400 text-sm text-center">Reviens plus tard, de nouvelles courses peuvent arriver à tout moment.</span>
    </div>
  </div>
  <div class="w-full max-w-6xl bg-gray-900 bg-opacity-80 border border-gray-800 rounded-2xl shadow-2xl p-3 flex items-center justify-center h-[520px]">
    <GoogleMap
      :from="null"
      :to="null"
      :follow="false"
      :zoomToCurrent="zoomOnMe"
      :currentPosition="myPosition"
    />
  </div>
</template>

    <template v-else-if="displayedCourses.length">
      <div class="w-full flex flex-col md:flex-row items-center justify-center gap-8">
        <div class="w-full max-w-sm p-7 rounded-2xl shadow-2xl bg-gray-900 bg-opacity-90 border border-gray-800 mb-8 md:mb-0">
          <template v-for="c in displayedCourses" :key="c.id">
            <h2 class="text-white text-2xl font-bold mb-6">Course</h2>
            <div class="bg-gradient-to-r from-gray-800/90 via-gray-900/95 to-violet-950/70 rounded-xl px-6 py-5 mb-3 shadow-inner border border-gray-800 space-y-4">
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Départ</span>
                <span class="text-emerald-400 font-mono text-base font-semibold text-right">
                  {{ acceptedAddresses[c.id]?.start || 'Chargement...' }}
                </span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Arrivée</span>
                <span class="text-violet-400 font-mono text-base font-semibold text-right">
                  {{ acceptedAddresses[c.id]?.end || 'Chargement...' }}
                </span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Distance</span>
                <span class="font-mono text-blue-200 text-base font-semibold text-right">
                  <span class="inline-block bg-blue-900/70 px-2 py-1 rounded-md">{{ c.distance_km }} km</span>
                </span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Prix</span>
                <span class="font-mono text-amber-300 text-lg font-bold text-right">
                  <span class="inline-block bg-amber-800/80 px-3 py-1 rounded-lg">{{ c.amount }} €</span>
                </span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Statut</span>
                <span>
                  <span v-if="c.status === 'accepted'" class="inline-block px-3 py-1 rounded-full bg-emerald-600 text-white text-xs font-semibold">Prise en charge</span>
                  <span v-else-if="c.status === 'in_progress'" class="inline-block px-3 py-1 rounded-full bg-blue-600 text-white text-xs font-semibold">En cours</span>
                  <span v-else-if="c.status === 'completed'" class="inline-block px-3 py-1 rounded-full bg-violet-700 text-white text-xs font-semibold">Terminée</span>
                  <span v-else-if="c.status === 'cancelled'" class="inline-block px-3 py-1 rounded-full bg-rose-700 text-white text-xs font-semibold">Annulée</span>
                </span>
              </div>
            </div>
            <div class="mt-4 flex gap-2">
              <button v-if="c.status === 'accepted' && c.driver_id === driverId"
                @click="startCourse(c.id)"
                class="w-full py-3 rounded-lg bg-blue-600 hover:bg-blue-700 text-white text-lg font-bold shadow transition">
                Commencer la course
              </button>
              <button v-if="c.status === 'in_progress' && c.driver_id === driverId"
                @click="completeCourse(c.id)"
                class="w-full py-3 rounded-lg bg-violet-700 hover:bg-violet-800 text-white text-lg font-bold shadow transition">
                Terminer la course
              </button>
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
  </div>
</template>




<script setup>
import { ref, onMounted, onUnmounted, watch, computed } from 'vue'
import GoogleMap from '../components/GoogleMap.vue'

const courses = ref([])
const activeCourse = ref(null)
const driverId = ref('')

const myPosition = ref(null)
const zoomOnMe = ref(false)
let geoWatchId = null

const activeStatuses = ['accepted', 'in_progress']
const displayedCourses = computed(() => {
  const active = courses.value.find(
    c => activeStatuses.includes(c.status) && c.driver_id === driverId.value
  )
  if (active) {
    activeCourse.value = active
    return [active]
  }
  activeCourse.value = null
  return courses.value.filter(c => c.status === 'requested')
})

const requestedAddresses = ref({})
const acceptedAddresses = ref({})
const previewCourseId = ref(null)

const previewFrom = computed(() => {
  const c = displayedCourses.value.find(c => c.id === previewCourseId.value)
  return c ? { lat: c.start_lat, lng: c.start_lng } : null
})
const previewTo = computed(() => {
  const c = displayedCourses.value.find(c => c.id === previewCourseId.value)
  return c ? { lat: c.end_lat, lng: c.end_lng } : null
})

const mapFrom = computed(() => {
  if (previewCourseId.value) return previewFrom.value
  if (activeCourse.value) return { lat: activeCourse.value.start_lat, lng: activeCourse.value.start_lng }
  return null
})
const mapTo = computed(() => {
  if (previewCourseId.value) return previewTo.value
  if (activeCourse.value) return { lat: activeCourse.value.end_lat, lng: activeCourse.value.end_lng }
  return null
})

const GOOGLE_API_KEY = import.meta.env.VITE_GOOGLE_MAPS_API_KEY
const addressCache = new Map()

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
  } catch {}
  addressCache.set(key, `${lat}, ${lng}`)
  return `${lat}, ${lng}`
}

onMounted(() => {
  fetchCourses()
  if (navigator.geolocation) {
    geoWatchId = navigator.geolocation.watchPosition(
      pos => {
        myPosition.value = {
          lat: pos.coords.latitude,
          lng: pos.coords.longitude
        }

      },
      err => {
        console.warn("Erreur géoloc :", err)
      },
      { enableHighAccuracy: true }
    )
  }
})

onUnmounted(() => {
  if (geoWatchId) navigator.geolocation.clearWatch(geoWatchId)
})

watch(displayedCourses, async (val) => {
  if (!val.length) return
  if (val[0].status === "requested") {
    for (const c of val) {
      if (!requestedAddresses.value[c.id]) {
        const start = await reverseGeocode(c.start_lat, c.start_lng)
        const end = await reverseGeocode(c.end_lat, c.end_lng)
        requestedAddresses.value = { ...requestedAddresses.value, [c.id]: { start, end } }
      }
    }
  } else {
    for (const c of val) {
      if (!acceptedAddresses.value[c.id]) {
        const start = await reverseGeocode(c.start_lat, c.start_lng)
        const end = await reverseGeocode(c.end_lat, c.end_lng)
        acceptedAddresses.value = { ...acceptedAddresses.value, [c.id]: { start, end } }
      }
    }
  }
}, { immediate: true })

watch(() => activeCourse.value?.status, (status) => {
  if (status === 'in_progress') {
    zoomOnMe.value = true
  } else {
    zoomOnMe.value = false
  }
})

async function fetchCourses() {
  const token = localStorage.getItem('token')
  const payload = JSON.parse(atob(token.split('.')[1]))
  driverId.value = payload.sub
  const res = await fetch(`http://localhost:8006/courses?driverView=1&driver_id=${payload.sub}`, {
    headers: { Authorization: 'Bearer ' + token }
  })
  if (res.ok) {
    const data = await res.json()
    courses.value = data
  } else {
    courses.value = []
  }
}

function accept(id) {
  const token = localStorage.getItem('token')
  const payload = JSON.parse(atob(token.split('.')[1]))
  fetch(`http://localhost:8006/courses/${id}/accept`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json', Authorization: 'Bearer ' + token },
    body: JSON.stringify({ driver_id: payload.sub })
  }).then(async res => {
    if (res.ok) {
      fetchCourses()
    } else {
      alert('Erreur lors de l’acceptation')
    }
  })
}

function startCourse(id) {
  const token = localStorage.getItem('token')
  fetch(`http://localhost:8006/courses/${id}/start`, {
    method: 'PATCH',
    headers: { Authorization: 'Bearer ' + token }
  }).then(async res => {
    if (res.ok) {
      fetchCourses()
    } else {
      alert('Erreur lors du démarrage')
    }
  })
}
const isIdle = computed(() => {
  return (
    courses.value.length === 0 || (
      displayedCourses.value.length === 0 &&
      !courses.value.find(
        c => ['accepted', 'in_progress'].includes(c.status) && c.driver_id === driverId.value
      )
    )
  )
})
function completeCourse(id) {
  const token = localStorage.getItem('token')
  fetch(`http://localhost:8006/courses/${id}/complete`, {
    method: 'PATCH',
    headers: { Authorization: 'Bearer ' + token }
  }).then(async res => {
    if (res.ok) {
      fetchCourses()
    } else {
      alert('Erreur lors de la complétion')
    }
  })
}
</script>