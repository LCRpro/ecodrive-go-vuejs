<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-900 via-gray-950 to-violet-900 pb-16">
    <header class="w-full flex items-center justify-center pt-12 pb-20 px-4 relative">
      <div class="absolute inset-0 flex items-center justify-center pointer-events-none select-none">
        <div class="w-72 h-72 md:w-96 md:h-96 bg-violet-700 opacity-30 rounded-full blur-3xl animate-pulse"></div>
      </div>
      <div class="relative z-10 flex flex-col items-center text-center">
        <div class="relative group">
          <img
            src="https://images.unsplash.com/photo-1596032457104-baa1f5f9372a?q=80&w=654&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
            alt="Ecodrive voiture"
            class="w-40 h-40 md:w-48 md:h-48 object-cover rounded-full border-4 border-violet-500 shadow-2xl bg-gray-900 group-hover:scale-105 group-hover:shadow-violet-800 transition-transform"
          />
          <div class="absolute bottom-0 right-2 bg-emerald-600 text-white rounded-xl px-3 py-1 text-xs font-bold shadow -mb-2">Eco</div>
        </div>
        <h1 class="text-5xl md:text-6xl font-extrabold mt-7 mb-4 text-white tracking-tight drop-shadow-[0_2px_12px_rgba(115,85,220,0.6)]">
          Bienvenue sur <span class="text-violet-400">EcoDrive</span> !
        </h1>
        <p class="text-lg md:text-2xl text-gray-200 mb-7 max-w-2xl font-medium">
          La mobilité locale, éthique et responsable. <br class="hidden md:inline">Chaque course, un pas vers une ville plus verte.
        </p>
        <router-link
          to="/course"
          class="inline-block px-8 py-3 rounded-xl bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-lg shadow-xl transition"
        >
          Commander une course
        </router-link>
        <div class="text-sm text-gray-400 mt-5">
Adoptez EcoDrive en toute liberté
        </div>
      </div>
    </header>


<section class="max-w-6xl mx-auto my-16 flex flex-col md:flex-row items-center md:items-stretch gap-8 px-4">
  <div class="flex-1 flex flex-col justify-center">
    <div class="mb-3">
      <span class="px-4 py-1 rounded-full bg-gradient-to-r from-violet-700 to-emerald-500 text-white font-bold text-xs shadow uppercase tracking-widest inline-block">Nouveau&nbsp;!</span>
    </div>
    <h2 class="text-4xl font-extrabold mb-7 bg-gradient-to-r from-emerald-400 via-violet-400 to-emerald-400 bg-clip-text text-transparent drop-shadow text-left">
      <span class="text-white">Qu’est-ce que </span>EcoDrive<span class="text-white"> ?</span>
    </h2>
    <div class="rounded-2xl bg-gradient-to-br from-gray-800/70 via-gray-900/80 to-violet-900/70 shadow-2xl p-8 border border-violet-900 mb-2 transition hover:shadow-violet-700/40 flex flex-col gap-6">
      <p class="text-lg text-gray-200 leading-relaxed font-medium">
        EcoDrive est la première plateforme française de transport local, durable et solidaire.<br />
        Ici, passagers et conducteurs se rencontrent pour partager des trajets moins chers, réduire leur impact environnemental, et booster la vie locale.
      </p>
      <ul class="flex flex-col gap-4">
        <li class="flex items-center gap-3 text-emerald-300 font-semibold text-lg">
          <svg class="w-6 h-6 flex-shrink-0" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <circle cx="12" cy="12" r="10" stroke-width="2.5"/>
            <path d="M8.5 12.5l2.5 2.5 4.5-4.5" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          Trajets avec des conducteurs vérifiés
        </li>
        <li class="flex items-center gap-3 text-violet-300 font-semibold text-lg">
          <svg class="w-6 h-6 flex-shrink-0" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <rect x="4" y="4" width="16" height="16" rx="4" />
            <path d="M8 12h8" stroke-linecap="round"/>
          </svg>
          Tarifs transparents &amp; sans frais cachés
        </li>
        <li class="flex items-center gap-3 text-green-300 font-semibold text-lg">
          <svg class="w-6 h-6 flex-shrink-0" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path d="M12 21C7 21 2 16.5 2 12V7l10-5 10 5v5c0 4.5-5 9-10 9z"/>
            <circle cx="12" cy="12" r="3"/>
          </svg>
          Contribution directe à la réduction du CO₂
        </li>
      </ul>
    </div>
  </div>
  <div class="w-full  md:w-[410px] flex-shrink-0 flex flex-col items-center md:items-end justify-end mb-3">
    <div class="relative w-full h-64 md:h-80">
      <img
        src="https://plus.unsplash.com/premium_photo-1712496473320-af3aa50cd5e5?q=80&w=1770&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
        alt="Présentation Ecodrive"
        class="w-full h-full object-cover rounded-2xl shadow-2xl border-4 border-emerald-500"
      />
      <span class="absolute top-4 left-4 bg-emerald-600/90 text-white px-4 py-1 rounded-xl text-xs font-bold shadow-md">Plateforme 100 % française</span>
    </div>
  </div>
