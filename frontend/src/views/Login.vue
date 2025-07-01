<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 via-gray-950 to-violet-900">
    <div class="w-full max-w-sm p-8 bg-gray-900 bg-opacity-90 rounded-2xl shadow-2xl border border-gray-800">
      <h2 class="text-3xl font-bold text-white mb-2 text-center">Connexion</h2>
      <p class="text-gray-400 mb-6 text-center">Accédez à votre espace sécurisé EcoDrive</p>
      
      <button
        @click="googleLogin"
        class="w-full flex items-center justify-center gap-2 px-4 py-3 rounded-xl bg-white text-gray-900 font-semibold shadow hover:bg-gray-100 transition"
      >
        <svg class="w-6 h-6" viewBox="0 0 24 24">
          <g>
            <path fill="#4285F4" d="M21.805 10.023h-9.6v3.954h5.55a4.747 4.747 0 01-2.056 3.115l-.018.12 2.99 2.319.207.021c1.899-1.753 2.994-4.339 2.994-7.235z"/>
            <path fill="#34A853" d="M12.204 21c2.7 0 4.966-.896 6.621-2.428l-3.157-2.447c-.879.619-2.011.99-3.464.99-2.656 0-4.907-1.795-5.713-4.197l-.117.01-3.09 2.412-.04.113C4.26 19.146 7.027 21 12.204 21z"/>
            <path fill="#FBBC05" d="M6.491 13.117A5.845 5.845 0 016.203 12c0-.392.068-.771.117-1.117l-.001-.115-3.126-2.426-.103.048A8.765 8.765 0 003.204 12c0 1.433.348 2.792.97 3.983l3.126-2.427z"/>
            <path fill="#EA4335" d="M12.204 6.955c1.52 0 2.545.646 3.131 1.187l2.287-2.248C16.964 4.2 14.905 3 12.204 3 7.027 3 4.26 6.854 3.204 8.957l3.287 2.56c.447-1.349 1.718-4.562 5.713-4.562z"/>
          </g>
        </svg>
        Se connecter avec Google
      </button>
      
      <div v-if="error" class="mt-6 text-red-500 text-center">{{ error }}</div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const error = ref('')

function googleLogin() {
  window.location.href = 'http://localhost:8000/auth/google'
}
</script>