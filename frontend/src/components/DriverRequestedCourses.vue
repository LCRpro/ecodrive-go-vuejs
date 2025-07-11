<template>
  <div class="w-full max-w-5xl flex flex-wrap gap-6 justify-center mb-8">
    <div v-for="c in displayedCourses" :key="c.id"
      class="w-full md:w-[340px] bg-gradient-to-r from-gray-800/80 via-gray-900/80 to-violet-950/40 rounded-xl px-6 py-5 shadow-xl border border-gray-800 space-y-3 flex flex-col hover:scale-[1.015] transition duration-300 divide-y divide-gray-700">
      <div class="flex justify-between items-center pb-3">
        <span class="text-xs text-gray-300 uppercase tracking-widest font-semibold">Passager</span>
        <span class="font-mono text-emerald-400 text-lg font-semibold">{{ c.passenger_name }}</span>
      </div>
      <div class="flex justify-between items-center py-3">
        <span class="text-xs text-gray-300 uppercase tracking-widest font-semibold">Départ</span>
        <span class=" pl-4 font-mono text-violet-400 text-xs font-semibold">
          {{ requestedAddresses[c.id]?.start || 'Chargement...' }}
        </span>
      </div>
      <div class="flex justify-between items-center py-3">
        <span class="text-xs text-gray-300 uppercase tracking-widest font-semibold">Arrivée</span>
        <span class="font-mono text-violet-400 text-xs pl-4 font-semibold">
          {{ requestedAddresses[c.id]?.end || 'Chargement...' }}
        </span>
      </div>
      <div class="flex justify-between items-center py-3">
        <span class="text-xs text-gray-300 uppercase tracking-widest font-semibold">Distance</span>
        <span class="font-mono text-blue-300 text-xm font-semibold">
          <span class="inline-block bg-blue-900/70 px-2 py-1 rounded-md">{{ c.distance_km }} km</span>
        </span>
      </div>
      <div class="flex justify-between items-center py-3">
        <span class="text-xs text-gray-300 uppercase tracking-widest font-semibold">Prix</span>
        <span class="font-mono text-amber-300  font-bold">
          <span class="inline-block text-xm bg-amber-800/80 px-3 py-1 rounded-lg">{{ c.amount }} €</span>
        </span>
      </div>
      <div class="flex justify-between items-center pt-3">
        <span class="text-xs text-gray-300 uppercase tracking-widest font-semibold">Statut</span>
        <span class="inline-block px-3 py-1 rounded-full bg-yellow-500 text-white text-xs font-bold">En attente</span>
      </div>
      <div class="flex gap-2 mt-4">
        <button @click="$emit('preview', c.id)"
          :class="previewCourseId === c.id ? 'bg-blue-700' : 'bg-blue-600 hover:bg-blue-700'"
          class="flex-1 mt-2 py-2 text-xs rounded-lg text-white font-bold shadow transition">
          Voir sur la carte
        </button>
        <button @click="$emit('accept', c.id)"
          class="flex-1 mt-2 py-2 text-xs rounded-lg bg-emerald-600 hover:bg-emerald-700 text-white font-bold shadow transition">
          Accepter
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps(['displayedCourses', 'requestedAddresses', 'previewCourseId'])
defineEmits(['preview', 'accept'])
</script>