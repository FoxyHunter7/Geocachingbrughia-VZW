<script setup>
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import LanguageProvider from '@/services/LanguageService';
import { StaticContentProvider } from '@/services/StaticContentService';

const route = useRoute();
const orderId = route.query.order || '';
const lang = computed(() => LanguageProvider.CURR_LANG.value);
const dictionary = StaticContentProvider.DICTIONARY;

function t(key, fallback) {
    return dictionary[key]?.[lang.value] ?? fallback;
}

function tId(key, fallback) {
    const text = dictionary[key]?.[lang.value] ?? fallback;
    return text.split('///id///')[0] + ' ' + orderId + ' ' + text.split('///id///')[1];
}
</script>

<template>
    <main class="shop-result">
        <div class="result-card">
            <div class="result-icon result-cancel">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"/>
                    <line x1="15" y1="9" x2="9" y2="15"/>
                    <line x1="9" y1="9" x2="15" y2="15"/>
                </svg>
            </div>
            <h1>{{ t('ShopOrderCancelTitle', 'Betaling geannuleerd') }}</h1>
            <p v-if="orderId">{{ tId('ShopOrderCancelledId', 'Je bestelling #' + orderId + ' is geannuleerd.') }}</p>
            <p v-else>{{ t('ShopOrderCancelled', 'Je bestelling is geannuleerd.') }}</p>
            <p class="result-note">{{ t('ShopNoAmountCharged', 'Geen bedrag werd afgeschreven.') }}</p>
            <RouterLink to="/shop" class="result-back">{{ t('ShopBackToShop', 'Terug naar de webshop') }}</RouterLink>
        </div>
    </main>
</template>

<style scoped>
.shop-result {
    flex: 1 1 auto;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 2rem;
}

.result-card {
    text-align: center;
    max-width: 28rem;
    padding: 3rem 2rem;
}

.result-icon {
    width: 4rem;
    height: 4rem;
    margin: 0 auto 1.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
}

.result-cancel {
    background: var(--color-alert);
    color: var(--color-alert-dark);
}

.result-icon svg {
    width: 2.5rem;
    height: 2.5rem;
}

.result-card h1 {
    font-size: 1.5rem;
    font-weight: 700;
    margin: 0 0 0.75rem;
    color: var(--color-text);
}

.result-card p {
    color: var(--color-text);
    opacity: 0.75;
    margin: 0 0 0.25rem;
}

.result-note {
    font-size: 0.875rem;
    opacity: 0.6;
    margin-top: 0.5rem;
}

.result-back {
    display: inline-block;
    margin-top: 1.5rem;
    padding: 0.625rem 1.5rem;
    background: var(--color-accent-dark);
    color: var(--color-background);
    text-decoration: none;
    border-radius: 0.5rem;
    font-weight: 500;
    transition: filter 0.15s;
}

.result-back:hover {
    filter: brightness(1.15);
}
</style>