<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, RouterLink } from 'vue-router';
import { getGoldenKeyMonth } from '@/services/GoldenKeyMonthService';
import config from '@/data/config.js';

// TipTap rendering
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

function renderHint(content) {
    if (!content) return '';
    try {
        return generateHTML(JSON.parse(content), tiptapExtensions);
    } catch {
        return content;
    }
}

const route = useRoute();

const month = ref(null);
const locked = ref(false);
const loading = ref(true);

onMounted(async () => {
    const data = await getGoldenKeyMonth(route.params.id);
    if (!data || data.access_denied || (data.error === 'locked')) {
        locked.value = true;
    } else {
        month.value = data;
    }
    loading.value = false;
});
</script>

<template>
    <main class="gkm">
        <RouterLink to="/golden-key" class="gkm__back">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"
                stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                <path d="M19 12H5M12 5l-7 7 7 7"/>
            </svg>
            Terug naar overzicht
        </RouterLink>

        <!-- Loading -->
        <div v-if="loading" class="gkm__loading">
            <span>Laden&#8230;</span>
        </div>

        <!-- Locked -->
        <div v-else-if="locked" class="gkm__locked">
            <div class="gkm__lock-icon">&#128274;</div>
            <p class="gkm__lock-msg">Deze maand is nog niet beschikbaar.</p>
        </div>

        <!-- Content -->
        <template v-else-if="month">
            <header class="gkm__header">
                <p class="gkm__header-sub">The Golden Key</p>
                <h1 class="gkm__title">{{ month.month_name }}</h1>
            </header>

            <!-- Finder section — only shown when found -->
            <section v-if="month.state === 'found'" class="gkm__finder">
                <div class="gkm__finder-label">&#127942; Gevonden door</div>
                <img
                    v-if="month.finder_image"
                    :src="`${config.apiUrl}images/${month.finder_image}`"
                    :alt="month.finder_name"
                    class="gkm__finder-img"
                    loading="lazy"
                />
                <p class="gkm__finder-name">{{ month.finder_name }}</p>
                <p v-if="month.found_date" class="gkm__finder-date">
                    {{ new Date(month.found_date).toLocaleDateString('nl-BE', { day: 'numeric', month: 'long', year: 'numeric' }) }}
                </p>
            </section>

            <!-- Hints list — always shown for active / found months -->
            <section class="gkm__hints" v-if="month.hints && month.hints.length">
                <h2 class="gkm__hints-title">Hints</h2>
                <div
                    v-for="(hint, idx) in month.hints"
                    :key="hint.id"
                    class="gkm__hint"
                >
                    <div class="gkm__hint-num">{{ idx + 1 }}</div>
                    <div class="gkm__hint-body">
                        <div
                            class="gkm__hint-content"
                            v-html="renderHint(hint.content)"
                        ></div>
                        <img
                            v-if="hint.image_url"
                            :src="`${config.apiUrl}images/${hint.image_url}`"
                            alt=""
                            class="gkm__hint-img"
                            loading="lazy"
                        />
                    </div>
                </div>
            </section>

            <p v-else class="gkm__no-hints">Geen hints beschikbaar.</p>
        </template>
    </main>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Cinzel+Decorative:wght@700&display=swap');

/* Page background: dark gold-black */
.gkm {
    flex: 1 1 auto;
    min-height: calc(100vh - 4.5rem);
    overflow-y: auto;
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
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: clamp(1.5rem, 4vw, 2.5rem) clamp(1rem, 5vw, 2.5rem) 5rem;
    gap: clamp(1.5rem, 3vw, 2rem);
    color: #e8d5a3;
}

/* Back button */
.gkm__back {
    align-self: flex-start;
    display: inline-flex;
    align-items: center;
    gap: 0.45rem;
    padding: 0.55rem 1.1rem 0.55rem 0.85rem;
    border: 1px solid rgba(200, 134, 10, 0.4);
    border-radius: 6px;
    background: rgba(200, 134, 10, 0.07);
    color: #c8860a;
    font-size: clamp(0.85rem, 2vw, 0.95rem);
    font-weight: 500;
    text-decoration: none;
    letter-spacing: 0.03em;
    transition: background 0.15s, border-color 0.15s, color 0.15s;
}

