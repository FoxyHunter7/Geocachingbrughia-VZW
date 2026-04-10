<script setup>
import { ref, computed, onMounted } from 'vue';
import { RouterLink } from 'vue-router';
import AdminLayout from '@/components/admin/AdminLayout.vue';
import config from '@/data/config.js';
import { getAdminGoldenKeyMonths } from '@/services/GoldenKeyMonthService';

const loading = ref(true);
const saving = ref(false);
const successMsg = ref('');
const errorMsg = ref('');

// Stored activation time (UTC ISO from API)
const activationTimeUTC = ref('');
const bannerTexts = ref({});
const rulesTexts = ref({});

// Languages for per-language banner text
const languages = ref([]);
const showBannerModal = ref(false);
const draftTexts = ref({});
const modalSaving = ref(false);
const modalError = ref('');

const hasBannerText = computed(() =>
    Object.values(bannerTexts.value).some(t => t && t.trim() !== '')
);

function openBannerModal() {
    draftTexts.value = { ...bannerTexts.value };
    modalError.value = '';
    showBannerModal.value = true;
}

function closeBannerModal() {
    showBannerModal.value = false;
}

async function saveBannerModal() {
    modalError.value = '';
    modalSaving.value = true;
    try {
        const localDate = new Date(localDatetimeInput.value);
        const utcISO = isNaN(localDate.getTime())
            ? new Date(activationTimeUTC.value).toISOString()
            : localDate.toISOString();

        const res = await apiRequest('admin/golden-key', {
            method: 'PUT',
            body: JSON.stringify({ activation_time: utcISO, banner_text: draftTexts.value, rules: rulesTexts.value })
        });

        if (res?.ok) {
            const data = await res.json();
            activationTimeUTC.value = data.activation_time;
            localDatetimeInput.value = toLocalDatetimeInput(data.activation_time);
            bannerTexts.value = data.banner_text || {};
            rulesTexts.value = data.rules || {};
            showBannerModal.value = false;
        } else {
            modalError.value = 'Opslaan mislukt. Probeer opnieuw.';
        }
    } catch {
        modalError.value = 'Er is een fout opgetreden.';
    }
    modalSaving.value = false;
}

// The datetime-local input value (local time string, no timezone)
const localDatetimeInput = ref('');

const isActive = computed(() => {
    if (!activationTimeUTC.value) return false;
    return new Date() >= new Date(activationTimeUTC.value);
});

function getToken() {
    return localStorage.getItem('admin_token');
}

async function apiRequest(endpoint, options = {}) {
    const token = getToken();
    const headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        ...(token && { 'Authorization': `Bearer ${token}` }),
        ...options.headers
    };
    const response = await fetch(`${config.apiUrl}${endpoint}`, { ...options, headers });
    return response;
}

// Convert UTC ISO string to "YYYY-MM-DDTHH:mm" in local time for <input type="datetime-local">
function toLocalDatetimeInput(utcIso) {
    const d = new Date(utcIso);
    // Shift to local time by subtracting the UTC offset
    const local = new Date(d.getTime() - d.getTimezoneOffset() * 60000);
    return local.toISOString().slice(0, 16);
}

async function fetchSettings() {
    loading.value = true;
    try {
        const res = await apiRequest('admin/golden-key');
        if (res?.ok) {
            const data = await res.json();
            activationTimeUTC.value = data.activation_time;
            localDatetimeInput.value = toLocalDatetimeInput(data.activation_time);
            bannerTexts.value = data.banner_text || {};
            rulesTexts.value = data.rules || {};
        }
    } catch (err) {
        errorMsg.value = 'Kon de instellingen niet ophalen.';
    }
    loading.value = false;
}

async function saveSettings() {
    successMsg.value = '';
    errorMsg.value = '';
    saving.value = true;
    try {
        // Convert local datetime input back to UTC ISO (RFC3339)
        const localDate = new Date(localDatetimeInput.value);
        if (isNaN(localDate.getTime())) {
            errorMsg.value = 'Ongeldige datum/tijd.';
            saving.value = false;
            return;
        }
        const utcISO = localDate.toISOString();

        const res = await apiRequest('admin/golden-key', {
            method: 'PUT',
            body: JSON.stringify({ activation_time: utcISO, banner_text: bannerTexts.value, rules: rulesTexts.value })
        });

        if (res?.ok) {
            const data = await res.json();
            activationTimeUTC.value = data.activation_time;
            localDatetimeInput.value = toLocalDatetimeInput(data.activation_time);
            bannerTexts.value = data.banner_text || {};
            rulesTexts.value = data.rules || {};
            successMsg.value = 'Instellingen opgeslagen.';
        } else {
            errorMsg.value = 'Opslaan mislukt. Probeer opnieuw.';
        }
    } catch (err) {
        errorMsg.value = 'Er is een fout opgetreden.';
    }
    saving.value = false;
}

