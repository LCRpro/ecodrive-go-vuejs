<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-900 via-gray-950 to-violet-900 flex flex-col md:flex-row items-center justify-center px-2 py-10 gap-8">
    <CourseInfo
      :activeCourse="activeCourse"
      :startAddress="startAddress"
      :endAddress="endAddress"
      :duration="duration"
      :canCancel="canCancel"
      :cancelError="cancelError"
      v-model:search="search"
      :suggestions="suggestions"
      :destination="destination"
      :error="error"
      :price="price"
      :km="km"
      @refresh="fetchActiveCourse"
    />
    <CourseMap
      :from="mapFrom"
      :to="mapTo"
      :follow="activeCourse && activeCourse.status === 'in_progress'"
      :zoomToCurrent="zoomOnMe"
      :currentPosition="myPosition"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import GoogleMap from '../components/GoogleMap.vue'
import CourseInfo from '../components/CourseInfo.vue'
import CourseMap from '../components/CourseMap.vue'
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