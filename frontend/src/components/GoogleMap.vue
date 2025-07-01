<template>
  <div ref="mapDiv" style="width:100%;height:400px;border-radius:12px;"></div>
</template>

<script setup>
import { ref, watch, onMounted, onUnmounted } from 'vue'



const props = defineProps({
  from: Object,              
  to: Object,                
  follow: Boolean,           
  zoomToCurrent: Boolean,     
  currentPosition: Object    
})

const mapDiv = ref(null)
let map = null
let directionsRenderer = null
let directionsService = null
let myMarker = null
let watchId = null

function drawRoute(from, to) {
  if (!from || !to || !map || !window.google) return
  if (!directionsService) directionsService = new window.google.maps.DirectionsService()
  if (!directionsRenderer) {
    directionsRenderer = new window.google.maps.DirectionsRenderer({ suppressMarkers: false })
    directionsRenderer.setMap(map)
  }
  directionsService.route({
    origin: from,
    destination: to,
    travelMode: 'DRIVING'
  }, (result, status) => {
    if (status === 'OK') {
      directionsRenderer.setDirections(result)
    }
  })
}

function centerMap(coords, zoom=15) {
  if (map && coords) {
    map.setCenter(coords)
    map.setZoom(zoom)
  }
}

function updateMyMarker(coords) {
  if (!map || !window.google || !coords) return
  if (!myMarker) {
    myMarker = new window.google.maps.Marker({
      position: coords,
      map: map,
      icon: {
        url: 'https://maps.google.com/mapfiles/ms/icons/blue-dot.png',
        scaledSize: new window.google.maps.Size(38,38)
      }
    })
  } else {
    myMarker.setPosition(coords)
  }
}

onMounted(() => {
  function setupMap() {
    map = new window.google.maps.Map(mapDiv.value, {
      center: props.from || { lat: 48.86, lng: 2.35 },
      zoom: 13
    })
    if (props.from && props.to) drawRoute(props.from, props.to)
    else if (props.from) centerMap(props.from)

    if (props.follow && navigator.geolocation) {
      watchId = navigator.geolocation.watchPosition(
        pos => {
          const coords = { lat: pos.coords.latitude, lng: pos.coords.longitude }
          updateMyMarker(coords)
          centerMap(coords, 16)
        }
      )
    }
  }

  if (!window.google || !window.google.maps) {
    const script = document.createElement('script')
    script.src = `https://maps.googleapis.com/maps/api/js?key=${import.meta.env.VITE_GOOGLE_MAPS_API_KEY}&libraries=places`
    script.onload = setupMap
    document.body.appendChild(script)
  } else {
    setupMap()
  }
})

onUnmounted(() => {
  if (watchId) navigator.geolocation.clearWatch(watchId)
})

watch(() => [props.from, props.to], ([from, to]) => {
  if (from && to && map && window.google) drawRoute(from, to)
  else if (from && map && window.google) centerMap(from)
})

watch(() => props.currentPosition, (pos) => {
  if (pos && map && window.google) {
    updateMyMarker(pos)
  }
})

watch(
  () => props.zoomToCurrent,
  (val) => {
    if (val && props.currentPosition && map) {
      centerMap(props.currentPosition, 16)
    }
  }
)
</script>