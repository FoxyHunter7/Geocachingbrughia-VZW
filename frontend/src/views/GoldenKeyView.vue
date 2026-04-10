<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { RouterLink } from 'vue-router';
import { getGoldenKeySettings } from '@/services/GoldenKeyService';
import { getGoldenKeyMonths } from '@/services/GoldenKeyMonthService';
import LanguageProvider from '@/services/LanguageService';
import StaticContentProvider from '@/services/StaticContentService';

const lang = computed(() => LanguageProvider.CURR_LANG.value);
const dictionary = StaticContentProvider.DICTIONARY;

const ROMAN = ['I','II','III','IV','V','VI','VII','VIII','IX','X','XI','XII'];

const months = ref([]);

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

async function loadMonths() {
    const data = await getGoldenKeyMonths();
    if (Array.isArray(data)) months.value = data;
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
    loadMonths();
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

    <!-- Active state: month grid -->
    <main v-else class="gk-active">
        <header class="gk-hero">
            <p class="gk-hero__sub">The hunt for</p>
            <h1 class="gk-hero__title">The Golden Key</h1>
        </header>

        <div class="gk-rules-wrap">
            <RouterLink to="/golden-key/rules" class="gk-rules-btn">
                <img src="/assets/media/static/goldenkey-knop.png" alt="" class="gk-rules-btn__img" aria-hidden="true" />
                <span class="gk-rules-btn__text">{{ dictionary.GoldenKeyRulesBtn?.[lang] ?? 'Spelregels' }}</span>
            </RouterLink>
        </div>

        <section class="gk-months">
            <component
                v-for="month in months"
                :key="month.id"
                :is="month.state !== 'locked' ? RouterLink : 'div'"
                :to="month.state !== 'locked' ? `/golden-key/${month.id}` : undefined"
                class="gk-btn"
                :class="{
                    'gk-btn--locked': month.state === 'locked',
                    'gk-btn--found':  month.state === 'found',
                }"
            >
                <img
                    src="/assets/media/static/goldenkey-knop-maand.png"
                    :alt="month.month_name"
                    class="gk-btn__img"
                />
                <div class="gk-btn__overlay">
                    <span class="gk-btn__roman">{{ ROMAN[month.month_number - 1] }}</span>
                    <div class="gk-btn__labels">
                        <span class="gk-btn__label-top">The Golden Key</span>
                        <span class="gk-btn__label-name">{{ month.month_name }}</span>
                    </div>
                </div>
                <span v-if="month.state === 'found'" class="gk-btn__found-badge">FOUND</span>
            </component>
        </section>
    </main>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Cinzel+Decorative:wght@700&display=swap');

/* ===== Shared base ===== */
.gk-soon,
.gk-active {
    flex: 1 1 auto;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: calc(100vh - 4.5rem);
    padding: 2rem 1rem;
    gap: 2.5rem;
}

/* ===== Soon: plain black ===== */
.gk-soon {
    background: #000;
}

