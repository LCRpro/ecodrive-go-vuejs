<template>
  <div style="max-width:350px;margin:80px auto;">
    <h2>Connexion</h2>
    <button @click="googleLogin">Se connecter avec Google</button>
    <div v-if="error" style="color:red;">{{ error }}</div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const error = ref('')

function googleLogin() {
  window.location.href = 'http://localhost:8000/auth/google'
}
</script>