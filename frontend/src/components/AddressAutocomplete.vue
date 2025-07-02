<template>
  <div class="mb-4 relative">
    <label class="text-gray-300 mb-1 block">Destination :</label>
    <input
      :ref="inputRef"
      :value="search"
      @input="onInput"
      placeholder="Ex: 6 rue de Paris..."
      autocomplete="off"
      class="bg-gray-800 border border-gray-700 text-white placeholder-gray-400 rounded-lg px-4 py-3 mb-3 w-full focus:outline-none focus:ring-2 focus:ring-violet-500"
    />
    <ul v-if="suggestions.length" class="bg-gray-800 border border-gray-700 rounded-lg shadow-lg max-h-60 overflow-auto absolute z-50 w-full text-white">
      <li 
        v-for="addr in suggestions" 
        :key="addr.place_id" 
        @click="$emit('select', addr)" 
        class="cursor-pointer px-4 py-3 hover:bg-violet-700 transition"
      >{{ addr.description }}</li>
    </ul>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const props = defineProps(['search', 'suggestions', 'inputRef'])
const emit = defineEmits(['update:search', 'select', 'input'])

function onInput(e) {
  emit('update:search', e.target.value)
  emit('input', e)
}
</script>