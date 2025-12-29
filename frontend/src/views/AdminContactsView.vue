<template>
    <AdminLayout pageTitle="Contactberichten">
        <template #actions>
            <span class="contact-count">{{ totalContacts }} berichten totaal</span>
        </template>

        <!-- Filters -->
        <div class="admin-card" style="margin-bottom: 1.5rem;">
            <div class="admin-card-body">
                <div class="admin-filters">
                    <div class="admin-filter-group">
                        <label class="admin-filter-label">Status</label>
                        <select v-model="selectedStatus" @change="fetchContacts" class="admin-select">
                            <option value="">Alle Statussen</option>
                            <option value="new">Nieuw</option>
                            <option value="in_progress">In Behandeling</option>
                            <option value="resolved">Opgelost</option>
                            <option value="closed">Gesloten</option>
                        </select>
                    </div>
                </div>
            </div>
        </div>

        <!-- Loading State -->
        <div v-if="loading" class="admin-card">
            <div class="admin-card-body" style="text-align: center; padding: 3rem;">
                <div class="admin-spinner"></div>
                <p style="margin-top: 1rem; color: var(--admin-muted);">Berichten laden...</p>
            </div>
        </div>

        <!-- Empty State -->
        <div v-else-if="contacts.length === 0" class="admin-empty">
            <div class="admin-empty-icon">📭</div>
            <h3 class="admin-empty-title">Geen Berichten Gevonden</h3>
            <p class="admin-empty-description">
                {{ selectedStatus ? 'Geen berichten gevonden met de geselecteerde status.' : 'Nog geen contactformulier inzendingen.' }}
            </p>
        </div>

        <!-- Contacts Layout -->
        <div v-else class="contacts-layout">
            <!-- Contact List -->
            <div class="contacts-list" :class="{ 'has-selection': selectedContact }">
                <div class="admin-card">
                    <div class="contact-items">
                        <div
                            v-for="contact in contacts"
                            :key="contact.id"
                            class="contact-item"
                            :class="{ 
                                selected: selectedContact?.id === contact.id,
                                unread: contact.status === 'new'
                            }"
                            @click="selectContact(contact)"
                        >
                            <div class="contact-status-indicator" :class="`status-${contact.status}`"></div>
                            <div class="contact-info">
                                <div class="contact-header">
                                    <span class="contact-email">{{ contact.email }}</span>
                                    <span class="contact-time">{{ getDaysAgo(contact.created_at) }}</span>
                                </div>
                                <div class="contact-subject">{{ contact.subject }}</div>
                                <div class="contact-meta">
                                    <span class="admin-badge" :class="getStatusBadgeClass(contact.status)">
                                        {{ getStatusLabel(contact.status) }}
                                    </span>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- Pagination -->
                    <div v-if="totalPages > 1" class="contacts-pagination">
                        <button 
                            class="admin-btn admin-btn-secondary"
                            @click="changePage(-1)" 
                            :disabled="currentPage === 1"
                        >
                            ← Vorige
                        </button>
                        <span class="pagination-info">Pagina {{ currentPage }} van {{ totalPages }}</span>
                        <button 
                            class="admin-btn admin-btn-secondary"
                            @click="changePage(1)" 
                            :disabled="currentPage === totalPages"
                        >
                            Volgende →
                        </button>
                    </div>
                </div>
            </div>

            <!-- Contact Detail Panel -->
            <div class="contact-detail-panel" :class="{ visible: selectedContact }">
                <div v-if="selectedContact" class="admin-card">
                    <!-- Loading Details -->
                    <div v-if="loadingDetails" class="admin-card-body" style="text-align: center; padding: 3rem;">
                        <div class="admin-spinner"></div>
                    </div>

                    <!-- Contact Details -->
                    <template v-else-if="contactDetails">
                        <!-- Header -->
                        <div class="detail-header">
                            <button class="close-detail-btn" @click="selectedContact = null">×</button>
                            <h3 class="detail-subject">{{ contactDetails.subject }}</h3>
                            <div class="detail-meta">
                                <span class="detail-meta-item">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                                    {{ contactDetails.name }}
                                </span>
                                <span class="detail-meta-item">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="20" height="16" x="2" y="4" rx="2"/><path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"/></svg>
                                    {{ contactDetails.email }}
                                </span>
                                <span class="detail-meta-item">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" x2="16" y1="2" y2="6"/><line x1="8" x2="8" y1="2" y2="6"/><line x1="3" x2="21" y1="10" y2="10"/></svg>
                                    {{ formatDate(contactDetails.created_at) }}
                                </span>
                            </div>
                        </div>

                        <!-- Status Section -->
                        <div class="detail-status-section">
                            <label class="status-section-label">Status</label>
                            <div class="status-buttons">
                                <button 
                                    v-for="status in statusOptions"
                                    :key="status.value"
                                    class="status-btn"
                                    :class="{ active: contactDetails.status === status.value }"
                                    :data-status="status.value"
                                    @click="updateStatus(status.value)"
                                >
                                    {{ status.label }}
                                </button>
                            </div>
                        </div>

                        <!-- Message -->
                        <div class="detail-message-section">
                            <label class="message-section-label">Bericht</label>
                            <div class="message-content">{{ contactDetails.message }}</div>
                        </div>

                        <!-- Actions -->
                        <div class="detail-actions">
                            <button class="admin-btn admin-btn-primary" @click="openEmailClient">
                                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m22 2-7 20-4-9-9-4Z"/><path d="M22 2 11 13"/></svg>
                                Beantwoorden via E-mail
                            </button>
                            <button class="admin-btn admin-btn-danger" @click="deleteContact">
                                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/></svg>
                                Verwijderen
                            </button>
                        </div>

                        <!-- Notes Section -->
                        <div class="detail-notes-section">
                            <label class="notes-section-label">
                                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"/><path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4Z"/></svg>
                                Interne Notities
                            </label>

                            <div v-if="contactDetails.notes?.length" class="notes-list">
                                <div v-for="note in contactDetails.notes" :key="note.id" class="note-item">
                                    <div class="note-header">
                                        <span class="note-author">{{ note.admin_email }}</span>
                                        <span class="note-date">{{ formatDate(note.created_at) }}</span>
                                    </div>
                                    <p class="note-content">{{ note.content }}</p>
                                </div>
                            </div>
                            <p v-else class="no-notes">Nog geen notities</p>

                            <div class="add-note-form">
                                <textarea 
                                    v-model="newNote" 
                                    class="admin-textarea"
                                    placeholder="Voeg een interne notitie toe..."
                                    rows="2"
                                ></textarea>
                                <button 
                                    class="admin-btn admin-btn-secondary"
                                    @click="addNote"
                                    :disabled="!newNote.trim() || addingNote"
                                >
                                    {{ addingNote ? 'Toevoegen...' : 'Notitie Toevoegen' }}
                                </button>
                            </div>
                        </div>
                    </template>
                </div>

                <!-- No Selection State -->
                <div v-else class="no-selection-state">
                    <div class="no-selection-icon">
                        <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect width="20" height="16" x="2" y="4" rx="2"/><path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"/></svg>
                    </div>
                    <p>Selecteer een bericht om details te bekijken</p>
                </div>
            </div>
        </div>
    </AdminLayout>
