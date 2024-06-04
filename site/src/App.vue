<script setup>
  import { RouterView } from 'vue-router'
  import { watch } from 'vue'
  import { ref } from 'vue'
  import { StaticContentProvider } from './services/StaticContentService';

  const currLang = ref(StaticContentProvider.DICTIONARY)

  watch(() => StaticContentProvider.DICTIONARY, (newLang) => {
  currLang.value = newLang;
})
</script>

<template>
  <p>{{ currLang }}</p>
  <RouterView />
</template>

<style scoped>
</style>
