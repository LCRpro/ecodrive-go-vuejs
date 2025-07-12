<template>

  <div
    class="bg-gradient-to-br from-gray-900/95 via-gray-900/90 to-violet-950/80 rounded-2xl shadow-xl border border-gray-800 p-7 mb-4">
    <div class="flex items-center gap-3 mb-4">

      <h2
        class="text-2xl font-bold font-extrabold bg-gradient-to-r from-emerald-400 via-violet-400 to-emerald-400 bg-clip-text text-transparent drop-shadow-lg text-center">
        Devenir driver</h2>
    </div>
    <div v-if="driverRequest">
      <div class="flex items-center gap-2 mb-4">
        <span v-if="driverRequest.status === 'pending'"
          class="inline-block px-3 py-1 rounded-full bg-yellow-600/20 border border-yellow-500 text-yellow-300 font-semibold text-xs tracking-wide">En
          attente de validation</span>
        <span v-else-if="driverRequest.status === 'accepted'"
          class="inline-block px-3 py-1 rounded-full bg-emerald-700/20 border border-emerald-500 text-emerald-300 font-semibold text-xs tracking-wide">Validé</span>
        <span v-else-if="driverRequest.status === 'refused'"
          class="inline-block px-3 py-1 rounded-full bg-red-800/20 border border-red-500 text-red-300 font-semibold text-xs tracking-wide">Refusé</span>
      </div>
      <div class="mb-2 flex flex-col gap-3">
        <label class="text-violet-400 font-bold" for="car-input">Voiture</label>
        <input id="car-input" type="text" :value="driverRequest.car" readonly
          class="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-gray-200 placeholder-gray-400 font-mono cursor-not-allowed opacity-80" />
        <label class="text-violet-400 font-bold" for="plate-input">Plaque</label>
        <input id="plate-input" type="text" :value="driverRequest.plate" readonly
          class="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-gray-200 placeholder-gray-400 font-mono cursor-not-allowed opacity-80" />
      </div>
      <div v-if="driverRequest.status === 'pending'" class="text-xs text-yellow-300 mt-2">Un administrateur doit valider
        votre demande.</div>
      <div v-if="driverRequest.status === 'accepted'" class="text-sm text-emerald-400 mt-6 font-semibold">
        Si vous avez un problème, n'hésitez pas à contacter un administrateur ou à vous déconnecter puis vous reconnecter.
      </div>
      <div v-if="driverRequest.status === 'refused'" class="text-sm text-red-400 mt-2 font-semibold">Votre demande a été
        refusée.</div>
    </div>
    <button v-if="!driverRequest" @click="$emit('openModal')"
      class="mt-5 w-full py-3 rounded-xl bg-violet-600 hover:bg-violet-700 text-white font-semibold shadow-lg transition-all text-lg tracking-wide">Faire
      une demande driver</button>
  </div>
</template>
<script setup>
defineProps(['driverRequest', 'user'])
</script>