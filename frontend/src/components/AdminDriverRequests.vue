<template>
  <div class="mb-10">
    <h3 class="text-xl font-bold mb-4 text-white">Demandes driver</h3>
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
    class="ml-2 px-3 py-1 rounded bg-violet-700 hover:bg-violet-800 text-white font-semibold shadow transition">
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
    <div class="bg-gray-900 rounded-lg p-6 w-full max-w-md shadow-xl border border-violet-800">
      <h3 class="text-lg font-bold text-white mb-4">Modifier la demande #{{ editingRequest.ID }}</h3>
      <form @submit.prevent="saveEdit">
        <div class="mb-3">
          <label class="block text-gray-300 mb-1">Voiture</label>
          <input v-model="editingRequest.car" class="input" type="text" />
        </div>
        <div class="mb-3">
          <label class="block text-gray-300 mb-1">Plaque</label>
          <input v-model="editingRequest.plate" class="input" type="text" />
        </div>
        <div class="mb-3">
          <label class="block text-gray-300 mb-1">Statut</label>
          <select v-model="editingRequest.status" class="input">
            <option value="pending">En attente</option>
            <option value="accepted">Acceptée</option>
            <option value="refused">Refusée</option>
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
const res = await fetch(`http://localhost:8002/driver-requests/${editingRequest.value.ID}/edit`, {
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