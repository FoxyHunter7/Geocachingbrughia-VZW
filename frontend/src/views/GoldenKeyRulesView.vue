<script setup>
import { ref, computed, onMounted } from 'vue';
import { RouterLink } from 'vue-router';
import { getGoldenKeySettings } from '@/services/GoldenKeyService';
import LanguageProvider from '@/services/LanguageService';
import StaticContentProvider from '@/services/StaticContentService';
import { generateHTML } from '@tiptap/vue-3';
import Document from '@tiptap/extension-document';
import Paragraph from '@tiptap/extension-paragraph';
import Text from '@tiptap/extension-text';
import Bold from '@tiptap/extension-bold';
import Italic from '@tiptap/extension-italic';
import Underline from '@tiptap/extension-underline';
import Link from '@tiptap/extension-link';
import BulletList from '@tiptap/extension-bullet-list';
import ListItem from '@tiptap/extension-list-item';
import Strike from '@tiptap/extension-strike';
import Table from '@tiptap/extension-table';
import TableCell from '@tiptap/extension-table-cell';
import TableHeader from '@tiptap/extension-table-header';
import TableRow from '@tiptap/extension-table-row';

const tiptapExtensions = [
    Document, Paragraph, Text, Bold, Italic, Underline, Link,
    BulletList, ListItem, Strike, Table, TableRow, TableHeader, TableCell
];

function renderRules(content) {
    if (!content) return '';
    try {
        return generateHTML(JSON.parse(content), tiptapExtensions);
    } catch {
        return content;
    }
}

const lang = computed(() => LanguageProvider.CURR_LANG.value);
const dictionary = StaticContentProvider.DICTIONARY;

const rulesText = ref('');
const loading = ref(true);

onMounted(async () => {
    try {
        const settings = await getGoldenKeySettings();
        const rules = settings?.rules ?? {};
        rulesText.value = rules[lang.value] || rules['NL'] || '';
    } catch {
        rulesText.value = '';
    }
    loading.value = false;
});
</script>

<template>
    <main class="gkr">
        <RouterLink to="/golden-key" class="gkr__back">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"
                stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                <path d="M19 12H5M12 5l-7 7 7 7"/>
            </svg>
            {{ dictionary.ButtonBack?.[lang] ?? 'Terug' }}
        </RouterLink>

        <div v-if="loading" class="gkr__loading">Laden&#8230;</div>

        <template v-else>
            <header class="gkr__header">
                <p class="gkr__header-sub">The Golden Key</p>
                <h1 class="gkr__title">{{ dictionary.GoldenKeyRulesBtn?.[lang] ?? 'Spelregels' }}</h1>
            </header>

            <div v-if="rulesText" class="gkr__content" v-html="renderRules(rulesText)"></div>
            <p v-else class="gkr__empty">Geen spelregels beschikbaar.</p>
        </template>
    </main>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Cinzel+Decorative:wght@700&display=swap');

.gkr {
    flex: 1 1 auto;
    display: flex;
    flex-direction: column;
    align-items: center;
    min-height: calc(100vh - 4.5rem);
    overflow-y: auto;
    padding: clamp(2rem, 4vw, 3.5rem) clamp(1rem, 4vw, 3rem) 5rem;
    gap: clamp(1.5rem, 3vw, 2.5rem);

    background-color: #050505;
    background-image:
        radial-gradient(ellipse 90% 25% at 50% 0%, rgba(160, 110, 8, 0.10) 0%, transparent 100%),
        radial-gradient(ellipse 35% 70% at   0% 45%, rgba(110, 70, 5, 0.05) 0%, transparent 100%),
        radial-gradient(ellipse 35% 70% at 100% 55%, rgba(110, 70, 5, 0.05) 0%, transparent 100%),
        repeating-linear-gradient(
            128deg,
            transparent 0px, transparent 4px,
            rgba(190, 140, 20, 0.012) 4px, rgba(190, 140, 20, 0.012) 5px
        ),
        repeating-linear-gradient(
            52deg,
            transparent 0px, transparent 6px,
            rgba(190, 140, 20, 0.008) 6px, rgba(190, 140, 20, 0.008) 7px
        );
}

.gkr__back {
    align-self: flex-start;
    display: inline-flex;
    align-items: center;
    gap: 0.4em;
    color: #9a7230;
    text-decoration: none;
    font-size: 0.85rem;
    font-weight: 600;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    opacity: 0.75;
    transition: opacity 0.2s, color 0.2s;
}

.gkr__back:hover {
    opacity: 1;
    color: #f5d87a;
}

.gkr__loading {
    color: #9a7230;
    font-size: 0.95rem;
    margin-top: 3rem;
}

.gkr__header {
    text-align: center;
    font-family: 'Cinzel Decorative', serif;
}

.gkr__header-sub {
    font-size: clamp(0.75rem, 2vw, 1rem);
    font-weight: 700;
    letter-spacing: 0.3em;
    text-transform: uppercase;
    color: #9a7230;
    margin: 0 0 0.4rem;
}

.gkr__title {
    font-size: clamp(1.6rem, 5vw, 3rem);
    font-weight: 700;
    margin: 0;
    background: linear-gradient(160deg, #f5d87a 0%, #c8860a 45%, #f5d87a 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    letter-spacing: 0.03em;
    line-height: 1.1;
}

.gkr__content {
    width: 100%;
    max-width: 680px;
    color: #d4b06a;
    font-size: clamp(0.9rem, 2.2vw, 1.05rem);
    line-height: 1.75;
    word-break: break-word;
}

.gkr__content :deep(p) {
    margin: 0.5em 0;
}

.gkr__empty {
    color: #5a4220;
    font-size: 0.95rem;
    font-style: italic;
}
</style>