</template>

<script setup>
import { ref, onMounted, inject } from 'vue'
import { useRouter } from 'vue-router'
import AdminLayout from '@/components/admin/AdminLayout.vue'

const router = useRouter()
const refreshContactCount = inject('refreshContactCount', () => {})

// Auth & API helpers
const verifyLogin = () => {
    const token = localStorage.getItem('admin_token')
    if (!token) {
        router.push({ name: 'admin' })
        return false
    }
    return true
}

const apiRequest = async (endpoint, options = {}) => {
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`/api/admin${endpoint}`, {
        ...options,
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
            ...options.headers
        }
    })
    
    if (response.status === 401) {
        localStorage.removeItem('admin_token')
        router.push({ name: 'admin' })
        throw new Error('Session expired')
    }
    
    return response
}

// State
const contacts = ref([])
const loading = ref(true)
const selectedStatus = ref('')
const currentPage = ref(1)
const totalPages = ref(1)
const totalContacts = ref(0)
const pageSize = 20

const selectedContact = ref(null)
const contactDetails = ref(null)
const loadingDetails = ref(false)
const newNote = ref('')
const addingNote = ref(false)

const statusOptions = [
    { value: 'new', label: 'Nieuw' },
    { value: 'in_progress', label: 'In Behandeling' },
    { value: 'resolved', label: 'Opgelost' },
    { value: 'closed', label: 'Gesloten' }
]

// Status helpers
const getStatusLabel = (status) => {
    const labels = {
        new: 'Nieuw',
        in_progress: 'In Behandeling',
        resolved: 'Opgelost',
        closed: 'Gesloten'
    }
    return labels[status] || status
}

