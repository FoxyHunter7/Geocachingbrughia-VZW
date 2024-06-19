<script setup>
  import { computed, onMounted, ref } from 'vue';
  import { getHomePageEvents } from '@/services/EventService';
  import LanguageProvider from '@/services/LanguageService';

  const events = ref([]);
  const lang = computed(() => LanguageProvider.CURR_LANG.value);

  onMounted(async () => {
    events.value = await getHomePageEvents();
  });

</script>

<template>
  <main>
    <p>{{ lang }}</p>
  </main>
</template>

<style scoped>
</style>