async function fetchLanguages() {
    try {
        const res = await apiRequest('admin/languages');
        if (res?.ok) {
            const data = await res.json();
            languages.value = data.data || data || [];
        }
    } catch { /* languages stay empty */ }
}

// ---- Months (read-only list; editing happens on the month detail page) ----
const months = ref([]);
const fetchingMonths = ref(false);

function monthStateLabel(state) {
    if (state === 'found')  return '🏆 Gevonden';
    if (state === 'active') return '🔓 Actief';
    return '🔒 Vergrendeld';
}
function monthStateCss(state) {
    if (state === 'found')  return 'badge-found';
    if (state === 'active') return 'badge-active';
    return 'badge-locked';
}

async function fetchMonths() {
    fetchingMonths.value = true;
    const data = await getAdminGoldenKeyMonths();
    months.value = Array.isArray(data) ? data : [];
    fetchingMonths.value = false;
}

onMounted(async () => {
    await Promise.all([fetchSettings(), fetchLanguages(), fetchMonths()]);
});
</script>

<template>
    <AdminLayout pageTitle="Golden Key">
        <div class="gk-admin">
            <div class="page-header">
                <h1 class="page-title">Golden Key</h1>
                <p class="page-subtitle">Beheer de activatiestatus van het Golden Key systeem.</p>
            </div>

            <div v-if="loading" class="loading-state">
                <div class="spinner"></div>
                <p>Laden…</p>
            </div>

            <template v-else>
                <!-- Status card -->
                <div class="card status-card">
                    <h2 class="card-title">Status</h2>
                    <div class="status-row">
                        <div class="status-left">
                            <span class="status-label">Huidige status</span>
                            <span class="status-badge" :class="isActive ? 'status-active' : 'status-soon'">
                                {{ isActive ? 'Actief' : 'Binnenkort' }}
                            </span>
                        </div>
                        <div class="status-right">
                            <span v-if="hasBannerText" class="banner-set-indicator" title="Bannertekst ingesteld">
                                <svg width="13" height="13" viewBox="0 0 16 16" fill="currentColor" aria-hidden="true"><path d="M13.5 3.5a1.5 1.5 0 0 1 0 2.12l-7.25 7.25a1 1 0 0 1-.39.24l-3 1a.5.5 0 0 1-.63-.63l1-3a1 1 0 0 1 .24-.39L10.88 3.5a1.5 1.5 0 0 1 2.12 0z"/></svg>
                                Bannertekst ingesteld
                            </span>
                            <button class="admin-btn admin-btn-sm admin-btn-secondary" @click="openBannerModal">
                                Bannertekst bewerken
                            </button>
                        </div>
                    </div>
                </div>

                <!-- Banner text modal -->
                <div v-if="showBannerModal" class="admin-modal-overlay" @click.self="closeBannerModal">
                    <div class="admin-modal admin-modal-lg">
                        <div class="admin-modal-header">
                            <h2 class="admin-modal-title">Bannertekst bewerken</h2>
                            <button class="admin-modal-close" @click="closeBannerModal" aria-label="Sluiten">
                                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 6L6 18M6 6l12 12"/></svg>
                            </button>
                        </div>
                        <div class="admin-modal-body">
                            <p class="modal-hint">
                                Deze tekst verschijnt naast de afbeelding in de homepage-banner wanneer Golden Key actief is,
                                in de taal van de bezoeker. Laat een veld leeg voor talen zonder tekst.
                            </p>
                            <div
                                v-for="lang in languages"
                                :key="lang.code"
                                class="form-group"
                            >
                                <label :for="'draft-text-' + lang.code" class="form-label">
                                    {{ lang.name }} ({{ lang.code }})
                                </label>
                                <textarea
                                    :id="'draft-text-' + lang.code"
                                    v-model="draftTexts[lang.code]"
                                    class="form-input form-textarea"
                                    rows="3"
                                    :placeholder="'Bannertekst in ' + lang.name + '…'"
                                ></textarea>
                            </div>
                            <p v-if="languages.length === 0" class="modal-hint">Geen talen gevonden.</p>
                            <div v-if="modalError" class="alert alert-error">{{ modalError }}</div>
                        </div>
                        <div class="admin-modal-footer">
                            <button class="admin-btn admin-btn-secondary" @click="closeBannerModal">Annuleren</button>
                            <button class="admin-btn admin-btn-primary" @click="saveBannerModal" :disabled="modalSaving">
                                {{ modalSaving ? 'Opslaan…' : 'Opslaan' }}
                            </button>
                        </div>
                    </div>
                </div>

                <!-- Activation time card -->
                <div class="card">
                    <h2 class="card-title">Activatiedatum &amp; -tijd</h2>
                    <p class="card-hint">
                        Datum en tijd worden weergegeven in de lokale tijdzone van de browser.
                        De opgegeven tijd wordt intern als UTC opgeslagen.
                        <br>Huidige instelling: <strong>{{ new Date(activationTimeUTC).toLocaleString('nl-BE', { timeZoneName: 'short' }) }}</strong>
                    </p>

                    <div class="form-group">
                        <label for="activation-time" class="form-label">Activatie op</label>
                        <input
                            id="activation-time"
                            type="datetime-local"
                            v-model="localDatetimeInput"
                            class="form-input"
                        />
                    </div>
                </div>

                <div v-if="successMsg" class="alert alert-success">{{ successMsg }}</div>
                <div v-if="errorMsg" class="alert alert-error">{{ errorMsg }}</div>

                <div class="form-actions">
                    <button class="btn btn-primary" @click="saveSettings" :disabled="saving">
                        {{ saving ? 'Opslaan…' : 'Opslaan' }}
                    </button>
                </div>

                <!-- Months management -->
                <div class="card months-card">
                    <h2 class="card-title">Maanden</h2>
                    <p class="card-hint">
                        Klik op "Bewerken" om de live datum, vindersinformatie en hints van een maand te beheren.
                    </p>

                    <div v-if="fetchingMonths" class="loading-state">
                        <div class="spinner"></div>
                    </div>

                    <table v-else class="months-table">
                        <thead>
                            <tr>
                                <th>#</th>
                                <th>Maand</th>
                                <th>Live datum</th>
                                <th>Status</th>
                                <th></th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="month in months" :key="month.id">
                                <td class="month-col month-num">{{ month.month_number }}</td>
                                <td class="month-col month-name">{{ month.month_name }}</td>
                                <td class="month-col month-date">{{ new Date(month.live_date).toLocaleString('nl-BE', { timeZoneName: 'short' }) }}</td>
                                <td class="month-col">
                                    <span class="month-badge" :class="monthStateCss(month.state)">
                                        {{ monthStateLabel(month.state) }}
                                    </span>
                                </td>
                                <td class="month-col month-actions">
                                    <RouterLink
                                        :to="`/admin/golden-key/months/${month.id}`"
                                        class="admin-btn admin-btn-sm admin-btn-secondary"
                                    >
                                        Bewerken
                                    </RouterLink>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <!-- Rules card -->
                <div class="card">
                    <h2 class="card-title">Spelregels</h2>
                    <p class="card-hint">
                        De tekst die op de pagina "Spelregels" verschijnt, per taal van de bezoeker.
                        Eenvoudige HTML is toegestaan (bijv. &lt;b&gt;, &lt;ul&gt;, &lt;li&gt;).
                    </p>
                    <div
                        v-for="lang in languages"
                        :key="lang.code"
                        class="form-group"
                    >
                        <label :for="'rules-text-' + lang.code" class="form-label">
                            {{ lang.name }} ({{ lang.code }})
                        </label>
                        <textarea
                            :id="'rules-text-' + lang.code"
                            v-model="rulesTexts[lang.code]"
                            class="form-input form-textarea"
                            rows="6"
                            :placeholder="'Spelregels in ' + lang.name + '…'"
                        ></textarea>
                    </div>
                    <p v-if="languages.length === 0" class="card-hint">Geen talen gevonden.</p>
                </div>
            </template>
        </div>
    </AdminLayout>
