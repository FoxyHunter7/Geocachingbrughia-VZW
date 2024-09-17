<script setup>
    import { ref, computed } from 'vue';
    import { RouterLink } from 'vue-router';
    import LanguageProvider from '@/services/LanguageService';
    import config from '../data/config.js';
    import { StaticContentProvider as SCP } from '@/services/StaticContentService';

    const props = defineProps(["isMobile"]);
    const emits = defineEmits(['menuStateChange', 'langSelector']);

    const lang = computed(() => LanguageProvider.CURR_LANG.value);
    const isNavOpen = ref(false);

    const langInfo = computed(() => {
        return SCP.LANGUAGES.find(lang => lang.code === LanguageProvider.CURR_LANG.value);
    });

    const dictionary = SCP.DICTIONARY;
</script>

<template>
    <header>
        <img v-if="!isMobile" src="/assets/media/logo-full-black.webp" class="logo">
        <img v-if="isMobile" src="/assets/media/logo-head-black.webp" class="logo">
        <Teleport to="#side-menu" :disabled="!isMobile">
            <nav>
                <RouterLink @click="if (isMobile) {isNavOpen = !isNavOpen; $emit('menuStateChange', isNavOpen)}" :to="{ path: SCP.constructRoute(lang, 'NavHome') }">{{ dictionary.NavHome[lang] }}</RouterLink>
                <RouterLink @click="if (isMobile) {isNavOpen = !isNavOpen; $emit('menuStateChange', isNavOpen)}" :to="{ path: SCP.constructRoute(lang, 'NavEvents') }">{{ dictionary.NavEvents[lang] }}</RouterLink>
                <RouterLink @click="if (isMobile) {isNavOpen = !isNavOpen; $emit('menuStateChange', isNavOpen)}" :to="{ path: SCP.constructRoute(lang, 'NavGeocaches') }">{{ dictionary.NavGeocaches[lang] }}</RouterLink>
                <RouterLink @click="if (isMobile) {isNavOpen = !isNavOpen; $emit('menuStateChange', isNavOpen)}" :to="{ path: SCP.constructRoute(lang, 'NavShop') }">{{ dictionary.NavShop[lang] }}</RouterLink>
            </nav>
            <figure @click="$emit('langSelector')" id="lang-selector">
                <img :src="(langInfo.fallback) ? `/assets/media/${langInfo.imageUrl}` : `${config.apiUrl}images/${langInfo.imageUrl}`">
                <p>{{ langInfo.name }}</p>
            </figure>
        </Teleport>
        <div v-if="isMobile"></div>
        <div v-if="isMobile" class="nav" :class="{open: isNavOpen}" @click="isNavOpen = !isNavOpen; $emit('menuStateChange', isNavOpen)"></div>
    </header>
</template>

<style scoped>
    header {
        background-color: var(--color-primary);
        display: grid;
        grid-template-columns: 2fr 6fr 1.5fr;
        gap: 5rem;
        padding: 0.75rem 1rem;
    }

    header > img.logo {
        height: 3rem;
    }

    header > nav {
        display: flex;
        justify-content: flex-end;
        align-items: center;
        gap: 1.8rem;
    }

    #side-menu > nav {
        display: flex;
        flex-direction: column;
        align-items: flex-end;
        gap: 1.5rem;
        margin: 2rem 2.5rem 0 2.5rem;
    }

    #side-menu > nav a {
        color: var(--color-text);
    }

    #side-menu > nav a.router-link-active {
        border-bottom: 0.2rem solid var(--color-text);
    }

    header > .nav {
        margin: auto 0.5rem auto auto;
        height: 80%;
        aspect-ratio: 1 / 1;
        background-color: var(--color-text);
        mask: url(../assets/media/menu.svg);
        mask-size: contain;
    }

    header > .nav.open {
        mask: url(../assets/media/x.svg);
        mask-size: contain;
    }

    nav a {
        text-decoration: none;
        text-transform: capitalize;
        font-size: 1rem;
        font-weight: normal;
        color: var(--color-text);
        border-bottom: 0.2rem solid rgba(0, 0, 0, 0);
        transition: border-bottom, font-weight 0.1s;
    }

    nav a:hover {
        border-bottom: 0.2rem solid var(--color-accent-dark);
        transition: border-bottom, font-weight 0.1s;
    }

    nav a.router-link-active {
        border-bottom: 0.2rem solid var(--color-text);
        transition: border-bottom, font-weight 0.1s;
        font-weight: bold;
    }

    figure {
        height: 3rem;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: space-between;
        width: min-content;
        margin-left: auto;
        padding: 0.2rem;
        border-radius: 0.5rem;
        transform: scale(100%);
        transition: transform 0.1s;
    }

    figure:hover {
        transform: scale(102%);
        transition: transform 0.1s;
        cursor: pointer;
    }

    figure img {
        height: 50%;
    }

    figure p {
        color: var(--color-text);
        opacity: 75%;
        line-height: 100%;
    }

    #side-menu > figure {
        flex-direction: row;
        justify-content: flex-start;
        height: 3.5rem;
        gap: 1rem;
        margin: 0 2.5rem 2rem auto;
    }

    #side-menu > figure p {
        color: var(--color-text);
        font-weight: bold;
    }
</style>