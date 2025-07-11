<template>
  <div class="bg-gray-900 bg-opacity-90 rounded-2xl shadow-lg border border-gray-800 p-6 mb-4 mt-8">
    <div class="flex mb-6">
      <h2
        class="text-2xl font-bold font-extrabold bg-gradient-to-r from-emerald-400 via-violet-400 to-emerald-400 bg-clip-text text-transparent drop-shadow-lg text-center">
        Historique de mes transactions</h2>
    </div>
    <div v-if="loading" class="text-gray-400 mb-4">Chargement...</div>
    <div v-else>
      <div class="overflow-x-auto w-full">
        <table class="min-w-[520px] w-full text-sm border border-gray-700 rounded-xl overflow-hidden bg-gray-900">
          <thead class="bg-gray-800 text-violet-200 font-semibold text-center">
            <tr>
              <th class="border border-gray-700 px-3 py-2 md:px-4 md:py-3">Date</th>
              <th class="border border-gray-700 px-3 py-2 md:px-4 md:py-3">Type</th>
              <th class="border border-gray-700 px-3 py-2 md:px-4 md:py-3 text-center">Montant (€)</th>
              <th class="border border-gray-700 px-3 py-2 md:px-4 md:py-3 text-center">IBAN (retrait)</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="t in transactions" :key="t.id" class="even:bg-gray-900 hover:bg-gray-800/80">
              <td class="border border-gray-700 px-3 py-2 md:px-4 md:py-3 text-gray-200">{{ formatDate(t.created_at) }}
              </td>
              <td class="border border-gray-700 px-3 py-2 md:px-4 md:py-3 text-center text-gray-200">
                <span v-if="t.type === 'deposit'"
                  class="px-2 py-1 rounded-full bg-emerald-600 text-white text-xs font-semibold">Dépôt</span>
                <span v-else-if="t.type === 'withdrawal'"
                  class="px-2 py-1 rounded-full bg-rose-600 text-white text-xs font-semibold">Retrait</span>
                <span v-else class="text-gray-200">{{ t.type }}</span>
              </td>
              <td class="border border-gray-700 px-3 py-2 md:px-4 md:py-3 font-mono text-center text-gray-200"
                :class="t.type === 'deposit' ? 'text-emerald-400' : 'text-rose-400'">
                {{ t.amount.toFixed(2) }}
              </td>
              <td class="border border-gray-700 px-3 py-2 md:px-4 md:py-3 text-center text-gray-200">{{ t.iban_mask ||
                '-' }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-if="transactions.length === 0" class="mt-6 text-gray-400 text-center">Aucune transaction trouvée.</div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
const paiementServiceURL = import.meta.env.VITE_PAIEMENT_SERVICE_URL
const transactions = ref([])
const loading = ref(true)

function getGoogleIdFromToken(token) {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.sub
  } catch {
    return null
  }
}
function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleString('fr-FR')
}

onMounted(async () => {
  const token = localStorage.getItem('token')
  const googleId = getGoogleIdFromToken(token)
  if (!token || !googleId) return

  const res = await fetch(`${paiementServiceURL}/transactions/${googleId}`, {
    headers: { Authorization: 'Bearer ' + token }
  })
  if (res.ok) transactions.value = await res.json()
  loading.value = false
})
</script>