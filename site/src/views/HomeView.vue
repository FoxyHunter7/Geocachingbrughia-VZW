<script setup>
  import { computed, onMounted, ref, watch } from 'vue';
  import { getHomePageEvents } from '@/services/EventService';
  import { getAllMessages } from '@/services/MessageService';
  import LanguageProvider from '@/services/LanguageService';
  import Message from '@/components/Message.vue'
  import Event from '@/components/Event.vue'
  import Hero from '@/components/Hero.vue';

  const events = ref(["loading"]);
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
    <section v-show=" events.length !== 0 && events[0] === 'loading'" id=loading><p>Loading</p></section>
    <Hero v-show="events.length === 0" id="default-hero"/>
    <Event v-if="events.length !== 0 && events[0] !== 'loading'" v-for="event in events" :event="event" />
  </main>
</template>

<style scoped>
  main {
    flex: 1 1 auto;
    height: 100%;
    overflow-y: auto;
  }

  section {
    height: 100%;
    max-height: 50rem;
  }
</style>