<template>
  <div style="max-width:400px;margin:40px auto;">
    <h2>Retrait</h2>
    <div>
      <label>Montant (€)&nbsp;</label>
      <input type="number" step="0.01" v-model.number="amount" placeholder="Montant" />
    </div>
    <div style="margin: 10px 0;">
      <label>IBAN</label>
      <input v-model="iban" maxlength="34" placeholder="FR..." />
    </div>
    <button @click="submit" style="margin-top:18px;">Retirer</button>
    <div v-if="error" style="color:red; margin-top:10px;">{{ error }}</div>
    <div v-if="success" style="color:green; margin-top:10px;">{{ success }}</div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
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

async function submit() {
  error.value = ''
  success.value = ''
  if (!iban.value.match(/^[A-Z]{2}\d{2}[A-Z0-9]{11,30}$/)) {
    error.value = "IBAN invalide"
    return
  }
  const token = localStorage.getItem('token')
  const googleId = getGoogleIdFromToken(token)
  if (!googleId) {
    error.value = "Utilisateur inconnu"
    return
  }
  const res = await fetch('http://localhost:8004/withdraw', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Authorization: 'Bearer ' + token },
    body: JSON.stringify({
      google_id: googleId,   
      amount: amount.value,
      iban: iban.value
    })
  })
  if (res.ok) {
    success.value = 'Retrait effectué !'
    iban.value = ''
    amount.value = 0
  } else {
    const r = await res.json().catch(()=>'')
    error.value = r?.error || "Erreur lors du retrait"
  }
}
</script>