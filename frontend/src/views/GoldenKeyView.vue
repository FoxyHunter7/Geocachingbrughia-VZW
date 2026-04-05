<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { getGoldenKeySettings } from '@/services/GoldenKeyService';

const isActive = ref(false);
const activationTime = ref(null);
let ticker = null;

const timeLeft = ref({ days: 0, hours: 0, minutes: 0, seconds: 0 });

async function loadSettings() {
    try {
        const settings = await getGoldenKeySettings();
        activationTime.value = new Date(settings.activation_time);
    } catch {
        activationTime.value = new Date('2026-04-12T10:12:00Z');
    }
    tick();
}

function tick() {
    if (!activationTime.value) return;
    const now = new Date();
    const diff = activationTime.value - now;

    if (diff <= 0) {
        isActive.value = true;
        timeLeft.value = { days: 0, hours: 0, minutes: 0, seconds: 0 };
        return;
    }

    isActive.value = false;
    const totalSeconds = Math.floor(diff / 1000);
    timeLeft.value = {
        days:    Math.floor(totalSeconds / 86400),
        hours:   Math.floor((totalSeconds % 86400) / 3600),
        minutes: Math.floor((totalSeconds % 3600) / 60),
        seconds: totalSeconds % 60,
    };
}

function pad(n) {
    return String(n).padStart(2, '0');
}

onMounted(() => {
    loadSettings();
    ticker = setInterval(tick, 1000);
});

onUnmounted(() => {
    if (ticker) clearInterval(ticker);
});
</script>

<template>
    <!-- Soon state: black page with countdown + image -->
    <main v-if="!isActive" class="gk-soon">
        <div class="gk-countdown">
            <div class="gk-countdown__unit">
                <span class="gk-countdown__number">{{ pad(timeLeft.days) }}</span>
                <span class="gk-countdown__label">Days</span>
            </div>
            <span class="gk-countdown__sep">:</span>
            <div class="gk-countdown__unit">
                <span class="gk-countdown__number">{{ pad(timeLeft.hours) }}</span>
                <span class="gk-countdown__label">Hours</span>
            </div>
            <span class="gk-countdown__sep">:</span>
            <div class="gk-countdown__unit">
                <span class="gk-countdown__number">{{ pad(timeLeft.minutes) }}</span>
                <span class="gk-countdown__label">Minutes</span>
            </div>
            <span class="gk-countdown__sep">:</span>
            <div class="gk-countdown__unit">
                <span class="gk-countdown__number">{{ pad(timeLeft.seconds) }}</span>
                <span class="gk-countdown__label">Seconds</span>
            </div>
        </div>

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
@import url('https://fonts.googleapis.com/css2?family=Cinzel+Decorative:wght@700&display=swap');

.gk-soon,
.gk-active {
    flex: 1 1 auto;
    background: #000;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 2.5rem;
    min-height: calc(100vh - 4.5rem);
    padding: 2rem 1rem;
}

/* Countdown */
.gk-countdown {
    display: flex;
    align-items: flex-start;
    gap: 0.5rem;
    font-family: 'Cinzel Decorative', serif;
}

.gk-countdown__unit {
    display: flex;
    flex-direction: column;
    align-items: center;
    min-width: 4rem;
}

.gk-countdown__number {
    font-size: clamp(2.2rem, 6vw, 4rem);
    font-weight: 700;
    line-height: 1;
    background: linear-gradient(160deg, #f5d87a 0%, #c8860a 50%, #f5d87a 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    letter-spacing: 0.04em;
}

.gk-countdown__label {
    font-size: clamp(0.55rem, 1.5vw, 0.75rem);
    font-weight: 700;
    letter-spacing: 0.15em;
    text-transform: uppercase;
    color: #9a7230;
    margin-top: 0.35rem;
}

.gk-countdown__sep {
    font-family: 'Cinzel Decorative', serif;
    font-size: clamp(2rem, 5vw, 3.5rem);
    font-weight: 700;
    line-height: 1;
    color: #c8860a;
    padding-bottom: 0.6rem;
    align-self: flex-end;
}

.gk-soon-image {
    max-width: min(460px, 90vw);
    max-height: calc(100vh - 4.5rem - 8rem);
    object-fit: contain;
    display: block;
}
</style>
