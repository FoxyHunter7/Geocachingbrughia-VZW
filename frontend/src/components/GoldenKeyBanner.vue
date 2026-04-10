<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { RouterLink } from 'vue-router';
import { getGoldenKeySettings } from '@/services/GoldenKeyService';
import LanguageProvider from '@/services/LanguageService';

const isActive = ref(false);
const activationTime = ref(null);
const bannerTextMap = ref({});
const lang = computed(() => LanguageProvider.CURR_LANG.value);
const displayText = computed(() => bannerTextMap.value?.[lang.value] || '');
let ticker = null;

async function loadSettings() {
    try {
        const settings = await getGoldenKeySettings();
        activationTime.value = new Date(settings.activation_time);
        bannerTextMap.value = settings.banner_text || {};
    } catch {
        activationTime.value = new Date('2026-04-12T10:12:00Z');
    }
    checkActive();
}

function checkActive() {
    if (activationTime.value) {
        isActive.value = new Date() >= activationTime.value;
    }
}

onMounted(() => {
    loadSettings();
    ticker = setInterval(checkActive, 30000);
});

onUnmounted(() => {
    if (ticker) clearInterval(ticker);
});

const imageSrc = computed(() =>
    isActive.value
        ? '/assets/media/static/goldenkey-front-page-active.jpg'
        : '/assets/media/static/goldenkey-front-page-soon.jpg'
);
</script>

<template>
    <section class="gk-banner">
        <!-- Atmospheric golden glow -->
        <div class="gk-glow" aria-hidden="true"></div>

        <div class="gk-content" :class="{ 'gk-content--active': isActive && displayText }">
            <!-- Image -->
            <RouterLink
                class="gk-image-link"
                to="/golden-key"
                :title="isActive ? 'Golden Key' : 'Golden Key — Coming Soon'"
            >
                <img
                    :src="imageSrc"
                    :alt="isActive ? 'Golden Key' : 'Golden Key — Coming Soon'"
                    class="gk-image"
                />
            </RouterLink>

            <!-- Active-state text + button (only shown when active and text for current language is set) -->
            <div v-if="isActive && displayText" class="gk-text-block">
                <p class="gk-text">{{ displayText }}</p>
                <RouterLink to="/golden-key" class="gk-discover-btn">Ontdek</RouterLink>
            </div>
        </div>
    </section>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Cinzel+Decorative:wght@700&display=swap');

.gk-banner {
    position: relative;
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
    padding: 3rem 1rem;
    background: #000;
    border-bottom: 3px solid var(--color-primary);
    overflow: hidden;
}

/* Golden radial glow */
.gk-glow {
    position: absolute;
    inset: 0;
    background: radial-gradient(
        ellipse 70% 60% at 50% 52%,
        rgba(212, 175, 55, 0.16) 0%,
        rgba(180, 110, 10, 0.07) 45%,
        transparent 70%
    );
    pointer-events: none;
}

/* Content wrapper — soon: single centered child; active: side-by-side */
.gk-content {
    position: relative;
    z-index: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 3.5rem;
}

.gk-content--active {
    width: min(1100px, calc(100% - 2rem));
    justify-content: center;
}

@media (max-width: 1000px) {
    .gk-content--active {
        flex-direction: column;
        align-items: center;
        gap: 1.75rem;
        width: min(600px, calc(100% - 2rem));
    }
}

/* Image link */
.gk-image-link {
    display: block;
    cursor: pointer;
    max-width: min(400px, 86vw);
    flex-shrink: 0;
    transition: transform 0.35s ease, filter 0.35s ease;
}

.gk-content--active .gk-image-link {
    max-width: min(360px, 40vw);
    flex-shrink: 0;
}

@media (max-width: 1000px) {
    .gk-content--active .gk-image-link {
        max-width: min(320px, 86vw);
    }
}

.gk-image-link:hover {
    transform: translateY(-6px);
    filter: brightness(1.06) drop-shadow(0 0 32px rgba(212, 175, 55, 0.35));
}

.gk-image {
    display: block;
    width: 100%;
    height: auto;
}

/* Active text block */
.gk-text-block {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    gap: 1.5rem;
    flex: 1;
    min-width: 0;
}

@media (max-width: 1000px) {
    .gk-text-block {
        flex: none;
        width: 100%;
    }
}

.gk-text {
    font-family: 'Cinzel Decorative', serif;
    font-size: 0.85rem;
    font-weight: 700;
    line-height: 1.6;
    letter-spacing: 0.02em;
    background: linear-gradient(160deg, #f5d87a 0%, #c8860a 55%, #f5d87a 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    margin: 0;
    white-space: pre-wrap;
}

.gk-discover-btn {
    display: inline-block;
    font-family: 'Cinzel Decorative', serif;
    font-size: 0.75rem;
    font-weight: 700;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    color: #c8860a;
    border: 2px solid #c8860a;
    padding: 0.6rem 1.6rem;
    text-decoration: none;
    transition: background 0.25s ease, color 0.25s ease;
}

.gk-discover-btn:hover {
    background: #c8860a;
    color: #000;
}
</style>
