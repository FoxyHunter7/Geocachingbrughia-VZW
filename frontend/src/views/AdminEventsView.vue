<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';
import AdminLayout from '@/components/admin/AdminLayout.vue';
import TipTapEditor from '@/components/TipTapEditor.vue';
import config from '@/data/config.js';
import QRCode from 'qrcode';
import JSZip from 'jszip';

const router = useRouter();

// State
const loading = ref(true);
const events = ref([]);
const languages = ref([]);
const currPage = ref(1);
const lastPage = ref(1);
const totalItems = ref(0);
const search = ref('');
const statusFilter = ref('');

// Modal state
const showModal = ref(false);
const modalMode = ref('create'); // 'create' or 'edit'
const saving = ref(false);
const editingEvent = ref(null);
const editors = ref([]);

// File upload
const fileInput = ref(null);
const imagePreview = ref('');

// Form data
const formData = ref({
    state: 'draft',
    on_home: false,
    title: '',
    geolink: '',
    type: 'REGULAR',
    location: '',
    start_date: '',
    end_date: '',
    ticket_purchase_url: '',
    translations: []
});

// API helpers
function getToken() {
    return localStorage.getItem('admin_token');
}

async function apiRequest(endpoint, options = {}) {
    const token = getToken();
    const headers = {
        'Accept': 'application/json',
        ...(!(options.body instanceof FormData) && { 'Content-Type': 'application/json' }),
        ...(token && { 'Authorization': `Bearer ${token}` }),
        // Disable caching for admin requests
        'Cache-Control': 'no-cache',
        'Pragma': 'no-cache',
        ...options.headers
    };

    try {
        const response = await fetch(`${config.apiUrl}${endpoint}`, {
            ...options,
            headers,
            cache: 'no-store'
        });
        return response;
    } catch (err) {
        console.error('API request failed:', err);
        return null;
    }
}

// Fetch data
async function fetchEvents() {
    loading.value = true;
    try {
        const params = new URLSearchParams({ page: currPage.value });
        if (statusFilter.value) params.append('status', statusFilter.value);
        
        const response = await apiRequest(`admin/events?${params}`);
        if (response?.ok) {
            const data = await response.json();
            if (data.data) {
                events.value = data.data;
                currPage.value = data.current_page || 1;
                lastPage.value = data.last_page || 1;
                totalItems.value = data.total || data.data.length;
            } else if (Array.isArray(data)) {
                events.value = data;
                totalItems.value = data.length;
            }
        }
    } catch (err) {
        console.error('Failed to fetch events:', err);
    }
    loading.value = false;
}

async function fetchLanguages() {
    try {
        const response = await apiRequest('admin/languages');
        if (response?.ok) {
            const data = await response.json();
            languages.value = data.data || data || [];
        }
    } catch (err) {
        console.error('Failed to fetch languages:', err);
    }
}

// Computed
const filteredEvents = computed(() => {
    if (!search.value) return events.value;
    const q = search.value.toLowerCase();
    return events.value.filter(e => 
        e.title?.toLowerCase().includes(q) ||
        e.location?.toLowerCase().includes(q)
    );
});

// Modal handlers
function openCreateModal() {
    modalMode.value = 'create';
    editingEvent.value = null;
    imagePreview.value = '';
    
    formData.value = {
        state: 'draft',
        on_home: false,
        title: '',
        geolink: '',
        type: 'REGULAR',
        location: '',
        start_date: '',
        end_date: '',
        ticket_purchase_url: '',
        translations: languages.value.map(l => ({
            lang_code: l.code,
            description: ''
        }))
    };
    
    showModal.value = true;
}

function openEditModal(event) {
    modalMode.value = 'edit';
    editingEvent.value = event;
    
    // Format dates for datetime-local input
    const formatDate = (d) => {
        if (!d) return '';
        const date = new Date(d);
        return date.toISOString().slice(0, 16);
    };
    
    formData.value = {
        state: event.state || 'draft',
        on_home: event.on_home || false,
        title: event.title || '',
        geolink: event.geolink || '',
        type: event.type || 'REGULAR',
        location: event.location || '',
        start_date: formatDate(event.start_date),
        end_date: formatDate(event.end_date),
        ticket_purchase_url: event.ticket_purchase_url || '',
        translations: event.translations?.length 
            ? event.translations 
            : languages.value.map(l => ({ lang_code: l.code, description: '' }))
    };
    
    imagePreview.value = event.imageUrl ? `${config.apiUrl}images/${event.imageUrl}` : '';
    showModal.value = true;
}

