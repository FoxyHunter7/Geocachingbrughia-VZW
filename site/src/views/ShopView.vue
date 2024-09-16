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
    <p>{{ products }}</p>
  </main>
</template>

<style scoped>
</style>