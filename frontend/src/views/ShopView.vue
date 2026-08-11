<script setup>
import { ref, computed, onMounted } from 'vue';
import { getShopSettings, getShopItems, createCheckoutSession } from '@/services/ShopService';
import LanguageProvider from '@/services/LanguageService';
import { StaticContentProvider } from '@/services/StaticContentService';

const dictionary = StaticContentProvider.DICTIONARY;
const lang = computed(() => LanguageProvider.CURR_LANG.value);

const loading = ref(true);
const settings = ref({ stripe_publishable_key: '', pretix_widget_url: '', currency: 'EUR' });
const items = ref([]);
const error = ref('');
const pretixLoaded = ref(false);

function t(key, fallback) {
    return dictionary[key]?.[lang.value] ?? fallback;
}

async function loadData() {
    loading.value = true;
    try {
        const [s, i] = await Promise.all([getShopSettings(), getShopItems()]);
        settings.value = s || { stripe_publishable_key: '', pretix_widget_url: '', currency: 'EUR' };
        items.value = Array.isArray(i) ? i : [];
        if (settings.value.pretix_widget_url && !pretixLoaded.value) {
            loadPretix(settings.value.pretix_widget_url);
        }
    } catch { error.value = t('ShopLoadError', 'Kon de winkel niet laden.'); }
    loading.value = false;
}

function loadPretix(url) {
    const base = url.replace(/\/$/, '');
    const link = document.createElement('link');
    link.rel = 'stylesheet';
    link.href = `${base}/widget/v2.css`;
    document.head.appendChild(link);
    const script = document.createElement('script');
    script.src = `${base}/widget/v2.en.js`;
    script.async = true;
    script.onload = () => { pretixLoaded.value = true; };
    document.head.appendChild(script);
}

onMounted(loadData);

const checkoutItem = ref(null);
const form = ref({ quantity: 1, fulfillment_type: 'pickup', buyer_email: '', shipping_name: '', shipping_address: '', shipping_city: '', shipping_postal_code: '', shipping_country: 'BE' });
const checkoutError = ref('');
const checkoutLoading = ref(false);

function openCheckout(item) {
    checkoutItem.value = item;
    const defaultCountry = item.shipping_countries?.[0] || 'BE';
    form.value = { quantity: 1, fulfillment_type: item.allow_pickup ? 'pickup' : 'shipping', buyer_email: '', shipping_name: '', shipping_address: '', shipping_city: '', shipping_postal_code: '', shipping_country: defaultCountry };
    checkoutError.value = '';
}

function closeCheckout() { checkoutItem.value = null; checkoutError.value = ''; }

async function submitCheckout() {
    if (!checkoutItem.value) return;
    checkoutError.value = '';

    if (!form.value.buyer_email.trim()) {
        checkoutError.value = t('ShopErrorEmail', 'E-mailadres is verplicht');
        return;
    }
    if (form.value.quantity < 1) {
        checkoutError.value = t('ShopErrorQuantity', 'Aantal moet minimaal 1 zijn');
        return;
    }
    if (form.value.fulfillment_type === 'shipping') {
        if (!form.value.shipping_name.trim() || !form.value.shipping_address.trim() || !form.value.shipping_city.trim() || !form.value.shipping_postal_code.trim()) {
            checkoutError.value = t('ShopErrorShipping', 'Alle verzendvelden zijn verplicht');
            return;
        }
    }

    checkoutLoading.value = true;
    try {
        const res = await createCheckoutSession({ item_id: checkoutItem.value.id, quantity: form.value.quantity, fulfillment_type: form.value.fulfillment_type, buyer_email: form.value.buyer_email, shipping_name: form.value.shipping_name, shipping_address: form.value.shipping_address, shipping_city: form.value.shipping_city, shipping_postal_code: form.value.shipping_postal_code, shipping_country: form.value.shipping_country });
        if (res?.success && res.data?.checkout_url) { window.location.href = res.data.checkout_url; }
        else if (res?.error) { checkoutError.value = res.error; }
        else { checkoutError.value = t('ShopErrorStart', 'Betaling starten mislukt.'); }
    } catch { checkoutError.value = t('ShopErrorGeneric', 'Er is een fout opgetreden.'); }
    checkoutLoading.value = false;
}

