<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 via-gray-950 to-violet-900 px-2 py-10">
    <div class="w-full max-w-sm mx-auto p-7 rounded-2xl shadow-2xl bg-gray-900 bg-opacity-90 border border-gray-800">
      <div class="flex items-center gap-3 mb-6">
        <svg class="w-10 h-10 text-violet-400" fill="none" viewBox="0 0 24 24"><rect width="20" height="14" x="2" y="5" rx="3" fill="#fff"/><rect width="20" height="14" x="2" y="5" rx="3" stroke="#7c3aed" stroke-width="1.5"/><rect width="6" height="2" x="4" y="13" fill="#d1d5db"/><rect width="4" height="2" x="14" y="13" fill="#a78bfa"/></svg>
        <h2 class="text-2xl font-bold text-white tracking-tight">Dépôt d'argent</h2>
      </div>
      <form @submit.prevent="submit" autocomplete="off" class="space-y-5">
        <div>
          <label class="block text-gray-300 mb-1 font-medium">Montant (€)</label>
          <input
            type="number"
            min="10"
            step="1"
            v-model.number="amount"
            placeholder="Min 10€"
            class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-semibold transition-all duration-200 font-mono"
          />
        </div>
        <div>
          <label class="block text-gray-300 mb-1 font-medium">Numéro de carte</label>
          <input
            v-model="cardDisplay"
            maxlength="19"
            placeholder="1234 5678 9012 3456"
            inputmode="numeric"
            autocomplete="cc-number"
            @input="formatCard"
            class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-semibold transition-all duration-200 font-mono tracking-widest"
          />
        </div>
        <div class="flex gap-4">
          <div class="flex-1">
            <label class="block text-gray-300 mb-1 font-medium">Expiration</label>
            <input
              v-model="exp"
              type="month"
              placeholder="MM/AA"
              autocomplete="cc-exp"
              class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-semibold transition-all duration-200 font-mono"
            />
          </div>
          <div class="flex-1">
            <label class="block text-gray-300 mb-1 font-medium">CVC</label>
            <input
              v-model="cvc"
              maxlength="3"
              placeholder="123"
              autocomplete="cc-csc"
              class="w-full px-4 py-3 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow-md font-semibold transition-all duration-200 font-mono"
            />
          </div>
        </div>
        <button type="submit"
          class="w-full py-3 rounded-lg bg-emerald-600 hover:bg-emerald-700 text-white text-lg font-bold shadow transition">
          Déposer
        </button>
        <div v-if="error" class="bg-red-700/80 text-white text-center py-2 rounded-lg font-semibold">{{ error }}</div>
        <div v-if="success" class="bg-emerald-700/80 text-white text-center py-2 rounded-lg font-semibold">{{ success }}</div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const amount = ref(10)
const card = ref('')        
const cardDisplay = ref('') 
const exp = ref('')
const cvc = ref('')
const error = ref('')
const success = ref('')

function formatCard(e) {
  let value = e.target.value.replace(/\D/g, '').slice(0, 16)
  card.value = value
  cardDisplay.value = value.replace(/(.{4})/g, '$1 ').trim()
}

function getGoogleIdFromToken(token) {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.sub
  } catch {
    return null
  }
}

async function submit() {
  error.value = ''
  success.value = ''
  if (amount.value < 10) {
    error.value = 'Montant minimum : 10€'
    return
  }
  if (card.value.length !== 16) {
    error.value = "Numéro de carte incomplet (16 chiffres attendus)"
    return
  }
  if (!exp.value) {
    error.value = "Date d'expiration requise"
    return
  }
  if (!/^\d{3}$/.test(cvc.value)) {
    error.value = "CVC incomplet (3 chiffres)"
    return
  }
  const token = localStorage.getItem('token')
  const googleId = getGoogleIdFromToken(token)
  if (!googleId) {
    error.value = "Utilisateur inconnu"
    return
  }
  const res = await fetch('http://localhost:8004/deposit', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Authorization: 'Bearer ' + token },
    body: JSON.stringify({
      google_id: googleId,
      amount: amount.value,
      card_number: card.value,
      expiry: exp.value,
      cvv: cvc.value
    })
  })
  if (res.ok) {
    success.value = 'Dépôt effectué !'
    card.value = ''
    cardDisplay.value = ''
    exp.value = ''
    cvc.value = ''
    amount.value = 10
    setTimeout(() => window.location.reload(), 1000)
  } else {
    const r = await res.json().catch(()=>'')
    error.value = r?.error || "Erreur lors du dépôt"
  }
}
</script>