<script setup>
import { ref, computed, onMounted } from 'vue';
import AdminLayout from '@/components/admin/AdminLayout.vue';
import config from '@/data/config.js';

// State
const loading = ref(true);
const staticContent = ref([]);
const languages = ref([]);
const search = ref('');
const selectedLanguage = ref('');

// Modal state
const showModal = ref(false);
const modalMode = ref('create');
const saving = ref(false);
const editingContent = ref(null);

// Form data
const formData = ref({
    property: '',
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
        'Content-Type': 'application/json',
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
async function fetchStaticContent() {
    loading.value = true;
    try {
        const response = await apiRequest('admin/static');
        if (response?.ok) {
            const data = await response.json();
            staticContent.value = data || [];
        }
    } catch (err) {
        console.error('Failed to fetch static content:', err);
    }
    loading.value = false;
}

async function fetchLanguages() {
    try {
        const response = await apiRequest('admin/languages');
        if (response?.ok) {
            const data = await response.json();
            languages.value = data.data || data || [];
            if (languages.value.length > 0 && !selectedLanguage.value) {
                selectedLanguage.value = languages.value[0].code;
            }
        }
    } catch (err) {
        console.error('Failed to fetch languages:', err);
    }
}

// Computed
const filteredContent = computed(() => {
    if (!search.value) return staticContent.value;
    const q = search.value.toLowerCase();
    return staticContent.value.filter(c => 
        c.property?.toLowerCase().includes(q)
    );
});

// Get content preview for table
function getContentPreview(content) {
    const translation = content.contents?.find(t => t.lang_code === selectedLanguage.value);
    if (!translation?.content) return 'Geen inhoud';
    return translation.content.slice(0, 80) + (translation.content.length > 80 ? '...' : '');
}

// Modal handlers
function openCreateModal() {
    modalMode.value = 'create';
    editingContent.value = null;
    
    formData.value = {
        property: '',
        translations: languages.value.map(l => ({
            lang_code: l.code,
            content: ''
        }))
    };
    
    showModal.value = true;
}

function openEditModal(content) {
    modalMode.value = 'edit';
    editingContent.value = content;
    
    const existingTranslations = new Map(
        content.contents?.map(t => [t.lang_code, t.content]) || []
    );
    
    formData.value = {
        property: content.property || '',
        translations: languages.value.map(l => ({
            lang_code: l.code,
            content: existingTranslations.get(l.code) || ''
        }))
    };
    
    showModal.value = true;
}

function closeModal() {
    showModal.value = false;
    editingContent.value = null;
}

async function handleSave() {
    if (!formData.value.property) {
        window.$toast?.error('Eigenschap sleutel is verplicht');
        return;
    }
    
    saving.value = true;
    
    try {
        if (modalMode.value === 'create') {
            for (const t of formData.value.translations) {
                if (!t.content.trim()) continue;
                
                const payload = {
                    property: formData.value.property,
                    lang_code: t.lang_code,
                    content: t.content
                };
                
                const response = await apiRequest('admin/static', {
                    method: 'POST',
                    body: JSON.stringify(payload)
                });
                
                if (!response?.ok) {
                    const err = await response?.json();
                    window.$toast?.error(err?.error || 'Vertaling opslaan mislukt');
                }
            }
            window.$toast?.success('Inhoud aangemaakt!');
        } else {
            for (const t of formData.value.translations) {
                const payload = {
                    lang_code: t.lang_code,
                    content: t.content
                };
                
                await apiRequest(`admin/static/${encodeURIComponent(formData.value.property)}`, {
                    method: 'PUT',
                    body: JSON.stringify(payload)
                });
            }
            window.$toast?.success('Inhoud bijgewerkt!');
        }
        
        closeModal();
        fetchStaticContent();
    } catch (err) {
        console.error('Save failed:', err);
        window.$toast?.error('Er is een fout opgetreden bij het opslaan');
    }
    
    saving.value = false;
}

async function handleDelete() {
    if (!editingContent.value) return;
    
    if (!confirm(`Weet je zeker dat je "${editingContent.value.property}" wilt verwijderen?`)) return;
    
    saving.value = true;
    try {
        const response = await apiRequest(`admin/static/${encodeURIComponent(editingContent.value.property)}`, {
            method: 'DELETE'
        });
        
        if (response?.ok) {
            window.$toast?.success('Inhoud verwijderd');
            closeModal();
            fetchStaticContent();
        } else {
            window.$toast?.error('Verwijderen mislukt');
        }
    } catch (err) {
        window.$toast?.error('Er is een fout opgetreden');
    }
    saving.value = false;
}

function isMultiLine(property) {
    return property && (
        property.includes('Txt') || 
        property.includes('Body') || 
        property.includes('SubTxt') ||
        property.includes('Description')
    );
}

onMounted(async () => {
    await fetchLanguages();
    await fetchStaticContent();
});
</script>

<template>
    <AdminLayout pageTitle="Vertalingen">
        <template #actions>
            <button class="admin-btn admin-btn-primary" @click="openCreateModal">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 1rem; height: 1rem;">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                Nieuwe Vertaling
            </button>
        </template>

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
                        <input v-model="search" type="search" class="admin-input" placeholder="Zoeken op eigenschap sleutel...">
                    </div>
                    <div class="admin-filter-group">
                        <label class="admin-filter-label">Voorbeeld Taal:</label>
                        <select v-model="selectedLanguage" class="admin-select" style="width: auto;">
                            <option v-for="lang in languages" :key="lang.code" :value="lang.code">
                                {{ lang.code }} - {{ lang.name }}
                            </option>
                        </select>
                    </div>
                </div>
            </div>
        </div>

        <div class="admin-card">
            <div class="admin-table-wrapper">
                <table class="admin-table">
                    <thead>
                        <tr>
                            <th style="width: 15rem;">Eigenschap Sleutel</th>
                            <th>Inhoud Voorbeeld</th>
                            <th style="width: 5rem;"></th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-if="loading">
                            <td colspan="3" style="text-align: center; padding: 3rem;">
                                <div class="admin-spinner" style="margin: 0 auto;"></div>
                            </td>
                        </tr>
                        <tr v-else-if="filteredContent.length === 0">
                            <td colspan="3">
                                <div class="admin-empty">
                                    <div class="admin-empty-icon">
                                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                                            <circle cx="12" cy="12" r="10"/>
                                            <line x1="2" y1="12" x2="22" y2="12"/>
                                            <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>
                                        </svg>
                                    </div>
                                    <p class="admin-empty-title">Geen vertalingen gevonden</p>
                                    <p class="admin-empty-description">Maak je eerste vertaling om te beginnen</p>
                                    <button class="admin-btn admin-btn-primary" @click="openCreateModal">
                                        Nieuwe Vertaling
                                    </button>
                                </div>
                            </td>
                        </tr>
                        <tr v-else v-for="content in filteredContent" :key="content.property">
                            <td>
                                <code style="font-size: 0.875rem;">{{ content.property }}</code>
                            </td>
                            <td>
                                <span class="content-preview">{{ getContentPreview(content) }}</span>
                            </td>
                            <td>
                                <div class="admin-table-actions">
                                    <button class="admin-btn admin-btn-ghost admin-btn-icon admin-btn-sm" @click="openEditModal(content)" title="Edit">
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

        <Teleport to="body">
            <div v-if="showModal" class="admin-modal-overlay" @click.self="closeModal">
                <div class="admin-modal admin-modal-lg">
                    <div class="admin-modal-header">
                        <h2 class="admin-modal-title">
                            {{ modalMode === 'create' ? 'Nieuwe Vertaling' : 'Vertaling Bewerken' }}
                        </h2>
                        <button class="admin-modal-close" @click="closeModal">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <line x1="18" y1="6" x2="6" y2="18"/>
                                <line x1="6" y1="6" x2="18" y2="18"/>
                            </svg>
                        </button>
                    </div>
                    <div class="admin-modal-body">
                        <div class="admin-form-group">
                            <label class="admin-label">Eigenschap Sleutel *</label>
                            <input 
                                v-model="formData.property" 
                                type="text" 
                                class="admin-input" 
                                placeholder="NavHome, AboutText, etc."
                                :disabled="modalMode === 'edit'"
                                required
                            >
                            <span class="admin-form-hint">Gebruik PascalCase voor de sleutelnaam (bijv. NavHome, ButtonSubmit)</span>
                        </div>

                        <div class="translations-list">
                            <div v-for="(translation, index) in formData.translations" :key="translation.lang_code" class="translation-card">
                                <label class="translation-lang">{{ translation.lang_code }}</label>
                                <textarea 
                                    v-if="isMultiLine(formData.property)"
                                    v-model="formData.translations[index].content"
                                    class="admin-input translation-textarea"
                                    rows="3"
                                    :placeholder="`Vertaling in ${translation.lang_code}...`"
                                ></textarea>
                                <input 
                                    v-else
                                    v-model="formData.translations[index].content"
                                    type="text"
                                    class="admin-input"
                                    :placeholder="`Vertaling in ${translation.lang_code}...`"
                                >
                            </div>
                        </div>
                    </div>
                    <div class="admin-modal-footer">
                        <button v-if="modalMode === 'edit'" class="admin-btn admin-btn-danger" @click="handleDelete" :disabled="saving">
                            Verwijderen
                        </button>
                        <div style="flex: 1;"></div>
                        <button class="admin-btn admin-btn-secondary" @click="closeModal" :disabled="saving">
                            Annuleren
                        </button>
                        <button class="admin-btn admin-btn-primary" @click="handleSave" :disabled="saving">
                            {{ saving ? 'Opslaan...' : 'Opslaan' }}
                        </button>
                    </div>
                </div>
            </div>
        </Teleport>
    </AdminLayout>
</template>

<style scoped>
.content-preview {
    color: var(--admin-text-secondary);
    font-size: 0.875rem;
}

.translations-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    margin-top: 1rem;
}

.translation-card {
    background: var(--admin-bg);
    border-radius: var(--admin-radius);
    padding: 1rem;
    border: 1px solid var(--admin-border-light);
}

.translation-lang {
    display: block;
    font-size: 0.75rem;
    font-weight: 600;
    color: var(--admin-primary);
    margin-bottom: 0.5rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
}

.translation-textarea {
    min-height: 80px;
    resize: vertical;
}
</style>
