<!-- src/components/BalanceDisplay.vue -->
<template>
  <span v-if="balance !== null" style="margin-left:20px;">
    <strong>Solde :</strong> {{ balance.toFixed(2) }} €
  </span>
</template>
<script setup>
import { ref, onMounted, watch } from 'vue'
const balance = ref(null)
const googleId = (() => {
  try {
    return JSON.parse(atob(localStorage.getItem('token').split('.')[1])).sub
  } catch { return null }
})()
async function fetchBalance() {
  if (!googleId) return
  const res = await fetch(`http://localhost:8002/users/${googleId}`)
  if (res.ok) {
    const user = await res.json()
    balance.value = user.balance ?? 0
  }
}
onMounted(fetchBalance)
window.addEventListener('focus', fetchBalance) // refresh si retour focus
</script>