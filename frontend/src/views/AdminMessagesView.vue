<script setup>
import { ref, computed, onMounted } from 'vue';
import AdminLayout from '@/components/admin/AdminLayout.vue';
import config from '@/data/config.js';

// State
const loading = ref(true);
const messages = ref([]);
const languages = ref([]);
const search = ref('');
const statusFilter = ref('');

// Modal state
const showModal = ref(false);
const modalMode = ref('create');
const saving = ref(false);
const editingMessage = ref(null);

// Form data
const formData = ref({
    state: 'draft',
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
        ...options.headers
    };

    try {
        const response = await fetch(`${config.apiUrl}${endpoint}`, {
            ...options,
            headers
        });
        return response;
    } catch (err) {
        console.error('API request failed:', err);
        return null;
    }
}

// Fetch data
async function fetchMessages() {
    loading.value = true;
    try {
        const response = await apiRequest('admin/messages');
        if (response?.ok) {
            const data = await response.json();
            messages.value = data.data || data || [];
        }
    } catch (err) {
        console.error('Failed to fetch messages:', err);
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
const filteredMessages = computed(() => {
    let filtered = messages.value;
    
    if (statusFilter.value) {
        filtered = filtered.filter(m => 
            m.state?.toLowerCase() === statusFilter.value.toLowerCase()
        );
    }
    
    if (search.value) {
        const q = search.value.toLowerCase();
        filtered = filtered.filter(m => {
            const nlTranslation = m.translations?.find(t => t.lang_code === 'NL');
            return nlTranslation?.title?.toLowerCase().includes(q);
        });
    }
    
    return filtered;
});

// Modal handlers
function openCreateModal() {
    modalMode.value = 'create';
    editingMessage.value = null;
    
    formData.value = {
        state: 'draft',
        priority: 0,
        translations: languages.value.map(l => ({
            lang_code: l.code,
            title: '',
            content: ''
        }))
    };
    
    showModal.value = true;
}

function openEditModal(message) {
    modalMode.value = 'edit';
    editingMessage.value = message;
    
    formData.value = {
        state: message.state || 'draft',
        priority: message.priority || 0,
        translations: message.translations?.length 
            ? message.translations.map(t => ({ ...t }))
            : languages.value.map(l => ({ lang_code: l.code, title: '', content: '' }))
    };
    
    showModal.value = true;
}

function closeModal() {
    showModal.value = false;
    editingMessage.value = null;
}

async function handleSave(publish = false) {
    saving.value = true;
    
    try {
        // Determine state:
        // - If publish=true -> always 'published'
        // - If editing and publish=false -> keep current state
        // - If creating and publish=false -> 'draft'
        let messageState = 'draft';
        if (publish) {
            messageState = 'published';
        } else if (modalMode.value === 'edit' && editingMessage.value?.state) {
            messageState = editingMessage.value.state;
        }
        
        // Build JSON payload (backend expects JSON, not FormData)
        const payload = {
            state: messageState,
            priority: formData.value.priority || 0,
            translations: formData.value.translations
        };
        
        const endpoint = modalMode.value === 'create' 
            ? 'admin/messages'
            : `admin/messages/${editingMessage.value.id}`;
        
        const method = modalMode.value === 'create' ? 'POST' : 'PUT';
        
        const response = await apiRequest(endpoint, {
            method,
            body: JSON.stringify(payload)
        });
        
        if (response?.ok) {
            window.$toast?.success(modalMode.value === 'create' ? 'Bericht aangemaakt!' : 'Bericht bijgewerkt!');
            closeModal();
            fetchMessages();
        } else {
            const err = await response?.json();
            window.$toast?.error(err?.message || 'Opslaan mislukt');
        }
    } catch (err) {
        console.error('Save failed:', err);
        window.$toast?.error('Er is een fout opgetreden bij het opslaan');
    }
    
    saving.value = false;
}

async function handleDelete() {
    if (!editingMessage.value) return;
    
    const title = editingMessage.value.translations?.find(t => t.lang_code === 'NL')?.title || 'dit bericht';
    if (!confirm(`Weet je zeker dat je "${title}" wilt verwijderen?`)) return;
    
    saving.value = true;
    try {
        const response = await apiRequest(`admin/messages/${editingMessage.value.id}`, {
            method: 'DELETE'
        });
        
        if (response?.ok) {
            window.$toast?.success('Bericht verwijderd');
            closeModal();
            fetchMessages();
        } else {
            window.$toast?.error('Verwijderen mislukt');
        }
    } catch (err) {
        window.$toast?.error('Er is een fout opgetreden');
    }
    saving.value = false;
}

async function handleArchive() {
    if (!editingMessage.value) return;
    
    saving.value = true;
    try {
        const payload = {
            state: 'archived',
            priority: editingMessage.value.priority || 0,
            translations: editingMessage.value.translations
        };
        
        const response = await apiRequest(`admin/messages/${editingMessage.value.id}`, {
            method: 'PUT',
            body: JSON.stringify(payload)
        });
        
        if (response?.ok) {
            window.$toast?.success('Bericht gearchiveerd');
            closeModal();
            fetchMessages();
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
    return new Date(dateString).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
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

function getMessageTitle(message) {
    return message.translations?.find(t => t.lang_code === 'NL')?.title || 
           message.translations?.[0]?.title || 
           'Naamloos';
}

// Initialize
onMounted(async () => {
    await fetchLanguages();
    await fetchMessages();
});
</script>

<template>
    <AdminLayout pageTitle="Berichten">
        <template #actions>
            <button class="admin-btn admin-btn-primary" @click="openCreateModal">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 1rem; height: 1rem;">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                Nieuw Bericht
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
                        <input v-model="search" type="search" class="admin-input" placeholder="Zoeken...">
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
                            <th>Titel</th>
                            <th style="width: 10rem;">Laatst Bijgewerkt</th>
                            <th style="width: 6rem;">Status</th>
                            <th style="width: 5rem;"></th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-if="loading">
                            <td colspan="4" style="text-align: center; padding: 3rem;">
                                <div class="admin-spinner" style="margin: 0 auto;"></div>
                            </td>
                        </tr>
                        <tr v-else-if="filteredMessages.length === 0">
                            <td colspan="4">
                                <div class="admin-empty">
                                    <div class="admin-empty-icon">
                                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                                            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
                                        </svg>
                                    </div>
                                    <p class="admin-empty-title">Geen berichten gevonden</p>
                                    <p class="admin-empty-description">Maak je eerste bericht om aankondigingen te tonen</p>
                                    <button class="admin-btn admin-btn-primary" @click="openCreateModal">
                                        Nieuw Bericht
                                    </button>
                                </div>
                            </td>
                        </tr>
                        <tr v-else v-for="message in filteredMessages" :key="message.id">
                            <td>
                                <span style="font-weight: 500;">{{ getMessageTitle(message) }}</span>
                            </td>
                            <td>{{ formatDate(message.updated_at) }}</td>
                            <td>
                                <span :class="['admin-badge', `admin-badge-${getStateBadge(message.state)}`]">
                                    {{ message.state }}
                                </span>
                            </td>
                            <td>
                                <div class="admin-table-actions">
                                    <button class="admin-btn admin-btn-ghost admin-btn-icon admin-btn-sm" @click="openEditModal(message)" title="Edit">
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
        </div>

        <!-- Modal -->
        <Teleport to="body">
            <div v-if="showModal" class="admin-modal-overlay" @click.self="closeModal">
                <div class="admin-modal admin-modal-lg">
                    <div class="admin-modal-header">
                        <h2 class="admin-modal-title">
                            {{ modalMode === 'create' ? 'Nieuw Bericht' : 'Bericht Bewerken' }}
                        </h2>
                        <button class="admin-modal-close" @click="closeModal">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <line x1="18" y1="6" x2="6" y2="18"/>
                                <line x1="6" y1="6" x2="18" y2="18"/>
                            </svg>
                        </button>
                    </div>
                    <div class="admin-modal-body">
                        <p class="form-hint" style="margin-bottom: 1.5rem; color: var(--admin-text-muted); font-size: 0.875rem;">
                            Berichten worden getoond op de homepage. Voor langere teksten, gebruik zowel titel als inhoud. Voor kortere aankondigingen volstaat alleen een titel.
                        </p>
                        
                        <div class="translations-list">
                            <div v-for="(translation, index) in formData.translations" :key="translation.lang_code" class="translation-card">
                                <h4 class="translation-lang">{{ translation.lang_code }}</h4>
                                <div class="admin-form-group">
                                    <label class="admin-label">Titel *</label>
                                    <input v-model="formData.translations[index].title" type="text" class="admin-input" required>
                                </div>
                                <div class="admin-form-group">
                                    <label class="admin-label">Inhoud</label>
                                    <textarea v-model="formData.translations[index].content" class="admin-textarea" rows="3"></textarea>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="admin-modal-footer">
                        <button v-if="modalMode === 'edit'" class="admin-btn admin-btn-danger" @click="handleDelete" :disabled="saving">
                            Verwijderen
                        </button>
                        <button v-if="modalMode === 'edit' && (formData.state === 'published' || formData.state === 'ONLINE')" class="admin-btn admin-btn-secondary" @click="handleArchive" :disabled="saving">
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
                            {{ saving ? 'Opslaan...' : (modalMode === 'edit' && editingMessage?.state === 'published' ? 'Opslaan & Gepubliceerd' : 'Opslaan & Publiceren') }}
                        </button>
                    </div>
                </div>
            </div>
        </Teleport>
    </AdminLayout>
</template>

<style scoped>
.translations-list {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
}

.translation-card {
    background: var(--admin-bg);
    border-radius: var(--admin-radius);
    padding: 1.25rem;
    border: 1px solid var(--admin-border-light);
}

.translation-lang {
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--admin-primary);
    margin: 0 0 1rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
}
</style>