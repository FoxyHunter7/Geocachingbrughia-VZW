<template>
    <AdminLayout pageTitle="Gebruikers">
        <template #actions>
            <button class="admin-btn admin-btn-primary" @click="showAddModal = true">
                <span class="btn-icon">+</span>
                Nieuwe Gebruiker
            </button>
        </template>

        <!-- Loading State -->
        <div v-if="loading" class="admin-card">
            <div class="admin-card-body" style="text-align: center; padding: 3rem;">
                <div class="admin-spinner"></div>
                <p style="margin-top: 1rem; color: var(--admin-muted);">Gebruikers laden...</p>
            </div>
        </div>

        <!-- Empty State -->
        <div v-else-if="users.length === 0" class="admin-empty">
            <div class="admin-empty-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                    <circle cx="9" cy="7" r="4"/>
                    <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
                    <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
                </svg>
            </div>
            <h3 class="admin-empty-title">Geen gebruikers gevonden</h3>
            <p class="admin-empty-description">
                Voeg een nieuwe beheerder toe om te beginnen.
            </p>
        </div>

        <!-- Users Table -->
        <div v-else class="admin-card">
            <div class="admin-table-container">
                <table class="admin-table">
                    <thead>
                        <tr>
                            <th>Naam</th>
                            <th>Email</th>
                            <th>Status</th>
                            <th>Aangemaakt</th>
                            <th>Acties</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="user in users" :key="user.id">
                            <td>
                                <div class="user-name">{{ user.name }}</div>
                            </td>
                            <td>
                                <div class="user-email">{{ user.email }}</div>
                            </td>
                            <td>
                                <span 
                                    class="admin-badge"
                                    :class="user.needs_password_update ? 'admin-badge-warning' : 'admin-badge-success'"
                                >
                                    {{ user.needs_password_update ? 'Wacht op activatie' : 'Actief' }}
                                </span>
                            </td>
                            <td>
                                <span class="date-cell">{{ formatDate(user.created_at) }}</span>
                            </td>
                            <td>
                                <div class="action-buttons">
                                    <button 
                                        v-if="user.needs_password_update"
                                        class="admin-btn admin-btn-secondary admin-btn-sm"
                                        @click="confirmResendInvitation(user)"
                                        title="Uitnodiging opnieuw versturen"
                                    >
                                        Opnieuw versturen
                                    </button>
                                    <button 
                                        class="admin-btn admin-btn-danger admin-btn-sm"
                                        @click="confirmDelete(user)"
                                        title="Gebruiker verwijderen"
                                    >
                                        Verwijderen
                                    </button>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <!-- Add User Modal -->
        <div v-if="showAddModal" class="admin-modal-overlay" @click.self="closeAddModal">
            <div class="admin-modal">
                <div class="admin-modal-header">
                    <h2>Nieuwe Gebruiker Toevoegen</h2>
                    <button class="admin-modal-close" @click="closeAddModal">&times;</button>
                </div>
                <form @submit.prevent="handleAddUser">
                    <div class="admin-modal-body">
                        <div class="admin-form-group">
                            <label for="name" class="admin-label">Naam *</label>
                            <input 
                                v-model="newUser.name" 
                                type="text" 
                                id="name" 
                                class="admin-input"
                                placeholder="Volledige naam"
                                required
                            >
                            <p v-if="errors.name" class="admin-error">{{ errors.name[0] }}</p>
                        </div>
                        <div class="admin-form-group">
                            <label for="email" class="admin-label">Email *</label>
                            <input 
                                v-model="newUser.email" 
                                type="email" 
                                id="email" 
                                class="admin-input"
                                placeholder="email@voorbeeld.be"
                                required
                            >
                            <p v-if="errors.email" class="admin-error">{{ errors.email[0] }}</p>
                        </div>
                        <p v-if="formError" class="admin-error">{{ formError }}</p>
                        <div class="info-box">
                            <strong>Let op:</strong> Er wordt automatisch een uitnodigingsmail naar dit emailadres verstuurd.
                        </div>
                    </div>
                    <div class="admin-modal-footer">
                        <button type="button" class="admin-btn admin-btn-secondary" @click="closeAddModal">
                            Annuleren
                        </button>
                        <button type="submit" class="admin-btn admin-btn-primary" :disabled="submitting">
                            {{ submitting ? 'Bezig...' : 'Toevoegen & Uitnodiging Versturen' }}
                        </button>
                    </div>
                </form>
            </div>
        </div>

        <!-- Delete Confirmation Modal -->
        <div v-if="showDeleteModal" class="admin-modal-overlay" @click.self="showDeleteModal = false">
            <div class="admin-modal admin-modal-sm">
                <div class="admin-modal-header">
                    <h2>Gebruiker Verwijderen</h2>
                    <button class="admin-modal-close" @click="showDeleteModal = false">&times;</button>
                </div>
                <div class="admin-modal-body">
                    <p>Weet u zeker dat u <strong>{{ userToDelete?.name }}</strong> ({{ userToDelete?.email }}) wilt verwijderen?</p>
                    <p class="admin-text-muted">Deze actie kan niet ongedaan worden gemaakt.</p>
                </div>
                <div class="admin-modal-footer">
                    <button class="admin-btn admin-btn-secondary" @click="showDeleteModal = false">
                        Annuleren
                    </button>
                    <button class="admin-btn admin-btn-danger" @click="handleDelete" :disabled="deleting">
                        {{ deleting ? 'Bezig...' : 'Verwijderen' }}
                    </button>
                </div>
            </div>
        </div>

        <!-- Resend Invitation Confirmation Modal -->
        <div v-if="showResendModal" class="admin-modal-overlay" @click.self="showResendModal = false">
            <div class="admin-modal admin-modal-sm">
                <div class="admin-modal-header">
                    <h2>Uitnodiging Opnieuw Versturen</h2>
                    <button class="admin-modal-close" @click="showResendModal = false">&times;</button>
                </div>
                <div class="admin-modal-body">
                    <p>Weet u zeker dat u een nieuwe uitnodiging wilt versturen naar <strong>{{ userToResend?.name }}</strong>?</p>
                    <p class="admin-text-muted">Er wordt een nieuw tijdelijk wachtwoord gegenereerd en verzonden naar {{ userToResend?.email }}.</p>
                </div>
                <div class="admin-modal-footer">
                    <button class="admin-btn admin-btn-secondary" @click="showResendModal = false">
                        Annuleren
                    </button>
                    <button class="admin-btn admin-btn-primary" @click="handleResend" :disabled="resending">
                        {{ resending ? 'Bezig...' : 'Versturen' }}
                    </button>
                </div>
            </div>
        </div>

        <!-- Toast Notification -->
        <ToastNotification 
            v-if="toast.show" 
            :message="toast.message" 
            :type="toast.type" 
            @close="toast.show = false" 
        />
    </AdminLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import AdminLayout from '@/components/admin/AdminLayout.vue';
