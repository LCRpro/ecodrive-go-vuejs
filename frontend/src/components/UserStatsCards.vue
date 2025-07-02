<template>
  <div class="grid grid-cols-1 sm:grid-cols-3 gap-6 mb-10">
    <div
      class="rounded-2xl bg-gradient-to-br from-gray-900/90 via-gray-800/80 to-violet-950/80 shadow-lg border border-gray-800 p-6 flex flex-col items-center justify-center">
      <div class="text-sm text-gray-400 mb-1">Solde actuel</div>
      <div class="text-3xl font-bold text-emerald-400 font-mono">{{ balance }} €</div>
    </div>
    <div
      class="rounded-2xl bg-gradient-to-br from-gray-900/90 via-gray-800/80 to-violet-950/80 shadow-lg border border-gray-800 p-6 flex flex-col items-center justify-center">
      <div class="text-sm text-gray-400 mb-1">CO₂ économisé</div>
      <div class="text-3xl font-bold text-blue-400 font-mono">
        <template v-if="co2Passenger !== '0.00' && co2Driver !== '0.00'">
          {{ co2Passenger }} kg / {{ co2Driver }} kg
        </template>
        <template v-else-if="co2Passenger !== '0.00'">
          {{ co2Passenger }} kg
        </template>
        <template v-else>
          {{ co2Driver }} kg
        </template>
      </div>
    </div>
    <div
      class="rounded-2xl bg-gradient-to-br from-gray-900/90 via-gray-800/80 to-violet-950/80 shadow-lg border border-gray-800 p-6 flex flex-col items-center justify-center">
      <div class="text-sm text-gray-400 mb-1">Kilomètres parcourus</div>
      <div class="text-3xl font-bold text-violet-400 font-mono">
        <template v-if="kmPassenger !== '0.0' && kmDriver !== '0.0'">
          {{ kmPassenger }} km / {{ kmDriver }} km
        </template>
        <template v-else-if="kmPassenger !== '0.0'">
          {{ kmPassenger }} km
        </template>
        <template v-else>
          {{ kmDriver }} km
        </template>
      </div>
    </div>
  </div>
</template>
<script setup>
defineProps(['balance', 'co2Passenger', 'co2Driver', 'kmPassenger', 'kmDriver'])
</script>