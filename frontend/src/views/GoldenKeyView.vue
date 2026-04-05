<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { getGoldenKeySettings } from '@/services/GoldenKeyService';

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
</script>

<template>
    <!-- Soon state: black page with centered contained image -->
    <main v-if="!isActive" class="gk-soon">
        <img
            :src="'/assets/media/static/goldenkey-page-soon.jpeg'"
            alt="Golden Key — Coming Soon"
            class="gk-soon-image"
        />
    </main>

    <!-- Active state: content to be added in a future update -->
    <main v-else class="gk-active">
        <!-- TODO: Active Golden Key content will be implemented in a later update -->
    </main>
</template>

<style scoped>
.gk-soon,
.gk-active {
    flex: 1 1 auto;
    background: #000;
    display: flex;
    align-items: center;
    justify-content: center;
    /* fill whatever height remains after the header */
    min-height: calc(100vh - 4.5rem);
}

.gk-soon-image {
    max-width: 100%;
    max-height: calc(100vh - 4.5rem);
    object-fit: contain;
    display: block;
}
</style>
