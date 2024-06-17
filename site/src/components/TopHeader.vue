<script setup>
    import { ref } from 'vue';
    import { RouterLink } from 'vue-router';
    import LanguageProvider from '@/services/LanguageService';
    import { StaticContentProvider as SCP } from '@/services/StaticContentService';

    const props = defineProps(["isMobile"]);
    const emits = defineEmits(['menuStateChange']);

    const lang = ref(LanguageProvider.CURR_LANG);
    const isNavOpen = ref(false);

    const dictionary = SCP.DICTIONARY;
</script>

<template>
    <header>
        <img v-if="!isMobile" src="../assets/media/logo-full-black.webp" class="logo">
        <img v-if="isMobile" src="../assets/media/logo-head-black.webp" class="logo">
        <Teleport to="#side-menu" :disabled="!isMobile">
            <nav>
                <RouterLink @click="if (isMobile) {isNavOpen = !isNavOpen; $emit('menuStateChange', isNavOpen)}" :to="{ path: SCP.constructRoute(lang, 'NavHome') }">{{ dictionary.NavHome[lang] }}</RouterLink>
                <RouterLink @click="if (isMobile) {isNavOpen = !isNavOpen; $emit('menuStateChange', isNavOpen)}" :to="{ path: SCP.constructRoute(lang, 'NavEvents') }">{{ dictionary.NavEvents[lang] }}</RouterLink>
                <RouterLink @click="if (isMobile) {isNavOpen = !isNavOpen; $emit('menuStateChange', isNavOpen)}" :to="{ path: SCP.constructRoute(lang, 'NavGeocaches') }">{{ dictionary.NavGeocaches[lang] }}</RouterLink>
                <RouterLink @click="if (isMobile) {isNavOpen = !isNavOpen; $emit('menuStateChange', isNavOpen)}" :to="{ path: SCP.constructRoute(lang, 'NavShop') }">{{ dictionary.NavShop[lang] }}</RouterLink>
            </nav>
        </Teleport>
        <div v-if="isMobile" class="nav" :class="{open: isNavOpen}" @click="isNavOpen = !isNavOpen; $emit('menuStateChange', isNavOpen)"></div>
    </header>
</template>

<style scoped>
    header {
        background-color: var(--color-primary);
        display: grid;
        grid-template-columns: 2fr 6fr;
        gap: 5rem;
        padding: 1rem;
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
        margin: 2rem 2.5rem;
    }

    #side-menu > nav a {
        color: var(--color-text);
    }

    #side-menu > nav a.router-link-active {
        border-bottom: 0.2rem solid var(--color-text);
    }

    header > .nav {
        margin: auto 0 auto auto;
        height: 80%;
        aspect-ratio: 1 / 1;
        background-color: var(--color-text2);
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
        color: var(--color-text2);
        border-bottom: 0.2rem solid rgba(0, 0, 0, 0);
        transition: border-bottom, font-weight 0.1s;
    }

    nav a:hover {
        border-bottom: 0.2rem solid var(--color-secondary);
        transition: border-bottom, font-weight 0.1s;
    }

    nav a.router-link-active {
        border-bottom: 0.2rem solid var(--color-text2);
        transition: border-bottom, font-weight 0.1s;
        font-weight: bold;
    }
</style>