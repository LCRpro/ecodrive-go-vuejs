<template>
  <div class="mb-10">
    <h3 class="text-xl font-bold mb-4 text-white">Toutes les courses</h3>
    <div class="overflow-x-auto w-full">
      <table class="min-w-[850px] w-full text-sm border border-gray-700 rounded-xl overflow-hidden bg-gray-900">
        <thead>
          <tr>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              ID</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Départ</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Arrivée</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Passager</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Driver</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Distance</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Prix (€)</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              CO₂ (kg)</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Statut</th>
          <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
  Actions
</th>
            </tr>

        </thead>
        <tbody>
          <tr v-for="c in courses" :key="c.id" class="even:bg-gray-900 hover:bg-gray-800/80">
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ c.id }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ c.start_lat }}, {{
              c.start_lng }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ c.end_lat }}, {{
              c.end_lng }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{
              getUserByGoogleId(c.passenger_id)?.email || c.passenger_id }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{
              getUserByGoogleId(c.driver_id)?.email || c.driver_id || '-' }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{
              c.distance_km?.toFixed(2) ?? '-' }} km</td>
            <td class="text-gray-200 border border-gray-700 text-center font-mono px-3 py-2 md:px-4 md:py-3">{{
              c.amount?.toFixed(2) ?? '-' }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ c.co2 ?
              (c.co2 / 1000).toFixed(2) : '-' }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">
              <span v-if="c.status === 'requested'"
                class="px-2 py-1 rounded-full bg-amber-500 text-white text-xs font-semibold">En attente</span>
              <span v-else-if="c.status === 'accepted'"
                class="px-2 py-1 rounded-full bg-emerald-600 text-white text-xs font-semibold">Acceptée</span>
              <span v-else-if="c.status === 'in_progress'"
                class="px-2 py-1 rounded-full bg-blue-600 text-white text-xs font-semibold">En cours</span>
              <span v-else-if="c.status === 'completed'"
                class="px-2 py-1 rounded-full bg-violet-600 text-white text-xs font-semibold">Terminée</span>
              <span v-else-if="c.status === 'cancelled'"
                class="px-2 py-1 rounded-full bg-rose-600 text-white text-xs font-semibold">Annulée</span>
              <span v-else class="px-2 py-1 rounded-full bg-gray-700 text-gray-200 text-xs font-semibold">{{ c.status
                }}</span>
            </td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">
  <button
    v-if="isAdmin"
    class="px-2 py-1 rounded bg-violet-700 hover:bg-violet-800 text-white text-xs font-semibold"
    @click="openEditModal(c)"
  >
    Éditer
  </button>
</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
<template v-if="showEditModal">
  <div class="fixed inset-0 z-50 bg-black/60 flex items-center justify-center">
    <div class="bg-gray-900 rounded-lg p-6 w-full max-w-lg shadow-xl border border-violet-800">
      <h3 class="text-lg font-bold text-white mb-4">Modifier la course #{{ editingCourse.id }}</h3>
      <form @submit.prevent="saveEdit">
        <!-- Exemples de champs, adapte selon les infos de la course -->
        <div class="mb-3">
          <label class="block text-gray-300 mb-1">Départ (lat, lng)</label>
          <input v-model="editingCourse.start_lat" class="input" type="text" />
          <input v-model="editingCourse.start_lng" class="input" type="text" />
        </div>
        <div class="mb-3">
          <label class="block text-gray-300 mb-1">Arrivée (lat, lng)</label>
          <input v-model="editingCourse.end_lat" class="input" type="text" />
          <input v-model="editingCourse.end_lng" class="input" type="text" />
        </div>
        <div class="mb-3">
          <label class="block text-gray-300 mb-1">Passager</label>
          <input v-model="editingCourse.passenger_id" class="input" type="text" />
        </div>
        <div class="mb-3">
          <label class="block text-gray-300 mb-1">Driver</label>
          <input v-model="editingCourse.driver_id" class="input" type="text" />
        </div>
        <div class="mb-3">
          <label class="block text-gray-300 mb-1">Distance (km)</label>
          <input v-model.number="editingCourse.distance_km" class="input" type="number" step="0.01" />
        </div>
        <div class="mb-3">
          <label class="block text-gray-300 mb-1">Prix (€)</label>
          <input v-model.number="editingCourse.amount" class="input" type="number" step="0.01" />
        </div>
        <div class="mb-3">
          <label class="block text-gray-300 mb-1">CO₂ (g)</label>
          <input v-model.number="editingCourse.co2" class="input" type="number" />
        </div>
        <div class="mb-3">
          <label class="block text-gray-300 mb-1">Statut</label>
          <select v-model="editingCourse.status" class="input">
            <option value="requested">En attente</option>
            <option value="accepted">Acceptée</option>
            <option value="in_progress">En cours</option>
            <option value="completed">Terminée</option>
            <option value="cancelled">Annulée</option>
          </select>
        </div>
        <div class="flex justify-end gap-3 mt-6">
          <button type="button" @click="closeEditModal" class="btn-cancel">Annuler</button>
          <button type="submit" class="btn-save">Sauvegarder</button>
        </div>
      </form>
    </div>
  </div>
</template>
  
</template>
<script setup>
defineProps({
  courses: {
    type: Array,
    required: true,
    default: () => []
  },
  getUserByGoogleId: {
    type: Function,
    required: true
  }
})

import { ref } from 'vue'
const editingCourse = ref(null)
const showEditModal = ref(false)
const isAdmin = JSON.parse(localStorage.getItem('roles') || '[]').includes('ROLE_ADMIN')

function openEditModal(course) {
  editingCourse.value = { ...course } // Clone pour éviter d'éditer directement
  showEditModal.value = true
}
function closeEditModal() {
  showEditModal.value = false
}
async function saveEdit() {
  const token = localStorage.getItem('token')
const payload = {
  start_lat: editingCourse.value.start_lat,
  start_lng: editingCourse.value.start_lng,
  end_lat: editingCourse.value.end_lat,
  end_lng: editingCourse.value.end_lng,
  passenger_id: editingCourse.value.passenger_id,
  driver_id: editingCourse.value.driver_id,
  distance_km: editingCourse.value.distance_km,
  amount: editingCourse.value.amount,
  co2: editingCourse.value.co2,
  status: editingCourse.value.status
}
const res = await fetch(`http://localhost:8006/courses/${editingCourse.value.id}`, {
  method: 'PATCH',
  headers: { 'Content-Type': 'application/json', 'Authorization': 'Bearer ' + token },
  body: JSON.stringify(payload)
})
  if (res.ok) {
    showEditModal.value = false
    window.location.reload() // Ou un meilleur refresh
  } else {
    alert("Erreur lors de la modification")
  }
}</script>