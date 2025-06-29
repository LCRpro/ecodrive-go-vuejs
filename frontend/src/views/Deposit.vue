<template>
  <div style="max-width:400px;margin:40px auto;">
    <h2>Dépôt d'argent</h2>
    <div>
      <label>Montant (€)&nbsp;</label>
      <input type="number" min="10" step="0.01" v-model.number="amount" placeholder="Min 10€" />
    </div>
    <div style="margin: 10px 0;">
      <label>Numéro de carte</label>
      <input v-model="card" maxlength="16" placeholder="1234 5678 9012 3456" />
    </div>
    <div>
      <label>Expiration</label>
      <input v-model="exp" maxlength="5" placeholder="MM/YY" style="width:70px;" />
      &nbsp;
      <label>CVC</label>
      <input v-model="cvc" maxlength="4" placeholder="123" style="width:60px;" />
    </div>
    <button @click="submit" style="margin-top:18px;">Déposer</button>
    <div v-if="error" style="color:red; margin-top:10px;">{{ error }}</div>
    <div v-if="success" style="color:green; margin-top:10px;">{{ success }}</div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const amount = ref(10)
const card = ref('')
const exp = ref('')
const cvc = ref('')
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
  if (amount.value < 10) {
    error.value = 'Montant minimum : 10€'
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
    exp.value = ''
    cvc.value = ''
    amount.value = 10
  } else {
    const r = await res.json().catch(()=>'')
    error.value = r?.error || "Erreur lors du dépôt"
  }
}
</script>