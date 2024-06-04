<script setup>
  import { RouterView } from 'vue-router'
  import { watch } from 'vue'
  import { ref } from 'vue'
  import { StaticContentProvider } from './services/StaticContentService';

  const staticInitComplete = ref(StaticContentProvider.INIT_COMPLETE)

  watch(() => StaticContentProvider.INIT_COMPLETE, (newStaticInitComplete) => {
    staticInitComplete.value = newStaticInitComplete;
  })
</script>

<template>
  <p>Hello</p>
  <p>{{ staticInitComplete }}</p>
  <template v-if="staticInitComplete">
    <!--<RouterView />-->
  </template>
  <template v-if="!staticInitComplete">
    <p>Loading page content</p>
  </template>
</template>

<style scoped>
</style>
