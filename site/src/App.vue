<script setup>
  import { RouterView } from 'vue-router'
  import TopHeader from '@/components/TopHeader.vue'
</script>

<template>
  <TopHeader />
  <RouterView />
</template>

<style scoped>
  Header {
    height: 6rem;
  }
</style>