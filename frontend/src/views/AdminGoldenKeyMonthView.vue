<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter, RouterLink } from 'vue-router';
import AdminLayout from '@/components/admin/AdminLayout.vue';
import TipTapEditor from '@/components/TipTapEditor.vue';
import config from '@/data/config.js';
import {
    getAdminGoldenKeyMonth,
    updateGoldenKeyMonth,
    addGoldenKeyHint,
    updateGoldenKeyHint,
    deleteGoldenKeyHint,
} from '@/services/GoldenKeyMonthService';

const route  = useRoute();
const router = useRouter();

const monthId = parseInt(route.params.id, 10);

// ---- Page state ----
const loading   = ref(true);
const month     = ref(null);
const saving    = ref(false);
const saveMsg   = ref('');
const saveError = ref('');

// Live date form
const liveDateLocal = ref('');
const isFound       = ref(false);

// ---- Finder modal ----
const showFinderModal = ref(false);
const finderForm = ref({ name: '', found_date_local: '', image: null, image_preview: '', existing_image: '' });
const finderFileInput   = ref(null);
const finderSaving  = ref(false);
const finderError   = ref('');

// ---- Hint modal (add/edit) ----
const showHintModal  = ref(false);
const hintModalMode  = ref('add');   // 'add' | 'edit'
const editingHintId  = ref(null);
const hintEditorRef  = ref(null);
const hintEditorKey  = ref(0);
const hintInitContent = ref('');
const hintImageFile  = ref(null);
const hintFileInput  = ref(null);
const hintImagePreview = ref('');
const hintExistingImage = ref('');
const hintSaving     = ref(false);
const hintError      = ref('');

// ---- API helper ----
function getToken() { return localStorage.getItem('admin_token'); }

async function apiRequest(endpoint, options = {}) {
    const token = getToken();
    const headers = {
        'Accept': 'application/json',
        ...(!(options.body instanceof FormData) && { 'Content-Type': 'application/json' }),
        ...(token && { 'Authorization': `Bearer ${token}` }),
    };
    return fetch(`${config.apiUrl}${endpoint}`, { ...options, headers });
}

async function uploadFile(file) {
    const fd = new FormData();
    fd.append('image', file);
    const res = await apiRequest('admin/upload-image', { method: 'POST', body: fd });
    if (res?.ok) {
        const data = await res.json();
        return data.filename;
    }
    return null;
}

// ---- Helpers ----
function toLocalDatetimeInput(utcIso) {
    if (!utcIso) return '';
    const d = new Date(utcIso);
    const local = new Date(d.getTime() - d.getTimezoneOffset() * 60000);
    return local.toISOString().slice(0, 16);
}

// ---- Load ----
async function loadMonth() {
    loading.value = true;
    const data = await getAdminGoldenKeyMonth(monthId);
    if (!data || data.access_denied) {
        router.push('/admin/golden-key');
        return;
    }
    month.value = data;
    liveDateLocal.value = toLocalDatetimeInput(data.live_date);
    isFound.value = data.is_found;
    loading.value = false;
}

onMounted(loadMonth);

// ---- Save date/state ----
async function saveSettings() {
    saveMsg.value = '';
    saveError.value = '';
    saving.value = true;
    const d = new Date(liveDateLocal.value);
    if (isNaN(d)) {
        saveError.value = 'Ongeldige datum.';
        saving.value = false;
        return;
    }
    const res = await updateGoldenKeyMonth(monthId, {
        live_date:    d.toISOString(),
        is_found:     isFound.value,
        finder_name:  month.value.finder_name  || '',
        finder_image: month.value.finder_image || '',
        found_date:   month.value.found_date   || '',
    });
    if (res?.success) {
        month.value = res.data;
        saveMsg.value = 'Opgeslagen.';
    } else {
        saveError.value = 'Opslaan mislukt.';
    }
    saving.value = false;
}

// ---- Finder modal ----
function openFinderModal() {
    finderError.value = '';
    finderForm.value = {
        name:              month.value.finder_name  || '',
        found_date_local:  toLocalDatetimeInput(month.value.found_date) || '',
        image:             null,
        image_preview:     month.value.finder_image ? `${config.apiUrl}images/${month.value.finder_image}` : '',
        existing_image:    month.value.finder_image || '',
    };
    showFinderModal.value = true;
}

function closeFinderModal() {
    showFinderModal.value = false;
    if (finderFileInput.value) finderFileInput.value.value = '';
    finderForm.value.image = null;
}

function onFinderFileChange(e) {
    const file = e.target.files?.[0];
    if (!file) return;
    finderForm.value.image = file;
    const reader = new FileReader();
    reader.onload = ev => { finderForm.value.image_preview = ev.target.result; };
    reader.readAsDataURL(file);
}

