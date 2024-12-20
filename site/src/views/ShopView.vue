<script setup>
import { onMounted } from 'vue';

onMounted(() => {
  window.location.href = "https://ticketing.geocachingbrughia.be/geocachingbrughia";
});
</script>

<template>
  <main>
  </main>
</template>

<style scoped>
main {
    position: relative;
    flex: 1 1 auto;
    height: 100%;
    overflow-y: auto;
}
section {
  height: 70vh;
  max-width: 100rem;
  margin: 0 auto;
}
iframe {
  width: 100%;
  height: 100%;
}
</style>