<script setup>
    import { getProfileData, login, logout, changePassword } from '@/services/AdminService';
    import { onMounted, ref } from 'vue';
    import { useRouter } from 'vue-router';

    const router = useRouter();

    const doneChecking = ref(false);
    const loggedIn = ref(false);
    const needsPasswordUpdate = ref(false);
    const userProfile = ref({});

    async function setProfileData() {
        const response = await getProfileData();

        if (response && response.status) {
            userProfile.value = response.data;
            loggedIn.value = true;
            
            // Check if password update is required
            if (response.data.needs_password_update) {
                needsPasswordUpdate.value = true;
            } else {
                // Redirect to new dashboard when logged in and no password change needed
                router.push({ name: 'adminDashboard' });
            }
        }

        doneChecking.value = true;
    }

    // Login form fields
    const email = ref("");
    const password = ref("");
    const emailErrors = ref([]);
    const passwordErrors = ref([]);
    const loginFailedMessage = ref("");

    async function tryToLogin() {
        emailErrors.value = [];
        passwordErrors.value = [];
        loginFailedMessage.value = "";

        const response = await login(email.value, password.value);

        if (response.data && response.data.status && response.data.token) {
            // Store the JWT token
            localStorage.setItem('admin_token', response.data.token);
            userProfile.value = response.data.user;
            loggedIn.value = true;
            
            // Check if password update is required
            if (response.data.needs_password_update) {
                needsPasswordUpdate.value = true;
            } else {
                // Redirect to dashboard
                router.push({ name: 'adminDashboard' });
            }
        } else if (response.data && response.data.errors) {
            if (response.data.errors.email) {
                emailErrors.value = response.data.errors.email;
            }
            if (response.data.errors.password) {
                passwordErrors.value = response.data.errors.password;
            }
        } else if (response.data && response.data.message) {
            loginFailedMessage.value = response.data.message;
        } else {
            loginFailedMessage.value = "Er ging iets mis bij het inloggen.";
        }
    }

    // Password change form fields
    const currentPassword = ref("");
    const newPassword = ref("");
    const newPasswordConfirm = ref("");
    const passwordChangeErrors = ref({});
    const passwordChangeMessage = ref("");

    async function tryChangePassword() {
        passwordChangeErrors.value = {};
        passwordChangeMessage.value = "";

        // Validate
        if (!currentPassword.value) {
            passwordChangeErrors.value.current = ["Huidig wachtwoord is verplicht"];
        }
        if (newPassword.value.length < 8) {
            passwordChangeErrors.value.new = ["Nieuw wachtwoord moet minimaal 8 tekens zijn"];
        }
        if (newPassword.value !== newPasswordConfirm.value) {
            passwordChangeErrors.value.confirm = ["Wachtwoorden komen niet overeen"];
        }

        if (Object.keys(passwordChangeErrors.value).length > 0) {
            return;
        }

        const response = await changePassword(currentPassword.value, newPassword.value);

        if (response.data && response.data.status && response.data.token) {
            // Update the token with the new one (without needs_password_update)
            localStorage.setItem('admin_token', response.data.token);
            needsPasswordUpdate.value = false;
            // Redirect to dashboard
            router.push({ name: 'adminDashboard' });
        } else if (response.data && response.data.errors) {
            passwordChangeErrors.value = response.data.errors;
        } else if (response.data && response.data.message) {
            passwordChangeMessage.value = response.data.message;
        } else {
            passwordChangeMessage.value = "Er ging iets mis bij het wijzigen van het wachtwoord.";
        }
    }

    async function tryToLogout() {
        const response = await logout();

        // Clear token regardless of response
        localStorage.removeItem('admin_token');
        userProfile.value = {};
        loggedIn.value = false;
        needsPasswordUpdate.value = false;
        
        if (!response.success) {
            console.warn("Logout request failed, but token cleared locally");
        }
    }

    onMounted(setProfileData);
</script>

