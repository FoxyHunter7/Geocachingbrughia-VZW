<script setup>
import { onMounted, ref, computed } from 'vue';
  import { RouterView } from 'vue-router'
  import TopHeader from '@/components/TopHeader.vue'

  const innerWidth = ref(window.innerWidth);
  const isSideMenuOpen = ref(false);
  const isMobile = computed(() => {
    return innerWidth.value <= 850;
  });

  onMounted(() => {
        window.addEventListener('resize', () => {innerWidth.value = window.innerWidth});
  });
</script>

<template>
  <section id="side-menu" :class="{ open: isSideMenuOpen }"></section>
  <TopHeader :isMobile="isMobile" @menu-state-change="(state) => {isSideMenuOpen = state}"/>
  <RouterView />
</template>

<style scoped>
  Header {
    height: 6rem;
    box-sizing: border-box;
    max-width: 100dvw;
    overflow: hidden;
  }

  #side-menu {
    background-color: var(--color-secondary);
    position: absolute;
    top: 6rem;
    right: 0;
    width: 0;
    max-width: 80dvw;
    overflow-x: hidden;
    height: calc(100lvh - 6rem);
    transition: width ease 0.3s;
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 3rem;
  }
  
  #side-menu.open {
    width: 15rem;
    transition: width ease 0.3s;
  }
  
  #side-menu p {
    height: 1.5rem;
    text-overflow: ellipsis;
    overflow: hidden;
  }

  #side-menu.open p {
    height: fit-content;
    overflow: visible;
    text-overflow: default;
  }
</style>