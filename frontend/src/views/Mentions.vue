<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 via-gray-950 to-violet-900">
    <div class="max-w-xl w-full p-8 bg-gray-900 bg-opacity-90 rounded-2xl shadow-2xl border border-gray-800 text-left">
      <h1 class="text-3xl font-bold mb-5 text-white drop-shadow">Mentions légales</h1>
      <p class="text-gray-300 mb-4">
        <strong>EcoDrive</strong> est un service édité par la société LIAM Production.
      </p>
      <ul class="mb-6 text-gray-400 space-y-2">
        <li><span class="font-bold text-white">Éditeur :</span> LIAM Production</li>
        <li><span class="font-bold text-white">Adresse :</span> 22 Pl. Henri Gadeau de Kerville, 76100 Rouen France</li>
        <li><span class="font-bold text-white">SIRET :</span> 82392759500020</li>
        <li><span class="font-bold text-white">Contact :</span> contact[a]ecodrive.fr</li>
        <li><span class="font-bold text-white">Directeur de la publication :</span> Liam Cariou</li>
        <li><span class="font-bold text-white">Hébergement :</span> OVHcloud, 2 rue Kellermann, 59100 Roubaix</li>
      </ul>
      <p class="text-gray-500 text-xs mt-6">
        &copy; {{ year }} EcoDrive. Tous droits réservés.
      </p>
    </div>
  </div>
</template>

<script setup>
const year = new Date().getFullYear()
</script>