<template>
    <section v-if="!loggedIn && !doneChecking" id="loader">
        <p>Even geduld, we kijken of u al ingelogd bent...</p>
    </section>
    
    <!-- Password change required form -->
    <section v-if="loggedIn && needsPasswordUpdate" id="password-change">
        <h1>Wachtwoord Wijzigen Vereist</h1>
        <p class="password-intro">U moet uw wachtwoord wijzigen voordat u verder kunt.</p>
        <form method="post" @submit.prevent="tryChangePassword">
            <div>
                <label for="current-password">Huidig Wachtwoord</label>
                <input v-model="currentPassword" type="password" id="current-password" name="current-password" autocomplete="current-password" placeholder="Tijdelijk wachtwoord uit logs">
            </div>
            <div>
                <label for="new-password">Nieuw Wachtwoord</label>
                <input v-model="newPassword" type="password" id="new-password" name="new-password" autocomplete="new-password" placeholder="Minimaal 8 tekens">
            </div>
            <div>
                <label for="new-password-confirm">Bevestig Nieuw Wachtwoord</label>
                <input v-model="newPasswordConfirm" type="password" id="new-password-confirm" name="new-password-confirm" autocomplete="new-password" placeholder="Herhaal nieuw wachtwoord">
            </div>
            <input type="submit" value="Wachtwoord Wijzigen">
        </form>
        <div v-if="passwordChangeMessage || Object.keys(passwordChangeErrors).length > 0" id="errors">
            <p v-if="passwordChangeMessage">{{ passwordChangeMessage }}</p>
            <ul v-for="(errors, field) in passwordChangeErrors" :key="field">
                <li v-for="error in errors" :key="error">{{ error }}</li>
            </ul>
        </div>
        <button type="button" @click="tryToLogout" class="logout-btn">Uitloggen</button>
    </section>

    <!-- Regular login form -->
    <section v-if="!loggedIn && doneChecking" id="login">
        <h1>Admin Login</h1>
        <form method="post" @submit.prevent="tryToLogin">
            <div>
                <label for="email">Email / Gebruikersnaam</label>
                <input v-model="email" type="text" id="email" name="email" autocomplete="username">
            </div>
            <div>
                <label for="password">Wachtwoord</label>
                <input v-model="password" type="password" id="password" name="password" autocomplete="current-password">
            </div>
            <input type="submit" value="Inloggen">
        </form>
        <div v-if="loginFailedMessage || emailErrors.length > 0 || passwordErrors.length > 0" id="errors">
            <p v-if="loginFailedMessage">{{ loginFailedMessage }}</p>
            <ul v-if="emailErrors.length > 0">
                <li v-for="error in emailErrors" :key="error">{{ error }}</li>
            </ul>
            <ul v-if="passwordErrors.length > 0">
                <li v-for="error in passwordErrors" :key="error">{{ error }}</li>
            </ul>
        </div>
        <p id="login-cookie-notice">Door in te loggen in het admin-paneel gaat u ermee akkoord dat er cookies, essentieel voor de login-functionaliteit, gebruikt zullen worden.</p>
    </section>
    
    <header v-show="loggedIn && !needsPasswordUpdate">
        <img src="@/assets/media/logo-full-black.webp" class="logo">
        <article id="userinfo">
            <div>
                <p>{{ userProfile.name }}</p>
                <p>{{ userProfile.email }}</p>
            </div>
            <button type="button" @click="tryToLogout">Uitloggen</button>
        </article>
    </header>
    <main v-show="loggedIn && !needsPasswordUpdate" id="admin-dash">
        <section>
            <h2>website inhoud</h2>
            <div>
                <figure @click="() => router.push({ name: 'adminEvents' })">
                    <div class="icon icon-event"></div>
                    <p>Evenementen</p>
                </figure>
                <figure @click="() => router.push({ name: 'adminGeocaches' })">
                    <div class="icon icon-cache"></div>
                    <p>Geocaches</p>
                </figure>
                <figure @click="() => router.push({ name: 'adminSocials' })">
                    <div class="icon icon-social"></div>
                    <p>Sociale Media</p>
                </figure>
                <figure @click="() => router.push({ name: 'adminMessages' })">
                    <div class="icon icon-msg"></div>
                    <p>Berichten</p>
                </figure>
            </div>
        </section>
        <section>
            <h2>Technisch & Andere</h2>
            <div>
                <figure @click="() => router.push({ name: 'adminStatic' })">
                    <div class="icon icon-static"></div>
                    <p>Vertaaltabel</p>
                </figure>
                <figure @click="() => router.push({ name: 'adminLanguages' })">
                    <div class="icon icon-lang"></div>
                    <p>Talen</p>
                </figure>
                <figure @click="() => router.push({ name: 'adminContactForm' })">
                    <div class="icon icon-mail"></div>
                    <p>contact formulier</p>
                </figure>
            </div>
        </section>
    </main>