function fmt(item) { return item.price_display || `\u20ac ${(item.price_cents / 100).toFixed(2)}`; }
const hasItems = computed(() => items.value.length > 0);
const hasPretix = computed(() => !!settings.value.pretix_widget_url);
function imgUrl(url) { return url?.startsWith('http') ? url : `/api/images/${url}`; }

function countryName(code) {
    return dictionary['Country' + code]?.[lang.value] || code;
}

function countriesSummary(codes) {
    if (!codes || codes.length === 0) return '';
    if (codes.length <= 3) return codes.map(countryName).join(', ');
    return `${codes.length} ${t('ShopCountries', 'landen')}`;
}

const checkoutCountries = computed(() => {
    if (!checkoutItem.value?.shipping_countries) return [];
    return checkoutItem.value.shipping_countries.map(code => ({
        code,
        name: countryName(code),
    }));
});
</script>

<template>
    <main class="shop-page">
        <section v-if="loading" id="loading" class="shop-status">
            <div class="shop-loader"></div>
            <p>{{ t('UILoadingShop', 'Winkel laden') }}...</p>
        </section>

        <section v-else-if="error" class="shop-status shop-status-error">
            <p>{{ error }}</p>
        </section>

        <template v-else>
            <section v-if="hasPretix" class="shop-section">
                <h2 class="shop-section-title">{{ t('ShopTicketsTitle', 'Tickets') }}</h2>
                <p class="shop-section-sub">{{ t('ShopTicketsSubTxt', 'Bestel je tickets voor onze evenementen.') }}</p>
                <div class="pretix-container">
                    <pretix-widget :event="settings.pretix_widget_url"></pretix-widget>
                    <noscript><div class="pretix-noscript"><a :href="settings.pretix_widget_url" target="_blank" rel="noopener">{{ t('ShopTicketsGoTo', 'Ga naar de ticketshop') }}</a></div></noscript>
                </div>
            </section>

            <section v-if="hasItems" class="shop-section">
                <h2 class="shop-section-title">{{ t('ShopMerchTitle', 'Webshop') }}</h2>
                <p class="shop-section-sub">{{ t('ShopMerchSubTxt', 'Merchandise en andere items.') }}</p>
                <div class="shop-grid">
                    <article v-for="item in items" :key="item.id" class="shop-card">
                        <div class="shop-card-img">
                            <img v-if="item.image_url" :src="imgUrl(item.image_url)" :alt="item.title" loading="lazy" />
                            <div v-else class="shop-card-noimg"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg></div>
                        </div>
                        <div class="shop-card-body">
                            <h3 class="shop-card-title">{{ item.title }}</h3>
                            <p v-if="item.description" class="shop-card-desc">{{ item.description }}</p>
                            <div class="shop-card-meta">
                                <span class="shop-card-price">{{ fmt(item) }}</span>
                                <span class="shop-card-stock">{{ item.stock_quantity != null ? item.stock_quantity + ' ' + t('ShopStockCount', 'op voorraad') : t('ShopInStock', 'Op voorraad') }}</span>
                            </div>
                            <div class="shop-card-fulfillment">
                                <span v-if="item.allow_pickup" class="fb fb-pickup">{{ t('ShopPickup', 'Afhalen') }}{{ item.pickup_label ? ': ' + item.pickup_label : '' }}</span>
                                <span v-if="item.allow_shipping" class="fb fb-shipping">{{ t('ShopShipping', 'Verzenden') }}{{ item.shipping_countries ? ' (' + countriesSummary(item.shipping_countries) + ')' : '' }}</span>
                            </div>
                            <button class="shop-buy-btn" @click="openCheckout(item)">{{ t('ShopBuy', 'Kopen') }}</button>
                        </div>
                    </article>
                </div>
            </section>

            <section v-if="!hasPretix && !hasItems" class="shop-status">
                <p>{{ t('UIShopEmpty', 'De webshop is momenteel leeg.') }}</p>
                <p class="shop-status-sub">{{ t('UIShopEmptySubTxt', 'Kom later terug!') }}</p>
            </section>
        </template>

        <Teleport to="body">
            <div v-if="checkoutItem" class="shop-modal-overlay" @click.self="closeCheckout">
                <div class="shop-modal" role="dialog" aria-modal="true" :aria-label="t('ShopOrderTitle', 'Bestellen')">
                    <div class="shop-modal-header">
                        <h2 class="shop-modal-title">{{ t('ShopOrderTitle', 'Bestellen') }}: {{ checkoutItem.title }}</h2>
                        <button class="shop-modal-close" @click="closeCheckout" :aria-label="t('ShopClose', 'Sluiten')"><svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg></button>
                    </div>
                    <div class="shop-modal-body">
                        <div class="shop-form-group"><label class="shop-label">{{ t('ShopEmail', 'E-mailadres') }} *</label><input v-model="form.buyer_email" type="email" class="shop-input" required :placeholder="t('ShopEmail', 'E-mailadres')" /></div>
                        <div class="shop-form-group"><label class="shop-label">{{ t('ShopQuantity', 'Aantal') }} *</label><input v-model.number="form.quantity" type="number" class="shop-input" min="1" :max="checkoutItem.stock_quantity || 999" required /></div>
                        <div v-if="checkoutItem.allow_pickup && checkoutItem.allow_shipping" class="shop-form-group">
                            <label class="shop-label">{{ t('ShopDelivery', 'Levering') }} *</label>
                            <div class="fc">
                                <label class="fc-opt"><input type="radio" v-model="form.fulfillment_type" value="pickup" /><span>{{ t('ShopPickup', 'Afhalen') }}{{ checkoutItem.pickup_label ? ' (' + checkoutItem.pickup_label + ')' : '' }}</span></label>
                                <label class="fc-opt"><input type="radio" v-model="form.fulfillment_type" value="shipping" /><span>{{ t('ShopShipping', 'Verzenden') }}{{ checkoutItem.shipping_countries ? ' (' + countriesSummary(checkoutItem.shipping_countries) + ')' : '' }}</span></label>
                            </div>
                        </div>
                        <template v-if="form.fulfillment_type === 'shipping'">
                            <div class="shop-form-group"><label class="shop-label">{{ t('ShopName', 'Naam') }} *</label><input v-model="form.shipping_name" type="text" class="shop-input" required /></div>
                            <div class="shop-form-group"><label class="shop-label">{{ t('ShopAddress', 'Adres') }} *</label><input v-model="form.shipping_address" type="text" class="shop-input" required /></div>
                            <div class="form-row">
                                <div class="shop-form-group"><label class="shop-label">{{ t('ShopCity', 'Stad') }} *</label><input v-model="form.shipping_city" type="text" class="shop-input" required /></div>
                                <div class="shop-form-group"><label class="shop-label">{{ t('ShopPostalCode', 'Postcode') }} *</label><input v-model="form.shipping_postal_code" type="text" class="shop-input" required /></div>
                            </div>
                            <div class="shop-form-group"><label class="shop-label">{{ t('ShopCountry', 'Land') }} *</label><select v-model="form.shipping_country" class="shop-input"><option v-for="c in checkoutCountries" :key="c.code" :value="c.code">{{ c.name }}</option></select></div>
                        </template>
                        <div v-if="checkoutError" class="shop-alert"><span>{{ checkoutError }}</span></div>
                        <div class="checkout-sum">
                            <div class="checkout-sum-row"><span>{{ t('ShopTotal', 'Totaal:') }}</span><strong>{{ fmt({ price_cents: checkoutItem.price_cents * form.quantity, price_display: '' }) }}</strong></div>
                            <p class="checkout-sum-note">{{ t('ShopStripeNote', 'Je wordt doorgestuurd naar de beveiligde betaalomgeving van Stripe.') }}</p>
                        </div>
                    </div>
                    <div class="shop-modal-footer">
                        <button class="shop-btn shop-btn-secondary" @click="closeCheckout" :disabled="checkoutLoading">{{ t('ButtonCancel', 'Annuleren') }}</button>
                        <button class="shop-btn shop-btn-primary" @click="submitCheckout" :disabled="checkoutLoading">{{ checkoutLoading ? t('ShopProcessing', 'Bezig...') : t('ShopPay', 'Betalen') }}</button>
                    </div>
                </div>
            </div>
        </Teleport>
    </main>
