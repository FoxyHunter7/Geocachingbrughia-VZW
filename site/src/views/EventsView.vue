<script setup>
  import { getAllEvents } from '@/services/EventService';
  import LanguageProvider from '@/services/LanguageService';
  import { computed, onMounted, ref, watch } from 'vue';
  import Event from '@/components/Event.vue';

  const events = ref(["loading"]);
  const loaderActive = ref(false);
  const lang = computed(() => LanguageProvider.CURR_LANG.value);

  const search = ref("");

  async function fetchData() {
    const pagedResponse = await getAllEvents(search.value);
    if (pagedResponse.data) {
      events.value = pagedResponse.data;
    } else {
      events.value = pagedResponse;
    }
  }

  async function initLoaderOnDelay() {
    await new Promise(r => setTimeout(r, 2500));
    loaderActive.value = true;
  }

  onMounted(fetchData);
  watch(lang, async () => await fetchData());

  initLoaderOnDelay();
</script>

<template>
  <main>
    <section v-show="events.length !== 0 && events[0] === 'loading'" id=loading>
      <div v-show="loaderActive" class="initial-loader"></div>
    </section>
    <section v-show="events.length === 0"><p>No events Found</p></section>
    <Event v-if="events.length !== 0 && events[0] !== 'loading'" v-for="event in events" :event="event" />
  </main>
</template>

<style scoped>
  main {
    flex: 1 1 auto;
    gap: 5rem;
    height: 100%;
    overflow-y: auto;
  }

  section {
    height: 100%;
    max-height: 40rem;
    max-width: 100rem;
    margin: 0 auto;
  }

  #loading {
    display: flex;
    justify-content: center;
    align-items: center;
  }
</style>