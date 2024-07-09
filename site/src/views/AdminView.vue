<script setup>
    import { getProfileData, login, logout} from '@/services/AdminService';
    import { onMounted, ref } from 'vue';
    import { useRouter } from 'vue-router';

    const router = useRouter();

    const doneChecking = ref(false);
    const loggedIn = ref(false);
    const userProfile = ref({});

    async function setProfileData() {
        const response = await getProfileData();

        if (response.status) {
            userProfile.value = response.data;
            loggedIn.value = true;
        }

        doneChecking.value = true;
    }

    const email = ref("");
    const password = ref("");
    const emailErrors = ref([]);
    const passwordErrors = ref([]);
    const loginFailedMessage = ref("");

    async function tryToLogin() {
        const response = await login(email.value, password.value);

        if (response.data.status) {
            setProfileData();
        } else if (response.data && response.data.errors) {
            if (response.data.errors.email) {
                emailErrors.value = response.data.errors.email;
            }
            if (response.data.errors.password) {
                passwordErrors.value = response.data.errors.password;
            }
        } else if (response.data.message) {
            loginFailedMessage.value = response.data.message;
        }
    }

    async function tryToLogout() {
        const response = await logout();

        if (response.status) {
            userProfile.value = {};
            loggedIn.value = false;
        } else {
            window.alert("Er is iets foutgelopen bij het uitloggen (geen antwoord teruggekregen van de server). \n Refresh de pagina.");
        }
    }

    onMounted(setProfileData);
</script>

<template>
    <section v-if="!loggedIn && !doneChecking" id="loader">
        <p>Even geduld, we kijken of u al ingelogd bent...</p>
    </section>
    <section v-if="!loggedIn && doneChecking" id="login">
        <h1>Admin Login</h1>
        <form method="post" @submit.prevent="tryToLogin">
            <div>
                <label for="email">Email</label>
                <input v-model="email" type="email" id="email" name="email" autocomplete="email">
            </div>
            <div>
                <label for="password">Password</label>
                <input v-model="password" type="password" id="password" name="password" autocomplete="current-password">
            </div>
            <input type="submit" value="Inloggen">
        </form>
        <div id="errors">
            <p v-show="loginFailedMessage">{{ loginFailedMessage }}</p>
            <ul v-show="emailErrors.length > 0">
                <li v-for="error in emailErrors">{{ error }}</li>
            </ul>
            <ul v-show="passwordErrors.length> 0">
                <li v-for="error in passwordErrors">{{ error }}</li>
            </ul>
        </div>
        <p id="login-cookie-notice">Door in te loggen in het admin-paneel gaat u ermee akkoord dat er cookies, essentieel voor de login-functionaliteit, gebruikt zullen worden.</p>
    </section>
    <header v-show="loggedIn">
        <picture>
            <source srcset="../assets/media/logo-full-black.webp" media="(prefers-color-scheme: light)" class="logo">
            <source srcset="../assets/media/logo-full-white.webp" media="(prefers-color-scheme: dark)" class="logo">
            <img src="../assets/media/logo-full-black.webp" class="logo">
        </picture>
        <article id="userinfo">
            <div>
                <p>{{ userProfile.name }}</p>
                <p>{{ userProfile.email }}</p>
            </div>
            <button type="button" @click="tryToLogout">Uitloggen</button>
        </article>
    </header>
    <main v-show="loggedIn" id="admin-dash">
        <section>
            <h2>website inhoud</h2>
            <div>
                <figure @click="() => router.push('adminEvents')">
                    <div class="icon icon-event"></div>
                    <p>Evenementen</p>
                </figure>
                <figure @click="() => router.push('adminGeocaches')">
                    <div class="icon icon-cache"></div>
                    <p>Geocaches</p>
                </figure>
                <figure @click="() => router.push('adminSocials')">
                    <div class="icon icon-social"></div>
                    <p>Sociale Media</p>
                </figure>
                <figure @click="() => router.push('adminMessages')">
                    <div class="icon icon-msg"></div>
                    <p>Berichten</p>
                </figure>
                <!--<figure @click="() => router.push('adminShop')">
                    <div class="icon icon-shop"></div>
                    <p>Webshop</p>
                </figure>-->
            </div>
        </section>
        <section>
            <h2>Technisch</h2>
            <div>
                <figure @click="() => router.push('adminStatic')">
                    <div class="icon icon-static"></div>
                    <p>Vertaaltabel</p>
                </figure>
                <figure @click="() => router.push('adminLanguages')">
                    <div class="icon icon-lang"></div>
                    <p>Talen</p>
                </figure>
            </div>
        </section>
    </main>
