<script setup>
    import { computed } from 'vue';
    import LanguageProvider from '@/services/LanguageService.js';
    import StaticContentProvider from '@/services/StaticContentService.js';
    import config from '../data/config.js';

    const props = defineProps({
        socials: Array
    });

    const lang = computed(() => LanguageProvider.CURR_LANG.value);

    function openLink(url) {
        window.open(url, '_blank');
    }
</script>

<template>
    <section>
        <h2>{{ StaticContentProvider.DICTIONARY.SocialsFollowTxt[lang] }}</h2>
        <div>
            <figure v-for="social in socials" @click="openLink(social.url)">
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
        font-weight: bold;
        margin-bottom: 3rem;
    }

    h2::first-letter {
        text-transform: uppercase;
    }

    figure {
        padding: 1.5rem;
        border-radius: 1rem;
        height: 10rem;
        width: 10rem;
        background-color: var(--color-accent-light);
        transform: scale(100%);
        transition: transform 0.15s;
        display: flex;
        flex-direction: column;
        justify-content: space-between;
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
        user-select: none;
        pointer-events: none;
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

    @media screen and (max-width: 50rem) {
        figure {
            height: 8rem;
            width: 8rem;
            padding: 1.5rem 1rem;
        }

        section {
            padding: 1rem;
        }

        section div {
            gap: 2rem;
        }
    }

    @media screen and (max-width: 30rem) {
        figure {
            height: 6rem;
            width: 6rem;
            padding: 1rem 0.5rem;
        }
    }
</style>