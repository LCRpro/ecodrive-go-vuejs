<template>
  <div class="flex" style="gap:32px;">
    <!-- Colonne gauche -->
    <div style="flex:1;max-width:350px;">
      <h2>Demandes de courses</h2>
      <div v-for="c in courses" :key="c.id" style="background:#eee;margin:10px 0;padding:10px;border-radius:8px;">
        <div>Passager : {{ c.passenger_email }}</div>
        <div>Nom : {{ c.passenger_name }}</div>
        <div>De : {{ c.start_lat }},{{ c.start_lng }}</div>
        <div>À : {{ c.end_lat }},{{ c.end_lng }}</div>
        <div>Distance : {{ c.distance_km }} km</div>
        <div>Prix : {{ c.amount }} €</div>
        <div>Status : <b>{{ displayStatus(c.status) }}</b></div>
        <div v-if="c.status === 'requested'">
          <button @click="accept(c.id)">Accepter</button>
        </div>
        <div v-else-if="c.status === 'accepted' && c.driver_id === driverId">
          <button @click="startCourse(c.id)">Commencer la course</button>
        </div>
        <div v-else-if="c.status === 'in_progress' && c.driver_id === driverId">
          <button @click="completeCourse(c.id)">Terminer la course</button>
        </div>
        <div v-else-if="c.status === 'completed'">
          <span>Course terminée</span>
        </div>
        <div v-else>
          <span>En attente...</span>
        </div>
      </div>
    </div>
    <!-- Colonne droite -->
    <div style="flex:2;">
      <GoogleMap
        :from="mapFrom"
        :to="mapTo"
        :follow="activeCourse && activeCourse.status === 'in_progress'"
        :zoomToCurrent="zoomOnMe"
        :currentPosition="myPosition"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import GoogleMap from '../components/GoogleMap.vue'

const myPosition = ref(null)
const zoomOnMe = ref(false)
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

const courses = ref([])
const activeCourse = ref(null)
const driverId = ref('')

function updateActiveCourse() {
  const mine = courses.value.find(c =>
    (c.status === 'accepted' || c.status === 'in_progress') && c.driver_id === driverId.value
  )
  activeCourse.value = mine || null
}

const mapFrom = computed(() => {
  if (activeCourse.value) return { lat: activeCourse.value.start_lat, lng: activeCourse.value.start_lng }
  return null
})
const mapTo = computed(() => {
  if (activeCourse.value) return { lat: activeCourse.value.end_lat, lng: activeCourse.value.end_lng }
  return null
})

watch(() => activeCourse.value?.status, (status) => {
  if (status === 'in_progress') {
    zoomOnMe.value = true
    startGeoloc()
  } else {
    zoomOnMe.value = false
    stopGeoloc()
  }
})

function displayStatus(status) {
  switch (status) {
    case 'requested': return "En attente"
    case 'accepted': return "Course acceptée"
    case 'in_progress': return "En cours"
    case 'completed': return "Terminée"
    case 'cancelled': return "Annulée"
    default: return status
  }
}

async function fetchCourses() {
  const token = localStorage.getItem('token')
  const payload = JSON.parse(atob(token.split('.')[1]))
  driverId.value = payload.sub
  const res = await fetch('http://localhost:8006/courses?role=driver&id=' + payload.sub, {
    headers: { Authorization: 'Bearer ' + token }
  })
  courses.value = res.ok ? await res.json() : []
  updateActiveCourse()
}

onMounted(fetchCourses)

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