</template>

<style scoped>
.shop-page {
    position: relative;
    flex: 1 1 auto;
    height: 100%;
    overflow-y: auto;
    padding: 1rem;
    max-width: 90rem;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    gap: 2.5rem;
}

.shop-status {
    min-height: 60vh;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    gap: 0.75rem;
    text-align: center;
    color: var(--color-text);
    font-size: 1.1rem;
}

.shop-status-error {
    color: var(--color-alert-dark);
    font-weight: bold;
}

.shop-status-sub {
    color: var(--color-text);
    opacity: 0.7;
    font-size: 0.95rem;
}

.shop-loader {
    width: 2.5rem;
    height: 2.5rem;
    border-radius: 50%;
    border: 0.25rem solid var(--color-primary);
    border-top-color: var(--color-accent-dark);
    animation: shop-spin 0.8s linear infinite;
}

@keyframes shop-spin {
    to { transform: rotate(360deg); }
}

.shop-section {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
}

.shop-section-title {
    font-size: 1.75rem;
    font-weight: 700;
    color: var(--color-text);
    margin: 0;
}

.shop-section-sub {
    color: var(--color-text);
    opacity: 0.7;
    font-size: 0.95rem;
    margin: 0;
}

.pretix-container {
    min-height: 200px;
    border: 1px solid var(--color-border);
    border-radius: 0.5rem;
    overflow: hidden;
    background: var(--color-background);
}

