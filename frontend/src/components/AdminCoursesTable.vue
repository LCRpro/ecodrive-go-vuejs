<template>
  <div class="mb-10">
    <div class="flex flex-col items-center mb-10">
      <h3
        class="text-3xl md:text-4xl font-extrabold bg-gradient-to-r from-emerald-400 via-violet-400 to-emerald-400 bg-clip-text text-transparent drop-shadow-lg text-center">
        Toutes les courses
      </h3>
    </div>
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
              getUserByGoogleId(c.passenger_id)?.name || c.passenger_id }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{
              getUserByGoogleId(c.driver_id)?.name || c.driver_id || '-' }}</td>
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
              <button v-if="isAdmin" class="px-2 py-1 rounded text-white text-xs font-semibold"
                @click="openEditModal(c)">
                ✏️
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
  <template v-if="showEditModal">
    <div
      class="fixed inset-0 z-50 flex items-center justify-center  from-gray-900 via-gray-950 to-violet-900 bg-opacity-90 px-2 py-10">
      <div
        class="w-full max-w-2xl mx-auto p-8 rounded-2xl shadow-2xl bg-gray-900 bg-opacity-95 border border-violet-800">
        <div class="flex items-center gap-3 mb-6">
          <svg class="w-10 h-10 text-violet-400" fill="none" viewBox="0 0 24 24">
            <rect width="20" height="14" x="2" y="5" rx="3" fill="#fff" />
            <rect width="20" height="14" x="2" y="5" rx="3" stroke="#7c3aed" stroke-width="1.5" />
            <rect width="6" height="2" x="4" y="13" fill="#d1d5db" />
            <rect width="4" height="2" x="14" y="13" fill="#a78bfa" />
          </svg>
          <h2 class="text-2xl font-bold text-white tracking-tight">Modifier la course #{{ editingCourse.id }}</h2>
        </div>
        <form @submit.prevent="saveEdit" autocomplete="off" class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label class="block text-gray-300 mb-1 font-medium">Départ - Latitude</label>
            <input v-model="editingCourse.start_lat"
              class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white
            focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-mono font-semibold transition-all duration-200" type="text"
              placeholder="Latitude" />
          </div>
          <div>
            <label class="block text-gray-300 mb-1 font-medium">Départ - Longitude</label>
            <input v-model="editingCourse.start_lng"
              class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white
            focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-mono font-semibold transition-all duration-200" type="text"
              placeholder="Longitude" />
          </div>
          <div>
            <label class="block text-gray-300 mb-1 font-medium">Arrivée - Latitude</label>
            <input v-model="editingCourse.end_lat"
              class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white
            focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-mono font-semibold transition-all duration-200"
              type="text" placeholder="Latitude" />
          </div>
          <div>
            <label class="block text-gray-300 mb-1 font-medium">Arrivée - Longitude</label>
            <input v-model="editingCourse.end_lng"
              class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white
            focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-mono font-semibold transition-all duration-200"
              type="text" placeholder="Longitude" />
          </div>
          <div>
            <label class="block text-gray-300 mb-1 font-medium">Passager (ID)</label>
            <input v-model="editingCourse.passenger_id"
              class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white
            focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-mono font-semibold transition-all duration-200" type="text" />
          </div>
          <div>
            <label class="block text-gray-300 mb-1 font-medium">Driver (ID)</label>
            <input v-model="editingCourse.driver_id"
              class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white
            focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-mono font-semibold transition-all duration-200"
              type="text" />
          </div>
          <div>
            <label class="block text-gray-300 mb-1 font-medium">Distance (km)</label>
            <input v-model.number="editingCourse.distance_km"
              class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white
            focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-mono font-semibold transition-all duration-200" type="number"
              step="0.01" />
          </div>
          <div>
            <label class="block text-gray-300 mb-1 font-medium">Prix (€)</label>
            <input v-model.number="editingCourse.amount"
              class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white
            focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-mono font-semibold transition-all duration-200" type="number"
              step="0.01" />
          </div>
          <div>
            <label class="block text-gray-300 mb-1 font-medium">CO₂ (g)</label>
            <input v-model.number="editingCourse.co2"
              class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white
            focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-mono font-semibold transition-all duration-200"
              type="number" />
          </div>
          <div>
            <label class="block text-gray-300 mb-1 font-medium">Statut</label>
            <select v-model="editingCourse.status"
              class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white
            focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-semibold transition-all duration-200">
              <option value="requested">En attente</option>
              <option value="accepted">Acceptée</option>
              <option value="in_progress">En cours</option>
              <option value="completed">Terminée</option>
              <option value="cancelled">Annulée</option>
            </select>
          </div>
          <div class="md:col-span-2 flex justify-end gap-3 mt-2">
            <button type="button" @click="closeEditModal"
              class="px-5 py-2 rounded-lg bg-gray-700 hover:bg-gray-800 text-white font-semibold transition">Annuler</button>
            <button type="submit"
              class="px-5 py-2 rounded-lg bg-violet-700 hover:bg-violet-800 text-white font-bold shadow transition">Sauvegarder</button>
          </div>
          <div v-if="error"
            class="md:col-span-2 bg-red-700/80 text-white text-center py-2 rounded-lg font-semibold mt-2">{{ error }}
          </div>
        </form>
      </div>
    </div>
  </template>

</template>
<script setup>

const driverServiceURL = import.meta.env.VITE_DRIVER_SERVICE_URL


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
  editingCourse.value = { ...course }
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
  const res = await fetch(`${driverServiceURL}/courses/${editingCourse.value.id}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json', 'Authorization': 'Bearer ' + token },
    body: JSON.stringify(payload)
  })
  if (res.ok) {
    showEditModal.value = false
    window.location.reload()
  } else {
    alert("Erreur lors de la modification")
  }
}</script>