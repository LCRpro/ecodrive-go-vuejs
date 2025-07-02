<template>
  <div class="mb-10">
    <h3 class="text-xl font-bold mb-4 text-white">Toutes les courses</h3>
    <div class="overflow-x-auto w-full">
      <table class="min-w-[850px] w-full text-sm border border-gray-700 rounded-xl overflow-hidden bg-gray-900">
        <thead>
          <tr>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              ID</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Départ</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Arrivée</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Passager</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Driver</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Distance</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Prix (€)</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              CO₂ (kg)</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Statut</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="c in courses" :key="c.id" class="even:bg-gray-900 hover:bg-gray-800/80">
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ c.id }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ c.start_lat }}, {{
              c.start_lng }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ c.end_lat }}, {{
              c.end_lng }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{
              getUserByGoogleId(c.passenger_id)?.email || c.passenger_id }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{
              getUserByGoogleId(c.driver_id)?.email || c.driver_id || '-' }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{
              c.distance_km?.toFixed(2) ?? '-' }} km</td>
            <td class="text-gray-200 border border-gray-700 text-center font-mono px-3 py-2 md:px-4 md:py-3">{{
              c.amount?.toFixed(2) ?? '-' }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ c.co2 ?
              (c.co2 / 1000).toFixed(2) : '-' }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">
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
defineProps(['courses', 'getUserByGoogleId'])
</script>