</section>

    <section class="max-w-6xl mx-auto grid grid-cols-1 md:grid-cols-3 gap-8 mb-12 px-4">
      <div class="glass-card p-8 flex flex-col items-center border border-gray-700 hover:shadow-violet-500/60 transition group">
        <div class="text-5xl mb-2 text-violet-400 group-hover:scale-110 transition drop-shadow">🚗</div>
        <div class="text-4xl text-white font-extrabold group-hover:text-violet-300 transition">{{ countCourses ?? '...' }}</div>
        <div class="text-gray-300 mt-2 text-lg font-medium">Courses réalisées</div>
      </div>
      <div class="glass-card p-8 flex flex-col items-center border border-gray-700 hover:shadow-emerald-500/60 transition group">
        <div class="text-5xl mb-2 text-green-400 group-hover:scale-110 transition drop-shadow">🌱</div>
        <div class="text-4xl text-white font-extrabold group-hover:text-emerald-200 transition">{{ totalCO2 ?? '...' }} kg</div>
        <div class="text-gray-300 mt-2 text-lg font-medium">CO₂ économisé</div>
      </div>
      <div class="glass-card p-8 flex flex-col items-center border border-gray-700 hover:shadow-amber-400/60 transition group">
        <div class="text-5xl mb-2 text-amber-300 group-hover:scale-110 transition drop-shadow">🛣️</div>
        <div class="text-4xl text-white font-extrabold group-hover:text-amber-100 transition">{{ totalKm ?? '...' }} km</div>
        <div class="text-gray-300 mt-2 text-lg font-medium">Kilomètres parcourus</div>
      </div>
    </section>

    <section class="w-full flex items-center justify-center mb-16 px-4">
      <div class="w-full max-w-3xl bg-gradient-to-r from-emerald-700 via-violet-700 to-violet-900 rounded-2xl shadow-2xl p-8 border border-emerald-800 text-center relative overflow-hidden">
        <div class="absolute inset-0 bg-white/10 blur-2xl opacity-70 pointer-events-none"></div>
        <div class="relative z-10">
          <h3 class="text-2xl md:text-3xl font-bold text-white mb-4 drop-shadow">Prêt à rejoindre le mouvement ?</h3>
          <p class="text-gray-100 mb-6 text-lg font-medium">
            Testez EcoDrive et devenez acteur du changement !<br />
            Commander une course, c’est soutenir un modèle plus juste, pour vous et pour la planète.
          </p>
          <router-link
            to="/course"
            class="inline-block px-10 py-3 rounded-xl bg-emerald-500 hover:bg-emerald-600 text-white font-bold text-lg shadow-xl transition"
          >Commander maintenant</router-link>
        </div>
      </div>
    </section>

