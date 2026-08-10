<script setup>
import { ref, onMounted } from 'vue';
import AdminLayout from '@/components/admin/AdminLayout.vue';
import config from '@/data/config.js';

// ---- State ----
const loadingSettings = ref(true);
const loadingItems = ref(true);
const savingSettings = ref(false);
const savingItem = ref(false);

// Settings
const settings = ref({
    stripe_secret_key: '',
    stripe_publishable_key: '',
    stripe_webhook_secret: '',
    pretix_widget_url: ''
});

// Items
const items = ref([]);

// Item modal
const showModal = ref(false);
const modalMode = ref('create');
const editingItem = ref(null);
const fileInput = ref(null);
const imagePreview = ref('');
const selectedFile = ref(null);

const itemForm = ref({
    title: '',
    description: '',
    priceEuros: '',
    image_url: '',
    stock: '',
    allow_pickup: false,
    pickup_label: '',
    allow_shipping: false,
    shipping_regions: '',
    auto_confirm: false,
    is_active: true,
    sort_order: 0
});

// ---- API helper ----
function getToken() {
    return localStorage.getItem('admin_token');
}

async function apiRequest(endpoint, options = {}) {
    const token = getToken();
    const headers = {
        'Accept': 'application/json',
        ...(!(options.body instanceof FormData) && { 'Content-Type': 'application/json' }),
        ...(token && { 'Authorization': `Bearer ${token}` }),
        ...options.headers
    };
    const response = await fetch(`${config.apiUrl}${endpoint}`, { ...options, headers });
    if (!response.ok) {
        let msg = `Request failed (${response.status})`;
        try {
            const err = await response.json();
            msg = err.error || err.message || msg;
        } catch { /* response had no JSON body */ }
        throw new Error(msg);
    }
    return response;
}

// ---- Helpers ----
function centsToEuros(cents) {
    if (cents === null || cents === undefined) return '';
    return (Number(cents) / 100).toFixed(2);
}

function eurosToCents(euros) {
    if (euros === '' || euros === null || euros === undefined) return 0;
    const n = parseFloat(euros);
    return isNaN(n) ? 0 : Math.round(n * 100);
}

function imageUrl(filename) {
    return filename ? `${config.apiUrl}images/${filename}` : '';
}

function stockLabel(stock) {
    if (stock === null || stock === undefined) return 'Onbeperkt';
    return stock;
}

// ---- Settings ----
async function fetchSettings() {
    loadingSettings.value = true;
    try {
        const res = await apiRequest('admin/shop/settings');
        const data = await res.json();
        settings.value = {
            stripe_secret_key: data.stripe_secret_key || '',
            stripe_publishable_key: data.stripe_publishable_key || '',
            stripe_webhook_secret: data.stripe_webhook_secret || '',
            pretix_widget_url: data.pretix_widget_url || ''
        };
    } catch {
        window.$toast?.error('Instellingen laden mislukt');
    }
    loadingSettings.value = false;
}

async function saveSettings() {
    savingSettings.value = true;
    try {
        await apiRequest('admin/shop/settings', {
            method: 'PUT',
            body: JSON.stringify({ ...settings.value, currency: 'EUR' })
        });
        window.$toast?.success('Instellingen opgeslagen');
    } catch {
        window.$toast?.error('Instellingen opslaan mislukt');
    }
    savingSettings.value = false;
}

// ---- Items ----
async function fetchItems() {
    loadingItems.value = true;
    try {
        const res = await apiRequest('admin/shop/items');
        const data = await res.json();
        items.value = data.data || data || [];
    } catch {
        window.$toast?.error('Items laden mislukt');
    }
    loadingItems.value = false;
}