</template>

<style scoped>
#loader {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    height: 100vh;
    background: var(--admin-bg, #f8fafc);
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    color: var(--admin-text-secondary, #64748b);
}

#loader::before {
    content: '';
    width: 2rem;
    height: 2rem;
    border: 3px solid var(--admin-border, #e2e8f0);
    border-top-color: var(--admin-primary, #0d9488);
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 1rem;
}

@keyframes spin {
    to { transform: rotate(360deg); }
}

#login, #password-change {
    max-width: 420px;
    margin: 0 auto;
    padding: 2.5rem;
    background: var(--admin-surface, #ffffff);
    border-radius: var(--admin-radius-lg, 0.75rem);
    box-shadow: var(--admin-shadow-lg, 0 10px 15px -3px rgb(0 0 0 / 0.1));
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
}

#login h1, #password-change h1 {
    text-align: center;
    margin-bottom: 2rem;
    font-size: 1.5rem;
    font-weight: 600;
    color: var(--admin-text, #1e293b);
}

#password-change h1 {
    font-size: 1.25rem;
}

.password-intro {
    text-align: center;
    margin-bottom: 1rem;
    color: var(--admin-text-secondary, #64748b);
    font-size: 0.9375rem;
}

.password-hint {
    font-size: 0.8125rem;
    color: var(--admin-text-secondary, #64748b);
    background: var(--admin-warning-bg, #fffbeb);
    padding: 0.75rem 1rem;
    border-radius: var(--admin-radius, 0.5rem);
    margin-bottom: 1.5rem;
    text-align: center;
    border: 1px solid rgba(245, 158, 11, 0.2);
}

form div {
    margin-bottom: 1.25rem;
}

form label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    font-size: 0.875rem;
    color: var(--admin-text, #1e293b);
}

form input[type="text"],
form input[type="email"],
form input[type="password"] {
    width: 100%;
    padding: 0.75rem 1rem;
    border: 1px solid var(--admin-border, #e2e8f0);
    border-radius: var(--admin-radius, 0.5rem);
    font-size: 0.9375rem;
    background: var(--admin-surface, #ffffff);
    color: var(--admin-text, #1e293b);
    transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

form input[type="text"]:focus,
form input[type="email"]:focus,
form input[type="password"]:focus {
    outline: none;
    border-color: var(--admin-primary, #0d9488);
    box-shadow: 0 0 0 3px var(--admin-primary-bg, rgba(13, 148, 136, 0.1));
}

form input[type="text"]::placeholder,
form input[type="email"]::placeholder,
form input[type="password"]::placeholder {
    color: var(--admin-text-muted, #94a3b8);
}

form input[type="submit"] {
    width: 100%;
    padding: 0.875rem 1rem;
    background: var(--admin-primary, #0d9488);
    color: white;
    border: none;
    border-radius: var(--admin-radius, 0.5rem);
    font-size: 0.9375rem;
    font-weight: 500;
    cursor: pointer;
    transition: background 0.15s ease;
    margin-top: 0.5rem;
}

form input[type="submit"]:hover {
    background: var(--admin-primary-dark, #0f766e);
}

#errors {
    margin-top: 1rem;
    padding: 0.75rem 1rem;
    background: var(--admin-danger-bg, #fef2f2);
    border: 1px solid rgba(239, 68, 68, 0.2);
    border-radius: var(--admin-radius, 0.5rem);
    color: var(--admin-danger, #ef4444);
    font-size: 0.875rem;
}

#errors:empty {
    display: none;
}

#errors ul {
    list-style: none;
    padding: 0;
    margin: 0;
}

#errors ul li {
    padding: 0.25rem 0;
}

#login-cookie-notice {
    margin-top: 1.5rem;
    font-size: 0.75rem;
    color: var(--admin-text-muted, #94a3b8);
    text-align: center;
    line-height: 1.5;
}

.logout-btn {
    width: 100%;
    margin-top: 1.5rem;
    padding: 0.75rem 1rem;
    background: var(--admin-surface, #ffffff);
    color: var(--admin-text-secondary, #64748b);
    border: 1px solid var(--admin-border, #e2e8f0);
    border-radius: var(--admin-radius, 0.5rem);
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.15s ease;
}

.logout-btn:hover {
    background: var(--admin-surface-hover, #f1f5f9);
    border-color: var(--admin-text-muted, #94a3b8);
}

/* Background for auth pages */
#login, #password-change {
    background: var(--admin-surface, #ffffff);
}

body:has(#login), body:has(#password-change) {
    background: var(--admin-bg, #f8fafc);
}

header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 2rem;
    background: var(--admin-surface, #ffffff);
    box-shadow: var(--admin-shadow, 0 1px 3px 0 rgb(0 0 0 / 0.1));
    border-bottom: 1px solid var(--admin-border, #e2e8f0);
}

header .logo {
    height: 50px;
}

#userinfo {
    display: flex;
    align-items: center;
    gap: 1rem;
}

#userinfo div p {
    margin: 0;
    color: var(--admin-text, #1e293b);
}

#userinfo div p:last-child {
    font-size: 0.875rem;
    color: var(--admin-text-muted, #94a3b8);
}

#userinfo button {
    padding: 0.5rem 1rem;
    background: var(--admin-danger, #ef4444);
    color: white;
    border: none;
    border-radius: var(--admin-radius, 0.5rem);
    font-size: 0.875rem;
    cursor: pointer;
    transition: background 0.15s ease;
}

#userinfo button:hover {
    background: #dc2626;
}

#admin-dash {
    padding: 2rem;
    background: var(--admin-bg, #f8fafc);
    min-height: calc(100vh - 82px);
}

#admin-dash section {
    margin-bottom: 2rem;
}

#admin-dash h2 {
    text-transform: uppercase;
    color: var(--admin-text-secondary, #64748b);
    margin-bottom: 1rem;
    font-size: 0.75rem;
    font-weight: 600;
    letter-spacing: 0.05em;
}

#admin-dash section > div {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
}

#admin-dash figure {
    width: 150px;
    padding: 1.5rem;
    background: var(--admin-surface, #ffffff);
    border-radius: var(--admin-radius-lg, 0.75rem);
    box-shadow: var(--admin-shadow, 0 1px 3px 0 rgb(0 0 0 / 0.1));
    border: 1px solid var(--admin-border, #e2e8f0);
    text-align: center;
    cursor: pointer;
    transition: transform 0.2s, box-shadow 0.2s, border-color 0.2s;
    margin: 0;
}

#admin-dash figure:hover {
    transform: translateY(-2px);
    box-shadow: var(--admin-shadow-md, 0 4px 6px -1px rgb(0 0 0 / 0.1));
    border-color: var(--admin-primary, #0d9488);
}

#admin-dash figure .icon {
    width: 50px;
    height: 50px;
    margin: 0 auto 1rem;
    background: var(--admin-primary, #0d9488);
    border-radius: var(--admin-radius, 0.5rem);
}

#admin-dash figure p {
    margin: 0;
    font-size: 0.875rem;
    color: var(--admin-text, #1e293b);
    font-weight: 500;
}
</style>
