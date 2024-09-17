<script setup>
  import { getAllProducts } from '@/services/ShopService';
  import LanguageProvider from '@/services/LanguageService';
  import StaticContentProvider from '@/services/StaticContentService';
  import { computed, onMounted, ref, watch } from 'vue';

  const lang = computed(() => LanguageProvider.CURR_LANG.value);
  const dictionary = StaticContentProvider.DICTIONARY;

  const products = ref(["loading"]);

  async function fetchData() {
    const response = await getAllProducts();
    if (response.data) {
      products.value = response.data;
    } else {
      products.value = response;
    }
  }

  const search = ref("");

  const filteredProducts = computed(() => {
    if (!search.value) {
      return products.value;
    }

    return products.value.filter(product =>
      product.name.toLowerCase().includes(search.value.toLocaleLowerCase())
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
    <section v-show="products.length !== 0 && products[0] === 'loading'" id="loading">
      <div v-show="loaderActive" class="initial-loader"></div>
    </section>
    <section id="no-products" v-show="products.length === 0">
      <h2>{{ dictionary.UINoStoreItems[lang] }}</h2>
      <p>{{ dictionary.UINoStoreItemsSubTxt[lang] }}</p>
    </section>
    <form v-show="products.length > 0 && products[0] !== 'loading'" method="post" @submit.prevent="">
      <div>
        <input v-model="search" type="search" id="search" name="search" autocomplete="search" required :placeholder="dictionary.FormSearch[lang]">
      </div>
    </form>
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
    background-color: var(--color-background);
    color: var(--color-text)
  }

  #no-products {
    display: flex;
    flex-direction: column;
    justify-content: center;
    height: 100%;
    padding: 0 1rem;
  }

  #no-products h2 {
    text-align: center;
    font-weight: bold;
    font-size: 1.2rem;
  }

  #no-products p {
    text-align: center;
  }

  #no-products a {
    color: var(--color-quaternary);
    text-decoration: none;
  }

  #no-products a:hover {
    text-decoration: underline;
    cursor: pointer;
  }

  @media screen and (max-width: 1000px) {
    form {
      justify-content: center;
    }
  }
</style>