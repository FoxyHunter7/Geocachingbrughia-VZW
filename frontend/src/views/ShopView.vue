<script setup>
import { ref, computed, onMounted } from 'vue';
import { getShopSettings, getShopItems, createCheckoutSession } from '@/services/ShopService';
import LanguageProvider from '@/services/LanguageService';
import { StaticContentProvider } from '@/services/StaticContentService';

const dictionary = StaticContentProvider.DICTIONARY;

const loading = ref(true);
const settings = ref({ stripe_publishable_key: '', pretix_widget_url: '', currency: 'EUR' });
const items = ref([]);
const error = ref('');
const pretixLoaded = ref(false);
const lang = computed(() => LanguageProvider.CURR_LANG.value);

async function loadData() {
    loading.value = true;
    try {
        const [s, i] = await Promise.all([getShopSettings(), getShopItems()]);
        settings.value = s || { stripe_publishable_key: '', pretix_widget_url: '', currency: 'EUR' };
        items.value = Array.isArray(i) ? i : [];
        if (settings.value.pretix_widget_url && !pretixLoaded.value) {
            loadPretix(settings.value.pretix_widget_url);
        }
    } catch { error.value = 'Kon de winkel niet laden.'; }
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
        checkoutError.value = 'E-mailadres is verplicht';
        return;
    }
    if (form.value.quantity < 1) {
        checkoutError.value = 'Aantal moet minimaal 1 zijn';
        return;
    }
    if (form.value.fulfillment_type === 'shipping') {
        if (!form.value.shipping_name.trim() || !form.value.shipping_address.trim() || !form.value.shipping_city.trim() || !form.value.shipping_postal_code.trim()) {
            checkoutError.value = 'Alle verzendvelden zijn verplicht';
            return;
        }
    }

    checkoutLoading.value = true;
    try {
        const res = await createCheckoutSession({ item_id: checkoutItem.value.id, quantity: form.value.quantity, fulfillment_type: form.value.fulfillment_type, buyer_email: form.value.buyer_email, shipping_name: form.value.shipping_name, shipping_address: form.value.shipping_address, shipping_city: form.value.shipping_city, shipping_postal_code: form.value.shipping_postal_code, shipping_country: form.value.shipping_country });
        if (res?.success && res.data?.checkout_url) { window.location.href = res.data.checkout_url; }
        else if (res?.error) { checkoutError.value = res.error; }
        else { checkoutError.value = 'Betaling starten mislukt.'; }
    } catch { checkoutError.value = 'Er is een fout opgetreden.'; }
    checkoutLoading.value = false;
}

function fmt(item) { return item.price_display || `\u20ac ${(item.price_cents / 100).toFixed(2)}`; }
const hasItems = computed(() => items.value.length > 0);
const hasPretix = computed(() => !!settings.value.pretix_widget_url);
function imgUrl(url) { return url?.startsWith('http') ? url : `/api/images/${url}`; }
const euCountryNames = {
    BE: 'België', NL: 'Nederland', LU: 'Luxemburg', FR: 'Frankrijk',
    DE: 'Duitsland', GB: 'Verenigd Koninkrijk', IE: 'Ierland', ES: 'Spanje',
    PT: 'Portugal', IT: 'Italië', AT: 'Oostenrijk', CH: 'Zwitserland',
    DK: 'Denemarken', SE: 'Zweden', NO: 'Noorwegen', FI: 'Finland',
    PL: 'Polen', CZ: 'Tsjechië', SK: 'Slowakije', HU: 'Hongarije',
    RO: 'Roemenië', BG: 'Bulgarije', HR: 'Kroatië', SI: 'Slovenië',
    EE: 'Estland', LV: 'Letland', LT: 'Litouwen', GR: 'Griekenland',
};

function countriesSummary(codes) {
    if (!codes || codes.length === 0) return '';
    if (codes.length <= 3) return codes.map(c => euCountryNames[c] || c).join(', ');
    return `${codes.length} landen`;
}

const checkoutCountries = computed(() => {
    if (!checkoutItem.value?.shipping_countries) return [];
    return checkoutItem.value.shipping_countries.map(code => ({
        code,
        name: euCountryNames[code] || code,
    }));
});
</script>

