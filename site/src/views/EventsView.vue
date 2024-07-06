<script setup>
  import { getAllEvents } from '@/services/EventService';
  import LanguageProvider from '@/services/LanguageService';
  import StaticContentProvider from '@/services/StaticContentService';
  import { computed, onMounted, ref, watch } from 'vue';
  import Event from '@/components/Event.vue';

  const lang = computed(() => LanguageProvider.CURR_LANG.value);
  const dictionary = StaticContentProvider.DICTIONARY;

  const events = ref(["loading"]);
  const currPage = ref(1);
  const lastPage = ref(1);

  async function fetchData() {
    const pagedResponse = await getAllEvents("", "", "", "", currPage.value);
    if (pagedResponse.data) {
      events.value = pagedResponse.data;
      currPage.value = pagedResponse.current_page;
      lastPage.value = pagedResponse.last_page;
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

  function nextPage() {
    if (currPage.value < lastPage.value) {
      currPage.value++;
    }
  }

  function prevPage() {
    if (currPage > 1) {
      currPage--;
    }
  }

  const loaderActive = ref(false);

  async function initLoaderOnDelay() {
    await new Promise(r => setTimeout(r, 2500));
    loaderActive.value = true;
  }

  onMounted(fetchData);
  watch(lang, async () => await fetchData());
  watch(currPage, async () => await fetchData());

  initLoaderOnDelay();
</script>

<template>
  <main>
    <section v-show="events.length !== 0 && events[0] === 'loading'" id=loading>
      <div v-show="loaderActive" class="initial-loader"></div>
    </section>
    <section id="no-events" v-show="events.length === 0">
      <h2>{{ dictionary.UINoEvents[lang] }}</h2>
      <p>{{ dictionary.UINoEventsSubTxt[lang] }}</p>
    </section>
    <form v-show="events.length > 0 && events[0] !== 'loading'" method="post" @submit.prevent="" >
      <div>
        <input v-model="search" type="search" id="search" name="search" autocomplete="search" required :placeholder="dictionary.FormSearch[lang]">
      </div>
    </form>
    <Event v-if="events.length !== 0 && events[0] !== 'loading'" v-for="event in filteredEvents" :event="event" />
    <div id="pager" v-show="events.length > 0 && events[0] !== 'loading'">
      <div class="prev pagerNavBtn" :class="{ disabled: currPage === 1 }" @click="prevPage"></div>
      <p>{{ currPage }} / {{ lastPage }}</p>
      <div class="next pagerNavBtn" :class="{ disabled: currPage === lastPage}" @click="nextPage"></div>
    </div>
  </main>
</template>

<style scoped>
  main {
    position: relative;
    flex: 1 1 auto;
    height: 100%;
    overflow-y: auto;
    padding-bottom: 1rem;
  }

  section {
    height: 70vh;
    max-height: 40rem;
    max-width: 100rem;
    margin: 0 auto;
  }

  #loading {
    display: flex;
    justify-content: center;
    align-items: center;
  }

  #pager {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 2rem;
    margin-top: 1rem;
  }

  #pager p {
    text-align: center;
    font-size: 1rem;
    font-weight: bold;
    letter-spacing: 0.5rem;
  }

  #pager .pagerNavBtn {
    width: 2rem;
    height: 2rem;
    background-color: var(--color-text);
    cursor: pointer;
  }

  #pager .pagerNavBtn.next {
    mask: url(../assets/media/chevron-right.svg);
    mask-size: contain;
  }

  #pager .pagerNavBtn.prev {
    mask: url(../assets/media/chevron-left.svg);
    mask-size: contain;
  }

  #pager .pagerNavBtn.disabled {
    filter: opacity(30%);
    cursor: auto;
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

  #no-events {
    display: flex;
    flex-direction: column;
    justify-content: center;
    height: 100%;
    padding: 0 1rem;
  }

  #no-events h2 {
    text-align: center;
    font-weight: bold;
    font-size: 1.2rem;
  }

  #no-events p {
    text-align: center;
  }

  #no-events a {
    color: var(--color-quaternary);
    text-decoration: none;
  }

  #no-events a:hover {
    text-decoration: underline;
    cursor: pointer;
  }

  @media screen and (max-width: 1000px) {
    form {
      justify-content: center;
    }
  }
</style>