<script setup>
  import { getAllEvents } from '@/services/EventService';
  import LanguageProvider from '@/services/LanguageService';
  import StaticContentProvider from '@/services/StaticContentService';
  import { computed, onMounted, ref, watch } from 'vue';
  import Event from '@/components/Event.vue';

  const lang = computed(() => LanguageProvider.CURR_LANG.value);
  const dictionary = StaticContentProvider.DICTIONARY;

  const events = ref(["loading"]);

  async function fetchData() {
    const pagedResponse = await getAllEvents();
    if (pagedResponse.data) {
      events.value = pagedResponse.data;
    } else {
      events.value = pagedResponse;
    }
  }

  const search = ref("");

  const filteredEvents = computed(() => {
    if (!search.value) {
      return events.value;
    }

    return events.value.filter(event =>
      event.title.toLowerCase().includes(search.value.toLowerCase())
    );
  });

  const loaderActive = ref(false);

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
    <form method="post">
      <div>
        <input v-model="search" type="search" id="search" name="search" autocomplete="search" required :placeholder="dictionary.FormSearch[lang]">
      </div>
    </form>
    <Event v-if="events.length !== 0 && events[0] !== 'loading'" v-for="event in filteredEvents" :event="event" />
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

  form {
    display: flex;
    justify-content: flex-end;
    padding: 1rem;
    gap: 1rem;
  }

  form > div input {
    width: 20rem;
    max-width: 95vw;
  }

  form > div {
    position: relative;
  }

  form > div::before {
    content: '';
    height: 1.3rem;
    width: 1.3rem;
    position: absolute;
    left: 0.5rem;
    top: 0.3rem;
    background-color: var(--color-text);
    mask: url(../assets/media/search.svg);
    mask-size: contain;
  }

  form input, form textarea {
    height: 2rem;
    border-radius: 0.3rem;
    border: solid 0.1rem var(--color-text);
    outline: none;
    font-family: inherit;
    font-size: 1rem;
    box-sizing: border-box;
    padding: 0 0.5rem 0 2rem;
    line-height: 2rem;
    text-transform: capitalize;
  }

  @media screen and (max-width: 1000px) {
    form {
      justify-content: center;
    }
  }
</style>