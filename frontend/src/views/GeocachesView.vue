<script setup>
import { getAllGeocaches } from '@/services/GeocacheService';
import LanguageProvider from '@/services/LanguageService';
import StaticContentProvider from '@/services/StaticContentService';
import { computed, onMounted, ref, watch } from 'vue';

const lang = computed(() => LanguageProvider.CURR_LANG.value);
const dictionary = StaticContentProvider.DICTIONARY;

const geocaches = ref([]);
const loading = ref(true);
const currPage = ref(1);
const lastPage = ref(1);

async function fetchData() {
    loading.value = true;
    const pagedResponse = await getAllGeocaches(false, currPage.value);
    if (pagedResponse.data) {
        geocaches.value = pagedResponse.data;
        currPage.value = pagedResponse.current_page;
        lastPage.value = pagedResponse.last_page;
    } else if (Array.isArray(pagedResponse)) {
        geocaches.value = pagedResponse;
    }
    loading.value = false;
}

const search = ref("");

const filteredGeocaches = computed(() => {
    if (!search.value) return geocaches.value;
    return geocaches.value.filter(geocache =>
        (geocache.title || geocache.name)?.toLowerCase().includes(search.value.toLowerCase())
    );
});

function nextPage() {
    if (currPage.value < lastPage.value) currPage.value++;
}

function prevPage() {
    if (currPage.value > 1) currPage.value--;
}

function extractGeocode(url) {
    if (!url) return '';
    return url.split("/").pop();
}

// Map cache type to image filename
function getCacheTypeImage(type) {
    const typeMap = {
        'traditional': 'TRADITIONAL',
        'multi': 'MULTI',
        'mystery': 'MYSTERY',
        'letterbox': 'LETTERBOX',
        'wherigo': 'WHEREIGO',
        'earthcache': 'EARTH',
        'virtual': 'VIRTUAL',
        'event': 'TRADITIONAL',
        'cito': 'TRADITIONAL',
        'mega': 'TRADITIONAL',
        'giga': 'TRADITIONAL',
        'lab': 'LAB',
        'webcam': 'WEBCAM'
    };
    return typeMap[type?.toLowerCase()] || 'TRADITIONAL';
}

onMounted(fetchData);
watch(lang, fetchData);
watch(currPage, fetchData);
</script>

<template>
    <main>
        <section v-if="loading" id="loading">
            <div class="initial-loader"></div>
        </section>

        <section v-else-if="geocaches.length === 0" id="no-geocaches">
            <h2>{{ dictionary.UINoGeocaches?.[lang] }}</h2>
            <p>{{ dictionary.UINoGeocachesSubTxt?.[lang] }}</p>
        </section>

        <template v-else>
            <form method="post" @submit.prevent="">
                <div>
                    <input v-model="search" type="search" id="search" name="search" autocomplete="off" :placeholder="dictionary.FormSearch?.[lang]">
                </div>
            </form>

            <div id="geocaches">
                <article v-for="geocache in filteredGeocaches" :key="geocache.id" class="geocache-card">
                    <div class="card-header">
                        <img :src="`/assets/media/cachetypes/${getCacheTypeImage(geocache.type)}.png`" :alt="geocache.type" class="cache-icon" />
                        <h3>{{ geocache.title || geocache.name }}</h3>
                    </div>
                    
                    <div class="rating-row">
                        <div class="rating-item">
                            <span class="rating-label">{{ dictionary.GeocacheDifficulty?.[lang] }}</span>
                            <span class="rating-value">{{ geocache.difficulty }}</span>
                        </div>
                        <div class="rating-item">
                            <span class="rating-label">{{ dictionary.GeocacheTerrain?.[lang] }}</span>
                            <span class="rating-value">{{ geocache.terrain }}</span>
                        </div>
                    </div>

                    <a v-if="geocache.geolink" :href="geocache.geolink" target="_blank" class="geocache-link">
                        {{ extractGeocode(geocache.geolink) }}
                    </a>
                </article>
            </div>

            <div v-if="filteredGeocaches.length === 0" id="no-results">
                <p>{{ dictionary.UINoGeocachesSearchResults?.[lang] }}</p>
            </div>

            <div id="pager" v-if="lastPage > 1">
                <div class="prev pagerNavBtn" :class="{ disabled: currPage === 1 }" @click="prevPage"></div>
                <p>{{ currPage }} / {{ lastPage }}</p>
                <div class="next pagerNavBtn" :class="{ disabled: currPage === lastPage }" @click="nextPage"></div>
            </div>
        </template>
    </main>
</template>

<style scoped>
main {
    position: relative;
    flex: 1 1 auto;
    height: 100%;
    overflow-y: auto;
    padding: 1rem;
}

section {
    height: 70vh;
    max-height: 40rem;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
}

#loading {
    display: flex;
    justify-content: center;
    align-items: center;
}

#no-geocaches h2, #no-results p {
    text-align: center;
    font-weight: bold;
    font-size: 1.2rem;
}

#no-geocaches p {
    text-align: center;
    margin-top: 0.5rem;
}

form {
    display: flex;
    justify-content: flex-end;
    padding: 1rem;
}

form > div {
    position: relative;
}

form > div::before {
    content: '';
    height: 1.3rem;
    width: 1.3rem;
    position: absolute;
    left: 0.5rem;
    top: 0.3rem;
    background-color: var(--color-text);
    mask: url(@/assets/media/search.svg);
    mask-size: contain;
    mask-repeat: no-repeat;
    mask-position: center;
}

form > div input {
    width: 20rem;
    max-width: 95vw;
    padding-left: 2.2rem;
}

#geocaches {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 1.5rem;
    padding: 1rem;
    max-width: 80rem;
    margin: 0 auto;
}

.geocache-card {
    background: var(--color-background);
    border-radius: 0.5rem;
    padding: 1rem;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.card-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 1rem;
}

.cache-icon {
    width: 2rem;
    height: 2rem;
    flex-shrink: 0;
}

.card-header h3 {
    font-size: 1rem;
    font-weight: 600;
    color: var(--color-text);
    margin: 0;
}

.rating-row {
    display: flex;
    gap: 1.5rem;
    margin-bottom: 1rem;
}

.rating-item {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
}

.rating-label {
    font-size: 0.75rem;
    color: var(--color-text);
    opacity: 0.6;
}

.rating-value {
    font-size: 1rem;
    font-weight: 600;
    color: var(--color-text);
}

.geocache-link {
    display: inline-block;
    color: var(--color-accent-dark);
    text-decoration: underline;
    font-size: 0.875rem;
}

.geocache-link:hover {
    font-weight: bold;
}

#no-results {
    text-align: center;
    padding: 2rem;
}

#pager {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 2rem;
    margin-top: 2rem;
}

#pager p {
    font-size: 1rem;
    font-weight: bold;
    letter-spacing: 0.5rem;
}

#pager .pagerNavBtn {
    width: 2rem;
    height: 2rem;
    background-color: var(--color-text);
    cursor: pointer;
}

#pager .pagerNavBtn.next {
    mask: url(@/assets/media/chevron-right.svg);
    mask-size: contain;
    mask-repeat: no-repeat;
    mask-position: center;
}

#pager .pagerNavBtn.prev {
    mask: url(@/assets/media/chevron-left.svg);
    mask-size: contain;
    mask-repeat: no-repeat;
    mask-position: center;
}

#pager .pagerNavBtn.disabled {
    filter: opacity(30%);
    cursor: auto;
}

@media screen and (max-width: 600px) {
    form {
        justify-content: center;
    }
    
    #geocaches {
        grid-template-columns: 1fr;
    }
}
</style>