import ToastNotification from '@/components/admin/ToastNotification.vue';
import { fetchUsers, createUser, deleteUser, resendInvitation } from '@/services/AdminService';

const loading = ref(true);
const users = ref([]);

// Add user modal state
const showAddModal = ref(false);
const newUser = ref({ name: '', email: '' });
const errors = ref({});
const formError = ref('');
const submitting = ref(false);

// Delete modal state
const showDeleteModal = ref(false);
const userToDelete = ref(null);
const deleting = ref(false);

// Resend modal state
const showResendModal = ref(false);
const userToResend = ref(null);
const resending = ref(false);

// Toast state
const toast = ref({ show: false, message: '', type: 'success' });

function showToast(message, type = 'success') {
    toast.value = { show: true, message, type };
    setTimeout(() => { toast.value.show = false; }, 4000);
}

function formatDate(dateString) {
    if (!dateString) return '-';
    const date = new Date(dateString);
    return date.toLocaleDateString('nl-BE', {
        day: 'numeric',
        month: 'short',
        year: 'numeric'
    });
}

async function loadUsers() {
    loading.value = true;
    try {
        const response = await fetchUsers();
        if (response && response.data) {
            users.value = response.data;
        }
    } catch (err) {
        console.error('Failed to load users:', err);
        showToast('Kon gebruikers niet laden', 'error');
    }
    loading.value = false;
}

function closeAddModal() {
    showAddModal.value = false;
    newUser.value = { name: '', email: '' };
    errors.value = {};
    formError.value = '';
}

async function handleAddUser() {
    errors.value = {};
    formError.value = '';
    submitting.value = true;

    try {
        const response = await createUser(newUser.value.name, newUser.value.email);
        
        if (response && response.data && response.data.status) {
            showToast('Gebruiker toegevoegd en uitnodiging verstuurd!', 'success');
            closeAddModal();
            await loadUsers();
        } else if (response && response.data && response.data.errors) {
            errors.value = response.data.errors;
        } else if (response && response.data && response.data.message) {
            formError.value = response.data.message;
        } else {
            formError.value = 'Er ging iets mis bij het toevoegen van de gebruiker.';
        }
    } catch (err) {
        console.error('Failed to create user:', err);
        formError.value = 'Er ging iets mis bij het toevoegen van de gebruiker.';
    }

    submitting.value = false;
}

