<template>
  <div style="margin-top:32px;">
    <h3>Historique de mes transactions</h3>
    <div v-if="loading">Chargement...</div>
    <table v-else border="1" cellpadding="5" style="width:100%; margin-top:10px;">
      <thead>
        <tr>
          <th>Date</th>
          <th>Type</th>
          <th>Montant (€)</th>
          <th>IBAN (retrait)</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="t in transactions" :key="t.id">
          <td>{{ formatDate(t.created_at) }}</td>
          <td>
            <span v-if="t.type === 'deposit'">Dépôt</span>
            <span v-else-if="t.type === 'withdrawal'">Retrait</span>
            <span v-else>{{ t.type }}</span>
          </td>
          <td :style="{color: t.type === 'deposit' ? 'green' : 'red'}">
            {{ t.amount.toFixed(2) }}
          </td>
          <td>{{ t.iban_mask || '-' }}</td>
        </tr>
      </tbody>
    </table>
    <div v-if="!loading && transactions.length === 0" style="margin-top:10px;">Aucune transaction trouvée.</div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

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

  const res = await fetch(`http://localhost:8004/transactions/${googleId}`, {
    headers: { Authorization: 'Bearer ' + token }
  })
  if (res.ok) transactions.value = await res.json()
  loading.value = false
})
</script>