<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import config from '@/data/config.js';
import { LanguageProvider } from '@/services/LanguageService';
import { StaticContentProvider as SCP } from '@/services/StaticContentService.js';

// TipTap for rendering descriptions
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
    Document,
    Paragraph,
    Text,
    Bold,
    Italic,
    Underline,
    Link,
    BulletList,
    ListItem,
    Strike,
    Table,
    TableRow,
    TableHeader,
    TableCell
];

const route = useRoute();
const router = useRouter();

const dictionary = SCP.DICTIONARY;
const event = ref(null);
const loading = ref(true);
const error = ref(null);
const isPreview = computed(() => route.query.preview === 'true');
const allLanguages = ref([]);

const currentLang = computed(() => LanguageProvider.CURR_LANG.value);

const missingTranslations = computed(() => {
    if (!event.value?.translations || !allLanguages.value.length) return [];
    const translatedLangs = event.value.translations.map(t => t.lang_code);
    return allLanguages.value
        .filter(lang => lang.active && !translatedLangs.includes(lang.code))
        .map(lang => lang.name);
});

const isDraft = computed(() => event.value?.state === 'draft');

async function fetchLanguages() {
    try {
        const response = await fetch(`${config.apiUrl}languages`);
        if (response.ok) {
            allLanguages.value = await response.json();
        }
    } catch (err) {
        console.error('Failed to fetch languages:', err);
    }
}

async function fetchEvent() {
    loading.value = true;
    error.value = null;
    
    try {
        const uuid = route.params.uuid;
        const previewParam = isPreview.value ? '&preview=true' : '';
        
        // Build request options - include auth header for preview mode
        const options = {};
        if (isPreview.value) {
            const token = localStorage.getItem('admin_token');
            if (token) {
                options.headers = { 'Authorization': `Bearer ${token}` };
            }
        }
        
        const response = await fetch(`${config.apiUrl}events/${uuid}?lang=${currentLang.value}${previewParam}`, options);
        
        if (response.ok) {
            event.value = await response.json();
        } else if (response.status === 404) {
            error.value = 'Event not found';
        } else if (response.status === 401) {
            error.value = 'Niet ingelogd - preview niet beschikbaar';
        } else {
            error.value = 'Failed to load event';
        }
    } catch (err) {
        console.error('Error fetching event:', err);
        error.value = 'Failed to load event';
    }
    
    loading.value = false;
}

function getDescription() {
    if (!event.value?.translations?.length) return '';
    const translation = event.value.translations.find(t => t.lang_code === currentLang.value);
    if (!translation?.description) return '';
    
    try {
        const parsed = JSON.parse(translation.description);
        return generateHTML(parsed, tiptapExtensions);
    } catch {
        return translation.description;
    }
}

function formatDate(dateString) {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString(currentLang.value, {
        weekday: 'short',
        day: 'numeric',
        month: 'short',
        year: 'numeric'
    }) + ' ' + date.toLocaleTimeString(currentLang.value, {
        hour: '2-digit',
        minute: '2-digit'
    });
}

function getEventTypeImage(type) {
    return `/assets/media/eventtypes/${type || 'REGULAR'}.png`;
}

function goBack() {
    // If coming from preview (admin), go back in history
    // Otherwise navigate to events page
    if (isPreview.value) {
        router.back();
    } else {
        router.push({ name: 'events' });
    }
}

onMounted(() => {
    fetchEvent();
    if (isPreview.value) {
        fetchLanguages();
    }
});

watch(() => route.params.uuid, fetchEvent);
watch(currentLang, fetchEvent);
</script>

