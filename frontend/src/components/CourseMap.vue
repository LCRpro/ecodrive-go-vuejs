<template>
  <div class="w-full max-w-4xl bg-gray-900 bg-opacity-80 border border-gray-800 rounded-2xl shadow-2xl p-3 h-full flex items-center justify-center">
    <GoogleMap
      :from="from"
      :to="to"
      :follow="follow"
      :zoomToCurrent="zoomToCurrent"
      :currentPosition="currentPosition"
      class="w-full h-[500px] rounded-2xl"
    />
  </div>
</template>
<script setup>
import GoogleMap from './GoogleMap.vue'
defineProps(['from', 'to', 'follow', 'zoomToCurrent', 'currentPosition'])
</script>