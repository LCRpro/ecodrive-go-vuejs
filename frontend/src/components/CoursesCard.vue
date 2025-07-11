<template>
  <div class="bg-gray-900 bg-opacity-90 rounded-2xl shadow-lg border border-gray-800 p-6 mb-4">
    <h3 class="text-xl font-bold mb-3" :class="mode === 'passenger' ? 'text-emerald-400' : 'text-violet-400'">
      {{ mode === 'passenger' ? 'Mes courses passager' : 'Mes courses en tant que driver' }}
    </h3>
    <div class="overflow-x-auto w-full">
      <table class="min-w-[700px] w-full text-sm border border-gray-700 rounded-xl overflow-hidden bg-gray-900">
        <thead>
          <tr>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Départ</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Arrivée</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Distance</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Prix</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              CO₂</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              {{ mode === 'passenger' ? 'Driver' : 'Passager' }}</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Statut</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="c in courses" :key="c.id" class="even:bg-gray-900 hover:bg-gray-800/80">
            <td class="text-gray-200 border border-gray-700 px-3 py-2 md:px-4 md:py-3 text-center">{{ c.start_address }}
            </td>
            <td class="text-gray-200 border border-gray-700 px-3 py-2 md:px-4 md:py-3 text-center">{{ c.end_address }}
            </td>
            <td class="text-gray-200 border border-gray-700 px-3 py-2 md:px-4 md:py-3 text-center">{{
              c.distance_km.toFixed(2) }} km</td>
            <td class="text-gray-200 border border-gray-700 px-3 py-2 md:px-4 md:py-3 text-center">
              <span class="font-mono" :class="c.amount > 0 ? 'text-green-400' : 'text-red-400'">
                {{ c.amount?.toFixed(2) ?? '?' }} €
              </span>
            </td>
            <td class="text-gray-200 border border-gray-700 px-3 py-2 md:px-4 md:py-3 text-center">{{ (c.co2 /
              1000).toFixed(2) }} kg</td>
            <td class="text-gray-200 border border-gray-700 px-3 py-2 md:px-4 md:py-3 text-center">
              {{ names[mode === 'passenger' ? c.driver_id : c.passenger_id] ?? '-' }}
            </td>
            <td class="text-gray-200 border border-gray-700 px-3 py-2 md:px-4 md:py-3 text-center">
              <span v-if="c.status === 'requested'"
                class="px-2 py-1 rounded-full bg-amber-500 text-white text-xs font-semibold">En attente</span>
              <span v-else-if="c.status === 'accepted'"
                class="px-2 py-1 rounded-full bg-emerald-600 text-white text-xs font-semibold">Acceptée</span>
              <span v-else-if="c.status === 'in_progress'"
                class="px-2 py-1 rounded-full bg-blue-600 text-white text-xs font-semibold">En cours</span>
              <span v-else-if="c.status === 'completed'"
                class="px-2 py-1 rounded-full bg-violet-600 text-white text-xs font-semibold">Terminée</span>
              <span v-else-if="c.status === 'cancelled'"
                class="px-2 py-1 rounded-full bg-rose-600 text-white text-xs font-semibold">Annulée</span>
              <span v-else class="px-2 py-1 rounded-full bg-gray-700 text-gray-200 text-xs font-semibold">{{ c.status
                }}</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
<script setup>
defineProps(['courses', 'names', 'mode'])
</script>