<template>
    <!-- Preview Warning Banner (Admin only) -->
    <div v-if="isPreview && isDraft" class="preview-banner">
        <div class="preview-banner-content">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="preview-icon">
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                <circle cx="12" cy="12" r="3"/>
            </svg>
            <div class="preview-text">
                <strong>Conceptweergave</strong>
                <span>Dit evenement is nog niet gepubliceerd en alleen zichtbaar voor beheerders.</span>
            </div>
        </div>
        <div v-if="missingTranslations.length > 0" class="translation-warning">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="warning-icon">
                <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>
                <line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/>
            </svg>
            <span>Ontbrekende vertalingen: {{ missingTranslations.join(', ') }}</span>
        </div>
    </div>
    
    <main class="event-detail">
        <!-- Loading State -->
        <div v-if="loading" class="loading-container">
            <div class="loading-spinner"></div>
            <p>{{ dictionary.UILoading?.[currentLang] }}</p>
        </div>
        
        <!-- Error State -->
        <div v-else-if="error" class="error-container">
            <div class="error-icon">⚠️</div>
            <h2>{{ dictionary.UIError?.[currentLang] }}</h2>
            <p>{{ error }}</p>
            <button @click="goBack" class="back-btn">← {{ dictionary.UIBack?.[currentLang] }}</button>
        </div>
        
        <!-- Event Content -->
        <article v-else-if="event" class="event-article">
            <!-- Back link -->
            <router-link :to="{ name: 'events' }" class="back-link">
                ← {{ dictionary.NavEvents?.[currentLang] }}
            </router-link>
            
            <!-- Main content: Image left, details right -->
            <div class="event-layout">
                <!-- Left: Image -->
                <div class="event-image-section">
                    <img 
                        v-if="event.imageUrl" 
                        :src="`${config.apiUrl}images/${event.imageUrl}`" 
                        :alt="event.title"
                        class="event-image"
                    />
                    <img 
                        v-else 
                        :src="getEventTypeImage(event.type)" 
                        :alt="event.type"
                        class="event-type-fallback"
                    />
                </div>
                
                <!-- Right: Details -->
                <div class="event-details">
                    <div class="event-header">
                        <img :src="getEventTypeImage(event.type)" :alt="event.type" class="event-type-icon" />
                        <h1 class="event-title">{{ event.title }}</h1>
                    </div>
                    
                    <div class="event-meta">
                        <p class="event-datetime">
                            <strong>{{ dictionary.UIStart?.[currentLang] }}:</strong> {{ formatDate(event.start_date) }}
                            &nbsp;—&nbsp;
                            <strong>{{ dictionary.UIEnd?.[currentLang] }}:</strong> {{ formatDate(event.end_date) }}
                        </p>
                        
                        <p class="event-location" v-if="event.location">
                            <strong>{{ dictionary.UILocation?.[currentLang] }}:</strong> {{ event.location }}
                        </p>
                    </div>
                    
                    <!-- Action Buttons -->
                    <div class="event-actions">
                        <a v-if="event.geolink" :href="event.geolink" target="_blank" rel="noopener" class="action-btn primary">
                            <svg class="btn-icon-svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <circle cx="12" cy="12" r="10"></circle>
                                <path d="M12 2a14.5 14.5 0 0 0 0 20 14.5 14.5 0 0 0 0-20"></path>
                                <path d="M2 12h20"></path>
                            </svg>
                            {{ dictionary.ButtonViewOnGC?.[currentLang] }}
                        </a>
                        <a v-if="event.ticket_purchase_url" :href="event.ticket_purchase_url" target="_blank" rel="noopener" class="action-btn secondary">
                            <svg class="btn-icon-svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <path d="M2 9a3 3 0 0 1 0 6v2a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2v-2a3 3 0 0 1 0-6V7a2 2 0 0 0-2-2H4a2 2 0 0 0-2 2Z"></path>
                                <path d="M13 5v2"></path>
                                <path d="M13 17v2"></path>
                                <path d="M13 11v2"></path>
                            </svg>
                            {{ dictionary.ButtonPurchase?.[currentLang] }}
                        </a>
                    </div>
                </div>
            </div>
            
            <!-- Description (full width below) -->
            <div class="event-description" v-if="getDescription()" v-html="getDescription()"></div>
        </article>
    </main>
</template>

<style scoped>
.event-detail {
    background: var(--color-background);
    padding: 2rem 3rem;
    min-height: calc(100vh - 4.5rem);
    max-height: calc(100vh - 4.5rem);
    overflow-y: auto;
}

.loading-container,
.error-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 50vh;
    text-align: center;
    padding: 2rem;
}

.loading-spinner {
    width: 3rem;
    height: 3rem;
    border: 3px solid var(--color-background-2);
    border-top-color: var(--color-primary);
    border-radius: 50%;
    animation: spin 1s linear infinite;
}

@keyframes spin {
    to { transform: rotate(360deg); }
}

.error-icon {
    font-size: 4rem;
    margin-bottom: 1rem;
}

.back-btn {
    margin-top: 1.5rem;
    padding: 0.5rem 1rem;
    background: var(--color-accent-light);
    color: var(--color-text);
    border: none;
    border-radius: 0.4rem;
    cursor: pointer;
    font-weight: 600;
    transition: transform 0.15s;
}

.back-btn:hover {
    transform: scale(1.02);
}

.event-article {
    max-width: 1400px;
    margin: 0 auto;
}

.back-link {
    display: inline-block;
    margin-bottom: 1rem;
    color: var(--color-accent-dark);
    background: none;
    border: none;
    cursor: pointer;
    font-size: 0.9rem;
    padding: 0;
    transition: transform 0.2s;
}

.back-link:hover {
    transform: translateX(-4px);
}

/* Main layout: image left, details right - matching homepage Event.vue */
.event-layout {
    display: grid;
    grid-template-columns: 1fr 2fr;
    gap: 3rem;
    align-items: start;
}

.event-image-section {
    display: flex;
    justify-content: center;
    align-items: flex-start;
    max-height: 35rem;
    overflow: hidden;
}

.event-image {
    max-width: 100%;
    max-height: 90%;
    object-fit: contain;
    border-radius: 0.5rem;
}

.event-type-fallback {
    width: 8rem;
    height: 8rem;
    object-fit: contain;
    padding: 1.5rem;
    background: var(--color-background-2);
    border-radius: 0.5rem;
}

