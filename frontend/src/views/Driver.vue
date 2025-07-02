<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-900 via-gray-950 to-violet-900 flex flex-col items-center px-2 py-10 gap-8">
    <DriverRequestedCourses
      v-if="displayedCourses.length && displayedCourses[0].status === 'requested'"
      :displayedCourses="displayedCourses"
      :requestedAddresses="requestedAddresses"
      :previewCourseId="previewCourseId"
      @preview="previewCourseId = $event"
      @accept="accept"
    />
    <DriverIdleMessage v-else-if="isIdle" />
    <DriverActiveCourse
      v-else-if="displayedCourses.length"
      :displayedCourses="displayedCourses"
      :driverId="driverId"
      @start="startCourse"
      @complete="completeCourse"
    >
      <template #map>
        <DriverMapContainer
          :from="mapFrom"
          :to="mapTo"
          :follow="activeCourse && activeCourse.status === 'in_progress'"
          :zoomToCurrent="zoomOnMe"
          :currentPosition="myPosition"
        />
      </template>
      <template #default="{ course }">
        <!-- Affiche les infos de la course (Départ/Arrivée/etc), réutilise ce que tu avais -->
        <div class="flex items-center justify-between">
          <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Départ</span>
          <span class="text-emerald-400 font-mono text-base font-semibold text-right">
            {{ acceptedAddresses[course.id]?.start || 'Chargement...' }}
          </span>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Arrivée</span>
          <span class="text-violet-400 font-mono text-base font-semibold text-right">
            {{ acceptedAddresses[course.id]?.end || 'Chargement...' }}
          </span>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Distance</span>
          <span class="font-mono text-blue-200 text-base font-semibold text-right">
            <span class="inline-block bg-blue-900/70 px-2 py-1 rounded-md">{{ course.distance_km }} km</span>
          </span>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Prix</span>
          <span class="font-mono text-amber-300 text-lg font-bold text-right">
            <span class="inline-block bg-amber-800/80 px-3 py-1 rounded-lg">{{ course.amount }} €</span>
          </span>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Statut</span>
          <span>
            <span v-if="course.status === 'accepted'" class="inline-block px-3 py-1 rounded-full bg-emerald-600 text-white text-xs font-semibold">Prise en charge</span>
            <span v-else-if="course.status === 'in_progress'" class="inline-block px-3 py-1 rounded-full bg-blue-600 text-white text-xs font-semibold">En cours</span>
            <span v-else-if="course.status === 'completed'" class="inline-block px-3 py-1 rounded-full bg-violet-700 text-white text-xs font-semibold">Terminée</span>
            <span v-else-if="course.status === 'cancelled'" class="inline-block px-3 py-1 rounded-full bg-rose-700 text-white text-xs font-semibold">Annulée</span>
          </span>
        </div>
      </template>
    </DriverActiveCourse>
    <!-- Map container alone for requested and idle (si besoin) -->
    <DriverMapContainer
      v-if="displayedCourses.length && displayedCourses[0].status === 'requested' || isIdle"
      :from="mapFrom"
      :to="mapTo"
      :follow="activeCourse && activeCourse.status === 'in_progress'"
      :zoomToCurrent="zoomOnMe"
      :currentPosition="myPosition"
    />
  </div>
</template>



<script setup>
import { ref, onMounted, onUnmounted, watch, computed } from 'vue'
import GoogleMap from '../components/GoogleMap.vue'
import DriverRequestedCourses from '../components/DriverRequestedCourses.vue'
import DriverActiveCourse from '../components/DriverActiveCourse.vue'
import DriverIdleMessage from '../components/DriverIdleMessage.vue'
import DriverMapContainer from '../components/DriverMapContainer.vue'

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