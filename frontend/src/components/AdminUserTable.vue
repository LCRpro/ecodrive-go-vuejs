<template>
  <div class="mb-10">
    <h3 class="text-xl font-bold mb-4 text-white">Utilisateurs</h3>
    <div class="overflow-x-auto w-full">
      <table class="min-w-[700px] w-full text-sm border border-gray-700 rounded-xl overflow-hidden bg-gray-900">
        <thead>
          <tr>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              ID</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Email</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Nom</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Rôles</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Voiture</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Plaque</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="u in users" :key="u.id" class="even:bg-gray-900 hover:bg-gray-800/80">
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ u.id }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ u.email }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ u.name }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">
              <span v-for="role in JSON.parse(u.roles || '[]')" :key="role"
                class="inline-block px-2 py-1 mx-1 rounded-full bg-violet-700 text-white text-xs font-semibold">{{ role
                }}</span>
            </td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ u.car || '-' }}
            </td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ u.plate || '-' }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
<script setup>
defineProps(['users'])
</script>