async function saveFinder() {
    finderError.value = '';
    finderSaving.value = true;

    let imageFilename = finderForm.value.existing_image;
    if (finderForm.value.image) {
        const uploaded = await uploadFile(finderForm.value.image);
        if (!uploaded) {
            finderError.value = 'Afbeelding uploaden mislukt.';
            finderSaving.value = false;
            return;
        }
        imageFilename = uploaded;
    }

    const foundDate = finderForm.value.found_date_local
        ? new Date(finderForm.value.found_date_local).toISOString()
        : '';

    const res = await updateGoldenKeyMonth(monthId, {
        live_date:    month.value.live_date,
        is_found:     true,
        finder_name:  finderForm.value.name,
        finder_image: imageFilename,
        found_date:   foundDate,
    });

    if (res?.success) {
        month.value = res.data;
        isFound.value = true;
        showFinderModal.value = false;
        if (finderFileInput.value) finderFileInput.value.value = '';
    } else {
        finderError.value = 'Opslaan mislukt.';
    }
    finderSaving.value = false;
}

// ---- Hint modal ----
function openAddHintModal() {
    hintModalMode.value = 'add';
    editingHintId.value = null;
    hintInitContent.value = '';
    hintImageFile.value = null;
    hintImagePreview.value = '';
    hintExistingImage.value = '';
    hintError.value = '';
    hintEditorKey.value++;
    showHintModal.value = true;
}

function openEditHintModal(hint) {
    hintModalMode.value = 'edit';
    editingHintId.value = hint.id;
    hintInitContent.value = hint.content || '';
    hintImageFile.value = null;
    hintExistingImage.value = hint.image_url || '';
    hintImagePreview.value = hint.image_url ? `${config.apiUrl}images/${hint.image_url}` : '';
    hintError.value = '';
    hintEditorKey.value++;
    showHintModal.value = true;
}

function closeHintModal() {
    showHintModal.value = false;
    hintImageFile.value = null;
    if (hintFileInput.value) hintFileInput.value.value = '';
}

function onHintFileChange(e) {
    const file = e.target.files?.[0];
    if (!file) return;
    hintImageFile.value = file;
    const reader = new FileReader();
    reader.onload = ev => { hintImagePreview.value = ev.target.result; };
    reader.readAsDataURL(file);
}

function clearHintImage() {
    hintImageFile.value = null;
    hintImagePreview.value = '';
    hintExistingImage.value = '';
    if (hintFileInput.value) hintFileInput.value.value = '';
}

async function saveHint() {
    hintError.value = '';
    hintSaving.value = true;

    const content = hintEditorRef.value?.getContent()?.description ?? '';

    let imageFilename = hintExistingImage.value;
    if (hintImageFile.value) {
        const uploaded = await uploadFile(hintImageFile.value);
        if (!uploaded) {
            hintError.value = 'Afbeelding uploaden mislukt.';
            hintSaving.value = false;
            return;
        }
        imageFilename = uploaded;
    }

    let res;
    if (hintModalMode.value === 'add') {
        res = await addGoldenKeyHint(monthId, { content, image_url: imageFilename });
    } else {
        res = await updateGoldenKeyHint(editingHintId.value, { content, image_url: imageFilename });
    }

    if (res?.success) {
        // Refresh hints list
        const fresh = await getAdminGoldenKeyMonth(monthId);
        if (fresh) month.value.hints = fresh.hints || [];
        showHintModal.value = false;
        if (hintFileInput.value) hintFileInput.value.value = '';
    } else {
        hintError.value = 'Opslaan mislukt.';
    }
    hintSaving.value = false;
}

async function deleteHint(hintId) {
    if (!confirm('Deze hint verwijderen?')) return;
    const res = await deleteGoldenKeyHint(hintId);
    if (res?.success) {
        month.value.hints = (month.value.hints || []).filter(h => h.id !== hintId);
    }
}

function stateLabel(state) {
    if (state === 'found')  return '🏆 Gevonden';
    if (state === 'active') return '🔓 Actief';
    return '🔒 Vergrendeld';
}
</script>

