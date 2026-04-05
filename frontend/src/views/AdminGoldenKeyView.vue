<script setup>
import { ref, computed, onMounted } from 'vue';
import AdminLayout from '@/components/admin/AdminLayout.vue';
import config from '@/data/config.js';

const loading = ref(true);
const saving = ref(false);
const successMsg = ref('');
const errorMsg = ref('');

// Stored activation time (UTC ISO from API)
const activationTimeUTC = ref('');

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
            body: JSON.stringify({ activation_time: utcISO })
        });

        if (res?.ok) {
            const data = await res.json();
            activationTimeUTC.value = data.activation_time;
            localDatetimeInput.value = toLocalDatetimeInput(data.activation_time);
            successMsg.value = 'Instellingen opgeslagen.';
        } else {
            errorMsg.value = 'Opslaan mislukt. Probeer opnieuw.';
        }
    } catch (err) {
        errorMsg.value = 'Er is een fout opgetreden.';
    }
    saving.value = false;
}

onMounted(fetchSettings);
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
                        <span class="status-label">Huidige status</span>
                        <span class="status-badge" :class="isActive ? 'status-active' : 'status-soon'">
                            {{ isActive ? 'Actief' : 'Binnenkort' }}
                        </span>
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

                    <div v-if="successMsg" class="alert alert-success">{{ successMsg }}</div>
                    <div v-if="errorMsg" class="alert alert-error">{{ errorMsg }}</div>

                    <div class="form-actions">
                        <button class="btn btn-primary" @click="saveSettings" :disabled="saving">
                            {{ saving ? 'Opslaan…' : 'Opslaan' }}
                        </button>
                    </div>
                </div>
            </template>
        </div>
    </AdminLayout>
</template>

<style scoped>
.gk-admin {
    max-width: 640px;
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
.status-card .status-row {
    display: flex;
    align-items: center;
    gap: 1rem;
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
</style>
