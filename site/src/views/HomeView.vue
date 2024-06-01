<script setup>
import { useRoute } from 'vue-router'
import { watch } from 'vue'
import { ref } from 'vue'
import LanguageSelector from '../components/LanguageSelector.vue'
import Events from '../components/Events.vue'
//import uiTranslations from '../data/uiTranslations.json'
import { getLanguageFromPath } from '../services/LanguageService.js'

const route = useRoute()
const lang = ref(getLanguageFromPath(route.path))

watch(() => route.path, (newPath) => {
  lang.value = getLanguageFromPath(newPath).data
})

components: {
  LanguageSelector
  Events
}
</script>

<template>
  <main>
    <h1> {{uiTranslations[lang].tagline}} </h1>
    <p>Help it's 21:30 I'm running out of time aaaaaaaaaaaaa</p>
    <language-selector :lang="lang" :pathName="route.name" />
    <Suspense>
      <template #default>
        <events :lang="lang" />
      </template>
      <template #fallback>
        <p>loading events...</p>
      </template>
    </Suspense>
  </main>
</template>