function openCreateModal() {
    modalMode.value = 'create';
    editingItem.value = null;
    selectedFile.value = null;
    imagePreview.value = '';
    itemForm.value = {
        title: '',
        description: '',
        priceEuros: '',
        image_url: '',
        stock: '',
        allow_pickup: false,
        pickup_label: '',
        allow_shipping: false,
        shipping_regions: '',
        auto_confirm: false,
        is_active: true,
        sort_order: 0
    };
    showModal.value = true;
}

function openEditModal(item) {
    modalMode.value = 'edit';
    editingItem.value = item;
    selectedFile.value = null;
    itemForm.value = {
        title: item.title || '',
        description: item.description || '',
        priceEuros: centsToEuros(item.price_cents),
        image_url: item.image_url || '',
        stock: item.stock_quantity === null || item.stock_quantity === undefined ? '' : item.stock_quantity,
        allow_pickup: !!item.allow_pickup,
        pickup_label: item.pickup_label || '',
        allow_shipping: !!item.allow_shipping,
        shipping_regions: item.shipping_regions || '',
        auto_confirm: !!item.auto_confirm,
        is_active: !!item.active,
        sort_order: item.sort_order || 0
    };
    imagePreview.value = imageUrl(item.image_url);
    showModal.value = true;
}

function closeModal() {
    if (savingItem.value) return;
    showModal.value = false;
    editingItem.value = null;
    selectedFile.value = null;
    imagePreview.value = '';
    if (fileInput.value) fileInput.value.value = '';
}

function handleImageChange(e) {
    const file = e.target.files?.[0];
    if (!file) return;
    selectedFile.value = file;
    const reader = new FileReader();
    reader.onload = (ev) => { imagePreview.value = ev.target.result; };
    reader.readAsDataURL(file);
}

async function uploadImage(file) {
    const formData = new FormData();
    formData.append('image', file);
    const res = await apiRequest('admin/upload-image', {
        method: 'POST',
        body: formData
    });
    const data = await res.json();
    return data.filename;
}

async function handleSaveItem() {
    const errors = [];

    if (!itemForm.value.title.trim()) {
        errors.push('Titel is verplicht');
    }

    const priceCents = eurosToCents(itemForm.value.priceEuros);
    if (itemForm.value.priceEuros === '' || itemForm.value.priceEuros === null || itemForm.value.priceEuros === undefined) {
        errors.push('Prijs is verplicht');
    } else if (isNaN(parseFloat(itemForm.value.priceEuros)) || parseFloat(itemForm.value.priceEuros) < 0) {
        errors.push('Prijs moet een positief getal zijn');
    } else if (priceCents <= 0) {
        errors.push('Prijs moet groter zijn dan 0');
    }

    if (!itemForm.value.allow_pickup && !itemForm.value.allow_shipping) {
        errors.push('Minstens één leveringsoptie (afhalen of verzenden) is verplicht');
    }

    if (itemForm.value.allow_pickup && !itemForm.value.pickup_label.trim()) {
        errors.push('Afhaal-label is verplicht wanneer afhalen is ingeschakeld');
    }

    if (itemForm.value.allow_shipping && !itemForm.value.shipping_regions.trim()) {
        errors.push('Verzendregio\'s zijn verplicht wanneer verzenden is ingeschakeld');
    }

    if (errors.length > 0) {
        window.$toast?.error(errors.join(', '));
        return;
    }

    savingItem.value = true;
    try {
        let finalImageUrl = itemForm.value.image_url || '';
        if (selectedFile.value) {
            const uploaded = await uploadImage(selectedFile.value);
            if (!uploaded) {
                window.$toast?.error('Afbeelding uploaden mislukt');
                savingItem.value = false;
                return;
            }
            finalImageUrl = uploaded;
        }

        const payload = {
            title: itemForm.value.title,
            description: itemForm.value.description || '',
            price_cents: priceCents,
            image_url: finalImageUrl,
            stock_quantity: itemForm.value.stock === '' ? null : parseInt(itemForm.value.stock, 10),
            allow_pickup: itemForm.value.allow_pickup,
            pickup_label: itemForm.value.pickup_label || '',
            allow_shipping: itemForm.value.allow_shipping,
            shipping_regions: itemForm.value.shipping_regions || '',
            auto_confirm: itemForm.value.auto_confirm,
            active: itemForm.value.is_active,
            sort_order: parseInt(itemForm.value.sort_order, 10) || 0
        };

        const endpoint = modalMode.value === 'create'
            ? 'admin/shop/items'
            : `admin/shop/items/${editingItem.value.id}`;
        const method = modalMode.value === 'create' ? 'POST' : 'PUT';

        await apiRequest(endpoint, { method, body: JSON.stringify(payload) });
        window.$toast?.success(modalMode.value === 'create' ? 'Item aangemaakt' : 'Item bijgewerkt');
        closeModal();
        fetchItems();
    } catch (err) {
        window.$toast?.error(err.message || 'Opslaan mislukt');
    }
    savingItem.value = false;
}

