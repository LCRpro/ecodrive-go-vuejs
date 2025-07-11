<template>
  <div class="mb-10">
    <div class="flex flex-col items-center mb-10">
      <h3
        class="text-3xl md:text-4xl font-extrabold bg-gradient-to-r from-emerald-400 via-violet-400 to-emerald-400 bg-clip-text text-transparent drop-shadow-lg text-center">
        Transactions
      </h3>
    </div>
    <div class="overflow-x-auto w-full">
      <table class="min-w-[600px] w-full text-sm border border-gray-700 rounded-xl overflow-hidden bg-gray-900">
        <thead>
          <tr>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Date</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Utilisateur</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Type</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Montant (€)</th>
            <th
              class="bg-gray-800 text-violet-200 font-semibold text-center border border-gray-700 px-3 py-2 md:px-4 md:py-3">
              Détails</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="tx in transactions" :key="tx.id" class="even:bg-gray-900 hover:bg-gray-800/80">
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">{{ new
              Date(tx.created_at).toLocaleString() }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">
              {{ getUserByGoogleId(tx.google_id)?.email || tx.google_id }}
            </td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">
              <span v-if="tx.type === 'deposit'"
                class="px-2 py-1 rounded-full bg-emerald-600 text-white text-xs font-semibold">Dépôt</span>
              <span v-else-if="tx.type === 'withdrawal'"
                class="px-2 py-1 rounded-full bg-rose-600 text-white text-xs font-semibold">Retrait</span>
              <span v-else>{{ tx.type }}</span>
            </td>
            <td class="text-gray-200 border border-gray-700 text-center font-mono px-3 py-2 md:px-4 md:py-3">{{
              tx.amount.toFixed(2) }}</td>
            <td class="text-gray-200 border border-gray-700 text-center px-3 py-2 md:px-4 md:py-3">
              <span v-if="tx.type === 'deposit'">Carte</span>
              <span v-else>IBAN : {{ tx.iban_mask }}</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
<script setup>
defineProps(['transactions', 'getUserByGoogleId'])
</script>