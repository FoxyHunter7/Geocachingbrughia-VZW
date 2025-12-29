<script setup>
import { getAllEvents } from '@/services/EventService';
import LanguageProvider from '@/services/LanguageService';
import StaticContentProvider from '@/services/StaticContentService';
import { computed, onMounted, ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import config from '@/data/config.js';

const router = useRouter();
const lang = computed(() => LanguageProvider.CURR_LANG.value);
const dictionary = StaticContentProvider.DICTIONARY;

const events = ref([]);
const loading = ref(true);
const currPage = ref(1);
const lastPage = ref(1);

async function fetchData() {
    loading.value = true;
    const pagedResponse = await getAllEvents(false, currPage.value);
    if (pagedResponse.data) {
        events.value = pagedResponse.data;
        currPage.value = pagedResponse.current_page;
        lastPage.value = pagedResponse.last_page;
    } else if (Array.isArray(pagedResponse)) {
        events.value = pagedResponse;
    }
    loading.value = false;
}

const search = ref("");

const filteredEvents = computed(() => {
    if (!search.value) return events.value;
    return events.value.filter(event =>
        event.title?.toLowerCase().includes(search.value.toLowerCase())
    );
});

function nextPage() {
    if (currPage.value < lastPage.value) currPage.value++;
}

function prevPage() {
    if (currPage.value > 1) currPage.value--;
}

function goToEvent(event) {
    if (event.uuid) {
        router.push({ name: 'eventDetail', params: { uuid: event.uuid } });
    }
}

function formatDate(dateString) {
    if (!dateString) return '';
    return new Date(dateString).toLocaleDateString(lang.value, {
        month: 'short',
        day: 'numeric',
        year: 'numeric'
    });
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

        <section v-else-if="events.length === 0" id="no-events">
            <h2>{{ dictionary.UINoEvents?.[lang] }}</h2>
            <p>{{ dictionary.UINoEventsSubTxt?.[lang] }}</p>
        </section>

        <template v-else>
            <form method="post" @submit.prevent="">
                <div>
                    <input v-model="search" type="search" id="search" name="search" autocomplete="off" :placeholder="dictionary.FormSearch?.[lang]">
                </div>
            </form>

            <div id="events">
                <article v-for="event in filteredEvents" :key="event.id || event.uuid" class="event-card" @click="goToEvent(event)">
                    <div class="event-image">
                        <img v-if="event.imageUrl" :src="`${config.apiUrl}images/${event.imageUrl}`" :alt="event.title" loading="lazy" />
                        <img v-else :src="`/assets/media/eventtypes/${event.type || 'REGULAR'}.png`" :alt="event.type" class="type-icon" />
                    </div>
                    <div class="event-info">
                        <h3>{{ event.title }}</h3>
                        <p class="event-date">{{ formatDate(event.start_date) }}</p>
                        <p v-if="event.location" class="event-location">{{ event.location }}</p>
                    </div>
                </article>
            </div>

            <div v-if="filteredEvents.length === 0" id="no-results">
                <p>{{ dictionary.UINoEventsSearchResults?.[lang] }}</p>
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

#no-events h2, #no-results p {
    text-align: center;
    font-weight: bold;
    font-size: 1.2rem;
}

#no-events p {
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

#events {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 2rem;
    padding: 1rem;
    max-width: 90rem;
    margin: 0 auto;
}

.event-card {
    background: var(--color-background);
    border-radius: 0.5rem;
    overflow: hidden;
    cursor: pointer;
    transition: transform 0.2s, box-shadow 0.2s;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    display: flex;
    flex-direction: column;
}

.event-card:hover {
    transform: translateY(-4px);
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.event-image {
    aspect-ratio: 4/3;
    overflow: hidden;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--color-background-2);
}

.event-image img {
    max-width: 100%;
    max-height: 100%;
    object-fit: contain;
}

.event-image img.type-icon {
    width: 5rem;
    height: 5rem;
    object-fit: contain;
}

.event-info {
    padding: 1rem 1.25rem;
    flex: 1;
    display: flex;
    flex-direction: column;
}

.event-info h3 {
    font-size: 1.1rem;
    font-weight: 600;
    margin-bottom: 0.5rem;
    color: var(--color-text);
}

.event-date, .event-location {
    font-size: 0.875rem;
    color: var(--color-text);
    opacity: 0.7;
    margin: 0.25rem 0;
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
    
    #events {
        grid-template-columns: 1fr;
        gap: 1.5rem;
    }
    
    .event-image {
        aspect-ratio: 16/9;
    }
}
</style>
