<template>
  <div class="w-full flex flex-col md:flex-row items-center justify-center gap-8">
    <div
      class="w-full max-w-sm p-7 rounded-2xl shadow-2xl bg-gray-900 bg-opacity-90 border border-gray-800 mb-8 md:mb-0">
      <template v-for="c in displayedCourses" :key="c.id">
        <h2
          class="text-white text-2xl font-extrabold bg-gradient-to-r from-emerald-400 via-violet-400 to-emerald-400 bg-clip-text text-transparent drop-shadow-lg text-center mb-6">
          Course</h2>
        <div
          class="bg-gradient-to-r from-gray-800/90 via-gray-900/95 to-violet-950/70 rounded-xl px-6 py-5 mb-3 shadow-inner border border-gray-800 space-y-4">
          <slot :course="c" />
        </div>
        <div class="mt-4 flex gap-2">
          <button v-if="c.status === 'accepted' && c.driver_id === driverId" @click="$emit('start', c.id)"
            class="w-full py-3 rounded-lg bg-blue-600 hover:bg-blue-700 text-white text-lg font-bold shadow transition">
            Commencer la course
          </button>
          <button v-if="c.status === 'in_progress' && c.driver_id === driverId" @click="$emit('complete', c.id)"
            class="w-full py-3 rounded-lg bg-violet-700 hover:bg-violet-800 text-white text-lg font-bold shadow transition">
            Terminer la course
          </button>
        </div>
      </template>
    </div>
    <slot name="map" />
  </div>
</template>
<script setup>
defineProps(['displayedCourses', 'driverId'])
defineEmits(['start', 'complete'])
</script>