<template>
    <main class="shop-page">
        <div v-if="loading" class="shop-loading"><div class="admin-spinner"></div><p>Winkel laden...</p></div>
        <div v-else-if="error" class="shop-error"><p>{{ error }}</p></div>
        <template v-else>
            <section v-if="hasPretix" class="shop-section">
                <h2 class="shop-section-title">Tickets</h2>
                <p class="shop-section-sub">Bestel je tickets voor onze evenementen.</p>
                <div class="pretix-container">
                    <pretix-widget :event="settings.pretix_widget_url"></pretix-widget>
                    <noscript><div class="pretix-noscript"><a :href="settings.pretix_widget_url" target="_blank" rel="noopener">Ga naar de ticketshop</a></div></noscript>
                </div>
            </section>
            <section v-if="hasItems" class="shop-section">
                <h2 class="shop-section-title">Webshop</h2>
                <p class="shop-section-sub">Merchandise en andere items.</p>
                <div class="shop-grid">
                    <div v-for="item in items" :key="item.id" class="shop-card">
                        <div class="shop-card-img">
                            <img v-if="item.image_url" :src="imgUrl(item.image_url)" :alt="item.title" loading="lazy" />
                            <div v-else class="shop-card-noimg"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg></div>
                        </div>
                        <div class="shop-card-body">
                            <h3 class="shop-card-title">{{ item.title }}</h3>
                            <p v-if="item.description" class="shop-card-desc">{{ item.description }}</p>
                            <div class="shop-card-meta">
                                <span class="shop-card-price">{{ fmt(item) }}</span>
                                <span class="shop-card-stock">{{ item.stock_quantity != null ? item.stock_quantity + ' op voorraad' : 'Op voorraad' }}</span>
                            </div>
                            <div class="shop-card-fulfillment">
                                <span v-if="item.allow_pickup" class="fb fb-pickup">Afhalen{{ item.pickup_label ? ': ' + item.pickup_label : '' }}</span>
                                <span v-if="item.allow_shipping" class="fb fb-shipping">Verzenden{{ item.shipping_countries ? ' (' + countriesSummary(item.shipping_countries) + ')' : '' }}</span>
                            </div>
                            <button class="shop-buy-btn" @click="openCheckout(item)">Kopen</button>
                        </div>
                    </div>
                </div>
            </section>
            <section v-if="!hasPretix && !hasItems" class="shop-empty">
                <p>{{ dictionary.UIShopEmpty?.[lang] ?? 'De webshop is momenteel leeg.' }}</p>
                <p class="shop-empty-sub">{{ dictionary.UIShopEmptySubTxt?.[lang] ?? 'Kom later terug!' }}</p>
            </section>
        </template>
        <Teleport to="body">
            <div v-if="checkoutItem" class="admin-modal-overlay" @click.self="closeCheckout">
                <div class="admin-modal admin-modal-lg" role="dialog" aria-modal="true" aria-label="Bestellen">
                    <div class="admin-modal-header">
                        <h2 class="admin-modal-title">Bestellen: {{ checkoutItem.title }}</h2>
                        <button class="admin-modal-close" @click="closeCheckout" aria-label="Sluiten"><svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg></button>
                    </div>
                    <div class="admin-modal-body">
                        <div class="admin-form-group"><label class="admin-label">E-mailadres *</label><input v-model="form.buyer_email" type="email" class="admin-input" required placeholder="jouw@email.be" /></div>
                        <div class="admin-form-group"><label class="admin-label">Aantal *</label><input v-model.number="form.quantity" type="number" class="admin-input" min="1" :max="checkoutItem.stock_quantity || 999" required /></div>
                        <div v-if="checkoutItem.allow_pickup && checkoutItem.allow_shipping" class="admin-form-group">
                            <label class="admin-label">Levering *</label>
                            <div class="fc">
                                <label class="fc-opt"><input type="radio" v-model="form.fulfillment_type" value="pickup" /><span>Afhalen{{ checkoutItem.pickup_label ? ' (' + checkoutItem.pickup_label + ')' : '' }}</span></label>
                                <label class="fc-opt"><input type="radio" v-model="form.fulfillment_type" value="shipping" /><span>Verzenden{{ checkoutItem.shipping_countries ? ' (' + countriesSummary(checkoutItem.shipping_countries) + ')' : '' }}</span></label>
                            </div>
                        </div>
                        <template v-if="form.fulfillment_type === 'shipping'">
                            <div class="admin-form-group"><label class="admin-label">Naam *</label><input v-model="form.shipping_name" type="text" class="admin-input" required /></div>
                            <div class="admin-form-group"><label class="admin-label">Adres *</label><input v-model="form.shipping_address" type="text" class="admin-input" required /></div>
                            <div class="form-row">
                                <div class="admin-form-group"><label class="admin-label">Stad *</label><input v-model="form.shipping_city" type="text" class="admin-input" required /></div>
                                <div class="admin-form-group"><label class="admin-label">Postcode *</label><input v-model="form.shipping_postal_code" type="text" class="admin-input" required /></div>
                            </div>
                            <div class="admin-form-group"><label class="admin-label">Land *</label><select v-model="form.shipping_country" class="admin-select"><option v-for="c in checkoutCountries" :key="c.code" :value="c.code">{{ c.name }}</option></select></div>
                        </template>
                        <div v-if="checkoutError" class="admin-alert admin-alert-danger"><span>{{ checkoutError }}</span></div>
                        <div class="checkout-sum"><div class="checkout-sum-row"><span>Totaal:</span><strong>{{ fmt({ price_cents: checkoutItem.price_cents * form.quantity, price_display: '' }) }}</strong></div><p class="checkout-sum-note">Je wordt doorgestuurd naar de beveiligde betaalomgeving van Stripe.</p></div>
                    </div>
                    <div class="admin-modal-footer">
                        <button class="admin-btn admin-btn-secondary" @click="closeCheckout" :disabled="checkoutLoading">Annuleren</button>
                        <button class="admin-btn admin-btn-primary" @click="submitCheckout" :disabled="checkoutLoading">{{ checkoutLoading ? 'Bezig...' : 'Betalen' }}</button>
                    </div>
                </div>
            </div>
        </Teleport>
    </main>
