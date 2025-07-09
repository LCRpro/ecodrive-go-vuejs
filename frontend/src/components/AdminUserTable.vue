<template>
  <div class="mb-10">
    <h3 class="text-xl font-bold mb-4 text-white">Utilisateurs</h3>
    <div class="overflow-x-auto w-full">
      <table class="min-w-[700px] w-full text-sm border border-gray-700 rounded-xl overflow-hidden bg-gray-900">
        <thead>
          <tr>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              ID</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Email</th>
                <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Numero de téléphone</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Nom</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Rôles</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Voiture</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Plaque</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Actions</th>

          </tr>
        </thead>
        <tbody>
          <tr v-for="u in users" :key="u.id" class="even:bg-gray-900 hover:bg-gray-800/80">
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ u.id }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ u.email }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ u.phone }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ u.name }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">
              <span v-for="role in JSON.parse(u.roles || '[]')" :key="role"
                class="inline-block px-2 py-1 mx-1 rounded-full bg-violet-700 text-white text-xs font-semibold">{{ role
                }}</span>
            </td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ u.car || '-' }}
            </td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ u.plate || '-' }}
            </td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">
              <button @click="openEdit(u)" class="text-violet-400 hover:text-violet-200">✏️</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>

  <div v-if="showModal" class="modal-bg">
    <div class="modal">
      <h4 class="text-lg mb-3">Modifier utilisateur #{{ editingUser?.id }}</h4>
      <div class="mb-2">
        <label>Nom :</label>
        <input v-model="editForm.name" class="input" />
      </div>
      <div class="mb-2">
        <label>Email :</label>
        <input v-model="editForm.email" class="input" />
      </div>
       <div class="mb-2">
        <label>Numero de téléphone :</label>
        <input v-model="editForm.phone" class="input" />
      </div>
      <div class="mb-2">
        <label>Voiture :</label>
        <input v-model="editForm.car" class="input" />
      </div>
      <div class="mb-2">
        <label>Plaque :</label>
        <input v-model="editForm.plate" class="input" />
      </div>
      <div class="mb-2">
        <label>Solde (€) :</label>
        <input type="number" v-model.number="editForm.balance" min="0" step="0.01" class="input" />
      </div>
      <div class="mb-2">
        <label>Rôles (JSON) :</label>
        <input v-model="editForm.roles" class="input" />
      </div>
      <div class="flex gap-3 mt-4">
        <button @click="saveEdit" class="btn bg-violet-700 text-white px-4 py-2 rounded">Enregistrer</button>
        <button @click="closeModal" class="btn bg-gray-700 text-white px-4 py-2 rounded">Annuler</button>
      </div>
      <div v-if="error" class="text-red-500 mt-2">{{ error }}</div>
    </div>
  </div>
</template>
<script setup>
import { ref, reactive, watch } from 'vue'

const props = defineProps(['users'])
const emit = defineEmits(['refresh'])

const showModal = ref(false)
const editingUser = ref(null)
const editForm = reactive({
  name: '',
  email: '',
  car: '',
  plate: '',
  balance: 0,
  roles: '',
  phone:'',
})
const error = ref('')

function openEdit(user) {
  editingUser.value = user
  editForm.name = user.name || ''
  editForm.email = user.email || ''
  editForm.car = user.car || ''
  editForm.phone = user.phone || ''
  editForm.plate = user.plate || ''
  editForm.balance = user.balance ?? 0
  editForm.roles = user.roles || '["ROLE_PASSAGER"]'
  showModal.value = true
  error.value = ''
}

function closeModal() {
  showModal.value = false
  editingUser.value = null
}

async function saveEdit() {
  error.value = ''
  let rolesJson = editForm.roles
  try { JSON.parse(rolesJson) } catch { error.value = 'Rôles doit être un JSON valide'; return }
  const googleId = editingUser.value.google_id
  const res = await fetch(`http://localhost:8002/users/${googleId}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      name: editForm.name,
      email: editForm.email,
      car: editForm.car,
      plate: editForm.plate,
      balance: editForm.balance,
      roles: editForm.roles,
      phone: editForm.phone
    })
  })
  if (res.ok) {
    showModal.value = false
    emit('refresh')
  } else {
    error.value = 'Erreur lors de la sauvegarde'
  }
}
</script>