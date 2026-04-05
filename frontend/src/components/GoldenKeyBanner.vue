<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import { getGoldenKeySettings } from '@/services/GoldenKeyService';

const router = useRouter();
const isActive = ref(false);
const activationTime = ref(null);
let ticker = null;

async function loadSettings() {
    try {
        const settings = await getGoldenKeySettings();
        activationTime.value = new Date(settings.activation_time);
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

function navigate() {
    router.push('/golden-key');
}
</script>

<template>
    <section class="gk-banner">
        <!-- Atmospheric golden glow, matches the light in the image itself -->
        <div class="gk-glow" aria-hidden="true"></div>

        <a class="gk-image-link" @click.prevent="navigate" href="/golden-key"
           :title="isActive ? 'Golden Key' : 'Golden Key — Coming Soon'">
            <img
                :src="imageSrc"
                :alt="isActive ? 'Golden Key' : 'Golden Key — Coming Soon'"
                class="gk-image"
            />
        </a>
    </section>
</template>

<style scoped>
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

/* Golden radial glow that mirrors the light emanating from the image */
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

.gk-image-link {
    position: relative;
    z-index: 1;
    display: block;
    cursor: pointer;
    max-width: min(460px, 86vw);
    transition: transform 0.35s ease, filter 0.35s ease;
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
</style>