async function handleDeleteItem(item) {
    if (!item) return;
    if (!confirm(`Weet je zeker dat je "${item.title}" wilt verwijderen?`)) return;
    savingItem.value = true;
    try {
        await apiRequest(`admin/shop/items/${item.id}`, { method: 'DELETE' });
        window.$toast?.success('Item verwijderd');
        fetchItems();
    } catch (err) {
        window.$toast?.error(err.message || 'Verwijderen mislukt');
    }
    savingItem.value = false;
}

onMounted(async () => {
    await Promise.all([fetchSettings(), fetchItems()]);
});
</script>

<template>
    <AdminLayout pageTitle="Webshop">
        <template #actions>
            <button class="admin-btn admin-btn-primary" @click="openCreateModal">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 1rem; height: 1rem;">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                Nieuw Item
            </button>
        </template>

        <!-- Section 1: Settings -->
        <div class="admin-card" style="margin-bottom: 1.5rem;">
            <div class="admin-card-header">
                <h2 class="admin-card-title">Webshop Instellingen</h2>
                <div style="display: flex; gap: 0.75rem; align-items: center;">
                    <a
                        href="https://ticketing.geocachingbrughia.be/control/"
                        target="_blank"
                        rel="noopener"
                        class="admin-btn admin-btn-secondary"
                    >
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 1rem; height: 1rem;">
                            <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/>
                            <polyline points="15 3 21 3 21 9"/>
                            <line x1="10" y1="14" x2="21" y2="3"/>
                        </svg>
                        Pretix Beheer
                    </a>
                    <button
                        class="admin-btn admin-btn-primary"
                        @click="saveSettings"
                        :disabled="savingSettings || loadingSettings"
                    >
                        {{ savingSettings ? 'Opslaan...' : 'Opslaan' }}
                    </button>
                </div>
            </div>
            <div class="admin-card-body">
                <div v-if="loadingSettings" style="text-align: center; padding: 2rem;">
                    <div class="admin-spinner" style="margin: 0 auto;"></div>
                </div>
                <template v-else>
                    <div class="admin-form-group">
                        <label class="admin-label" for="stripe-secret-key">Stripe Secret Key</label>
                        <input
                            id="stripe-secret-key"
                            v-model="settings.stripe_secret_key"
                            type="password"
                            class="admin-input"
                            placeholder="sk_live_..."
                            autocomplete="off"
                        >
                    </div>
                    <div class="admin-form-group">
                        <label class="admin-label" for="stripe-publishable-key">Stripe Publishable Key</label>
                        <input
                            id="stripe-publishable-key"
                            v-model="settings.stripe_publishable_key"
                            type="text"
                            class="admin-input"
                            placeholder="pk_live_..."
                        >
                    </div>
                    <div class="admin-form-group">
                        <label class="admin-label" for="stripe-webhook-secret">Stripe Webhook Secret</label>
                        <input
                            id="stripe-webhook-secret"
                            v-model="settings.stripe_webhook_secret"
                            type="text"
                            class="admin-input"
                            placeholder="whsec_..."
                        >
                    </div>
                    <div class="admin-form-group">
                        <label class="admin-label" for="pretix-widget-url">Pretix Widget URL</label>
                        <input
                            id="pretix-widget-url"
                            v-model="settings.pretix_widget_url"
                            type="text"
                            class="admin-input"
                            placeholder="https://ticketing.geocachingbrughia.be/..."
                        >
                        <p class="admin-form-hint">
                            URL van de Pretix widget voor het insluiten van ticketverkoop op de website.
                        </p>
                    </div>
                </template>
            </div>
        </div>

        <!-- Section 2: Items table -->
        <div class="admin-card">
            <div class="admin-card-header">
                <h2 class="admin-card-title">Shop Items</h2>
            </div>
            <div class="admin-table-wrapper">
                <table class="admin-table">
                    <thead>
                        <tr>
                            <th style="width: 4rem;">Afbeelding</th>
                            <th>Titel</th>
                            <th style="width: 7rem;">Prijs</th>
                            <th style="width: 6rem;">Voorraad</th>
                            <th style="width: 6rem;">Actief</th>
                            <th style="width: 8rem;">Levering</th>
                            <th style="width: 6rem; text-align: right;">Acties</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-if="loadingItems">
                            <td colspan="7" style="text-align: center; padding: 3rem;">
                                <div class="admin-spinner" style="margin: 0 auto;"></div>
                            </td>
                        </tr>
                        <tr v-else-if="items.length === 0">
                            <td colspan="7">
                                <div class="admin-empty">
                                    <div class="admin-empty-icon">
                                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                                            <path d="M3 9l1-5h16l1 5"/>
                                            <path d="M4 9h16v11H4z"/>
                                            <path d="M9 13h6"/>
                                        </svg>
                                    </div>
                                    <p class="admin-empty-title">Geen items gevonden</p>
                                    <p class="admin-empty-description">Maak je eerste shop item aan om te beginnen.</p>
                                    <button class="admin-btn admin-btn-primary" @click="openCreateModal">
                                        Nieuw Item
                                    </button>
                                </div>
                            </td>
                        </tr>
                        <tr v-else v-for="item in items" :key="item.id">
                            <td>
                                <img
                                    v-if="item.image_url"
                                    :src="imageUrl(item.image_url)"
                                    :alt="item.title"
                                    class="shop-thumb"
                                >
                                <span v-else class="admin-text-muted">&mdash;</span>
                            </td>
                            <td>
                                <span style="font-weight: 500;">{{ item.title }}</span>
                            </td>
                            <td>{{ item.price_display || ('€ ' + centsToEuros(item.price_cents)) }}</td>
                            <td>{{ stockLabel(item.stock_quantity) }}</td>
                            <td>
                                <span
                                    :class="['admin-badge', item.active ? 'admin-badge-success' : 'admin-badge-neutral']"
                                >
                                    {{ item.active ? 'Actief' : 'Inactief' }}
                                </span>
                            </td>
                            <td>
                                <div class="fulfillment-badges">
                                    <span v-if="item.allow_pickup" class="admin-badge admin-badge-info">Afhalen</span>
                                    <span v-if="item.allow_shipping" class="admin-badge admin-badge-success">Verzending</span>
                                    <span
                                        v-if="!item.allow_pickup && !item.allow_shipping"
                                        class="admin-text-muted"
                                    >&mdash;</span>
                                </div>
                            </td>
                            <td>
                                <div class="admin-table-actions">
                                    <button
                                        class="admin-btn admin-btn-ghost admin-btn-icon admin-btn-sm"
                                        @click="openEditModal(item)"
                                        title="Bewerken"
                                    >
                                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 1rem; height: 1rem;">
                                            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                                            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                                        </svg>
                                    </button>
                                    <button
                                        class="admin-btn admin-btn-ghost admin-btn-icon admin-btn-sm"
                                        @click="handleDeleteItem(item)"
                                        title="Verwijderen"
                                    >
                                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 1rem; height: 1rem;">
                                            <polyline points="3 6 5 6 21 6"/>
                                            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                                        </svg>
                                    </button>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <!-- Item Modal -->
        <Teleport to="body">
            <div v-if="showModal" class="admin-modal-overlay" @click.self="closeModal">
                <div class="admin-modal admin-modal-lg">
                    <div class="admin-modal-header">
                        <h2 class="admin-modal-title">
                            {{ modalMode === 'create' ? 'Nieuw Item' : 'Item Bewerken' }}
                        </h2>
                        <button class="admin-modal-close" @click="closeModal" aria-label="Sluiten">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <line x1="18" y1="6" x2="6" y2="18"/>
                                <line x1="6" y1="6" x2="18" y2="18"/>
                            </svg>
                        </button>
                    </div>

                    <div class="admin-modal-body">
                        <div class="admin-form-group">
                            <label class="admin-label" for="item-title">Titel *</label>
                            <input
                                id="item-title"
                                v-model="itemForm.title"
                                type="text"
                                class="admin-input"
                                required
                            >
                        </div>

                        <div class="admin-form-group">
                            <label class="admin-label" for="item-description">Beschrijving</label>
                            <textarea
                                id="item-description"
                                v-model="itemForm.description"
                                class="admin-textarea"
                                rows="3"
                            ></textarea>
                        </div>

                        <div class="admin-form-row">
                            <div class="admin-form-group">
                                <label class="admin-label" for="item-price">Prijs (&euro;) *</label>
                                <input
                                    id="item-price"
                                    v-model="itemForm.priceEuros"
                                    type="number"
                                    step="0.01"
                                    min="0"
                                    class="admin-input"
                                    placeholder="0.00"
                                    required
                                >
                            </div>
                            <div class="admin-form-group">
                                <label class="admin-label" for="item-stock">Voorraad</label>
                                <input
                                    id="item-stock"
                                    v-model="itemForm.stock"
                                    type="number"
                                    min="0"
                                    class="admin-input"
                                    placeholder="Leeg = onbeperkt"
                                >
                                <p class="admin-form-hint">Leeg laten voor onbeperkte voorraad.</p>
                            </div>
                        </div>

                        <div class="admin-form-group">
                            <label class="admin-label">Afbeelding</label>
                            <div class="image-upload">
                                <div v-if="imagePreview" class="image-preview">
                                    <img :src="imagePreview" alt="Voorbeeld">
                                </div>
                                <input
                                    ref="fileInput"
                                    type="file"
                                    accept="image/*"
                                    class="admin-input"
                                    @change="handleImageChange"
                                >
                            </div>
                        </div>

                        <div class="form-section-title">Levering <span class="form-required-hint">minimaal 1 vereist</span></div>
                        <div class="admin-form-group">
                            <label class="admin-checkbox">
                                <input
                                    type="checkbox"
                                    v-model="itemForm.allow_pickup"
                                >
                                Afhalen toestaan
                            </label>
                        </div>
                        <div v-if="itemForm.allow_pickup" class="admin-form-group">
                            <label class="admin-label" for="item-pickup-label">Afhaal-label</label>
                            <input
                                id="item-pickup-label"
                                v-model="itemForm.pickup_label"
                                type="text"
                                class="admin-input"
                                placeholder="bijv. Afhalen in Brugge"
                            >
                        </div>

                        <div class="admin-form-group">
                            <label class="admin-checkbox">
                                <input
                                    type="checkbox"
                                    v-model="itemForm.allow_shipping"
                                >
                                Verzending toestaan
                            </label>
                        </div>
                        <div v-if="itemForm.allow_shipping" class="admin-form-group">
                            <label class="admin-label" for="item-shipping-regions">Verzendregio's</label>
                            <input
                                id="item-shipping-regions"
                                v-model="itemForm.shipping_regions"
                                type="text"
                                class="admin-input"
                                placeholder="bijv. België, Nederland, Luxemburg"
                            >
                        </div>

                        <div class="form-section-title">Opties</div>
                        <div class="admin-form-row">
                            <div class="admin-form-group">
                                <label class="admin-checkbox">
                                    <input
                                        type="checkbox"
                                        v-model="itemForm.auto_confirm"
                                    >
                                    Auto-bevestigen
                                </label>
                                <p class="admin-form-hint">Bestellingen worden automatisch bevestigd na betaling. Schakel uit als je bestellingen manueel wilt controleren voor verzending/afhalen.</p>
                            </div>
                            <div class="admin-form-group">
                                <label class="admin-checkbox">
                                    <input
                                        type="checkbox"
                                        v-model="itemForm.is_active"
                                    >
                                    Actief
                                </label>
                                <p class="admin-form-hint">Alleen actieve items zijn zichtbaar in de webshop. Schakel uit om een item tijdelijk te verbergen zonder het te verwijderen.</p>
                            </div>
                        </div>

                        <div class="admin-form-group">
                            <label class="admin-label" for="item-sort-order">Sorteervolgorde</label>
                            <input
                                id="item-sort-order"
                                v-model="itemForm.sort_order"
                                type="number"
                                class="admin-input"
                            >
                            <p class="admin-form-hint">Lagere waarden worden eerst getoond.</p>
                        </div>
                    </div>

                    <div class="admin-modal-footer">
                        <button
                            v-if="modalMode === 'edit'"
                            class="admin-btn admin-btn-danger"
                            @click="handleDeleteItem(editingItem)"
                            :disabled="savingItem"
                        >
                            Verwijderen
                        </button>
                        <div style="flex: 1;"></div>
                        <button
                            class="admin-btn admin-btn-secondary"
                            @click="closeModal"
                            :disabled="savingItem"
                        >
                            Annuleren
                        </button>
                        <button
                            class="admin-btn admin-btn-primary"
                            @click="handleSaveItem"
                            :disabled="savingItem"
                        >
                            {{ savingItem ? 'Opslaan...' : 'Opslaan' }}
                        </button>
                    </div>
                </div>
            </div>
        </Teleport>
    </AdminLayout>
