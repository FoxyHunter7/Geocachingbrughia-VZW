<script setup>
    import { computed } from 'vue';
    import LanguageProvider from '@/services/LanguageService.js';
    import StaticContentProvider from '@/services/StaticContentService.js';
    import config from '../data/config.json';

    const props = defineProps({
        socials: Array
    });

    const lang = computed(() => LanguageProvider.CURR_LANG.value);
</script>

<template>
    <section>
        <h2>{{ StaticContentProvider.DICTIONARY.SocialsFollowTxt[lang] }}</h2>
        <div>
            <figure v-for="social in socials">
                <img :src="`${config.apiUrl}images/${social.imageUrl}`">
                <p><a :href="social.url" target="_blank">{{ social.name }}</a></p>
            </figure>
        </div>
    </section>
</template>

<style scoped>
    section {
        height: auto !important;
        padding: 3rem;
        background-color: var(--color-background2)
    }

    section div {
        display: flex;
        gap: 3rem;
        justify-content: space-evenly;
        flex-wrap: wrap;
    }

    h2 {
        text-align: center;
        text-transform: uppercase;
        margin-bottom: 3rem;
        font-weight: bold;
    }

    figure {
        padding: 1.5rem;
        border-radius: 0.5rem;
        height: 10rem;
        background-color: var(--color-tertiary);
        transform: scale(100%);
        transition: transform 0.3s;
    }

    figure:hover {
        cursor: pointer;
        transform: scale(105%);
        transition: transform 0.3s;
        text-decoration: underline;
        color: var(--color-text2);
    }

    figure img {
        max-height: 70%;
        margin-bottom: 0.5rem;
    }

    figure p {
        text-align: center
    }

    figure a {
        text-decoration: none;
        text-transform: capitalize;
        color: var(--color-text2);
    }
</style>