.pretix-noscript {
    padding: 2rem;
    text-align: center;
}

.pretix-noscript a {
    display: inline-block;
    padding: 0.75rem 1.5rem;
    background: var(--color-accent-dark);
    color: var(--color-background);
    text-decoration: none;
    border-radius: 0.5rem;
    font-weight: 500;
}

.shop-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 1.5rem;
    margin-top: 0.5rem;
}

.shop-card {
    background: var(--color-background);
    border-radius: 0.5rem;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    transition: transform 0.2s, box-shadow 0.2s;
}

.shop-card:hover {
    transform: translateY(-4px);
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.shop-card-img {
    width: 100%;
    aspect-ratio: 1;
    background: var(--color-background-2);
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
}

.shop-card-img img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.shop-card-noimg {
    color: var(--color-accent-dark);
    opacity: 0.4;
}

.shop-card-noimg svg {
    width: 3rem;
    height: 3rem;
}

.shop-card-body {
    padding: 1.25rem;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    flex: 1;
}

.shop-card-title {
    font-size: 1.125rem;
    font-weight: 600;
    color: var(--color-text);
    margin: 0;
}

.shop-card-desc {
    font-size: 0.875rem;
    color: var(--color-text);
    opacity: 0.7;
    margin: 0;
    line-height: 1.5;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
}

.shop-card-meta {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 0.5rem;
}

.shop-card-price {
    font-size: 1.25rem;
    font-weight: 700;
    color: var(--color-accent-dark);
}

.shop-card-stock {
    font-size: 0.75rem;
    color: var(--color-text);
    opacity: 0.6;
}

.shop-card-fulfillment {
    display: flex;
    gap: 0.5rem;
    flex-wrap: wrap;
}

.fb {
    font-size: 0.6875rem;
    padding: 0.125rem 0.5rem;
    border-radius: 9999px;
    font-weight: 500;
}

.fb-pickup {
    background: var(--color-primary);
    color: var(--color-text);
}

.fb-shipping {
    background: var(--color-accent-light);
    color: var(--color-text);
}

.shop-buy-btn {
    margin-top: 0.75rem;
    padding: 0.625rem 1.25rem;
    background: var(--color-accent-dark);
    color: var(--color-background);
    border: none;
    border-radius: 0.5rem;
    font-size: 0.9375rem;
    font-weight: 600;
    font-family: inherit;
    cursor: pointer;
    transition: filter 0.15s;
}

.shop-buy-btn:hover {
    filter: brightness(1.15);
}

.shop-modal-overlay {
    position: fixed;
    inset: 0;
    z-index: 1000;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 1rem;
}

.shop-modal {
    background: var(--color-background);
    color: var(--color-text);
    border-radius: 0.75rem;
    max-width: 34rem;
    width: 100%;
    max-height: 90vh;
    display: flex;
    flex-direction: column;
    box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
}

.shop-modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.25rem 1.5rem;
    border-bottom: 1px solid var(--color-border);
}