</template>

<style scoped>
.shop-page { flex: 1 1 auto; max-width: 100rem; margin: 0 auto; padding: 2rem 1rem; display: flex; flex-direction: column; gap: 3rem; }
.shop-loading { display: flex; flex-direction: column; align-items: center; gap: 1rem; padding: 4rem; color: #64748b; }
.shop-error { text-align: center; padding: 4rem; color: #ef4444; }
.shop-section { display: flex; flex-direction: column; gap: 1rem; }
.shop-section-title { font-size: 1.75rem; font-weight: 700; color: #1e293b; margin: 0; }
.shop-section-sub { color: #64748b; font-size: 0.95rem; margin: 0; }
.pretix-container { min-height: 200px; border: 1px solid #e2e8f0; border-radius: 0.75rem; overflow: hidden; }
.pretix-noscript { padding: 2rem; text-align: center; }
.pretix-noscript a { display: inline-block; padding: 0.75rem 1.5rem; background: #0d9488; color: white; text-decoration: none; border-radius: 0.5rem; font-weight: 500; }
.shop-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 1.5rem; }
.shop-card { background: white; border: 1px solid #e2e8f0; border-radius: 0.75rem; overflow: hidden; display: flex; flex-direction: column; transition: box-shadow 0.15s; }
.shop-card:hover { box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1); }
.shop-card-img { width: 100%; aspect-ratio: 1; background: #f8fafc; display: flex; align-items: center; justify-content: center; overflow: hidden; }
.shop-card-img img { width: 100%; height: 100%; object-fit: cover; }
.shop-card-noimg { color: #cbd5e1; }
.shop-card-noimg svg { width: 3rem; height: 3rem; }
.shop-card-body { padding: 1.25rem; display: flex; flex-direction: column; gap: 0.5rem; flex: 1; }
.shop-card-title { font-size: 1.125rem; font-weight: 600; color: #1e293b; margin: 0; }
.shop-card-desc { font-size: 0.875rem; color: #64748b; margin: 0; line-height: 1.5; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }
.shop-card-meta { display: flex; justify-content: space-between; align-items: center; margin-top: 0.5rem; }
.shop-card-price { font-size: 1.25rem; font-weight: 700; color: #0d9488; }
.shop-card-stock { font-size: 0.75rem; color: #94a3b8; }
.shop-card-fulfillment { display: flex; gap: 0.5rem; flex-wrap: wrap; }
.fb { font-size: 0.6875rem; padding: 0.125rem 0.5rem; border-radius: 9999px; font-weight: 500; }
.fb-pickup { background: #f0fdfa; color: #0d9488; }
.fb-shipping { background: #eff6ff; color: #3b82f6; }
.shop-buy-btn { margin-top: 0.75rem; padding: 0.625rem 1.25rem; background: #0d9488; color: white; border: none; border-radius: 0.5rem; font-size: 0.9375rem; font-weight: 600; cursor: pointer; transition: background 0.15s; }
.shop-buy-btn:hover { background: #0f766e; }
.shop-empty { text-align: center; padding: 4rem; color: #64748b; font-size: 1.125rem; }
.shop-empty-sub { margin: 0.5rem 0 0; font-size: 0.95rem; }
.fc { display: flex; gap: 1rem; }
.fc-opt { display: flex; align-items: center; gap: 0.5rem; cursor: pointer; padding: 0.75rem 1rem; border: 1px solid #e2e8f0; border-radius: 0.5rem; flex: 1; }
.fc-opt input { accent-color: #0d9488; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
.checkout-sum { margin-top: 1rem; padding-top: 1rem; border-top: 1px solid #e2e8f0; }
.checkout-sum-row { display: flex; justify-content: space-between; align-items: center; font-size: 1.125rem; }
.checkout-sum-row strong { color: #0d9488; }
.checkout-sum-note { font-size: 0.8125rem; color: #94a3b8; margin: 0.5rem 0 0; }
</style>
