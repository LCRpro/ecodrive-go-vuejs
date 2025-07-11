<template>
  <div class="mb-10">
  <div class="flex flex-col items-center mb-10">
  <h3
    class="text-3xl md:text-4xl font-extrabold bg-gradient-to-r from-emerald-400 via-violet-400 to-emerald-400 bg-clip-text text-transparent drop-shadow-lg text-center"
  >
    Demande driver
  </h3>
</div>
    <div class="overflow-x-auto w-full">
      <table class="min-w-[700px] w-full text-sm border border-gray-700 rounded-xl overflow-hidden bg-gray-900">
        <thead>
          <tr>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              ID</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Nom</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Email</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Voiture</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Plaque</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Statut</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Action</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="req in driverRequests" :key="req.ID" class="even:bg-gray-900 hover:bg-gray-800/80">
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ req.ID }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{
              getUserByGoogleId(req.google_id || req.GoogleID)?.name || '-' }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{
              getUserByGoogleId(req.google_id || req.GoogleID)?.email || '-' }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ req.car }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ req.plate }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">
              <span v-if="req.status === 'pending'"
                class="px-2 py-1 rounded-full bg-amber-500 text-white text-xs font-semibold">En attente</span>
              <span v-else-if="req.status === 'accepted'"
                class="px-2 py-1 rounded-full bg-emerald-600 text-white text-xs font-semibold">Acceptée</span>
              <span v-else-if="req.status === 'refused'"
                class="px-2 py-1 rounded-full bg-rose-600 text-white text-xs font-semibold">Refusée</span>
              <span v-else class="px-2 py-1 rounded-full bg-gray-700 text-gray-200 text-xs font-semibold">{{ req.status
                }}</span>
            </td>
<td class="text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
  <button
    v-if="req.status === 'pending'"
    @click="$emit('handleRequest', req.ID, 'accept')"
    class="px-3 py-1 rounded bg-emerald-600 hover:bg-emerald-700 text-white font-semibold shadow transition">
    Accepter
  </button>
  <button
    v-if="req.status === 'pending'"
    @click="$emit('handleRequest', req.ID, 'refuse')"
    class="px-3 py-1 rounded bg-rose-600 hover:bg-rose-700 text-white font-semibold shadow transition ml-2">
    Refuser
  </button>
  <button
    v-if="isAdmin"
    @click="openEditModal(req)"
    class="ml-2 px-3 py-1 rounded   text-white font-semibold shadow transition">
    ✏️
  </button>
</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
<template v-if="showEditModal">
  <div class="fixed inset-0 z-50 flex items-center justify-center  from-gray-900 via-gray-950 to-violet-900 bg-opacity-90 px-2 py-10">
    <div class="w-full max-w-sm mx-auto p-7 rounded-2xl shadow-2xl bg-gray-900 bg-opacity-95 border border-violet-800">
      <div class="flex items-center gap-3 mb-6">
        <svg class="w-10 h-10 text-violet-400" fill="none" viewBox="0 0 24 24">
          <rect width="20" height="14" x="2" y="5" rx="3" fill="#fff"/>
          <rect width="20" height="14" x="2" y="5" rx="3" stroke="#7c3aed" stroke-width="1.5"/>
          <rect width="6" height="2" x="4" y="13" fill="#d1d5db"/>
          <rect width="4" height="2" x="14" y="13" fill="#a78bfa"/>
        </svg>
        <h2 class="text-2xl font-bold text-white tracking-tight">Modification demande #{{ editingRequest.ID }}</h2>
      </div>
      <form @submit.prevent="saveEdit" autocomplete="off" class="space-y-5">
        <div>
          <label class="block text-gray-300 mb-1 font-medium">Voiture</label>
          <input
            v-model="editingRequest.car"
            class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400
                   focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-semibold transition-all duration-200 font-mono"
            type="text"
            placeholder="Ex: Renault Clio"
          />
        </div>
        <div>
          <label class="block text-gray-300 mb-1 font-medium">Plaque</label>
          <input
            v-model="editingRequest.plate"
            class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400
                   focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-semibold transition-all duration-200 font-mono"
            type="text"
            placeholder="Ex: AB-123-CD"
          />
        </div>
        <div>
          <label class="block text-gray-300 mb-1 font-medium">Statut</label>
          <select
            v-model="editingRequest.status"
            class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white
                   focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-semibold transition-all duration-200"
          >
            <option value="pending">En attente</option>
            <option value="accepted">Acceptée</option>
            <option value="refused">Refusée</option>
          </select>
        </div>
        <div class="flex justify-end gap-3 mt-8">
          <button type="button"
                  @click="closeEditModal"
                  class="px-5 py-2 rounded-lg bg-gray-700 hover:bg-gray-800 text-white font-semibold transition">Annuler</button>
          <button type="submit"
                  class="px-5 py-2 rounded-lg bg-violet-700 hover:bg-violet-800 text-white font-bold shadow transition">Sauvegarder</button>
        </div>
        <div v-if="error" class="bg-red-700/80 text-white text-center py-2 rounded-lg font-semibold mt-2">{{ error }}</div>
        <div v-if="success" class="bg-emerald-700/80 text-white text-center py-2 rounded-lg font-semibold mt-2">{{ success }}</div>
      </form>
    </div>
  </div>
</template>
</template>
<script setup>
defineProps(['driverRequests', 'getUserByGoogleId'])
defineEmits(['handleRequest'])
const isAdmin = JSON.parse(localStorage.getItem('roles') || '[]').includes('ROLE_ADMIN')
import { ref } from 'vue'
const editingRequest = ref(null)
const showEditModal = ref(false)

function openEditModal(req) {
  editingRequest.value = { ...req }
  showEditModal.value = true
}
function closeEditModal() {
  showEditModal.value = false
}
async function saveEdit() {
  const token = localStorage.getItem('token')
const payload = {
  google_id: editingRequest.value.google_id,
  car: editingRequest.value.car,
  plate: editingRequest.value.plate,
  status: editingRequest.value.status,
}
const res = await fetch(`https://user-ecodrive.liamcariou.fr/driver-requests/${editingRequest.value.ID}/edit`, {
  method: 'PATCH',
  headers: { 'Content-Type': 'application/json', 'Authorization': 'Bearer ' + token },
  body: JSON.stringify(payload),
})
  if (res.ok) {
    showEditModal.value = false
    window.location.reload()
  } else {
    alert("Erreur lors de la modification")
  }
}
</script>