</template>

<style scoped>
/* ---- Months card ---- */
.gk-admin {
    max-width: 900px;
}

.page-header {
    margin-bottom: 2rem;
}

.page-title {
    font-size: 1.75rem;
    font-weight: 700;
    color: var(--admin-text);
    margin-bottom: 0.25rem;
}

.page-subtitle {
    color: var(--admin-text-secondary);
    font-size: 0.95rem;
}

/* Cards */
.card {
    background: var(--admin-surface);
    border: 1px solid var(--admin-border);
    border-radius: var(--admin-radius-lg, 0.75rem);
    padding: 1.5rem;
    margin-bottom: 1.5rem;
}

.card-title {
    font-size: 1rem;
    font-weight: 600;
    color: var(--admin-text);
    margin-bottom: 1rem;
}

.card-hint {
    font-size: 0.85rem;
    color: var(--admin-text-secondary);
    margin-bottom: 1.25rem;
    line-height: 1.5;
}

/* Status */
.status-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.75rem;
    flex-wrap: wrap;
}

.status-left {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    flex-wrap: wrap;
}

.status-right {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    flex-wrap: wrap;
}

.status-label {
    font-size: 0.9rem;
    color: var(--admin-text-secondary);
}

.status-badge {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    padding: 0.3rem 0.85rem;
    border-radius: 9999px;
    font-size: 0.875rem;
    font-weight: 600;
}