<template>
    <AdminLayout :pageTitle="month ? `Golden Key – ${month.month_name}` : 'Golden Key'">
        <div class="gkm-admin">

            <!-- Back -->
            <RouterLink to="/admin/golden-key" class="gkm-admin__back">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"
                    stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                    <path d="M19 12H5M12 5l-7 7 7 7"/>
                </svg>
                Terug naar Golden Key
            </RouterLink>

            <!-- Loading -->
            <div v-if="loading" class="loading-state">
                <div class="spinner"></div>
                <p>Laden…</p>
            </div>

            <template v-else-if="month">

                <!-- Header -->
                <div class="page-header">
                    <h1 class="page-title">{{ month.month_name }}</h1>
                    <span class="state-badge" :class="`state-${month.state}`">{{ stateLabel(month.state) }}</span>
                </div>

                <!-- Settings card -->
                <div class="card">
                    <h2 class="card-title">Live datum &amp; status</h2>
                    <div class="form-row">
                        <div class="admin-form-group">
                            <label class="admin-label">Live datum (lokale tijd)</label>
                            <input type="datetime-local" v-model="liveDateLocal" class="admin-input" />
                        </div>
                        <div class="admin-form-group" style="justify-content: flex-end; padding-top: 1.6rem;">
                            <label class="toggle-label">
                                <input type="checkbox" v-model="isFound" class="toggle-check" />
                                <span class="toggle-text">Gevonden</span>
                            </label>
                        </div>
                    </div>
                    <div v-if="saveError" class="alert alert-error">{{ saveError }}</div>
                    <div v-if="saveMsg"   class="alert alert-success">{{ saveMsg }}</div>
                    <button class="admin-btn admin-btn-primary" @click="saveSettings" :disabled="saving">
                        {{ saving ? 'Opslaan…' : 'Opslaan' }}
                    </button>
                </div>

                <!-- Finder card -->
                <div class="card">
                    <h2 class="card-title">Vinder</h2>
                    <div v-if="month.finder_name" class="finder-info">
                        <img
                            v-if="month.finder_image"
                            :src="`${config.apiUrl}images/${month.finder_image}`"
                            alt=""
                            class="finder-thumb"
                        />
                        <div class="finder-details">
                            <p class="finder-name">{{ month.finder_name }}</p>
                            <p v-if="month.found_date" class="finder-date">
                                {{ new Date(month.found_date).toLocaleString('nl-BE', { timeZoneName: 'short' }) }}
                            </p>
                        </div>
                    </div>
                    <p v-else class="card-hint">Nog geen vinder ingesteld.</p>
                    <button class="admin-btn admin-btn-secondary" style="margin-top:0.75rem;" @click="openFinderModal">
                        {{ month.finder_name ? 'Vinder bewerken' : 'Vinder instellen' }}
                    </button>
                </div>

                <!-- Hints card -->
                <div class="card">
                    <div class="card-top-row">
                        <h2 class="card-title" style="margin-bottom:0;">Hints</h2>
                        <button class="admin-btn admin-btn-primary admin-btn-sm" @click="openAddHintModal">
                            + Hint toevoegen
                        </button>
                    </div>
                    <p class="card-hint" style="margin-top:0.4rem;">
                        Elke hint bestaat uit een rijke tekst en een optionele afbeelding.
                    </p>

                    <div v-if="!month.hints || month.hints.length === 0" class="card-hint">
                        Nog geen hints.
                    </div>

                    <div v-else class="hints-list">
                        <div
                            v-for="(hint, idx) in month.hints"
                            :key="hint.id"
                            class="hint-row"
                        >
                            <span class="hint-num">{{ idx + 1 }}</span>
                            <div class="hint-preview">
                                <p class="hint-preview-text">
                                    {{ hint.content ? '(rijke tekst)' : '(leeg)' }}
                                </p>
                                <img
                                    v-if="hint.image_url"
                                    :src="`${config.apiUrl}images/${hint.image_url}`"
                                    alt=""
                                    class="hint-thumb"
                                />
                            </div>
                            <div class="hint-actions">
                                <button class="admin-btn admin-btn-sm admin-btn-secondary" @click="openEditHintModal(hint)">
                                    Bewerken
                                </button>
                                <button class="admin-btn admin-btn-sm admin-btn-danger" @click="deleteHint(hint.id)">
                                    Verwijder
                                </button>
                            </div>
                        </div>
                    </div>
                </div>

            </template>

        </div>

        <!-- ===== Finder modal ===== -->
        <Teleport to="body">
            <div v-if="showFinderModal" class="admin-modal-overlay" @click.self="closeFinderModal">
                <div class="admin-modal admin-modal-lg">
                    <div class="admin-modal-header">
                        <h2 class="admin-modal-title">Vinder instellen</h2>
                        <button class="admin-modal-close" @click="closeFinderModal" aria-label="Sluiten">
                            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 6L6 18M6 6l12 12"/></svg>
                        </button>
                    </div>
                    <div class="admin-modal-body">
                        <div class="admin-form-group">
                            <label class="admin-label">Naam vinder *</label>
                            <input type="text" v-model="finderForm.name" class="admin-input" placeholder="Naam van de vinder" />
                        </div>
                        <div class="admin-form-group">
                            <label class="admin-label">Gevonden op</label>
                            <input type="datetime-local" v-model="finderForm.found_date_local" class="admin-input" />
                        </div>
                        <div class="admin-form-group">
                            <label class="admin-label">Foto vinder (optioneel)</label>
                            <div class="image-upload-area">
                                <img v-if="finderForm.image_preview" :src="finderForm.image_preview" class="image-preview" />
                                <div v-else class="image-placeholder">
                                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                                        <rect x="3" y="3" width="18" height="18" rx="2"/>
                                        <circle cx="8.5" cy="8.5" r="1.5"/>
                                        <polyline points="21 15 16 10 5 21"/>
                                    </svg>
                                    <span>Geen afbeelding geselecteerd</span>
                                </div>
                            </div>
                            <input
                                ref="finderFileInput"
                                type="file"
                                accept="image/*"
                                @change="onFinderFileChange"
                                class="admin-input"
                                style="margin-top:0.6rem;"
                            />
                        </div>
                        <div v-if="finderError" class="alert alert-error">{{ finderError }}</div>
                    </div>
                    <div class="admin-modal-footer">
                        <button class="admin-btn admin-btn-secondary" @click="closeFinderModal" :disabled="finderSaving">Annuleren</button>
                        <button class="admin-btn admin-btn-primary"   @click="saveFinder"       :disabled="finderSaving">
                            {{ finderSaving ? 'Opslaan…' : 'Opslaan' }}
                        </button>
                    </div>
                </div>
            </div>
        </Teleport>

        <!-- ===== Hint modal ===== -->
        <Teleport to="body">
            <div v-if="showHintModal" class="admin-modal-overlay" @click.self="closeHintModal">
                <div class="admin-modal admin-modal-lg">
                    <div class="admin-modal-header">
                        <h2 class="admin-modal-title">{{ hintModalMode === 'add' ? 'Hint toevoegen' : 'Hint bewerken' }}</h2>
                        <button class="admin-modal-close" @click="closeHintModal" aria-label="Sluiten">
                            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 6L6 18M6 6l12 12"/></svg>
                        </button>
                    </div>
                    <div class="admin-modal-body">
                        <div class="admin-form-group">
                            <label class="admin-label">Tekst (rijke tekst)</label>
                            <TipTapEditor
                                :key="hintEditorKey"
                                :content="hintInitContent"
                                :editable="true"
                                ref="hintEditorRef"
                            />
                        </div>
                        <div class="admin-form-group">
                            <label class="admin-label">Afbeelding (optioneel)</label>
                            <div class="image-upload-area">
                                <img v-if="hintImagePreview" :src="hintImagePreview" class="image-preview" />
                                <div v-else class="image-placeholder">
                                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                                        <rect x="3" y="3" width="18" height="18" rx="2"/>
                                        <circle cx="8.5" cy="8.5" r="1.5"/>
                                        <polyline points="21 15 16 10 5 21"/>
                                    </svg>
                                    <span>Geen afbeelding geselecteerd</span>
                                </div>
                            </div>
                            <div class="upload-row">
                                <input
                                    ref="hintFileInput"
                                    type="file"
                                    accept="image/*"
                                    @change="onHintFileChange"
                                    class="admin-input"
                                />
                                <button
                                    v-if="hintImagePreview"
                                    type="button"
                                    class="admin-btn admin-btn-sm admin-btn-ghost"
                                    @click="clearHintImage"
                                >
                                    Verwijder afbeelding
                                </button>
                            </div>
                        </div>
                        <div v-if="hintError" class="alert alert-error">{{ hintError }}</div>
                    </div>
                    <div class="admin-modal-footer">
                        <button class="admin-btn admin-btn-secondary" @click="closeHintModal" :disabled="hintSaving">Annuleren</button>
                        <button class="admin-btn admin-btn-primary"   @click="saveHint"       :disabled="hintSaving">
                            {{ hintSaving ? 'Opslaan…' : 'Opslaan' }}
                        </button>
                    </div>
                </div>
            </div>
        </Teleport>

    </AdminLayout>
