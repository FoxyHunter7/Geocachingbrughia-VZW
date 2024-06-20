<script setup>
  import { onMounted, ref, computed } from 'vue';
  import { RouterView } from 'vue-router';
  import { StaticContentProvider as SCP } from '@/services/StaticContentService';
  import config from '@/data/config.json'
  import TopHeader from '@/components/TopHeader.vue'
  import WarningBanner from './components/WarningBanner.vue';
  import LanguageSelector from './components/LanguageSelector.vue';

  const scsErrors = SCP.ERRORS;

  const innerWidth = ref(window.innerWidth);
  const isMobile = computed(() => {
    return innerWidth.value <= 850;
  });

  const isSideMenuOpen = ref(false);

  const keepPopupMenuOpen = ref(false);
  const popupMenuStates = ref({
    "languageSelector": false
  });
  const isPopupOpen = computed(() => {
    return keepPopupMenuOpen.value || Object.values(popupMenuStates.value).some(state => state === true);
  });
  const popupMenuQueue = [];

  function openInPopup(componentName) {
    console.log(isPopupOpen.value);
    if (isPopupOpen.value) {
      popupMenuQueue.push(componentName);
    } else {
      popupMenuStates.value[componentName] = true;
    }
  }

  function closeOutPopup(componentName = "") {
    if (isPopupOpen.value) {
      if (popupMenuQueue.length > 0) {
        keepPopupMenuOpen.value = true;
        clearFromPopup(componentName);
        popupMenuStates.value[popupMenuQueue.shift()] = true;
        keepPopupMenuOpen.value = false;
      } else {
        clearFromPopup(componentName);
      }
    }

    function clearFromPopup(componentName) {
      if (componentName) {
        popupMenuStates.value[componentName] = false;
      } else {
        Object.keys(popupMenuStates.value).forEach(componentName => popupMenuStates.value[componentName] = false);
      }
    }
  }

  onMounted(() => {
        window.addEventListener('resize', () => {innerWidth.value = window.innerWidth});
  });
</script>

<template>
  <section id="side-menu" :class="{ open: isSideMenuOpen }"></section>
  <TopHeader :isMobile="isMobile" @menu-state-change="(state) => { isSideMenuOpen = state }" @lang-selector="openInPopup('languageSelector')"/>
  <div id="messages">
    <WarningBanner v-if="scsErrors" :error="scsErrors" :date="config.fallbackLastUpdated"></WarningBanner>
  </div>
  <div id="popup-menu" v-show="isPopupOpen">
    <LanguageSelector v-show="popupMenuStates.languageSelector" @close="closeOutPopup('languageSelector')"/>
  </div>
  <div id="overlay" v-show="isPopupOpen" @click="closeOutPopup()"></div>
  <RouterView />
</template>

<style scoped>
  Header {
    height: 5rem;
    box-sizing: border-box;
    max-width: 100vw;
    overflow: hidden;
  }

  #side-menu {
    background-color: var(--color-secondary);
    position: absolute;
    top: 5rem;
    right: 0;
    width: 0;
    bottom: 0;
    max-width: 80dvw;
    overflow-x: hidden;
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

  #messages {
    position: relative;
    height: fit-content;
    top: 0;
    left: 0;
    padding: 1rem 0.8rem;
    z-index: -1;
  }

  #messages:empty {
    padding: 0;
  }

  #popup-menu {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
      align-items: center;
      background-color: var(--color-background);
      border-radius: 0.5rem;
      z-index: 20;
      max-width: 60vw;
  }

  #overlay {
    position: absolute;
    top: 0;
    right: 0;
    bottom: 0;
    left: 0;
    z-index: 10;
    background-color: rgba(0, 0, 0, 0.5)
  }
</style>