<template>
  <div>
    <h2>Mon profil</h2>
    <div v-if="user">
      <div>Email: {{ user.email }}</div>
      <div>Nom: {{ user.name }}</div>
      <div>Naissance: <input v-model="user.birthdate" /></div>
      <div>Genre: <input v-model="user.gender" /></div>
      <div>Voiture: <input v-model="user.car" /></div>
      <div>Plaque: <input v-model="user.plate" /></div>
      <button @click="updateProfile">Mettre à jour</button>
    </div>
    <div v-else>Chargement...</div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
const user = ref(null)
const router = useRouter()

function getGoogleIdFromToken(token) {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.sub
  } catch {
    return null
  }
}

onMounted(async () => {
  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
    return
  }
  const googleId = getGoogleIdFromToken(token)
  if (!googleId) {
    router.push('/login')
    return
  }

  const res = await fetch(`http://localhost:8002/users/${googleId}`)
  if (res.ok) user.value = await res.json()
  else router.push('/login')
})

async function updateProfile() {
  const token = localStorage.getItem('token')
  const googleId = getGoogleIdFromToken(token)
  const res = await fetch(`http://localhost:8002/users/${googleId}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      Birthdate: user.value.birthdate,
      Gender: user.value.gender,
      Car: user.value.car,
      Plate: user.value.plate,
    })
  })
  if (res.ok) alert('Mise à jour OK')
}
</script>