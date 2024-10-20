<script setup>
  import { getAllProducts } from '@/services/ShopService';
  import LanguageProvider from '@/services/LanguageService';
  import StaticContentProvider from '@/services/StaticContentService';
  import { computed, onMounted, ref, watch } from 'vue';

  const lang = computed(() => LanguageProvider.CURR_LANG.value);
  const dictionary = StaticContentProvider.DICTIONARY;

  const products = ref(["loading"]);
  const detailsOf = ref(-1);

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
    <form v-show="products.length > 0 && products[0] !== 'loading' && detailsOf === -1" method="post" @submit.prevent="">
      <div>
        <input v-model="search" type="search" id="search" name="search" autocomplete="search" required :placeholder="dictionary.FormSearch[lang]">
      </div>
    </form>
    <section id="store-items" v-show="detailsOf === -1">
      <article v-if="products.length > 0 && products[0] !== 'loading'" v-for="(product, index) in filteredProducts" @click="detailsOf = index">
        <div>
          <img :src="product.images[0]" data-cookiescript="denied" data-category="functionality">
        </div>
        <p>{{ product.name }} <span>{{ product.price.amount / 100 }} {{ product.price.currency.toUpperCase() }}</span></p>
      </article>
    </section>
    <section v-if="detailsOf !== -1" id="product-details">
      <div>
        <img :src="filteredProducts[detailsOf].images[0]">
      </div>
      <div>
        <h2>{{ filteredProducts[detailsOf].name }} <span>{{ filteredProducts[detailsOf].price.amount / 100 }} {{ filteredProducts[detailsOf].price.currency.toUpperCase() }}</span></h2>
        <p>{{ filteredProducts[detailsOf].description }}</p>
      </div>
    </section>
  </main>
</template>

<style scoped>
  #product-details {
    display: grid;
    grid-template-columns: 1fr 2fr;
    align-items: center;
    width: 70%;
    max-width: 80rem;
    gap: 5rem;
    margin: auto;
    height: 100%;
  }

  #product-details h2 {
    font-weight: bold;
  }

  #product-details h2 span {
    background-color: var(--color-accent-light);
    padding: 0.1rem 0.5rem;
    border-radius: 0.4rem;
    font-size: 1rem;
    margin-left: 1rem;
    transform: translateY(-0.1rem);
    display: inline-block;
  }

  #product-details div {
    height: 30rem;
  }

  #product-details div:first-child {
    display: flex;
    justify-content: end;
  }

  #product-details div img {
    max-height: 80%;
    max-width: 100%;
    object-fit: contain;
    object-position: top;
    border-radius: 0.5rem;
  }

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