.event-details {
    display: flex;
    flex-direction: column;
    gap: 1rem;
}

.event-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
}

.event-type-icon {
    height: 2.5rem;
    object-fit: contain;
}

.event-title {
    font-size: 1.5rem;
    font-weight: bold;
    color: var(--color-text);
    line-height: 1.3;
    margin: 0;
}

.event-title::first-letter {
    text-transform: uppercase;
}

.event-meta {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    font-size: 0.9rem;
}

.event-datetime,
.event-location {
    margin: 0;
    color: var(--color-text);
}

.event-datetime strong,
.event-location strong {
    font-weight: normal;
    color: var(--color-text-2);
}

/* Action buttons - matching homepage Event.vue */
.event-actions {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
    margin-top: 0.5rem;
    font-size: 0.8rem;
}

.action-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    border-radius: 0.4rem;
    font-weight: normal;
    text-decoration: none;
    background-color: var(--color-accent-light);
    color: var(--color-text);
    transition: transform 0.15s;
}

.action-btn::first-letter {
    text-transform: capitalize;
}

.action-btn:hover {
    transform: scale(1.02);
}

.action-btn .btn-icon-svg {
    width: 1rem;
    height: 1rem;
    flex-shrink: 0;
}

.action-btn.primary {
    background: var(--color-accent-light);
    color: var(--color-text);
}

.action-btn.secondary {
    background: var(--color-accent-light);
    color: var(--color-text);
}

/* Description with TipTap styling from main.css */
.event-description {
    margin-top: 2rem;
    padding-top: 1.5rem;
    border-top: 1px solid var(--color-background-2);
    line-height: 1.6;
    color: var(--color-text);
}

.event-description :deep(p) {
    padding: 0.3rem 0;
    margin: 0;
}

.event-description :deep(a) {
    color: var(--color-accent-dark);
    text-decoration: underline;
}

.event-description :deep(ul) {
    margin: 0.5rem 0;
    padding-left: 1.5rem;
}

.event-description :deep(li) {
    margin-bottom: 0.25rem;
}

/* TipTap table styling - matching main.css */
.event-description :deep(table) {
    border-collapse: collapse;
    margin: 1rem 0;
    overflow: hidden;
    table-layout: fixed;
    width: 100%;
    max-width: 100%;
}

.event-description :deep(td),
.event-description :deep(th) {
    border: 1px solid var(--color-accent-dark);
    box-sizing: border-box;
    min-width: 1em;
    padding: 0.25rem 0.5rem;
    vertical-align: top;
}

.event-description :deep(th) {
    background-color: var(--color-primary);
    font-weight: bold;
    text-align: left;
}

.event-description :deep(strong) {
    font-weight: bold;
}

.event-description :deep(em) {
    font-style: italic;
}

.event-description :deep(u) {
    text-decoration: underline;
}

/* Mobile responsive */
@media (max-width: 1000px) {
    .event-detail {
        padding: 1rem;
        max-height: none;
    }
    
    .event-layout {
        grid-template-columns: 1fr;
        gap: 1.5rem;
    }
    
    .event-image-section {
        max-height: 40rem;
        justify-content: center;
        align-items: center;
    }
    
    .event-image {
        max-height: 100%;
    }
    
    .event-title {
        font-size: 1.3rem;
    }
    
    .event-details {
        text-align: center;
    }
    
    .event-header {
        justify-content: center;
    }
    
    .event-actions {
        justify-content: center;
    }
    
    .event-description :deep(table) {
        table-layout: auto !important;
        width: auto !important;
    }
}

/* Preview Banner */
.preview-banner {
    background: linear-gradient(135deg, var(--color-alert) 0%, var(--color-alert-dark, #c0392b) 100%);
    color: white;
    padding: 0.75rem 2rem;
}

.preview-banner-content {
    display: flex;
    align-items: center;
    gap: 1rem;
    max-width: 1400px;
    margin: 0 auto;
}

.preview-icon {
    width: 1.25rem;
    height: 1.25rem;
    flex-shrink: 0;
}

.preview-text {
    display: flex;
    flex-direction: column;
    gap: 0.1rem;
}

.preview-text strong {
    font-weight: 600;
    font-size: 0.9rem;
}

.preview-text span {
    font-size: 0.8rem;
    opacity: 0.9;
}

.translation-warning {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-top: 0.5rem;
    padding: 0.5rem 0.75rem;
    background: rgba(0, 0, 0, 0.2);
    border-radius: 0.4rem;
    max-width: 1400px;
    margin-left: auto;
    margin-right: auto;
}

.warning-icon {
    width: 1rem;
    height: 1rem;
    flex-shrink: 0;
}

.translation-warning span {
    font-size: 0.8rem;
}

@media (max-width: 768px) {
    .preview-banner {
        padding: 0.75rem 1rem;
    }
    
    .preview-banner-content {
        flex-direction: column;
        align-items: flex-start;
        gap: 0.5rem;
    }
}
</style>
