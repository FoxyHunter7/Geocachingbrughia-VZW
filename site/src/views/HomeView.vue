<script setup>
  import { computed, onMounted, ref, watch } from 'vue';
  import { getHomePageEvents } from '@/services/EventService';
  import { getAllMessages } from '@/services/MessageService';
  import LanguageProvider from '@/services/LanguageService';
  import Message from '@/components/Message.vue'

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
  <Teleport to="#messages">
    <Message v-for="message in messages" :message="message" />
  </Teleport>
  <main>
    <section></section>
    <section></section>
  </main>
</template>

<style scoped>
  main {
    flex: 1 1 auto;
    height: 100%;
    background-color: hotpink;
  }

  section {
    border: 1rem solid black;
    height: 100%;
  }
</style>