const getStatusBadgeClass = (status) => {
    const classes = {
        new: 'admin-badge-danger',
        in_progress: 'admin-badge-warning',
        resolved: 'admin-badge-success',
        closed: 'admin-badge-neutral'
    }
    return classes[status] || 'admin-badge-neutral'
}

// Fetch contacts
const fetchContacts = async () => {
    if (!verifyLogin()) return
    loading.value = true
    
    try {
        let url = `/contacts?page=${currentPage.value}&limit=${pageSize}`
        if (selectedStatus.value) {
            url += `&status=${selectedStatus.value}`
        }
        
        const response = await apiRequest(url)
        const data = await response.json()
        
        contacts.value = data.data || []
        totalContacts.value = data.total || 0
        totalPages.value = Math.ceil(totalContacts.value / pageSize)
    } catch (error) {
        console.error('Error fetching contacts:', error)
        window.$toast?.error('Berichten laden mislukt')
    } finally {
        loading.value = false
    }
}

// Select contact
const selectContact = async (contact) => {
    selectedContact.value = contact
    await fetchContactDetails(contact.id)
}

// Fetch contact details
const fetchContactDetails = async (id) => {
    loadingDetails.value = true
    
    try {
        const response = await apiRequest(`/contacts/${id}`)
        contactDetails.value = await response.json()
        
        // Update list item status
        const index = contacts.value.findIndex(c => c.id === id)
        if (index !== -1) {
            contacts.value[index].status = contactDetails.value.status
        }
    } catch (error) {
        console.error('Error fetching details:', error)
        window.$toast?.error('Berichtdetails laden mislukt')
    } finally {
        loadingDetails.value = false
    }
}

// Update status
const updateStatus = async (status) => {
    if (!contactDetails.value || contactDetails.value.status === status) return
    
    try {
        await apiRequest(`/contacts/${contactDetails.value.id}/status`, {
            method: 'PUT',
            body: JSON.stringify({ status })
        })
        
        contactDetails.value.status = status
        
        const index = contacts.value.findIndex(c => c.id === contactDetails.value.id)
        if (index !== -1) {
            contacts.value[index].status = status
        }
        
        refreshContactCount()
        window.$toast?.success('Status bijgewerkt')
    } catch (error) {
        console.error('Error updating status:', error)
        window.$toast?.error('Status bijwerken mislukt')
    }
}

// Add note
const addNote = async () => {
    if (!newNote.value.trim() || !contactDetails.value) return
    
    addingNote.value = true
    
    try {
        const response = await apiRequest(`/contacts/${contactDetails.value.id}/notes`, {
            method: 'POST',
            body: JSON.stringify({ content: newNote.value.trim() })
        })
        
        const note = await response.json()
        
        if (!contactDetails.value.notes) {
            contactDetails.value.notes = []
        }
        contactDetails.value.notes.push(note)
        newNote.value = ''
        
        window.$toast?.success('Notitie toegevoegd')
    } catch (error) {
        console.error('Error adding note:', error)
        window.$toast?.error('Notitie toevoegen mislukt')
    } finally {
        addingNote.value = false
    }
}

// Delete contact
const deleteContact = async () => {
    if (!contactDetails.value) return
    
    if (!confirm('Weet je zeker dat je dit bericht wilt verwijderen? Deze actie kan niet ongedaan worden gemaakt.')) return
    
    try {
        await apiRequest(`/contacts/${contactDetails.value.id}`, {
            method: 'DELETE'
        })
        
        contacts.value = contacts.value.filter(c => c.id !== contactDetails.value.id)
        selectedContact.value = null
        contactDetails.value = null
        totalContacts.value--
        
        refreshContactCount()
        window.$toast?.success('Bericht verwijderd')
    } catch (error) {
        console.error('Error deleting contact:', error)
        window.$toast?.error('Bericht verwijderen mislukt')
    }
}

// Open email client
const openEmailClient = () => {
    if (!contactDetails.value) return
    
    const subject = encodeURIComponent(`Re: ${contactDetails.value.subject}`)
    const body = encodeURIComponent(`\n\n--- Original Message ---\nFrom: ${contactDetails.value.name} <${contactDetails.value.email}>\nDate: ${formatDate(contactDetails.value.created_at)}\n\n${contactDetails.value.message}`)
    
    window.location.href = `mailto:${contactDetails.value.email}?subject=${subject}&body=${body}`
}