function confirmDelete(user) {
    userToDelete.value = user;
    showDeleteModal.value = true;
}

async function handleDelete() {
    if (!userToDelete.value) return;
    
    deleting.value = true;
    try {
        const response = await deleteUser(userToDelete.value.id);
        
        if (response && response.success !== false) {
            showToast('Gebruiker verwijderd', 'success');
            showDeleteModal.value = false;
            userToDelete.value = null;
            await loadUsers();
        } else {
            showToast(response?.data?.message || 'Kon gebruiker niet verwijderen', 'error');
        }
    } catch (err) {
        console.error('Failed to delete user:', err);
        showToast('Kon gebruiker niet verwijderen', 'error');
    }
    deleting.value = false;
}

function confirmResendInvitation(user) {
    userToResend.value = user;
    showResendModal.value = true;
}

async function handleResend() {
    if (!userToResend.value) return;
    
    resending.value = true;
    try {
        const response = await resendInvitation(userToResend.value.id);
        
        if (response && response.data && response.data.status) {
            showToast('Uitnodiging opnieuw verstuurd!', 'success');
            showResendModal.value = false;
            userToResend.value = null;
        } else {
            showToast(response?.data?.message || 'Kon uitnodiging niet versturen', 'error');
        }
    } catch (err) {
        console.error('Failed to resend invitation:', err);
        showToast('Kon uitnodiging niet versturen', 'error');
    }
    resending.value = false;
}

onMounted(loadUsers);
</script>

<style scoped>
.user-name {
    font-weight: 500;
}

.user-email {
    color: var(--admin-muted);
}

.date-cell {
    color: var(--admin-muted);
    font-size: 0.875rem;
}

.action-buttons {
    display: flex;
    gap: 0.5rem;
}

.info-box {
    background: var(--admin-info-bg);
    border: 1px solid var(--admin-info);
    border-radius: var(--admin-radius);
    padding: 1rem;
    margin-top: 1rem;
    font-size: 0.875rem;
    color: var(--admin-info);
}

.admin-text-muted {
    color: var(--admin-muted);
    font-size: 0.875rem;
    margin-top: 0.5rem;
}

.btn-icon {
    margin-right: 0.5rem;
    font-weight: bold;
}

/* Modal styles */
.admin-modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
}

.admin-modal {
    background: white;
    border-radius: 8px;
    max-width: 500px;
    width: 90%;
    max-height: 90vh;
    overflow: auto;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.admin-modal-sm {
    max-width: 400px;
}

.admin-modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 1.5rem;
    border-bottom: 1px solid var(--admin-border);
}

.admin-modal-header h2 {
    margin: 0;
    font-size: 1.25rem;
}

.admin-modal-close {
    background: none;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    color: var(--admin-muted);
    line-height: 1;
}

.admin-modal-close:hover {
    color: var(--admin-text);
}

.admin-modal-body {
    padding: 1.5rem;
}

.admin-modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
    padding: 1rem 1.5rem;
    border-top: 1px solid var(--admin-border);
    background: var(--admin-bg);
}

.admin-form-group {
    margin-bottom: 1rem;
}

.admin-label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
}

.admin-input {
    width: 100%;
    padding: 0.625rem 0.75rem;
    border: 1px solid var(--admin-border);
    border-radius: 6px;
    font-size: 1rem;
}

.admin-input:focus {
    outline: none;
    border-color: var(--admin-primary);
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.admin-error {
    color: var(--admin-danger);
    font-size: 0.875rem;
    margin-top: 0.25rem;
}

/* Table styles */
.admin-table-container {
    overflow-x: auto;
}

.admin-table {
    width: 100%;
    border-collapse: collapse;
}

.admin-table th,
.admin-table td {
    padding: 1rem;
    text-align: left;
    border-bottom: 1px solid var(--admin-border);
}

.admin-table th {
    background: var(--admin-bg);
    font-weight: 600;
    font-size: 0.875rem;
    color: var(--admin-muted);
    text-transform: uppercase;
}

.admin-table tbody tr:hover {
    background: var(--admin-bg);
}
</style>
