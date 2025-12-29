<script setup>
    import { computed } from 'vue';
    import { StaticContentProvider as SCP } from '@/services/StaticContentService';
    import LanguageProvider from '@/services/LanguageService';
    import config from '../data/config.js';

    const emits = defineEmits(['close']);
    const language = computed(() => LanguageProvider.CURR_LANG.value);

    function getFlagUrl(lang) {
        // Check if API returned a flag_url
        if (lang.flag_url) {
            // Local asset paths - use directly
            if (lang.flag_url.includes('/assets/') || lang.flag_url.includes('fallbackLangFlags')) {
                return lang.flag_url;
            }
            // Otherwise it's an uploaded image
            return `${config.apiUrl}images/${lang.flag_url}`;
        }
        // Check if fallback data has imageUrl  
        if (lang.imageUrl) {
            return `/assets/media/${lang.imageUrl}`;
        }
        // Default: use local flag based on language code
        return `/assets/media/fallbackLangFlags/${lang.code}.svg`;
    }
</script>

<template>
    <section>
        <div @click="$emit('close')" class="cross-btn"></div>
        <ul>
            <li v-for="lang in SCP.LANGUAGES" :key="lang.code" :class="{selected: lang.code === language}">
                <figure @click="LanguageProvider.CURR_LANG = lang.code; $emit('close')">
                    <img :src="getFlagUrl(lang)">
                    <p>{{ lang.name }}</p>
                </figure>
            </li>
        </ul>
    </section>
</template>

<style scoped>
    section {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 0.5rem;
    }

    .cross-btn {
        width: 1.5rem;
        aspect-ratio: 1 / 1;
        background-color: var(--color-text);
        mask: url(@/assets/media/x.svg);
        mask-size: contain;
        cursor: pointer;
        mask-repeat: no-repeat;
        mask-position: center;
        align-self: flex-end;
        margin: 0.3rem 0.3rem 0 0;
    }

    figure img {
        height: 3rem;
    }

    ul {
        list-style: none;
        margin: 0 1rem 1rem 1rem;
        padding: 0;
        display: flex;
        flex-wrap: wrap;
        justify-content: center;
        gap: 1rem;
        width: 100%;
    }

    li {
        box-sizing: border-box;
        padding: 0.5rem;
        border-radius: 0.2rem;
        text-align: center
    }

    li p {
        border-bottom: 0.2rem solid rgba(0, 0, 0, 0);
    }

    li.selected p {
        font-weight: bold;
        border-bottom: 0.2rem solid var(--color-text);
    }

    li:hover {
        cursor: pointer;
    }

    li:hover p {
        font-weight: bold;
    }
</style>