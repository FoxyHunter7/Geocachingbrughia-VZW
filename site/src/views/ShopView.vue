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
    await new Promise(r => setTimeout(r, 0));
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
    <section id="store-items">
      <article v-if="products.length > 0 && products[0] !== 'loading'" v-for="product in filteredProducts">
        <div>
          <img :src="product.images[0]" data-cookiescript="denied" data-category="functionality">
        </div>
        <p>{{ product.name }} <span>{{ product.price.amount / 100 }} {{ product.price.currency.toUpperCase() }}</span></p>
      </article>
    </section>
  </main>
</template>

<style scoped>
  #loading {
    height: 80vh;
    max-height: 40rem;
    max-width: 100rem;
    margin: 0 auto;
  }

  #store-items {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(15rem, 20rem));
    justify-items: center;
    align-items: center;
    gap: 3rem;
    padding: 3rem;
  }

  #store-items article {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-between;
    max-width: 85vw;
    width: 16rem;
    padding: 0.5rem;
    border-radius: 0.5rem;
    scale: 100%;
    transition: scale 0.1s;
  }

  #store-items article:hover {
    cursor: pointer;
    scale: 101%;
    transition: scale 0.1s;
  }

  #store-items article div {
    height: 15rem;
    width: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
    border-radius: 0.5rem;
  }

  #store-items article img {
    object-fit: contain;
    object-position: top;
    border-radius: 0.5rem;
    pointer-events: none;
    user-select: none;
    max-height: 100%;
    max-width: 100%;
  }

  #store-items article p {
    text-align: center;
    font-weight: bold;
    display: flex;
    justify-content: space-between;
    width: 100%;
    margin-top: 1.5rem;
  }

  #store-items article p span {
    background-color: var(--color-accent-light);
    padding: 0.1rem 0.5rem;
    border-radius: 0.4rem;
  }

  main {
    position: relative;
    flex: 1 1 auto;
    height: 100%;
    overflow-y: auto;
    padding-bottom: 1rem;
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

  @media screen and (max-width: 1000px) {
    form {
      justify-content: center;
    }
  }
</style>