// Pagination
const changePage = (delta) => {
    const newPage = currentPage.value + delta
    if (newPage >= 1 && newPage <= totalPages.value) {
        currentPage.value = newPage
        fetchContacts()
    }
}

// Helpers
const formatDate = (dateStr) => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        day: 'numeric',
        month: 'short',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    })
}

const getDaysAgo = (dateStr) => {
    const date = new Date(dateStr)
    const now = new Date()
    const diffTime = now - date
    const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24))
    
    if (diffDays === 0) return 'Vandaag'
    if (diffDays === 1) return 'Gisteren'
    if (diffDays < 7) return `${diffDays}d geleden`
    return formatDate(dateStr).split(',')[0]
}

onMounted(fetchContacts)
</script>

<style scoped>
.contact-count {
    color: var(--admin-muted);
    font-size: 0.875rem;
}

/* Contacts Layout */
.contacts-layout {
    display: grid;
    grid-template-columns: 400px 1fr;
    gap: 1.5rem;
    min-height: calc(100vh - 250px);
}

.contacts-list.has-selection {
    grid-column: 1;
}

/* Contact Items */
.contact-items {
    max-height: calc(100vh - 350px);
    overflow-y: auto;
}

.contact-item {
    display: flex;
    gap: 0.75rem;
    padding: 1rem 1.25rem;
    border-bottom: 1px solid var(--admin-border);
    cursor: pointer;
    transition: all 0.2s;
}

.contact-item:hover {
    background: var(--admin-hover);
}

.contact-item.selected {
    background: rgba(13, 148, 136, 0.08);
    border-left: 3px solid var(--admin-primary);
}

.contact-item.unread .contact-subject {
    font-weight: 600;
    color: var(--admin-text);
}

.contact-status-indicator {
    width: 10px;
    height: 10px;
    border-radius: 50%;
    flex-shrink: 0;
    margin-top: 6px;
}

.contact-status-indicator.status-new { background: var(--admin-danger); }
.contact-status-indicator.status-in_progress { background: var(--admin-warning); }
.contact-status-indicator.status-resolved { background: var(--admin-success); }
.contact-status-indicator.status-closed { background: var(--admin-text-muted); }

.contact-info {
    flex: 1;
    min-width: 0;
}

.contact-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.25rem;
}

.contact-email {
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--admin-text);
}

.contact-time {
    font-size: 0.75rem;
    color: var(--admin-muted);
}

.contact-subject {
    font-size: 0.8125rem;
    color: var(--admin-muted);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    margin-bottom: 0.5rem;
}

.contact-meta {
    display: flex;
    gap: 0.5rem;
}

/* Pagination */
.contacts-pagination {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 1rem;
    padding: 1rem;
    border-top: 1px solid var(--admin-border);
}

.pagination-info {
    font-size: 0.875rem;
    color: var(--admin-muted);
}

/* Detail Panel */
.contact-detail-panel {
    position: relative;
}

.contact-detail-panel .admin-card {
    position: sticky;
    top: 1rem;
}

.detail-header {
    padding: 1.25rem;
    border-bottom: 1px solid var(--admin-border);
    position: relative;
}

.close-detail-btn {
    position: absolute;
    top: 1rem;
    right: 1rem;
    width: 32px;
    height: 32px;
    border: none;
    background: var(--admin-hover);
    border-radius: 50%;
    font-size: 1.25rem;
    cursor: pointer;
    color: var(--admin-muted);
    display: none;
    transition: all 0.2s;
}

.close-detail-btn:hover {
    background: var(--admin-border);
    color: var(--admin-text);
}

.detail-subject {
    margin: 0 0 0.75rem 0;
    font-size: 1.125rem;
    font-weight: 600;
    color: var(--admin-text);
    padding-right: 2rem;
}

.detail-meta {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
}

.detail-meta-item {
    display: flex;
    align-items: center;
    gap: 0.375rem;
    font-size: 0.8125rem;
    color: var(--admin-muted);
}

.detail-meta-item svg {
    color: var(--admin-primary);
}

/* Status Section */
.detail-status-section {
    padding: 1rem 1.25rem;
    background: var(--admin-hover);
    border-bottom: 1px solid var(--admin-border);
}

.status-section-label {
    display: block;
    font-size: 0.75rem;
    font-weight: 600;
    color: var(--admin-muted);
    text-transform: uppercase;
    letter-spacing: 0.05em;
    margin-bottom: 0.75rem;
}

.status-buttons {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
}

