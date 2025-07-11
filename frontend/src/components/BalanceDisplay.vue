<template>
  <span
    v-if="balance !== null"
    class="inline-flex items-center gap-2 px-3 py-1 rounded-lg bg-gradient-to-r from-emerald-700/80 to-violet-800/80 border border-violet-600 text-white font-semibold shadow transition duration-200"
    style="margin-left: 20px; font-size: 1.07rem;"
  >
  
 <svg class="w-5 h-5 text-amber-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
  <circle cx="12" cy="12" r="10" class="stroke-amber-300" fill="#FFD700"/>
  <text x="12" y="16" text-anchor="middle" font-size="10" fill="#fff" font-weight="bold">€</text>
</svg>
    <span>Solde&nbsp;:</span>
    <span class="font-bold text-amber-300">{{ balance.toFixed(2) }} €</span>
  </span>
</template>
<script setup>
import { ref, onMounted, watch } from 'vue'
const userServiceURL = import.meta.env.VITE_USER_SERVICE_URL
const balance = ref(null)
const googleId = (() => {
  try {
    return JSON.parse(atob(localStorage.getItem('token').split('.')[1])).sub
  } catch { return null }
})()
async function fetchBalance() {
  if (!googleId) return
  const res = await fetch(`${userServiceURL}/users/${googleId}`)
  if (res.ok) {
    const user = await res.json()
    balance.value = user.balance ?? 0
  }
}
onMounted(fetchBalance)
window.addEventListener('focus', fetchBalance) 
</script>