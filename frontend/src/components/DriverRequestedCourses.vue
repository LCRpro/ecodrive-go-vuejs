<template>
  <div class="w-full max-w-5xl flex flex-wrap gap-6 justify-center mb-8">
    <div
      v-for="c in displayedCourses"
      :key="c.id"
      class="w-full md:w-[340px] bg-gradient-to-r from-gray-800/90 via-gray-900/90 to-violet-950/60 rounded-xl px-6 py-5 shadow-inner border border-gray-800 space-y-3 flex flex-col"
    >
      <div class="flex justify-between items-center">
        <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Passager</span>
        <span class="font-mono text-emerald-400 text-base font-semibold">{{ c.passenger_name }}</span>
      </div>
      <div class="flex justify-between items-center">
        <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Départ</span>
        <span class="font-mono text-violet-400 text-base font-semibold">
          {{ requestedAddresses[c.id]?.start || 'Chargement...' }}
        </span>
      </div>
      <div class="flex justify-between items-center">
        <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Arrivée</span>
        <span class="font-mono text-violet-400 text-base font-semibold">
          {{ requestedAddresses[c.id]?.end || 'Chargement...' }}
        </span>
      </div>
      <div class="flex justify-between items-center">
        <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Distance</span>
        <span class="font-mono text-blue-200 text-base font-semibold">
          <span class="inline-block bg-blue-900/70 px-2 py-1 rounded-md">{{ c.distance_km }} km</span>
        </span>
      </div>
      <div class="flex justify-between items-center">
        <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Prix</span>
        <span class="font-mono text-amber-300 text-lg font-bold">
          <span class="inline-block bg-amber-800/80 px-3 py-1 rounded-lg">{{ c.amount }} €</span>
        </span>
      </div>
      <div class="flex justify-between items-center">
        <span class="text-xs text-gray-400 uppercase tracking-widest font-bold">Statut</span>
        <span class="inline-block px-3 py-1 rounded-full bg-amber-600 text-white text-xs font-semibold">En attente</span>
      </div>
      <div class="flex gap-2 mt-2">
        <button
          @click="$emit('preview', c.id)"
          :class="previewCourseId === c.id ? 'bg-blue-700' : 'bg-blue-600 hover:bg-blue-700'"
          class="flex-1 py-2 rounded-lg text-white font-bold shadow transition"
        >Voir sur la carte</button>
        <button
          @click="$emit('accept', c.id)"
          class="flex-1 py-2 rounded-lg bg-emerald-600 hover:bg-emerald-700 text-white font-bold shadow transition"
        >Accepter</button>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps(['displayedCourses', 'requestedAddresses', 'previewCourseId'])
defineEmits(['preview', 'accept'])
</script>