function closeModal() {
    showModal.value = false;
    editingEvent.value = null;
    editors.value = [];
    if (fileInput.value) fileInput.value.value = '';
}

// Store file for upload
const selectedFile = ref(null);

function handleImageChange(e) {
    const file = e.target.files?.[0];
    if (file) {
        selectedFile.value = file;
        const reader = new FileReader();
        reader.onload = (e) => {
            imagePreview.value = e.target.result;
        };
        reader.readAsDataURL(file);
    }
}

async function uploadImage(file) {
    const formData = new FormData();
    formData.append('image', file);
    
    const response = await apiRequest('admin/upload-image', {
        method: 'POST',
        body: formData,
        headers: {} // Let browser set Content-Type with boundary for FormData
    });
    
    if (response?.ok) {
        const data = await response.json();
        return data.filename;
    }
    return null;
}

async function handleSave(publish = false) {
    saving.value = true;
    
    try {
        // Upload image if new file selected
        let imageUrl = editingEvent.value?.imageUrl || '';
        if (selectedFile.value) {
            const uploadedFilename = await uploadImage(selectedFile.value);
            if (uploadedFilename) {
                imageUrl = uploadedFilename;
            } else {
                window.$toast?.error('Afbeelding uploaden mislukt');
                saving.value = false;
                return;
            }
        }
        
        // Collect translations from editors
        const translations = editors.value.map(editor => editor.getContent());
        
        // Determine state:
        // - If publish=true -> always 'published'
        // - If editing and publish=false -> keep current state (draft stays draft, published stays published)
        // - If creating and publish=false -> 'draft'
        let eventState = 'draft';
        if (publish) {
            eventState = 'published';
        } else if (modalMode.value === 'edit' && editingEvent.value?.state) {
            eventState = editingEvent.value.state;
        }
        
        // Build JSON payload
        const payload = {
            state: eventState,
            on_home: formData.value.on_home,
            title: formData.value.title,
            type: formData.value.type,
            geolink: formData.value.geolink || '',
            location: formData.value.location || '',
            ticket_purchase_url: formData.value.ticket_purchase_url || '',
            start_date: formData.value.start_date ? new Date(formData.value.start_date).toISOString() : '',
            end_date: formData.value.end_date ? new Date(formData.value.end_date).toISOString() : '',
            translations: translations,
            imageUrl: imageUrl
        };
        
        const endpoint = modalMode.value === 'create' 
            ? 'admin/events'
            : `admin/events/${editingEvent.value.id}`;
        
        const method = modalMode.value === 'create' ? 'POST' : 'PUT';
        
        const response = await apiRequest(endpoint, {
            method,
            body: JSON.stringify(payload)
        });
        
        if (response?.ok) {
            window.$toast?.success(modalMode.value === 'create' ? 'Evenement aangemaakt!' : 'Evenement bijgewerkt!');
            selectedFile.value = null;
            closeModal();
            fetchEvents();
        } else {
            const err = await response?.json();
            window.$toast?.error(err?.message || err?.error || 'Opslaan mislukt');
        }
    } catch (err) {
        console.error('Save failed:', err);
        window.$toast?.error('Er is een fout opgetreden bij het opslaan');
    }
    
    saving.value = false;
}

async function handleDelete() {
    if (!editingEvent.value) return;
    
    if (!confirm(`Weet u zeker dat u "${editingEvent.value.title}" wilt verwijderen?`)) return;
    
    saving.value = true;
    try {
        const response = await apiRequest(`admin/events/${editingEvent.value.id}`, {
            method: 'DELETE'
        });
        
        if (response?.ok) {
            window.$toast?.success('Evenement verwijderd');
            closeModal();
            fetchEvents();
        } else {
            window.$toast?.error('Verwijderen mislukt');
        }
    } catch (err) {
        window.$toast?.error('Er is een fout opgetreden');
    }
    saving.value = false;
}

