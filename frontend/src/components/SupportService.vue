<template>
  <div
    class="min-h-screen flex flex-col items-center bg-gradient-to-br from-gray-900 via-gray-950 to-violet-900 px-4 py-10">
    <div class="w-full max-w-7xl">

      <div class="flex flex-col items-center mb-12 mt-6">
        <span
          class="mb-3 px-4 py-1 rounded-full bg-gradient-to-r from-violet-700 to-emerald-500 text-white font-bold text-xs shadow uppercase tracking-widest">
          Support &amp; assistance
        </span>
        <h2
          class="text-4xl md:text-5xl font-extrabold bg-gradient-to-r from-emerald-400 via-violet-400 to-emerald-400 bg-clip-text text-transparent drop-shadow-lg text-center">
          Centre de support EcoDrive
        </h2>
        <div class="text-gray-400 text-sm mt-2 font-medium">
          Créez vos tickets, suivez vos demandes et échangez avec notre équipe.
        </div>
      </div>

      <div class="flex flex-col items-center mb-10">
        <h3
          class="text-2xl md:text-3xl font-extrabold bg-gradient-to-r from-emerald-400 via-violet-400 to-emerald-400 bg-clip-text text-transparent drop-shadow text-center mb-2">
          Créer une demande de support
        </h3>
        <div class="text-sm text-gray-400 mb-5 text-center max-w-xl">
          Remplissez ce formulaire pour signaler un problème ou poser une question à notre équipe.
        </div>
      </div>

      <form @submit.prevent="submitTicket"
        class="bg-gray-800 p-6 rounded-lg mb-10 border border-violet-700 shadow-inner max-w-xl mx-auto">
        <label class="block text-gray-300 mb-3 font-medium">Catégorie :</label>
        <select v-model="form.category"
          class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-semibold transition-all duration-200 font-mono"
          required>
          <option disabled value="">Choisir une catégorie</option>
          <option value="compte">Compte</option>
          <option value="paiement">Paiement</option>
          <option value="course">Course</option>
        </select>

        <label class="block text-gray-300 mt-6 mb-3 font-medium">Décris ton problème :</label>
        <textarea v-model="form.message" rows="5"
          class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-semibold transition-all duration-200 font-mono"
          required></textarea>

        <button type="submit"
          class="btn mt-6 bg-violet-700 hover:bg-violet-800 text-white font-semibold px-6 py-2 rounded-lg shadow transition">
          Envoyer
        </button>

        <div v-if="formError" class="text-red-500 mt-3 font-medium">{{ formError }}</div>
        <div v-if="formSuccess" class="text-green-400 mt-3 font-medium">Ticket envoyé !</div>
      </form>

      <div class="flex flex-col items-center mb-10 mt-14">
        <h3
          class="text-3xl md:text-4xl font-extrabold bg-gradient-to-r from-emerald-400 via-violet-400 to-emerald-400 bg-clip-text text-transparent drop-shadow-lg text-center mb-2">
          Mes tickets de support
        </h3>
        <div class="text-sm text-gray-400 text-center max-w-xl">
          Retrouvez vos demandes, leur suivi et les réponses de l’équipe EcoDrive.
        </div>
      </div>
      <div class="overflow-x-auto rounded-lg border border-violet-700 shadow-lg mb-12">
        <table class="min-w-[700px] w-full text-sm border border-violet-700 rounded-xl overflow-hidden bg-gray-900">
          <thead>
            <tr>
              <th
                class="bg-gray-800 text-violet-200 font-semibold text-center border border-violet-700 px-3 py-2 md:px-4 md:py-3">
                ID</th>
              <th
                class="bg-gray-800 text-violet-200 font-semibold text-center border border-violet-700 px-3 py-2 md:px-4 md:py-3">
                Catégorie</th>
              <th
                class="bg-gray-800 text-violet-200 font-semibold text-left border border-violet-700 px-3 py-2 md:px-4 md:py-3">
                Message</th>
              <th
                class="bg-gray-800 text-violet-200 font-semibold text-center border border-violet-700 px-3 py-2 md:px-4 md:py-3">
                Statut</th>
              <th
                class="bg-gray-800 text-violet-200 font-semibold text-left border border-violet-700 px-3 py-2 md:px-4 md:py-3">
                Réponse admin</th>
              <th
                class="bg-gray-800 text-violet-200 font-semibold text-center border border-violet-700 px-3 py-2 md:px-4 md:py-3">
                Date</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="t in tickets" :key="t.id" class="even:bg-gray-900 hover:bg-gray-800/80 transition">
              <td class="text-gray-200 border border-violet-700 text-center px-3 py-2 md:px-4 md:py-3">{{ t.id }}</td>
              <td class="text-gray-200 border border-violet-700 text-center capitalize px-3 py-2 md:px-4 md:py-3">{{
                t.category }}</td>
              <td class="text-gray-200 border border-violet-700 text-left px-3 py-2 md:px-4 md:py-3">{{ t.message }}
              </td>
              <td class="text-gray-200 border border-violet-700 text-center capitalize px-3 py-2 md:px-4 md:py-3">{{
                t.status }}</td>
              <td class="text-gray-200 border border-violet-700 text-left px-3 py-2 md:px-4 md:py-3">{{ t.admin_reply ||
                '-' }}</td>
              <td class="text-gray-200 border border-violet-700 text-center px-3 py-2 md:px-4 md:py-3">{{ new
                Date(t.created_at).toLocaleString() }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="isAdmin" class="flex flex-col items-center mb-8 mt-20">
        <h3
          class="text-2xl md:text-3xl font-extrabold bg-gradient-to-r from-violet-400 via-emerald-400 to-violet-400 bg-clip-text text-transparent drop-shadow text-center mb-2">
          Gestion des tickets (admin)
        </h3>
      </div>
      <div v-if="isAdmin" class="overflow-x-auto rounded-lg border border-violet-700 shadow-lg">
        <table class="min-w-[700px] w-full text-sm border border-violet-700 rounded-xl overflow-hidden bg-gray-900">
          <thead>
            <tr>
              <th
                class="bg-gray-800 text-violet-200 font-semibold text-center border border-violet-700 px-3 py-2 md:px-4 md:py-3">
                ID</th>
              <th
                class="bg-gray-800 text-violet-200 font-semibold text-left border border-violet-700 px-3 py-2 md:px-4 md:py-3">
                Utilisateur</th>
              <th
                class="bg-gray-800 text-violet-200 font-semibold text-center border border-violet-700 px-3 py-2 md:px-4 md:py-3">
                Catégorie</th>
              <th
                class="bg-gray-800 text-violet-200 font-semibold text-left border border-violet-700 px-3 py-2 md:px-4 md:py-3">
                Message</th>
              <th
                class="bg-gray-800 text-violet-200 font-semibold text-center border border-violet-700 px-3 py-2 md:px-4 md:py-3">
                Statut</th>
              <th
                class="bg-gray-800 text-violet-200 font-semibold text-left border border-violet-700 px-3 py-2 md:px-4 md:py-3">
                Réponse admin</th>
              <th
                class="bg-gray-800 text-violet-200 font-semibold text-center border border-violet-700 px-3 py-2 md:px-4 md:py-3">
                Date</th>
              <th
                class="bg-gray-800 text-violet-200 font-semibold text-center border border-violet-700 px-3 py-2 md:px-4 md:py-3">
                Action</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="t in allTickets" :key="t.id" class="even:bg-gray-900 hover:bg-gray-800/80 transition">
              <td class="text-gray-200 border border-violet-700 text-center px-3 py-2 md:px-4 md:py-3">{{ t.id }}</td>
              <td class="text-gray-200 border border-violet-700 text-left px-3 py-2 md:px-4 md:py-3">
                {{ userNames[t.google_id] || t.google_id }}
              </td>
              <td class="text-gray-200 border border-violet-700 text-center capitalize px-3 py-2 md:px-4 md:py-3">{{
                t.category }}</td>
              <td class="text-gray-200 border border-violet-700 text-left px-3 py-2 md:px-4 md:py-3">{{ t.message }}
              </td>
              <td class="text-gray-200 border border-violet-700 text-center capitalize px-3 py-2 md:px-4 md:py-3">{{
                t.status }}</td>
              <td class="text-gray-200 border border-violet-700 text-left px-3 py-2 md:px-4 md:py-3">{{ t.admin_reply ||
                '-' }}</td>
              <td class="text-gray-200 border border-violet-700 text-center px-3 py-2 md:px-4 md:py-3">{{ new
                Date(t.created_at).toLocaleString() }}</td>
              <td class="text-gray-200 border border-violet-700 text-center px-3 py-2 md:px-4 md:py-3">
                <button v-if="!t.admin_reply" @click="replyTicket(t)"
                  class="btn bg-violet-700 hover:bg-violet-800 text-white font-semibold px-4 py-1 rounded shadow transition">
                  Répondre
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="modalTicket" class="fixed inset-0 bg-black/80 flex items-center justify-center z-50">
        <div class="bg-gray-900 p-6 rounded-lg border border-violet-700 shadow-lg max-w-md w-full">
          <h4 class="mb-4 text-lg font-semibold text-violet-300">Répondre au ticket #{{ modalTicket.id }}</h4>
          <textarea v-model="replyText" rows="5" class="input w-full mb-4 resize-y"
            placeholder="Tapez votre réponse ici..."></textarea>
          <div class="flex justify-end gap-3">
            <button @click="submitReply"
              class="btn bg-violet-700 hover:bg-violet-800 text-white font-semibold px-5 py-2 rounded shadow transition">
              Envoyer
            </button>
            <button @click="modalTicket = null"
              class="btn bg-gray-700 hover:bg-gray-600 text-white font-semibold px-5 py-2 rounded shadow transition">
              Annuler
            </button>
          </div>
          <div v-if="replyError" class="text-red-500 mt-3 font-medium">{{ replyError }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'


const supportServiceURL = import.meta.env.VITE_SUPPORT_SERVICE_URL
const userServiceURL = import.meta.env.VITE_USER_SERVICE_URL


const form = reactive({
  category: '',
  message: ''
})
const formError = ref('')
const formSuccess = ref(false)

const tickets = ref([])
const allTickets = ref([])
const modalTicket = ref(null)
const replyText = ref('')
const replyError = ref('')

let google_id = ''
const isAdmin = ref(false)

onMounted(() => {
  google_id = localStorage.getItem('google_id') || ''
  const roles = JSON.parse(localStorage.getItem('roles') || '[]')
  isAdmin.value = roles.includes('ROLE_ADMIN')
  fetchTickets()
  if (isAdmin.value) fetchAllTickets()
})

async function submitTicket() {
  formError.value = ''
  formSuccess.value = false
  if (!form.category || !form.message) {
    formError.value = 'Tous les champs sont obligatoires.'
    return
  }
  const res = await fetch(`${supportServiceURL}/support`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ google_id, category: form.category, message: form.message })
  })
  if (res.ok) {
    formSuccess.value = true
    form.category = ''
    form.message = ''
    fetchTickets()
    if (isAdmin.value) fetchAllTickets()
  } else {
    formError.value = 'Erreur lors de l\'envoi.'
  }
}

