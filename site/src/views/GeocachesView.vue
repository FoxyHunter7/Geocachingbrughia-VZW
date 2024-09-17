<script setup>
  import { getAllGeocaches } from '@/services/GeocacheService';
  import LanguageProvider from '@/services/LanguageService';
  import StaticContentProvider from '@/services/StaticContentService';
  import { computed, onMounted, ref, watch } from 'vue';

  const lang = computed(() => LanguageProvider.CURR_LANG.value);
  const dictionary = StaticContentProvider.DICTIONARY;

  const geocaches = ref(["loading"]);
  const currPage = ref(1);
  const lastPage = ref(1);

  async function fetchData() {
    const pagedResponse = await getAllGeocaches(false, currPage.value);
    if (pagedResponse.data) {
      geocaches.value = pagedResponse.data;
      currPage.value = pagedResponse.current_page;
      lastPage.value = pagedResponse.last_page;
    } else {
      geocaches.value = pagedResponse;
    }
  }

  const search = ref("");

  const filteredGeocaches = computed(() => {
    if (!search.value) {
      return geocaches.value;
    }

    return geocaches.value.filter(geocache =>
      geocache.title.toLowerCase().includes(search.value.toLowerCase())
    );
  });

  function nextPage() {
    if (currPage.value < lastPage.value) {
      currPage.value++;
    }
  }

  function prevPage() {
    if (currPage.value > 1) {
      currPage.value--;
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

  function extractGeocode(url) {
    return url.split("/").pop();
  }
</script>

<template>
  <main>
    <section v-show="geocaches.length !== 0 && geocaches[0] === 'loading'" id=loading>
      <div v-show="loaderActive" class="initial-loader"></div>
    </section>
    <section id="no-geocaches" v-show="geocaches.length === 0">
      <h2>{{ dictionary.UINoGeocaches[lang] }}</h2>
      <p>{{ dictionary.UINoGeocachesSubTxt[lang] }}</p>
    </section>
    <form v-show="geocaches.length > 0 && geocaches[0] !== 'loading'" method="post" @submit.prevent="" >
      <div>
        <input v-model="search" type="search" id="search" name="search" autocomplete="search" required :placeholder="dictionary.FormSearch[lang]">
      </div>
    </form>
    <div id="geocaches" v-if="geocaches.length !== 0 && geocaches[0] !== 'loading'">
      <table>
      <thead>
        <tr>
          <th scope="col"> </th>
          <th scope="col">{{ dictionary.GeocacheTitle[lang] }}</th>
          <th scope="col">{{ dictionary.GeocacheDifficulty[lang] }}</th>
          <th scope="col">{{ dictionary.GeocacheTerrain[lang] }}</th>
          <th scope="col">{{ dictionary.GeocacheGeoLink[lang] }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="geocache in filteredGeocaches">
          <td><div><img :src="`/assets/media/cachetypes/${geocache.type}.png`"></div></td>
          <td>{{ geocache.title }}</td>
          <td>{{ geocache.difficulty }}</td>
          <td>{{ geocache.terrain }}</td>
          <td><a :href="geocache.geolink">{{ extractGeocode(geocache.geolink) }}</a></td>
        </tr>
      </tbody>
    </table>
    </div>
    <div id="pager" v-show="geocaches.length > 0 && geocaches[0] !== 'loading'">
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

  a {
    text-decoration: underline;
    color: var(--color-accent-dark);
  }

  a:hover {
    font-weight: bold;
    cursor: pointer;
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
    mask-repeat: no-repeat;
    mask-position: center;
  }

  #pager .pagerNavBtn.prev {
    mask: url(../assets/media/chevron-left.svg);
    mask-size: contain;
    mask-repeat: no-repeat;
    mask-position: center;
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
    mask-repeat: no-repeat;
    mask-position: center;
  }

  #no-geocaches {
    display: flex;
    flex-direction: column;
    justify-content: center;
    height: 100%;
    padding: 0 1rem;
  }

  #no-geocaches h2 {
    text-align: center;
    font-weight: bold;
    font-size: 1.2rem;
  }

  #no-geocaches p {
    text-align: center;
  }

  #no-geocaches a {
    color: var(--color-accent-dark);
    text-decoration: none;
  }

  #no-geocaches a:hover {
    text-decoration: underline;
    cursor: pointer;
  }

  #geocaches {
    overflow-x: auto;
  }

  table {
    min-width: 40rem;
    width: 70vw;
    max-width: 80rem;
    margin: 4rem auto;
    border: 0.1rem solid var(--color-text);
    border-collapse: collapse;
  }

  thead {
    background-color: var(--color-primary);
    border: 0.1rem solid var(--color-text);
    height: 1.75rem;
  }

  td > div {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 0.2rem;
  }

  th:first-child {
    width: 3rem;
  }

  th:last-child {
    width: 10rem;
  }

  th:nth-child(3) {
    width: 6rem;
  }

  th:nth-child(4) {
    width: 6rem;
  }

  th {
    color: var(--color-text);
    padding: 0 0.5rem;
    font-weight: bold;
    text-transform: capitalize;
  }

  th, td, tr {
    border: 0.1rem solid var(--color-text);
    text-align: center;
  }

  tr img {
    max-height: 1.8rem;
  }

  tbody tr:nth-child(even) {
    background-color: var(--color-background-2);
  }

  @media screen and (max-width: 1000px) {
    form {
      justify-content: center;
    }
  }
</style>