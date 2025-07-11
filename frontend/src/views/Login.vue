<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 via-gray-950 to-violet-900">
    <div class="w-full max-w-sm p-8 bg-gray-900/90 rounded-2xl shadow-2xl border border-gray-800 relative overflow-hidden">
      <div class="absolute inset-0 bg-gradient-to-tr from-emerald-700/20 via-violet-700/10 to-yellow-500/10 blur-xl pointer-events-none" />
      <h2 class="text-3xl font-extrabold text-white mb-2 text-center drop-shadow bg-gradient-to-r from-violet-400 to-emerald-400 bg-clip-text text-transparent">
        Connexion
      </h2>
      <p class="text-gray-400 mb-5 text-center">Accédez à votre espace sécurisé EcoDrive</p>
      <div class="relative mb-8">
        <div class="bg-gradient-to-br from-emerald-700/40 via-violet-700/40 to-white/10 backdrop-blur-lg text-emerald-200 text-[15px] rounded-xl p-4 pt-6 border border-violet-700 shadow-lg text-center font-medium">
          Connectez-vous pour commander vos trajets et gérer votre solde.<br>
          <span class="block mt-2 text-violet-200 font-semibold">Conducteur&nbsp;?</span>
          <span class="text-white/90">Prenez des courses, touchez vos gains et accédez à vos statistiques.</span>
        </div>
      </div>
      <button
        @click="googleLogin"
        class="w-full flex items-center justify-center gap-2 px-4 py-3 rounded-xl bg-gradient-to-r from-emerald-400 to-violet-500 text-white font-bold shadow-lg hover:scale-[1.025] hover:shadow-emerald-600/40 active:scale-100 transition"
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
const authServiceURL = import.meta.env.VITE_AUTH_SERVICE_URL

function googleLogin() {
  window.location.href = authServiceURL + '/auth/google'
}
</script>