.gkm__back:hover {
    background: rgba(200, 134, 10, 0.14);
    border-color: rgba(200, 134, 10, 0.65);
    color: #f5d87a;
}

/* Loading */
.gkm__loading {
    color: #9a7230;
    font-style: italic;
}

/* Locked */
.gkm__locked {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    margin-top: 4rem;
}

.gkm__lock-icon {
    font-size: 3.5rem;
    filter: grayscale(1) opacity(0.5);
}

.gkm__lock-msg {
    font-family: 'Cinzel Decorative', serif;
    font-size: 1rem;
    color: #6b4c1e;
    text-align: center;
}

/* Header */
.gkm__header {
    text-align: center;
    font-family: 'Cinzel Decorative', serif;
}

.gkm__header-sub {
    font-size: clamp(0.75rem, 2vw, 0.95rem);
    font-weight: 700;
    letter-spacing: 0.25em;
    text-transform: uppercase;
    color: #9a7230;
    margin: 0 0 0.3rem;
}

.gkm__title {
    font-size: clamp(1.6rem, 5vw, 2.8rem);
    font-weight: 700;
    margin: 0;
    background: linear-gradient(160deg, #f5d87a 0%, #c8860a 50%, #f5d87a 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
}

/* Finder section */
.gkm__finder {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.75rem;
    padding: 1.5rem 2rem;
    border: 1px solid rgba(200, 134, 10, 0.35);
    border-radius: 8px;
    background: rgba(200, 134, 10, 0.06);
    max-width: 420px;
    width: 100%;
}

.gkm__finder-label {
    font-family: 'Cinzel Decorative', serif;
    font-size: 0.85rem;
    font-weight: 700;
    letter-spacing: 0.15em;
    color: #c8860a;
    text-transform: uppercase;
}

.gkm__finder-img {
    width: 120px;
    height: 120px;
    border-radius: 50%;
    object-fit: cover;
    border: 2px solid rgba(200, 134, 10, 0.5);
}

.gkm__finder-name {
    font-size: clamp(1rem, 3vw, 1.3rem);
    font-weight: 600;
    color: #f5d87a;
    margin: 0;
}

.gkm__finder-date {
    font-size: clamp(0.8rem, 2vw, 0.95rem);
    color: #9a7230;
    margin: 0;
    letter-spacing: 0.04em;
}

/* Hints */
.gkm__hints {
    width: 100%;
    max-width: 720px;
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
}

.gkm__hints-title {
    font-family: 'Cinzel Decorative', serif;
    font-size: 1.1rem;
    font-weight: 700;
    color: #c8860a;
    margin: 0 0 0.5rem;
    border-bottom: 1px solid rgba(200, 134, 10, 0.25);
    padding-bottom: 0.5rem;
}

.gkm__hint {
    display: flex;
    gap: 1rem;
    padding: 1.1rem 1.25rem;
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid rgba(200, 134, 10, 0.18);
    border-radius: 6px;
}

.gkm__hint-num {
    font-family: 'Cinzel Decorative', serif;
    font-size: 1.1rem;
    font-weight: 700;
    color: rgba(200, 134, 10, 0.6);
    min-width: 1.8rem;
    flex-shrink: 0;
    padding-top: 0.1em;
}

.gkm__hint-body {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
}

.gkm__hint-content {
    color: #e8d5a3;
    line-height: 1.7;
    font-size: clamp(0.9rem, 2vw, 1.05rem);
}

.gkm__hint-content :deep(a) {
    color: #f5d87a;
}

.gkm__hint-content :deep(strong) {
    color: #f5d87a;
}

.gkm__hint-img {
    max-width: 100%;
    border-radius: 4px;
    border: 1px solid rgba(200, 134, 10, 0.2);
}

.gkm__no-hints {
    color: #6b4c1e;
    font-style: italic;
    font-size: 0.9rem;
}
</style>
