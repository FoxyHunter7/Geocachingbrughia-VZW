<script setup>
  import { computed, onMounted, ref, watch } from 'vue';
  import { getHomePageEvents } from '@/services/EventService';
  import { getAllMessages } from '@/services/MessageService';
  import LanguageProvider from '@/services/LanguageService';

  const events = ref([]);
  const messages = ref([]);
  const lang = computed(() => LanguageProvider.CURR_LANG.value);

  async function fetchData() {
    events.value = await getHomePageEvents();
    messages.value = await getAllMessages();
  }

  onMounted(fetchData);
  watch(lang, async () => await fetchData());

</script>

<template>
  <div v-for="message in messages" >
    <h2>{{ message.title }}</h2>
    <p>{{ message.body }}</p>
  </div>
  <main>
  </main>
</template>

<style scoped>
</style>