<section class="max-w-3xl mx-auto my-20 flex flex-col items-center gap-10 px-4">
  <div class="w-full flex flex-col md:flex-row gap-6 justify-center items-center">
    <div class="relative w-48 h-48 md:w-56 md:h-56 rounded-2xl overflow-hidden shadow-2xl border-4 border-emerald-500 bg-gray-900 flex-shrink-0">
      <img
        src="https://plus.unsplash.com/premium_photo-1664283228853-cdefdd117b6f?q=80&w=2806&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
        alt="Conducteur VTC électrique"
        class="w-full h-full object-cover transition-transform duration-300 hover:scale-105"
      />
      <span class="absolute bottom-4 right-4 bg-violet-700/90 text-white px-4 py-1 rounded-xl text-xs font-bold shadow-md tracking-wide">VTC électrique</span>
    </div>
    <div class="relative w-48 h-48 md:w-56 md:h-56 rounded-2xl overflow-hidden shadow-2xl border-4 border-emerald-500 bg-gray-900 flex-shrink-0">
      <img
        src="https://plus.unsplash.com/premium_photo-1715789261470-fb25ffbf70d3?q=80&w=1470&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
        alt="Conductrice VTC"
        class="w-full h-full object-cover transition-transform duration-300 hover:scale-105"
      />
      <span class="absolute top-4 left-4 bg-emerald-600/90 text-white px-4 py-1 rounded-xl text-xs font-bold shadow-md tracking-wide">Mobilité propre</span>
    </div>
  </div>
  <div class="w-full text-center bg-gradient-to-br from-gray-900/80 via-violet-950/80 to-emerald-900/70 rounded-2xl shadow-xl border border-violet-800 px-7 py-9 flex flex-col gap-7 mt-3">
    <h2 class="text-3xl md:text-4xl font-bold bg-gradient-to-r from-violet-400 via-emerald-400 to-violet-400 bg-clip-text text-transparent drop-shadow mb-2">
      Les bienfaits d’EcoDrive
    </h2>
    <p class="text-lg text-gray-200 mb-1 leading-relaxed font-medium max-w-xl mx-auto">
      EcoDrive, c’est la garantie de trajets <span class="text-emerald-300 font-semibold">éco-responsables</span>, accessibles à tous, pour une ville plus verte et plus humaine.
    </p>
    <ul class="flex flex-col items-center gap-4 text-base md:text-lg font-medium mt-1">
      <li class="flex items-center gap-3 text-emerald-300">
        <svg class="w-7 h-7" fill="none" stroke="currentColor" stroke-width="2.2" viewBox="0 0 24 24">
          <circle cx="12" cy="12" r="10" /><path d="M8.5 12.5l2.5 2.5 4.5-4.5" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        Conducteurs certifiés et véhicules propres
      </li>
      <li class="flex items-center gap-3 text-violet-300">
        <svg class="w-7 h-7" fill="none" stroke="currentColor" stroke-width="2.2" viewBox="0 0 24 24">
          <circle cx="12" cy="12" r="10"/><path d="M12 8v4l3 3"/>
        </svg>
        Moins de CO₂, plus d’impact positif local
      </li>
      <li class="flex items-center gap-3 text-amber-300">
        <svg class="w-7 h-7" fill="none" stroke="currentColor" stroke-width="2.2" viewBox="0 0 24 24">
          <circle cx="12" cy="12" r="10"/><path d="M16 12v-2a4 4 0 00-8 0v2a4 4 0 108 0z"/>
        </svg>
        Rétribution juste pour les conducteurs locaux
      </li>
    </ul>
    <div class="text-base text-emerald-400 mt-2 font-semibold">
      <span class="inline-block animate-pulse">🌱</span> Vos trajets font avancer la transition !
    </div>
  </div>
</section>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const countCourses = ref(null)
const totalCO2 = ref(null)
const totalKm = ref(null)

async function fetchStats() {
  const res = await fetch('https://driver-ecodrive.liamcariou.fr/courses')
  if (res.ok) {
    const rawData = await res.json()
    const data = rawData.filter(c => c.status === 'completed')
    countCourses.value = data.length
    totalCO2.value = data.reduce((sum, c) => sum + (c.co2 || 0), 0).toFixed(1)
    totalKm.value = data.reduce((sum, c) => sum + (c.distance_km || 0), 0).toFixed(1)
  } else {
    countCourses.value = totalCO2.value = totalKm.value = '…'
  }
}
onMounted(fetchStats)
</script>

<style scoped>
.glass-card {
  background: linear-gradient(135deg, rgba(34,197,94,0.06) 0%, rgba(139,92,246,0.09) 80%);
  backdrop-filter: blur(7px) saturate(1.3);
  border-radius: 1.2rem;
  transition: box-shadow 0.2s, transform 0.18s;
}
</style>