/* ===== Active: gold-black premium leather feel ===== */
.gk-active {
    justify-content: flex-start;
    padding: clamp(2rem, 4vw, 3.5rem) clamp(1rem, 4vw, 3rem) 5rem;
    gap: clamp(1.5rem, 3vw, 2.5rem);

    /* Near-pure black base — warmth only in the gradients */
    background-color: #050505;
    background-image:
        /* Overhead gold light — the only real color in the page */
        radial-gradient(ellipse 90% 25% at 50% 0%, rgba(160, 110, 8, 0.10) 0%, transparent 100%),
        /* Barely-there warm depth at the sides */
        radial-gradient(ellipse 35% 70% at   0% 45%, rgba(110, 70, 5, 0.05) 0%, transparent 100%),
        radial-gradient(ellipse 35% 70% at 100% 55%, rgba(110, 70, 5, 0.05) 0%, transparent 100%),
        /* Micro leather grain — two crossing diagonals, nearly invisible */
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

/* ===== Hero header ===== */
.gk-hero {
    text-align: center;
    font-family: 'Cinzel Decorative', serif;
}

.gk-hero__sub {
    font-size: clamp(0.75rem, 2vw, 1rem);
    font-weight: 700;
    letter-spacing: 0.3em;
    text-transform: uppercase;
    color: #9a7230;
    margin: 0 0 0.4rem;
}

.gk-hero__title {
    font-size: clamp(2rem, 6vw, 4rem);
    font-weight: 700;
    margin: 0;
    background: linear-gradient(160deg, #f5d87a 0%, #c8860a 45%, #f5d87a 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    letter-spacing: 0.03em;
    line-height: 1.1;
}

/* ===== Rules button ===== */
.gk-rules-wrap {
    display: flex;
    justify-content: center;
}

.gk-rules-btn {
    position: relative;
    display: block;
    width: clamp(200px, 40vw, 380px);
    text-decoration: none;
    border-radius: 5px;
    overflow: hidden;
    cursor: pointer;
    transition: transform 0.15s ease, filter 0.15s ease;
}

.gk-rules-btn:hover {
    transform: scale(1.025);
    filter: brightness(1.12);
}

.gk-rules-btn__img {
    display: block;
    width: 100%;
    height: auto;
    pointer-events: none;
    user-select: none;
}

.gk-rules-btn__text {
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    font-family: 'Cinzel Decorative', serif;
    font-size: clamp(0.85rem, 3vw, 1.15rem);
    font-weight: 700;
    letter-spacing: 0.08em;
    background: linear-gradient(160deg, #f5d87a 0%, #c8860a 45%, #f5d87a 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
}

/* ===== Month grid ===== */
.gk-months {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: clamp(0.6rem, 1.5vw, 1.1rem);
    width: 100%;
    max-width: clamp(480px, 75vw, 1000px);
}

@media (max-width: 680px) {
    .gk-months {
        grid-template-columns: 1fr;
        max-width: 420px;
    }
}

/* ===== Button ===== */
.gk-btn {
    position: relative;
    display: block;
    width: 100%;
    text-decoration: none;
    border-radius: 5px;
    overflow: hidden;
    cursor: pointer;
    transition: transform 0.15s ease, filter 0.15s ease;
}

.gk-btn:not(.gk-btn--locked):hover {
    transform: scale(1.025);
    filter: brightness(1.12);
}

.gk-btn__img {
    display: block;
    width: 100%;
    height: auto;
    pointer-events: none;
    user-select: none;
}

/* ===== Text overlay (left ~62%, keyhole lives in right ~38%) ===== */
.gk-btn__overlay {
    position: absolute;
    top: 0;
    left: 0;
    bottom: 0;
    right: 38%;
    display: flex;
    align-items: center;
    padding-left: clamp(0.6rem, 4%, 1.2rem);
    gap: clamp(0.4rem, 3%, 0.85rem);
}

/* Roman numeral — scales with button width (which = ~45vw on single col, ~35vw on 2-col) */
.gk-btn__roman {
    font-family: 'Cinzel Decorative', serif;
    font-size: clamp(1.35rem, 5vw, 1.8rem);
    font-weight: 700;
    background: linear-gradient(160deg, #f5d87a 0%, #c8860a 55%, #f5d87a 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    flex-shrink: 0;
    min-width: 1.8em;
    text-align: center;
    line-height: 1;
}

/* Text labels */
.gk-btn__labels {
    display: flex;
    flex-direction: column;
    gap: 0.15em;
    min-width: 0;
    overflow: hidden;
}

.gk-btn__label-top {
    font-family: 'Cinzel Decorative', serif;
    font-size: clamp(0.72rem, 2.2vw, 0.88rem);
    font-weight: 700;
    letter-spacing: 0.06em;
    background: linear-gradient(160deg, #f5d87a 0%, #c8860a 65%, #f5d87a 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    line-height: 1.2;
}

.gk-btn__label-name {
    font-family: 'Cinzel Decorative', serif;
    font-size: clamp(1.05rem, 4vw, 1.3rem);
    font-weight: 700;
    background: linear-gradient(160deg, #fff4d6 0%, #e8cb82 70%, #fff4d6 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    line-height: 1.2;
}

/* ===== States ===== */
.gk-btn--locked {
    filter: grayscale(0.75) brightness(0.45);
    pointer-events: none;
    cursor: default;
}

/* Found badge — centered vertically between text area and keyhole slot */
.gk-btn__found-badge {
    position: absolute;
    top: 0;
    bottom: 0;
    left: 63%;
    right: 20%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-family: 'Cinzel Decorative', serif;
    font-size: clamp(0.85rem, 2.8vw, 1.1rem);
    font-weight: 700;
    letter-spacing: 0.04em;
    color: #f5d87a;
    text-shadow:
        0 0 10px rgba(200, 134, 10, 0.9),
        0 0 20px rgba(200, 134, 10, 0.5),
        0 1px 3px rgba(0, 0, 0, 0.9);
    pointer-events: none;
}

/* ============================================================
   Countdown (soon state)
   ============================================================ */
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