.status-btn {
    padding: 0.375rem 0.875rem;
    border: 2px solid var(--admin-border);
    border-radius: 9999px;
    background: var(--admin-card);
    font-size: 0.8125rem;
    cursor: pointer;
    transition: all 0.2s;
    color: var(--admin-muted);
}

.status-btn:hover {
    border-color: var(--admin-primary);
}

.status-btn.active[data-status="new"] {
    border-color: var(--admin-danger);
    background: var(--admin-danger-bg);
    color: var(--admin-danger);
}

.status-btn.active[data-status="in_progress"] {
    border-color: var(--admin-warning);
    background: var(--admin-warning-bg);
    color: var(--admin-warning);
}

.status-btn.active[data-status="resolved"] {
    border-color: var(--admin-success);
    background: var(--admin-success-bg);
    color: var(--admin-success);
}

.status-btn.active[data-status="closed"] {
    border-color: var(--admin-text-muted);
    background: var(--admin-surface-hover);
    color: var(--admin-text-secondary);
}

/* Message Section */
.detail-message-section {
    padding: 1.25rem;
    border-bottom: 1px solid var(--admin-border);
}

.message-section-label {
    display: block;
    font-size: 0.75rem;
    font-weight: 600;
    color: var(--admin-muted);
    text-transform: uppercase;
    letter-spacing: 0.05em;
    margin-bottom: 0.75rem;
}

.message-content {
    background: var(--admin-hover);
    padding: 1rem;
    border-radius: 0.5rem;
    white-space: pre-wrap;
    font-size: 0.875rem;
    line-height: 1.7;
    color: var(--admin-text);
}

/* Actions */
.detail-actions {
    padding: 1.25rem;
    display: flex;
    gap: 0.75rem;
    border-bottom: 1px solid var(--admin-border);
}

.detail-actions .admin-btn {
    display: flex;
    align-items: center;
    gap: 0.5rem;
}

/* Notes Section */
.detail-notes-section {
    padding: 1.25rem;
    background: var(--admin-bg);
}

.notes-section-label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--admin-text);
    margin-bottom: 1rem;
}

.notes-list {
    margin-bottom: 1rem;
}

.note-item {
    background: var(--admin-card);
    padding: 0.875rem;
    border-radius: 0.5rem;
    margin-bottom: 0.5rem;
    border: 1px solid var(--admin-border);
}

.note-header {
    display: flex;
    justify-content: space-between;
    margin-bottom: 0.375rem;
}

.note-author {
    font-size: 0.8125rem;
    font-weight: 500;
    color: var(--admin-primary);
}

.note-date {
    font-size: 0.75rem;
    color: var(--admin-muted);
}

.note-content {
    margin: 0;
    font-size: 0.8125rem;
    color: var(--admin-text);
    line-height: 1.5;
}

.no-notes {
    color: var(--admin-muted);
    font-size: 0.875rem;
    text-align: center;
    padding: 1rem;
    background: var(--admin-card);
    border-radius: 0.5rem;
    border: 1px dashed var(--admin-border);
    margin-bottom: 1rem;
}

.add-note-form {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
}

.add-note-form .admin-textarea {
    resize: vertical;
    min-height: 60px;
}

.add-note-form .admin-btn {
    align-self: flex-start;
}

/* No Selection State */
.no-selection-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    background: var(--admin-card);
    border-radius: var(--admin-radius);
    border: 1px solid var(--admin-border);
    padding: 4rem 2rem;
    text-align: center;
    min-height: 400px;
}

.no-selection-icon {
    color: var(--admin-border);
    margin-bottom: 1rem;
}

.no-selection-state p {
    color: var(--admin-muted);
    font-size: 0.875rem;
    margin: 0;
}

/* Responsive */
@media (max-width: 1024px) {
    .contacts-layout {
        grid-template-columns: 1fr;
    }

    .contacts-list.has-selection {
        display: none;
    }

    .contact-detail-panel {
        position: fixed;
        inset: 0;
        z-index: 100;
        background: var(--admin-bg);
        display: none;
        padding: 1rem;
        overflow-y: auto;
    }

    .contact-detail-panel.visible {
        display: block;
    }

    .close-detail-btn {
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .no-selection-state {
        display: none;
    }
}

/* Spinner */
.admin-spinner {
    width: 32px;
    height: 32px;
    border: 3px solid var(--admin-border);
    border-top-color: var(--admin-primary);
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin: 0 auto;
}

@keyframes spin {
    to { transform: rotate(360deg); }
}
</style>
