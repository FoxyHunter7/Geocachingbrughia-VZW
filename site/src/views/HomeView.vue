<script setup>
  import { computed, onMounted, ref, watch } from 'vue';
  import { getHomePageEvents } from '@/services/EventService';
  import { getAllMessages } from '@/services/MessageService';
  import { getAllSocials } from '@/services/SocialService';
  import LanguageProvider from '@/services/LanguageService';
  import Message from '@/components/Message.vue';
  import Event from '@/components/Event.vue';
  import Hero from '@/components/Hero.vue';
  import Socials from '@/components/Socials.vue';
  import Contact from '@/components/Contact.vue';

  const events = ref(["loading"]);
  const messages = ref([]);
  const socials = ref([]);
  const loaderActive = ref(false);
  const lang = computed(() => LanguageProvider.CURR_LANG.value);
  
  async function fetchData() {
    socials.value = await getAllSocials();
    messages.value = await getAllMessages();
    events.value = await getHomePageEvents();
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
  <Teleport to="#messages">
    <Message v-for="message in messages" :message="message" />
  </Teleport>
  <main>
    <section v-show="events.length !== 0 && events[0] === 'loading'" id=loading>
      <div v-show="loaderActive" class="initial-loader"></div>
    </section>
    <Hero v-show="events.length === 0" id="default-hero"/>
    <Event v-if="events.length !== 0 && events[0] !== 'loading'" v-for="event in events" :event="event" />
    <Socials v-if="socials.length !== 0 && events[0] !== 'loading'" :socials="socials"/>
    <Contact />
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
    max-width: 100rem;
    margin: 0 auto 5rem auto;
  }

  #loading {
    display: flex;
    justify-content: center;
    align-items: center;
  }
</style>