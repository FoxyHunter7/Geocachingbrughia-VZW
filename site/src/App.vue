<script setup>
  import { RouterView } from 'vue-router'
  import { watch } from 'vue'
  import { ref } from 'vue'
  import { LanguageProvider } from './services/LanguageService';

  const currLang = ref(LanguageProvider.DICTIONARY)

  watch(() => LanguageProvider.DICTIONARY, (newLang) => {
  currLang.value = newLang;
})
</script>

<template>
  <p>{{ currLang }}</p>
  <RouterView />
</template>

<style scoped>
</style>