</template>

<style scoped>
.gkm-admin {
    max-width: 720px;
    display: flex;
    flex-direction: column;
    gap: 0;
}

/* Back link */
.gkm-admin__back {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    font-size: 0.85rem;
    color: var(--admin-text-secondary);
    text-decoration: none;
    margin-bottom: 1.5rem;
    transition: color 0.15s;
}
.gkm-admin__back:hover { color: var(--admin-text); }

/* Page header */
.page-header {
    display: flex;
    align-items: center;
    gap: 1rem;
    flex-wrap: wrap;
    margin-bottom: 1.5rem;
}

.page-title {
    font-size: 1.75rem;
    font-weight: 700;
    color: var(--admin-text);
    margin: 0;
}

/* State badge */
.state-badge {
    display: inline-block;
    padding: 0.3rem 0.8rem;
    border-radius: 9999px;
    font-size: 0.8rem;
    font-weight: 600;
}
.state-locked  { background: color-mix(in srgb, #f59e0b 12%, transparent); color: #b45309; }
.state-active  { background: color-mix(in srgb, #22c55e 12%, transparent); color: #15803d; }
.state-found   { background: color-mix(in srgb, #f59e0b 18%, transparent); color: #92400e; }

/* Cards */
.card {
    background: var(--admin-surface);
    border: 1px solid var(--admin-border);
    border-radius: var(--admin-radius-lg, 0.75rem);
    padding: 1.5rem;
    margin-bottom: 1.5rem;
}

.card-top-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 0.75rem;
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
    line-height: 1.5;
    margin-bottom: 0.75rem;
}

/* Form rows */
.form-row {
    display: grid;
    grid-template-columns: 1fr auto;
    gap: 1rem;
    align-items: start;
    margin-bottom: 1rem;
}

/* Toggle */
.toggle-label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
    padding-top: 0.3rem;
}
.toggle-check { accent-color: var(--admin-primary); width: 1rem; height: 1rem; }
.toggle-text  { font-size: 0.9rem; font-weight: 500; color: var(--admin-text); }

/* Finder info display */
.finder-info {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1rem;
    background: var(--admin-surface-hover, rgba(0,0,0,0.03));
    border: 1px solid var(--admin-border);
    border-radius: var(--admin-radius);
}
.finder-thumb {
    width: 54px;
    height: 54px;
    border-radius: 50%;
    object-fit: cover;
    flex-shrink: 0;
}
.finder-details { flex: 1; }
.finder-name { font-weight: 600; color: var(--admin-text); margin: 0 0 0.15rem; }
.finder-date { font-size: 0.82rem; color: var(--admin-text-secondary); margin: 0; }

/* Hints list */
.hints-list {
    display: flex;
    flex-direction: column;
    gap: 0.6rem;
}

.hint-row {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem 1rem;
    background: var(--admin-surface-hover, rgba(0,0,0,0.03));
    border: 1px solid var(--admin-border);
    border-radius: var(--admin-radius);
}

.hint-num {
    font-size: 0.85rem;
    font-weight: 700;
    color: var(--admin-text-secondary);
    min-width: 1.5rem;
    flex-shrink: 0;
}

.hint-preview {
    flex: 1;
    min-width: 0;
    display: flex;
    align-items: center;
    gap: 0.75rem;
}

.hint-preview-text {
    font-size: 0.85rem;
    color: var(--admin-text-secondary);
    margin: 0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.hint-thumb {
    width: 36px;
    height: 36px;
    object-fit: cover;
    border-radius: 3px;
    flex-shrink: 0;
}

.hint-actions {
    display: flex;
    gap: 0.4rem;
    flex-shrink: 0;
}

/* Image upload */
.image-upload-area {
    border: 1px solid var(--admin-border);
    border-radius: var(--admin-radius);
    overflow: hidden;
    background: var(--admin-surface);
    min-height: 120px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.image-preview {
    width: 100%;
    max-height: 220px;
    object-fit: contain;
    display: block;
}

.image-placeholder {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
    padding: 1.25rem;
    color: var(--admin-text-secondary);
}
.image-placeholder svg { width: 2.5rem; height: 2.5rem; opacity: 0.4; }
.image-placeholder span { font-size: 0.82rem; }

.upload-row {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    flex-wrap: wrap;
    margin-top: 0.5rem;
}
.upload-row .admin-input { flex: 1; min-width: 0; margin: 0; }

/* Alerts */
.alert {
    padding: 0.65rem 0.9rem;
    border-radius: var(--admin-radius);
    font-size: 0.875rem;
    margin-top: 0.75rem;
}
.alert-success { background: color-mix(in srgb, #22c55e 12%, transparent); color: #15803d; border: 1px solid color-mix(in srgb, #22c55e 30%, transparent); }
.alert-error   { background: color-mix(in srgb, #ef4444 12%, transparent); color: #b91c1c; border: 1px solid color-mix(in srgb, #ef4444 30%, transparent); }

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
@keyframes spin { to { transform: rotate(360deg); } }
</style>