.shop-modal-title {
    font-size: 1.25rem;
    font-weight: 700;
    margin: 0;
}

.shop-modal-close {
    background: none;
    border: none;
    cursor: pointer;
    color: var(--color-text);
    opacity: 0.6;
    padding: 0.25rem;
    display: flex;
}

.shop-modal-close:hover {
    opacity: 1;
}

.shop-modal-body {
    padding: 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
    overflow-y: auto;
}

.shop-form-group {
    display: flex;
    flex-direction: column;
    gap: 0.375rem;
}

.shop-label {
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--color-text);
}

.shop-input {
    height: 2rem;
    border-radius: 0.3rem;
    border: 0.1rem solid transparent;
    outline: none;
    font-family: inherit;
    font-size: 1rem;
    box-sizing: border-box;
    padding: 0 0.5rem;
    background-color: var(--color-background-2);
    color: var(--color-text-2);
}

.shop-input:focus {
    border-color: var(--color-accent-dark);
}

select.shop-input {
    text-transform: none;
}

.fc {
    display: flex;
    gap: 1rem;
}

.fc-opt {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
    padding: 0.75rem 1rem;
    border: 1px solid var(--color-border);
    border-radius: 0.5rem;
    flex: 1;
    font-size: 0.875rem;
    background: var(--color-background-3);
}

.fc-opt input {
    accent-color: var(--color-accent-dark);
}

.form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
}

.shop-alert {
    padding: 0.625rem 0.875rem;
    border-radius: 0.5rem;
    background: var(--color-alert);
    color: var(--vt-c-white);
    font-size: 0.875rem;
}

.checkout-sum {
    margin-top: 0.25rem;
    padding-top: 1rem;
    border-top: 1px solid var(--color-border);
}

.checkout-sum-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 1.125rem;
}

.checkout-sum-row strong {
    color: var(--color-accent-dark);
}

.checkout-sum-note {
    font-size: 0.8125rem;
    color: var(--color-text);
    opacity: 0.7;
    margin: 0.5rem 0 0;
}

.shop-modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
    padding: 1.25rem 1.5rem;
    border-top: 1px solid var(--color-border);
}

.shop-btn {
    padding: 0.625rem 1.25rem;
    border: none;
    border-radius: 0.5rem;
    font-size: 0.9375rem;
    font-weight: 600;
    font-family: inherit;
    cursor: pointer;
    transition: filter 0.15s;
}

.shop-btn:disabled {
    opacity: 0.6;
    cursor: auto;
}

.shop-btn-primary {
    background: var(--color-accent-dark);
    color: var(--color-background);
}

.shop-btn-secondary {
    background: var(--color-background-2);
    color: var(--color-text);
}

.shop-btn:hover:not(:disabled) {
    filter: brightness(1.1);
}

@media (max-width: 768px) {
    .form-row {
        grid-template-columns: 1fr;
    }

    .fc {
        flex-direction: column;
    }
}
</style>