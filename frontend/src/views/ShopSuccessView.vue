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
            <div class="result-icon result-success">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
                    <polyline points="22 4 12 14.01 9 11.01"/>
                </svg>
            </div>
            <h1>{{ t('ShopOrderSuccessTitle', 'Bedankt voor je aankoop!') }}</h1>
            <p v-if="orderId">{{ tId('ShopOrderReceivedId', 'Je bestelling #' + orderId + ' is ontvangen.') }}</p>
            <p v-else>{{ t('ShopOrderReceived', 'Je bestelling is ontvangen.') }}</p>
            <p class="result-note">{{ t('ShopOrderNote', 'Je ontvangt een bevestigingsmail met verdere details.') }}</p>
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

.result-success {
    background: var(--color-primary);
    color: var(--color-accent-dark);
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