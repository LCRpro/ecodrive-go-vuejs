<template>
  <div>
    <template v-if="activeCourse">
      <h2 class="text-white text-2xl font-bold mb-6">Course</h2>
      <div class="bg-gradient-to-r from-gray-800/90 via-gray-900/95 to-violet-950/70 rounded-xl px-6 py-5 mb-3 shadow-inner border border-gray-800 space-y-4">
        <div class="flex items-center justify-between">
          <span class="text-xs text-gray-400 font-bold">Départ</span>
          <span class="text-emerald-400 font-mono text-base">{{ startAddress || 'Chargement...' }}</span>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-xs text-gray-400 font-bold">Arrivée</span>
          <span class="text-violet-400 font-mono text-base">{{ endAddress || 'Chargement...' }}</span>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-xs text-gray-400 font-bold">Distance</span>
          <span class="font-mono text-blue-200 text-base">
            <span class="inline-block bg-blue-900/70 px-2 py-1 rounded-md">{{ activeCourse.distance_km }} km</span>
          </span>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-xs text-gray-400 font-bold">Durée</span>
          <span class="font-mono text-cyan-200 text-base">
            <span class="inline-block bg-cyan-900/60 px-2 py-1 rounded-md">{{ duration || '...' }}</span>
          </span>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-xs text-gray-400 font-bold">Prix</span>
          <span class="font-mono text-amber-300 text-lg font-bold">
            <span class="inline-block bg-amber-800/80 px-3 py-1 rounded-lg">{{ activeCourse.amount }} €</span>
          </span>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-xs text-gray-400 font-bold">Statut</span>
          <span>
            <span
              :class="badgeClass(activeCourse.status)"
              class="inline-block px-3 py-1 rounded-full text-xs font-semibold"
            >{{ displayStatus(activeCourse.status) }}</span>
          </span>
        </div>
      </div>
      <button v-if="canCancel" @click="cancelCourse" class="w-full py-3 rounded-lg bg-red-600 hover:bg-red-700 text-white text-lg font-bold shadow transition mt-6">Annuler la course</button>
      <button @click="$emit('refresh')" class="w-full py-3 rounded-lg bg-violet-600 hover:bg-violet-700 text-white text-lg font-bold shadow transition mt-4">Rafraîchir</button>
      <div v-if="cancelError" class="mt-4 rounded-lg bg-red-700 bg-opacity-80 text-white p-3 text-center font-semibold">{{ cancelError }}</div>
    </template>
    <template v-else>
      <h2 class="text-white text-2xl font-bold mb-6">Commander une course</h2>
     <AddressAutocomplete
  :search="search"
  :suggestions="suggestions"
  @select="pickSuggestion"
  @update:search="val => search = val"
  @input="searchAddr"
/>
      <button @click="requestCourse" :disabled="!destination" class="w-full py-3 rounded-lg bg-emerald-600 hover:bg-emerald-700 text-white text-lg font-bold shadow transition disabled:opacity-50 disabled:cursor-not-allowed">Valider</button>
      <div v-if="error" class="mt-4 rounded-lg bg-red-700 bg-opacity-80 text-white p-3 text-center font-semibold">{{ error }}</div>
      <div v-if="price" class="mt-6 text-gray-300 space-y-1">
        <div class="font-bold text-white">Prix estimé : <span class="font-mono">{{ price }} €</span></div>
        <div v-if="km">Distance : <span class="font-mono">{{ km }} km</span></div>
        <div v-if="duration">Durée estimée : <span class="font-mono">{{ duration }}</span></div>
      </div>
    </template>
  </div>
</template>

<script setup>
import AddressAutocomplete from './AddressAutocomplete.vue'
defineProps([
  'activeCourse', 'startAddress', 'endAddress', 'duration', 'canCancel', 'cancelError',
  'search', 'suggestions', 'destination', 'error', 'price', 'km'
])
defineEmits(['refresh'])
</script>