async function fetchTickets() {
  if (!google_id) return
  const res = await fetch(`${supportServiceURL}/support/user/` + google_id)
  tickets.value = res.ok ? await res.json() : []
}

async function fetchAllTickets() {
  const res = await fetch(`${supportServiceURL}/support/all`)
  allTickets.value = res.ok ? await res.json() : []
  if (allTickets.value.length) {
    await resolveUserNames(allTickets.value)
  }
}

function replyTicket(t) {
  modalTicket.value = t
  replyText.value = ''
  replyError.value = ''
}

async function submitReply() {
  if (!replyText.value) {
    replyError.value = 'Réponse obligatoire'
    return
  }
  const res = await fetch(`${supportServiceURL}/support/reply/${modalTicket.value.id}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ admin_reply: replyText.value, status: 'fermé' })
  })
  if (res.ok) {
    modalTicket.value = null
    fetchAllTickets()
  } else {
    replyError.value = 'Erreur lors de la réponse'
  }
}

const userNames = ref({})
async function resolveUserNames(tickets) {
  const ids = [...new Set(tickets.map(t => t.google_id).filter(Boolean))]
  for (const id of ids) {
    if (!userNames.value[id]) {
      const res = await fetch(`${userServiceURL}/users/${id}`)
      if (res.ok) {
        const user = await res.json()
        userNames.value[id] = user.name
      } else {
        userNames.value[id] = id
      }
    }
  }
}
</script>

<style scoped>
.input {
  width: 100%;
  border-radius: 0.375rem;
  background: #1f2937;
  /* slate-800 */
  color: white;
  padding: 0.5rem 0.75rem;
  border: 1px solid #7c3aed;
  /* violet-600 */
  outline: none;
  transition: border-color 0.2s;
  font-size: 0.9rem;
}

.input:focus {
  border-color: #a78bfa;
  /* violet-400 */
  box-shadow: 0 0 5px #a78bfa;
}

.btn {
  font-weight: 600;
  transition: background-color 0.2s ease-in-out;
  cursor: pointer;
  user-select: none;
}
</style>