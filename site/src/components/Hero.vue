<script setup>
    import { computed } from 'vue';
    import LanguageProvider from '@/services/LanguageService';
    import StaticContentProvider from '@/services/StaticContentService';
    import config from '../data/config.json';

    const lang = computed(() => LanguageProvider.CURR_LANG.value);
    const dictOnFallback = computed(() => StaticContentProvider.DICT_ON_FALLBACK.value)
    const dictionary = StaticContentProvider.DICTIONARY;
</script>

<template>
    <section>
        <div>
            <img :src="(dictOnFallback) ? `/src/assets/media/static/bertje.jpg` : `${config.apiUrl}images/${dictionary.SplashImg[lang]}`">
        </div>
        <div>
            <h2>{{ dictionary.SplashTitle[lang] }}</h2>
            <p>{{ dictionary.SplashBody[lang] }}</p>
        </div>
    </section>
</template>

<style scoped>
    section {
        display: grid;
        grid-template-columns: 1fr 2fr;
        gap: 3rem;
        padding: 3rem;
        max-height: 40rem !important;
    }

    section > div:first-child {
        height: 100%;
        max-height: 40rem;
        display: flex;
        justify-content: center;
        align-items: flex-start;
        overflow: hidden;
    }

    section > div > img {
        max-width: 100%;
        max-height: 100%;
        object-fit: contain;
        border-radius: 0.5rem;
    }

    h2 {
        text-transform: uppercase;
    }

    p {
        margin-top: 1rem;
        font-size: 1rem;
        line-height: 170%;
        max-width: 60rem;
    }

    @media screen and (max-width: 800px) {
        section {
            height: fit-content !important;
            max-height: 80rem !important;
            grid-template-columns: 1fr;
            grid-template-rows: 1fr 1fr;
            justify-content: center;
            align-items: center;
        }

        h2 {
            margin-top: auto;
            text-align: center;
        }

        p {
            text-align: center;
        }

        section > div:first-child {
            align-items: center;
        }
    }
</style>