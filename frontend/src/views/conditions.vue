<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 via-gray-950 to-violet-900">
    <div class="max-w-2xl w-full p-8 bg-gray-900 bg-opacity-90 rounded-2xl shadow-2xl border border-gray-800 text-left">
      <h1 class="text-3xl font-bold mb-5 text-white drop-shadow">Conditions générales de vente</h1>
      <p class="text-gray-300 mb-6">
        Les présentes conditions générales régissent l’utilisation du service <strong>EcoDrive</strong> proposé par Mindset Production.
      </p>
      <ol class="list-decimal list-inside text-gray-400 space-y-3 mb-8">
        <li>
          <span class="font-bold text-white">Objet :</span> EcoDrive propose une plateforme de mise en relation entre passagers et conducteurs pour des trajets urbains et périurbains.
        </li>
        <li>
          <span class="font-bold text-white">Commande et paiement :</span> Toute commande de course s’effectue via l’application. Le paiement est sécurisé et doit être validé avant le début de la course.
        </li>
        <li>
          <span class="font-bold text-white">Annulation :</span> Les passagers peuvent annuler une course sans frais tant qu’elle n’est pas prise en charge par un conducteur. Passé ce délai, des frais peuvent s’appliquer.
        </li>
        <li>
          <span class="font-bold text-white">Responsabilité :</span> EcoDrive décline toute responsabilité en cas de retard ou d’annulation pour des raisons indépendantes de sa volonté (météo, panne, etc.).
        </li>
        <li>
          <span class="font-bold text-white">Protection des données :</span> Les données personnelles sont traitées conformément à la législation en vigueur et ne sont jamais revendues à des tiers.
        </li>
        <li>
          <span class="font-bold text-white">Contact :</span> Pour toute question, contactez-nous à <a href="mailto:contact@ecodrive.fr" class="text-violet-300 underline">contact@ecodrive.fr</a>.
        </li>
      </ol>
      <p class="text-gray-500 text-xs mt-4">
        Dernière mise à jour : {{ year }}<br />
        &copy; {{ year }} EcoDrive. Tous droits réservés.
      </p>
    </div>
  </div>
</template>

<script setup>
const year = new Date().getFullYear()
</script>