<template>
  <div
    class="bg-gradient-to-br from-gray-900/95 via-gray-900/90 to-violet-950/80 rounded-2xl shadow-xl border border-gray-800 p-7 mb-4">
    <h2 class="text-2xl font-bold mb-5 text-white">Mon profil</h2>
    <div class="space-y-5">
      <div>
        <label class="block text-gray-400 text-sm mb-1">Email</label>
        <input :value="user.email" disabled
          class="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-gray-400 font-mono placeholder-gray-400 cursor-not-allowed opacity-60" />
      </div>
      <div>
        <label class="block text-gray-400 text-sm mb-1">Nom</label>
        <input v-model="user.name"
          class="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow transition-all duration-200 font-semibold font-mono" />
      </div>
          <div>
        <label class="block text-gray-400 text-sm mb-1">Numéro de téléphone</label>
        <input v-model="user.phone"
          class="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow transition-all duration-200 font-semibold font-mono" />
      </div>
      <div>
        <label class="block text-gray-400 text-sm mb-1">Naissance</label>
        <input v-model="user.birthdate" type="date"
          class="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow transition-all duration-200 font-semibold font-mono" />
      </div>
      <div>
        <label class="block text-gray-400 text-sm mb-1">Genre</label>
        <select v-model="user.gender"
          class="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-violet-500 focus:bg-gray-900 shadow transition-all duration-200 font-semibold font-mono">
          <option value="">-</option>
          <option value="male">Homme</option>
          <option value="female">Femme</option>
          <option value="other">Autre</option>
        </select>
      </div>

    </div>
    <div class="flex justify-center mt-7">
      <button @click="$emit('updateProfile')"
        class="px-6 py-2 rounded-xl bg-violet-600 hover:bg-violet-700 text-white font-bold shadow transition">Mettre à
        jour</button>
    </div>
  </div>
</template>

<script setup>
defineProps(['user'])
</script>