async function handleArchive() {
    if (!editingEvent.value) return;
    
    saving.value = true;
    try {
        const response = await apiRequest(`admin/events/${editingEvent.value.id}`, {
            method: 'PUT',
            body: JSON.stringify({
                ...editingEvent.value,
                state: 'archived'
            })
        });
        
        if (response?.ok) {
            window.$toast?.success('Evenement gearchiveerd');
            closeModal();
            fetchEvents();
        } else {
            window.$toast?.error('Archiveren mislukt');
        }
    } catch (err) {
        window.$toast?.error('Er is een fout opgetreden');
    }
    saving.value = false;
}

// Formatting
function formatDate(dateString) {
    if (!dateString) return '-';
    return new Date(dateString).toLocaleDateString('nl-BE', {
        month: 'short',
        day: 'numeric',
        year: 'numeric'
    });
}

function getStateBadge(state) {
    const badges = {
        published: 'success',
        ONLINE: 'success',
        draft: 'warning',
        DRAFT: 'warning',
        archived: 'neutral',
        ARCHIVED: 'neutral'
    };
    return badges[state] || 'neutral';
}

function getStateLabel(state) {
    const labels = {
        published: 'Gepubliceerd',
        ONLINE: 'Gepubliceerd',
        draft: 'Concept',
        DRAFT: 'Concept',
        archived: 'Gearchiveerd',
        ARCHIVED: 'Gearchiveerd'
    };
    return labels[state] || state;
}

function getTypeIcon(type) {
    return `/assets/media/eventtypes/${type}.png`;
}

// QR Code generation
async function generateQRCodes(event) {
    if (!event.uuid) {
        window.$toast?.error('Event heeft geen UUID');
        return;
    }
    
    const eventUrl = `${window.location.origin}/event/${event.uuid}`;
    const zip = new JSZip();
    const qrSize = 1024; // High resolution QR codes
    
    try {
        // Generate QR codes in different color combinations
        const variants = [
            { name: 'zwart-op-wit', dark: '#000000', light: '#FFFFFF' },
            { name: 'wit-op-zwart', dark: '#FFFFFF', light: '#000000' },
            { name: 'zwart-op-transparant', dark: '#000000', light: '#00000000' },
            { name: 'wit-op-transparant', dark: '#FFFFFF', light: '#00000000' }
        ];
        
        for (const variant of variants) {
            const canvas = document.createElement('canvas');
            await QRCode.toCanvas(canvas, eventUrl, {
                width: qrSize,
                margin: 2,
                color: {
                    dark: variant.dark,
                    light: variant.light
                }
            });
            
            // Convert canvas to blob
            const blob = await new Promise(resolve => canvas.toBlob(resolve, 'image/png'));
            zip.file(`qr-${variant.name}.png`, blob);
        }
        
        // Generate and download ZIP
        const zipBlob = await zip.generateAsync({ type: 'blob' });
        const link = document.createElement('a');
        link.href = URL.createObjectURL(zipBlob);
        
        // Clean filename from event title
        const cleanTitle = (event.title || 'event').replace(/[^a-zA-Z0-9]/g, '-').toLowerCase();
        link.download = `qr-codes-${cleanTitle}.zip`;
        link.click();
        URL.revokeObjectURL(link.href);
        
        window.$toast?.success('QR codes gedownload!');
    } catch (err) {
        console.error('QR generation failed:', err);
        window.$toast?.error('QR codes genereren mislukt');
    }
}

// Initialize
onMounted(async () => {
    await fetchLanguages();
    await fetchEvents();
});

watch(currPage, fetchEvents);
watch(statusFilter, () => {
    currPage.value = 1;
    fetchEvents();
});
</script>

