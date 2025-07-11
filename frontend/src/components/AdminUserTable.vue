<template>
  <div class="mb-10">
<div class="flex flex-col items-center mb-10">
  <h3
    class="text-3xl md:text-4xl font-extrabold bg-gradient-to-r from-emerald-400 via-violet-400 to-emerald-400 bg-clip-text text-transparent drop-shadow-lg text-center"
  >
    Utilisateurs
  </h3>
</div>    <div class="overflow-x-auto w-full">
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
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ u.google_id }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ u.email }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ u.phone }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ u.name }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">
              <span
                v-for="role in JSON.parse(u.roles || '[]')"
                :key="role"
                :class="roleColorClass(role)"
                class="inline-block px-2 py-1 mx-1 rounded-full text-xs font-semibold capitalize"
              >
                {{ shortRole(role) }}
              </span>
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

<div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center  from-gray-900 via-gray-950 to-violet-900 bg-opacity-90 px-2 py-10">
  <div class="w-full max-w-2xl mx-auto p-8 rounded-2xl shadow-2xl bg-gray-900 bg-opacity-95 border border-violet-800">
    <div class="flex items-center gap-3 mb-6">
      <svg class="w-10 h-10 text-violet-400" fill="none" viewBox="0 0 24 24">
        <rect width="20" height="14" x="2" y="5" rx="3" fill="#fff"/>
        <rect width="20" height="14" x="2" y="5" rx="3" stroke="#7c3aed" stroke-width="1.5"/>
        <rect width="6" height="2" x="4" y="13" fill="#d1d5db"/>
        <rect width="4" height="2" x="14" y="13" fill="#a78bfa"/>
      </svg>
      <h2 class="text-2xl font-bold text-white tracking-tight">Modifier utilisateur #{{ editingUser?.id }}</h2>
    </div>
    <form @submit.prevent="saveEdit" autocomplete="off" class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div>
        <label class="block text-gray-300 mb-1 font-medium">Nom</label>
        <input v-model="editForm.name" class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400
          focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-semibold transition-all duration-200 font-mono" />
      </div>
      <div>
        <label class="block text-gray-300 mb-1 font-medium">Email</label>
        <input v-model="editForm.email" class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400
          focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-semibold transition-all duration-200 font-mono" />
      </div>
      <div>
        <label class="block text-gray-300 mb-1 font-medium">Numéro de téléphone</label>
        <input v-model="editForm.phone" class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400
          focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-semibold transition-all duration-200 font-mono" />
      </div>
      <div>
        <label class="block text-gray-300 mb-1 font-medium">Voiture</label>
        <input v-model="editForm.car" class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400
          focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-semibold transition-all duration-200 font-mono" />
      </div>
      <div>
        <label class="block text-gray-300 mb-1 font-medium">Plaque</label>
        <input v-model="editForm.plate" class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400
          focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-semibold transition-all duration-200 font-mono" />
      </div>
      <div>
        <label class="block text-gray-300 mb-1 font-medium">Solde (€)</label>
        <input type="number" v-model.number="editForm.balance" min="0" step="0.01"
          class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400
          focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-semibold transition-all duration-200 font-mono" />
      </div>
    <div class="md:col-span-2">
  <label class="block text-gray-300 mb-1 font-medium mb-2">Rôles</label>
  <div class="flex flex-wrap gap-2 mb-2">
    <span
      v-for="role in possibleRoles"
      :key="role.value"
      :class="[
        isRoleSelected(role.value)
          ? roleColorClass(role.value) + ' ring-2 ring-violet-500'
          : 'bg-gray-700 text-gray-300 hover:bg-violet-800 hover:text-white',
        'px-3 py-1 rounded-full cursor-pointer text-xs font-bold transition'
      ]"
      @click="toggleRole(role.value)"
    >
      {{ role.label }}
      <span v-if="isRoleSelected(role.value)" class="ml-1">×</span>
    </span>
  </div>
</div>
      <div class="md:col-span-2 flex justify-end gap-3 mt-2">
        <button type="submit"
                class="px-5 py-2 rounded-lg bg-violet-700 hover:bg-violet-800 text-white font-bold shadow transition">Enregistrer</button>
        <button type="button"
                @click="closeModal"
                class="px-5 py-2 rounded-lg bg-gray-700 hover:bg-gray-800 text-white font-semibold transition">Annuler</button>
      </div>
      <div v-if="error" class="md:col-span-2 bg-red-700/80 text-white text-center py-2 rounded-lg font-semibold mt-2">{{ error }}</div>
    </form>
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

function shortRole(role) {
  return (role.replace(/^ROLE_/, '').toLowerCase())
}

function roleColorClass(role) {
  switch (role) {
    case 'ROLE_ADMIN':
      return 'bg-red-700 text-white'
    case 'ROLE_DRIVER':
      return 'bg-green-700 text-white'
    case 'ROLE_PASSAGER':
      return 'bg-blue-700 text-white'
    default:
      return 'bg-gray-600 text-white'
  }
}


const possibleRoles = [
  { value: "ROLE_ADMIN", label: "Admin" },
  { value: "ROLE_DRIVER", label: "Driver" },
  { value: "ROLE_PASSAGER", label: "Passager" }
]

const rolesArray = ref([])

function isRoleSelected(role) {
  return rolesArray.value.includes(role)
}

function toggleRole(role) {
  if (isRoleSelected(role)) {
    rolesArray.value = rolesArray.value.filter(r => r !== role)
  } else {
    rolesArray.value = [...rolesArray.value, role]
  }
}

watch(
  () => editForm.roles,
  (val) => {
    try {
      rolesArray.value = JSON.parse(val || '[]')
    } catch {
      rolesArray.value = []
    }
  },
  { immediate: true }
)
watch(
  rolesArray,
  (val) => {
    editForm.roles = JSON.stringify(val)
  },
  { deep: true }
)
</script>