</template>

<style scoped>
    #loader, #login {
        height: 100vh;
        width: 100vw;
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
    }

    h1 {
        font-weight: bold;
        text-transform: capitalize;
        margin-bottom: 2dvh;
    }

    form {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 1rem;
    }

    form div {
        display: flex;
        flex-direction: column;
        gap: 0.2rem;
    }

    form label {
        text-transform: capitalize;
    }

    form > div input {
        width: 20rem;
        max-width: 90vw;
    }

    form input {
        height: 2rem;
        border-radius: 0.3rem;
        border: solid 0.1rem var(--color-text);
        outline: none;
        font-family: inherit;
        font-size: 0.85rem;
        box-sizing: border-box;
        background-color: var(--color-background);
        color: var(--color-text)
    }

    form input[type="submit"] {
        background-color: var(--color-secondary);
        border: none;
        width: 9rem;
        height: 2rem;
        font-family: inherit;
        border-radius: 0.4rem;
        box-shadow: var(--color-background2) 0.5rem 0.5rem;
        text-transform: capitalize;
        font-weight: bold;
        scale: 100%;
        transition: scale 0.15s;
        margin-top: 2rem;
    }

    form input[type="submit"]:hover {
        cursor: pointer;
        scale: 103%;
        transition: scale 0.25s;
    }

    #login-cookie-notice {
        position: absolute;
        bottom: 0.5rem;
        color: var(--color-text);
        opacity: 50%;
    }

    header {
        padding: 0.5rem;
        display: flex;
        justify-content: space-between
    }

    header .logo {
        height: 3rem;
    }

    header #userinfo {
        display: flex;
        align-items: center;
        gap: 1rem;
    }

    header #userinfo p {
        text-align: end;
    }

    header button {
        background-color: var(--color-primary);
        color: var(--color-text3);
        font-weight: bold;
        padding: 0.4rem 1.5rem;
        text-decoration: none;
        border-radius: 0.4rem;
        transform: scale(100%);
        transition: transform 0.15s;
    }

    header button:hover {
        cursor: pointer;
        transform: scale(103%);
        transition: transform 0.25s;
    }

    #admin-dash {
        height: 100%;
        margin: 2rem auto;
        width: fit-content;
        max-width: 90vw;
        display: flex;
        flex-direction: column;
        justify-content: center;
    }

    #admin-dash section {
        margin: 0 1rem 3rem 1rem;
    }

    #admin-dash section h2 {
        text-transform: capitalize;
        font-size: 1.5rem;
        font-weight: bold;
        margin-bottom: 0.5rem;
    }

    #admin-dash section div {
        display: flex;
        align-items: center;
        justify-content: flex-start;
        flex-wrap: wrap;
        gap: 2rem;
    }

    #admin-dash section figure {
        width: 10rem;
        display: flex;
        flex-direction: column;
        gap: 1rem;
        align-items: center;
        padding: 1rem;
        text-align: center;
    }

    #admin-dash section figure:hover {
        cursor: pointer;
    }

    #admin-dash section figure:hover div.icon {
        background-color: var(--color-primary);
        transform: scale(105%);
        transition: transform 0.2s;
    }

    #admin-dash div.icon {
        width: 3rem;
        height: 3rem;
        background-color: var(--color-text);
        transform: scale(100%);
        transition: transform 0.2s;
    }

    div.icon-event {
        mask: url(../assets/media/calendar.svg);
        mask-size: contain;
    }

    div.icon-cache {
        mask: url(../assets/media/box.svg);
        mask-size: contain;
    }

    div.icon-social {
        mask: url(../assets/media/share-2.svg);
        mask-size: contain;
    }

    div.icon-msg {
        mask: url(../assets/media/message-square.svg);
        mask-size: contain;
    }

    div.icon-shop {
        mask: url(../assets/media/shopping-cart.svg);
        mask-size: contain;
    }

    div.icon-static {
        mask: url(../assets/media/file-text.svg);
        mask-size: contain;
    }

    div.icon-lang {
        mask: url(../assets/media/book-open.svg);
        mask-size: contain;
    }
</style>