<template>
    <AdminLayout pageTitle="Evenementen">
        <template #actions>
            <button class="admin-btn admin-btn-primary" @click="openCreateModal">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 1rem; height: 1rem;">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                Nieuw Evenement
            </button>
        </template>

        <!-- Filters -->
        <div class="admin-card" style="margin-bottom: 1.5rem;">
            <div class="admin-card-body">
                <div class="admin-filters">
                    <div class="admin-search" style="flex: 1; max-width: 20rem;">
                        <span class="admin-search-icon">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 1rem; height: 1rem;">
                                <circle cx="11" cy="11" r="8"/>
                                <line x1="21" y1="21" x2="16.65" y2="16.65"/>
                            </svg>
                        </span>
                        <input v-model="search" type="search" class="admin-input" placeholder="Zoek evenementen...">
                    </div>
                    <div class="admin-filter-group">
                        <label class="admin-filter-label">Status:</label>
                        <select v-model="statusFilter" class="admin-select" style="width: auto;">
                            <option value="">Alle</option>
                            <option value="published">Gepubliceerd</option>
                            <option value="draft">Concept</option>
                            <option value="archived">Gearchiveerd</option>
                        </select>
                    </div>
                </div>
            </div>
        </div>

        <!-- Table -->
        <div class="admin-card">
            <div class="admin-table-wrapper">
                <table class="admin-table">
                    <thead>
                        <tr>
                            <th style="width: 3rem;">Type</th>
                            <th>Titel</th>
                            <th style="width: 7rem;">Startdatum</th>
                            <th style="width: 7rem;">Einddatum</th>
                            <th style="width: 5rem;">Homepage</th>
                            <th style="width: 6rem;">Status</th>
                            <th style="width: 5rem;"></th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-if="loading">
                            <td colspan="7" style="text-align: center; padding: 3rem;">
                                <div class="admin-spinner" style="margin: 0 auto;"></div>
                            </td>
                        </tr>
                        <tr v-else-if="filteredEvents.length === 0">
                            <td colspan="7">
                                <div class="admin-empty">
                                    <div class="admin-empty-icon">
                                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                                            <rect x="3" y="4" width="18" height="18" rx="2"/>
                                            <line x1="16" y1="2" x2="16" y2="6"/>
                                            <line x1="8" y1="2" x2="8" y2="6"/>
                                            <line x1="3" y1="10" x2="21" y2="10"/>
                                        </svg>
                                    </div>
                                    <p class="admin-empty-title">Geen evenementen gevonden</p>
                                    <p class="admin-empty-description">Maak uw eerste evenement aan om te beginnen</p>
                                    <button class="admin-btn admin-btn-primary" @click="openCreateModal">
                                        Nieuw Evenement
                                    </button>
                                </div>
                            </td>
                        </tr>
                        <tr v-else v-for="event in filteredEvents" :key="event.id">
                            <td>
                                <img :src="getTypeIcon(event.type)" :alt="event.type" style="width: 1.5rem; height: 1.5rem;">
                            </td>
                            <td>
                                <span style="font-weight: 500;">{{ event.title }}</span>
                            </td>
                            <td>{{ formatDate(event.start_date) }}</td>
                            <td>{{ formatDate(event.end_date) }}</td>
                            <td>
                                <span :class="['admin-badge', event.on_home ? 'admin-badge-success' : 'admin-badge-neutral']">
                                    {{ event.on_home ? 'Ja' : 'Nee' }}
                                </span>
                            </td>
                            <td>
                                <span :class="['admin-badge', `admin-badge-${getStateBadge(event.state)}`]">
                                    {{ getStateLabel(event.state) }}
                                </span>
                            </td>
                            <td>
                                <div class="admin-table-actions">
                                    <button v-if="event.uuid" class="admin-btn admin-btn-ghost admin-btn-icon admin-btn-sm" @click="generateQRCodes(event)" title="Download QR Codes">
                                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 1rem; height: 1rem;">
                                            <rect x="3" y="3" width="7" height="7"/>
                                            <rect x="14" y="3" width="7" height="7"/>
                                            <rect x="3" y="14" width="7" height="7"/>
                                            <rect x="14" y="14" width="3" height="3"/>
                                            <rect x="18" y="14" width="3" height="3"/>
                                            <rect x="14" y="18" width="3" height="3"/>
                                            <rect x="18" y="18" width="3" height="3"/>
                                        </svg>
                                    </button>
                                    <a v-if="event.uuid" :href="`/event/${event.uuid}?preview=true`" target="_blank" class="admin-btn admin-btn-ghost admin-btn-icon admin-btn-sm" title="Voorbeeld">
                                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 1rem; height: 1rem;">
                                            <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                                            <circle cx="12" cy="12" r="3"/>
                                        </svg>
                                    </a>
                                    <button class="admin-btn admin-btn-ghost admin-btn-icon admin-btn-sm" @click="openEditModal(event)" title="Bewerken">
                                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 1rem; height: 1rem;">
                                            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                                            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                                        </svg>
                                    </button>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <!-- Pagination -->
            <div v-if="lastPage > 1" class="admin-card-footer">
                <div class="admin-pagination">
                    <button 
                        class="admin-pagination-btn" 
                        :disabled="currPage === 1"
                        @click="currPage--"
                    >
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 1rem; height: 1rem;">
                            <polyline points="15 18 9 12 15 6"/>
                        </svg>
                    </button>
                    <span style="font-size: 0.875rem; color: var(--admin-text-secondary); padding: 0 0.5rem;">
                        Pagina {{ currPage }} van {{ lastPage }}
                    </span>
                    <button 
                        class="admin-pagination-btn" 
                        :disabled="currPage === lastPage"
                        @click="currPage++"
                    >
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 1rem; height: 1rem;">
                            <polyline points="9 18 15 12 9 6"/>
                        </svg>
                    </button>
                </div>
            </div>
        </div>

        <!-- Modal -->
        <Teleport to="body">
            <div v-if="showModal" class="admin-modal-overlay" @click.self="closeModal">
                <div class="admin-modal admin-modal-xl">
                    <div class="admin-modal-header">
                        <h2 class="admin-modal-title">
                            {{ modalMode === 'create' ? 'Nieuw Evenement' : 'Evenement Bewerken' }}
                        </h2>
                        <button class="admin-modal-close" @click="closeModal">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <line x1="18" y1="6" x2="6" y2="18"/>
                                <line x1="6" y1="6" x2="18" y2="18"/>
                            </svg>
                        </button>
                    </div>
                    <div class="admin-modal-body">
                        <div class="modal-grid">
                            <!-- Left Column - Basic Info -->
                            <div class="modal-section">
                                <h3 class="section-title">Basisinformatie</h3>
                                
                                <div class="admin-form-group">
                                    <label class="admin-label">Titel *</label>
                                    <input v-model="formData.title" type="text" class="admin-input" required>
                                </div>

                                <div class="form-row">
                                    <div class="admin-form-group">
                                        <label class="admin-label">Type Evenement *</label>
                                        <select v-model="formData.type" class="admin-select">
                                            <option value="REGULAR">Regular</option>
                                            <option value="CITO">CITO</option>
                                            <option value="MEGA">Mega</option>
                                            <option value="GIGA">Giga</option>
                                            <option value="BLOCK">Block Party</option>
                                        </select>
                                    </div>
                                    <div class="admin-form-group">
                                        <label class="admin-label">Toon op Homepage</label>
                                        <select v-model="formData.on_home" class="admin-select">
                                            <option :value="false">Nee</option>
                                            <option :value="true">Ja</option>
                                        </select>
                                    </div>
                                </div>

                                <div class="form-row">
                                    <div class="admin-form-group">
                                        <label class="admin-label">Startdatum *</label>
                                        <input v-model="formData.start_date" type="datetime-local" class="admin-input" required>
                                    </div>
                                    <div class="admin-form-group">
                                        <label class="admin-label">Einddatum *</label>
                                        <input v-model="formData.end_date" type="datetime-local" class="admin-input" required>
                                    </div>
                                </div>

                                <div class="admin-form-group">
                                    <label class="admin-label">Locatie</label>
                                    <input v-model="formData.location" type="text" class="admin-input" placeholder="N 34° 56.789 E 123° 45.678">
                                    <span class="admin-form-hint">Coördinaten in geocaching formaat</span>
                                </div>

                                <div class="admin-form-group">
                                    <label class="admin-label">Geocaching Link</label>
                                    <input v-model="formData.geolink" type="url" class="admin-input" placeholder="https://www.geocaching.com/geocache/...">
                                </div>

                                <div class="admin-form-group">
                                    <label class="admin-label">Ticket URL</label>
                                    <input v-model="formData.ticket_purchase_url" type="url" class="admin-input" placeholder="https://...">
                                </div>
                            </div>

                            <!-- Right Column - Image -->
                            <div class="modal-section">
                                <h3 class="section-title">Afbeelding</h3>
                                
                                <div class="admin-form-group">
                                    <div class="image-upload-area">
                                        <img v-if="imagePreview" :src="imagePreview" class="image-preview">
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
                                        ref="fileInput" 
                                        type="file" 
                                        accept="image/*" 
                                        @change="handleImageChange"
                                        class="admin-input"
                                        style="margin-top: 0.75rem;"
                                    >
                                </div>
                            </div>
                        </div>

                        <!-- Translations -->
                        <div class="modal-section" style="margin-top: 1.5rem;">
                            <h3 class="section-title">Vertalingen</h3>
                            <div class="translations-grid">
                                <div v-for="(translation, index) in formData.translations" :key="translation.lang_code" class="translation-item">
                                    <label class="admin-label">{{ translation.lang_code }} Beschrijving</label>
                                    <TipTapEditor 
                                        :content="translation.description" 
                                        :langCode="translation.lang_code" 
                                        :editable="true" 
                                        ref="editors"
                                    />
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="admin-modal-footer">
                        <button v-if="modalMode === 'edit'" class="admin-btn admin-btn-danger" @click="handleDelete" :disabled="saving">
                            Verwijderen
                        </button>
                        <button v-if="modalMode === 'edit' && editingEvent?.state !== 'archived'" class="admin-btn admin-btn-ghost" @click="handleArchive" :disabled="saving">
                            Archiveren
                        </button>
                        <div style="flex: 1;"></div>
                        <button class="admin-btn admin-btn-secondary" @click="closeModal" :disabled="saving">
                            Annuleren
                        </button>
                        <button class="admin-btn admin-btn-secondary" @click="handleSave(false)" :disabled="saving">
                            {{ saving ? 'Opslaan...' : 'Opslaan' }}
                        </button>
                        <button class="admin-btn admin-btn-primary" @click="handleSave(true)" :disabled="saving">
                            {{ saving ? 'Opslaan...' : (modalMode === 'edit' && editingEvent?.state === 'published' ? 'Opslaan & Gepubliceerd' : 'Opslaan & Publiceren') }}
                        </button>
                    </div>
                </div>
            </div>
        </Teleport>
    </AdminLayout>
</template>

<style scoped>
.modal-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 2rem;
}

@media (max-width: 900px) {
    .modal-grid {
        grid-template-columns: 1fr;
    }
}

.modal-section {
    display: flex;
    flex-direction: column;
}

.section-title {
    font-size: 0.9375rem;
    font-weight: 600;
    color: var(--admin-text);
    margin: 0 0 1rem;
    padding-bottom: 0.5rem;
    border-bottom: 1px solid var(--admin-border-light);
}

.form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
}

.image-upload-area {
    border: 2px dashed var(--admin-border);
    border-radius: var(--admin-radius-lg);
    overflow: hidden;
    background: var(--admin-bg);
}

.image-preview {
    width: 100%;
    height: 16rem;
    object-fit: contain;
    background: var(--admin-surface);
}

.image-placeholder {
    height: 16rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 0.75rem;
    color: var(--admin-text-muted);
}

.image-placeholder svg {
    width: 3rem;
    height: 3rem;
}

.translations-grid {
    display: grid;
    gap: 1.5rem;
}

.translation-item {
    background: var(--admin-bg);
    border-radius: var(--admin-radius);
    padding: 1rem;
}

.translation-item :deep(.tiptap-editor) {
    min-height: 8rem;
}
</style>