.status-badge::before {
    content: '';
    display: inline-block;
    width: 0.5rem;
    height: 0.5rem;
    border-radius: 50%;
    background: currentColor;
}

.status-soon {
    background: color-mix(in srgb, #f59e0b 15%, transparent);
    color: #b45309;
}

.status-active {
    background: color-mix(in srgb, #22c55e 15%, transparent);
    color: #15803d;
}

.banner-set-indicator {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    font-size: 0.8rem;
    color: var(--admin-text-secondary);
}

/* Modal */
.modal-hint {
    font-size: 0.85rem;
    color: var(--admin-text-secondary);
    margin-bottom: 1.25rem;
    line-height: 1.5;
}

/* Form */
.form-group {
    margin-bottom: 1.25rem;
}

.form-label {
    display: block;
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--admin-text);
    margin-bottom: 0.4rem;
}

.form-input {
    width: 100%;
    padding: 0.55rem 0.75rem;
    border: 1px solid var(--admin-border);
    border-radius: var(--admin-radius);
    background: var(--admin-surface);
    color: var(--admin-text);
    font-size: 0.9rem;
    transition: border-color 0.15s;
}

.form-input:focus {
    outline: none;
    border-color: var(--admin-primary);
}

.form-textarea {
    resize: vertical;
    min-height: 80px;
    font-family: inherit;
    line-height: 1.5;
}

.form-actions {
    display: flex;
    gap: 0.75rem;
    margin-top: 0.5rem;
}

/* Alerts */
.alert {
    padding: 0.65rem 0.9rem;
    border-radius: var(--admin-radius);
    font-size: 0.875rem;
    margin-bottom: 1rem;
}

.alert-success {
    background: color-mix(in srgb, #22c55e 12%, transparent);
    color: #15803d;
    border: 1px solid color-mix(in srgb, #22c55e 30%, transparent);
}

.alert-error {
    background: color-mix(in srgb, #ef4444 12%, transparent);
    color: #b91c1c;
    border: 1px solid color-mix(in srgb, #ef4444 30%, transparent);
}

/* Buttons */
.btn {
    padding: 0.55rem 1.25rem;
    border-radius: var(--admin-radius);
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    border: none;
    transition: opacity 0.15s;
}

.btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.btn-primary {
    background: var(--admin-primary);
    color: white;
}

.btn-primary:not(:disabled):hover {
    opacity: 0.88;
}

/* Loading */
.loading-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    padding: 3rem 0;
    color: var(--admin-text-secondary);
}

.spinner {
    width: 2rem;
    height: 2rem;
    border: 3px solid var(--admin-border);
    border-top-color: var(--admin-primary);
    border-radius: 50%;
    animation: spin 0.7s linear infinite;
}

@keyframes spin {
    to { transform: rotate(360deg); }
}

/* ---- Months card ---- */
.gk-admin {
    max-width: 900px;
}

.months-card {
    margin-top: 2rem;
}

.months-table {
    width: 100%;
    border-collapse: collapse;
    font-size: 0.875rem;
}

.months-table th {
    text-align: left;
    font-weight: 600;
    font-size: 0.8rem;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: var(--admin-text-secondary);
    padding: 0.5rem 0.75rem;
    border-bottom: 1px solid var(--admin-border);
}

.month-col {
    padding: 0.65rem 0.75rem;
    border-bottom: 1px solid var(--admin-border);
    vertical-align: middle;
}

.month-num {
    color: var(--admin-text-secondary);
    width: 2rem;
}

.month-name {
    font-weight: 500;
}

.month-date {
    color: var(--admin-text-secondary);
    font-size: 0.82rem;
}

.month-actions {
    text-align: right;
}

/* State badges */
.month-badge {
    display: inline-block;
    padding: 0.25rem 0.6rem;
    border-radius: 9999px;
    font-size: 0.78rem;
    font-weight: 600;
    white-space: nowrap;
}

.badge-locked {
    background: color-mix(in srgb, #f59e0b 12%, transparent);
    color: #b45309;
}

.badge-active {
    background: color-mix(in srgb, #22c55e 12%, transparent);
    color: #15803d;
}

.badge-found {
    background: color-mix(in srgb, #f59e0b 18%, transparent);
    color: #92400e;
}

</style>
