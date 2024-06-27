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
        border-radius: 1rem;
        height: 12rem;
        width: 12rem;
        background-color: var(--color-secondary);
        transform: scale(100%);
        transition: transform 0.15s;
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        box-shadow: var(--color-background2) 1rem 1rem;
    }

    figure:hover {
        cursor: pointer;
        transform: scale(105%);
        transition: transform 0.25s;
        text-decoration: underline;
        color: var(--color-text);
    }

    figure img {
        object-fit: contain;
        max-height: 50%;
        margin-bottom: 0.5rem;
    }

    figure p {
        text-align: center
    }

    figure a {
        font-weight: bold;
        text-decoration: none;
        text-transform: capitalize;
        color: var(--color-text);
    }
</style>