</template>

<style scoped>
.shop-thumb {
    width: 2.5rem;
    height: 2.5rem;
    object-fit: cover;
    border-radius: var(--admin-radius);
    border: 1px solid var(--admin-border);
}

.fulfillment-badges {
    display: flex;
    flex-wrap: wrap;
    gap: 0.375rem;
}

.admin-form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
}

.form-section-title {
    font-size: 0.75rem;
    font-weight: 600;
    color: var(--admin-text-muted);
    text-transform: uppercase;
    letter-spacing: 0.05em;
    margin: 1.5rem 0 0.75rem;
    padding-top: 1rem;
    border-top: 1px solid var(--admin-border-light);
    display: flex;
    align-items: center;
    gap: 0.5rem;
}

.form-required-hint {
    font-size: 0.6875rem;
    font-weight: 500;
    color: var(--admin-danger);
    text-transform: none;
    letter-spacing: 0;
}

.image-upload {
    display: flex;
    align-items: flex-start;
    gap: 1rem;
    flex-wrap: wrap;
}

.image-preview {
    flex-shrink: 0;
}

.image-preview img {
    width: 5rem;
    height: 5rem;
    object-fit: cover;
    border-radius: var(--admin-radius);
    border: 1px solid var(--admin-border);
}

.image-upload input[type="file"] {
    flex: 1;
    min-width: 12rem;
}

@media (max-width: 768px) {
    .admin-form-row {
        grid-template-columns: 1fr;
    }
}
</style>
