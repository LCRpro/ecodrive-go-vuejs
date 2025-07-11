<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 via-gray-950 to-violet-900 px-2 py-10">
    <div class="w-full max-w-sm mx-auto p-7 rounded-2xl shadow-2xl bg-gray-900 bg-opacity-90 border border-gray-800">
      <div class="flex items-center gap-3 mb-6">
        <svg class="w-10 h-10 text-violet-400" fill="none" viewBox="0 0 24 24">
          <rect width="20" height="14" x="2" y="5" rx="3" fill="#fff"/>
          <rect width="20" height="14" x="2" y="5" rx="3" stroke="#7c3aed" stroke-width="1.5"/>
          <rect width="6" height="2" x="4" y="13" fill="#d1d5db"/>
          <rect width="4" height="2" x="14" y="13" fill="#a78bfa"/>
        </svg>
        <h2 class="text-2xl font-bold font-extrabold bg-gradient-to-r from-emerald-400 via-violet-400 to-emerald-400 bg-clip-text text-transparent drop-shadow-lg text-center">Retrait</h2>
      </div>
      <form @submit.prevent="submit" autocomplete="off" class="space-y-5">
        <div>
          <label class="block text-gray-300 mb-1 font-medium">Montant (€)</label>
          <input
            type="number"
            min="1"
            step="1"
            v-model.number="amount"
            placeholder="Montant"
            class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-semibold transition-all duration-200 font-mono"
          />
        </div>
        <div>
          <label class="block text-gray-300 mb-1 font-medium">IBAN</label>
          <input
            v-model="iban"
            maxlength="34"
            placeholder="FR..."
            autocomplete="off"
            class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-semibold transition-all duration-200 font-mono"
          />
        </div>
        <button type="submit"
          class="w-full py-3 rounded-lg bg-violet-600 hover:bg-violet-700 text-white text-lg font-bold shadow transition">
          Retirer
        </button>
        <div v-if="error" class="bg-red-700/80 text-white text-center py-2 rounded-lg font-semibold">{{ error }}</div>
        <div v-if="success" class="bg-emerald-700/80 text-white text-center py-2 rounded-lg font-semibold">{{ success }}</div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const paiementServiceURL = import.meta.env.VITE_PAIEMENT_SERVICE_URL

const amount = ref(10)
const iban = ref('')
const error = ref('')
const success = ref('')

function getGoogleIdFromToken(token) {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.sub
  } catch {
    return null
  }
}

function isValidIban(iban) {
  return /^FR\d{2}[0-9A-Z]{11,27}$/.test(iban.replace(/\s+/g, '').toUpperCase())
}

async function submit() {
  error.value = ''
  success.value = ''
  if (!amount.value || amount.value <= 0) {
    error.value = "Le montant doit être supérieur à 0"
    return
  }
  if (!isValidIban(iban.value)) {
    error.value = "IBAN invalide (exemple : FR7612345678901234567890123)"
    return
  }
  const token = localStorage.getItem('token')
  const googleId = getGoogleIdFromToken(token)
  if (!googleId) {
    error.value = "Utilisateur inconnu"
    return
  }
  const res = await fetch(paiementServiceURL + '/withdraw', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Authorization: 'Bearer ' + token },
    body: JSON.stringify({
      google_id: googleId,
      amount: amount.value,
      iban: iban.value.replace(/\s+/g, '').toUpperCase()
    })
  })
  if (res.ok) {
    success.value = 'Retrait effectué !'
    iban.value = ''
    amount.value = 10
    setTimeout(() => window.location.reload(), 1000)
  } else {
    const r = await res.json().catch(()=>'')
    error.value = r?.error || "Erreur lors du retrait"
  }
}
</script>