<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 via-gray-950 to-violet-900">
    <div class="max-w-xl w-full p-8 bg-gray-900 bg-opacity-90 rounded-2xl shadow-2xl border border-gray-800 text-center">
      <h1 class="text-3xl font-bold mb-3 text-white drop-shadow">Bienvenue sur EcoDrive !</h1>
      <p v-if="isLoggedIn" class="text-green-400 text-lg font-medium">Vous êtes connecté 🎉</p>
      <p v-else class="text-gray-300 text-lg">Pour profiter de toutes les fonctionnalités, connectez-vous ou créez un compte.</p>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const isLoggedIn = computed(() => !!localStorage.getItem('token'))
</script>