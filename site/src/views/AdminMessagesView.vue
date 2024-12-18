<script setup>
    import { deleteMessage, fetchMessages, getProfileData, postMessage, updateMessage } from '@/services/AdminService';
    import { onMounted, ref, computed, toRaw } from 'vue';
    import { useRouter } from 'vue-router';
    import StaticContentProvider from '@/services/StaticContentService';

    const router = useRouter();

    async function verifyLogin() {
        const response = await getProfileData();

        if (response.status) {
            fetchData();
        } else {
            router.push({ name: "admin" });
        }
    }

    const messages = ref([]);
    const search = ref("");

    async function fetchData() {
        const response = await fetchMessages();
        if (response.access_denied) {
            window.alert(response.access_denied);
        } else {
            messages.value = response;
        }
    }

    const filteredMessages = computed(() => {
        if (!search.value) {
            return messages.value;
        }

        return messages.value.filter(message =>
            message.translations.find(translation => translation.lang_code === "NL").title.toLowerCase().includes(search.value.toLowerCase())
        );
    });

    function dateTimeFormatter(dateTimeString) {
        if (!dateTimeString) {
            "n/a";
        }

        const date = new Date(dateTimeString);
        if (isNaN(date)) {
            return "parsing error";
        }

        return date.toLocaleString("nl-be", {
            year: "numeric",
            month: "short",
            day: "2-digit",
            hour: "2-digit",
            minute: "2-digit"
        });
    }

    const currentlyEditing = ref(-1);
    const currentlyEditingData = ref({});

    function editMessage(message) {
        currentlyEditing.value = message.id;
        currentlyEditingData.value = structuredClone(toRaw(message));
        delete currentlyEditingData.value.id;
    }

    function createMessage() {
        const translations = [];

        StaticContentProvider.LANGUAGES.forEach(language => {
            translations.push({lang_code: language.code, title: "", body: ""});
        });

        currentlyEditingData.value = {
            state: "DRAFT",
            translations: translations
        }

        currentlyEditing.value = -2;
    }

    function stopEditing() {
        currentlyEditing.value = -1;
        currentlyEditingData.value = {};
    }

    const saving = ref(false);

    async function save(state) {
        if (!verifyInputs()) {
            return;
        }
        saving.value = true;

        currentlyEditingData.value.state = state;        
        const formData = setFormData();

        const response = await postMessage(formData);
        handleResponse(response);
    }

    async function update(state) {
        if (!verifyInputs()) {
            return;
        }
        saving.value = true;

        currentlyEditingData.value.state = state;
        const formData = setFormData();

        const response = await updateMessage(currentlyEditing.value, formData);
        handleResponse(response);
    }

    async function remove() {
        if (window.confirm(`Zeker dat je bericht: "${currentlyEditingData.value.translations.find(translation => translation.lang_code === "NL").title}" wilt verwijderen?`)) {
            saving.value = true;
            
            const response = await deleteMessage(currentlyEditing.value);

            if (!response.success) {
                window.alert(response.error);
                if (response.error.includes("Unautherized")) {
                    router.push({ name: "admin" });
                }
                saving.value = false;
            } else if (response.success && response.data && response.data.deleted) {
                currentlyEditingData.value = {};
                currentlyEditing.value = -1;
                saving.value = false;
                fetchData();
            } else {
                window.alert("onbekende fout");
                saving.value = false;
            }
        }
    }

    function verifyInputs() {
        const form = document.querySelector("#messageEdit");

        for (let element of form.elements) {
            if (!element.checkValidity()) {
                element.reportValidity();
                return false;
            }
        }

        if (StaticContentProvider.LANGUAGES.length !== currentlyEditingData.value.translations.length) {
            const translations = [];
            StaticContentProvider.LANGUAGES.forEach(language => {
                translations.push({lang_code: language.lang_code, title: "", body: ""});
            });
            currentlyEditingData.value.translations = translations;

            window.alert("Ontbrekende vertaling, talen zonder vertalingen gevonden, gelieve dit aan te vullen.");
            return false;
        }

        return true;
    }

    function setFormData() {
        const formData = new FormData();

        formData.append("translations", JSON.stringify(currentlyEditingData.value.translations));
        formData.append("state", currentlyEditingData.value.state);

        return formData;
    }

    function handleResponse(response) {
        if (!response.success) {
            window.alert(response.error);
            if (response.error.includes("Unautherized")) {
                router.push({ name: "admin" });
            }
            saving.value = false;
        } else if (response.success && response.data) {
            if (response.data.errors) {
                let errorsString = ""
                for (let field in response.data.errors) {
                    errorsString += `${field}: ${response.data.errors[field]}\n`;
                }
                window.alert(errorsString);
                saving.value = false;
            } else if (response.data.data) {
                currentlyEditingData.value = {};
                currentlyEditing.value = -1;
                saving.value = false;
                fetchData();
            }
        } else {
            window.alert("onbekende fout");
            saving.value = false;
        }
    }

    onMounted(verifyLogin);
</script>

<template>
    <header>
        <div>
            <button @click="() => (currentlyEditing === -1) ? router.push({ name: 'admin'}) : stopEditing()" type="button">Terug</button>
            <button v-show="currentlyEditing === -1" type="button" @click="createMessage">Nieuw Bericht</button>
        </div>
        <form v-show="messages.length > 0 && messages[0] !== 'loading' && currentlyEditing === -1" method="post" @submit.prevent="" id="search">
            <div>
                <input v-model="search" type="search" id="search" name="search" autocomplete="search" required placeholder="Zoeken">
            </div>
        </form>
    </header>
    <main>
        <table v-show="currentlyEditing === -1">
            <thead>
                <th>ID</th>
                <th>Status</th>
                <th>Titel</th>
                <th>Laatst Aangepast</th>
                <th></th>
            </thead>
            <tbody>
                <tr v-for="message in filteredMessages">
                    <td>{{ message.id }}</td>
                    <td :class="message.state">{{ message.state }}</td>
                    <td>{{ message.translations.find(translation => translation.lang_code === "NL").title }}</td>
                    <td>{{ dateTimeFormatter(message.updated_at) }}</td>
                    <td><div @click="() => editMessage(message)" class="icon-edit"></div></td>
                </tr>
            </tbody>
        </table>
        <form @submit.prevent="" method="post" id="messageEdit" v-show="currentlyEditing !== -1">
            <section class="general">
                <p><i>Bij een lang bericht is het aangeraden op te splitsen in titel en body, zo niet, is een titel alleen voldoende.</i></p>
                <div v-for="(translation, index) in currentlyEditingData.translations">
                    <label for="title">{{ `Titel - ${translation.lang_code}` }}<span>*</span></label>
                    <input v-model="currentlyEditingData.translations[index].title" type="text" id="title" name="title" max="200" required>
                    <label for="body">{{ `Body - ${translation.lang_code}` }}</label>
                    <textarea v-model="currentlyEditingData.translations[index].body" id="body" name="body" max="2000" required></textarea>
                </div>
            </section>
            <section class="buttons">
                <button class="btn-red" @click="remove()" type="button" v-if="currentlyEditing !== -2">Verwijderen</button>
                    <button class="btn-orange" @click="update('ARCHIVED')" type="button" v-if="currentlyEditingData.state === 'ONLINE'">Archiveren</button>
                    <button class="btn-orange" @click="() => (currentlyEditing === -2) ? save('DRAFT') : update('DRAFT')" type="button" v-if="currentlyEditingData.state !== 'ONLINE'">Opslaan Als Concept</button>
                    <button @click="() => (currentlyEditing === -2) ? save('ONLINE') : update('ONLINE')" type="button">Publiceren</button>
            </section>
        </form>
        <div id="overlay" v-show="saving"></div>
    </main>
</template>

<style scoped>
header {
        display: flex;
        justify-content: space-between;
    }

    main {
        overflow-y: auto;
    }

    button {
        color: var(--color-text);
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
        margin: 1rem;
    }

    button:hover {
        cursor: pointer;
        scale: 103%;
        transition: scale 0.25s;
    }

    #search {
        display: flex;
        justify-content: flex-end;
        padding: 1rem;
        gap: 1rem;
        background: none !important;
    }

    #search > div input {
        width: 20rem;
        max-width: 95vw;
    }

    #search > div {
        position: relative;
    }

    #search > div::before {
        content: '';
        height: 1.3rem;
        width: 1.3rem;
        position: absolute;
        left: 0.5rem;
        top: 0.3rem;
        background-color: var(--color-text);
        mask: url(@/assets/media/search.svg);
        mask-size: contain;
        mask-repeat: no-repeat;
        mask-position: center;
    }

    #search input {
        padding: 0 0.5rem 0 2rem;
        text-transform: capitalize;
    }

    form input, form textarea {
        height: 2rem;
        border-radius: 0.3rem;
        border: solid 0.1rem var(--color-text);
        outline: none;
        font-family: inherit;
        font-size: 1rem;
        box-sizing: border-box;
        line-height: 2rem;
        background-color: var(--color-background);
        color: var(--color-text)
    }

    table {
        width: 90vw;
        max-width: 80rem;
        margin: 3rem auto;
    }

    thead {
        background-color: var(--color-primary);
        color: var(--color-text);
    }

    th, td, tr {
        text-align: center;
        padding: 0 0.5rem;
    }

    tbody tr:nth-child(even) {
        background-color: var(--color-background2);
    }

    th:nth-child(3), td:nth-child(3) {
        text-align: start;
    }

    td {
        height: 1.8rem;
    }

    td div {
        display: flex;
        justify-content: center;
    }

    td img {
        height: 1.5rem;
    }

    th:nth-child(1) {
        width: 3rem;
    }

    th:nth-child(2) {
        width: 5rem;
    }

    td:nth-child(4) {
        width: 10rem;
    }

    th:nth-child(5) {
        width: 3rem;
    }

    td.ONLINE, td.true {
        background-color: rgb(125, 209, 0);
    }

    td.DRAFT {
        background-color: rgb(255, 136, 0);
        color: var(--color-text3);
    }

    td.ARCHIVED, td.false {
        background-color: red;
        color: var(--color-text4);
    }

    div.icon-edit {
        margin: auto;
        height: 60%;
        aspect-ratio: 1 / 1;
        mask: url(@/assets/media/edit.svg);
        mask-size: contain;
        background-color: var(--color-text);
        mask-repeat: no-repeat;
        mask-position: center;
    }

    div.icon-edit:hover {
        cursor: pointer;
        background-color: var(--color-primary);
    }

    #messageEdit {
        display: grid;
        grid-template-columns: 1fr;
        padding: 0 1rem 1rem 1rem;
    }

    #messageEdit  p {
        margin-bottom: 2rem;
        position: sticky;
        top: 0;
        padding: 0.5rem 0;
        background-color: var(--color-background)
    }

    #messageEdit section.general > div {
        display: flex;
        flex-direction: column;
        margin-bottom: 5rem;
    }

    #messageEdit section.general {
        display: flex;
        height: max-content;
        flex-direction: column;
    }

    #messageEdit section.image-upload {
        align-items: center;
    }

    #messageEdit label {
        margin-bottom: 0.3rem;
    }

    #messageEdit label:not(:first-child) {
        margin-top: 0.5rem;
    }

    #messageEdit label span {
        color: red;
    }

    #messageEdit label i {
        opacity: 60%;
    }

    #messageEdit input, #messageEdit textarea {
        width: 80%;
        font-size: 0.85rem;
    }

    #messageEdit textarea {
        resize: vertical;
        height: min-content;
    }

    #messageEdit section.general input:nth-of-type(4), #messageEdit section.general input:nth-of-type(5) {
        text-align: center;
        width: 12rem;
    }

    #messageEdit select {
        width: 12rem;
        height: 2rem;
        font: inherit;
        font-size: 0.9rem;
        border-radius: 0.3rem;
        border: solid 0.1rem var(--color-text);
        padding: 0.2rem 0;
        outline: none;
        background-color: var(--color-background);
        color: var(--color-text);
    }

    #messageEdit section.image-upload input[type="file"] {
        width: 18rem;
        padding-left: 0.5rem;
    }

    #messageEdit section.image-upload img {
        margin-top: 3rem;
        height: 30rem;
        max-width: 80%;
        object-fit: contain;
        object-position: 0;
    }

    #messageEdit section.translations {
        margin-top: 3rem;
        display: flex;
        justify-content: center;
        gap: 3rem;
        flex-wrap: wrap;
    }

    #messageEdit section.translations div {
        width: 40rem;
        border-radius: 0.3rem;
    }

    #messageEdit section.buttons {
        margin: 1rem;
        display: flex;
        gap: 1rem;
        justify-content: center;
    }

    #messageEdit button.btn-red {
        background-color: rgb(255, 66, 66);
        box-shadow: #e0b8aa 0.5rem 0.5rem;
    }

    #messageEdit button.btn-orange {
        background-color: rgb(255, 138, 70);
        box-shadow: #e9cdc2 0.5rem 0.5rem;
    }

    #overlay {
        position: absolute;
        top: 0;
        right: 0;
        bottom: 0;
        left: 0;
        z-index: 10;
        background-color: rgba(0, 0, 0, 0.5);
    }

    @media (prefers-color-scheme: dark) {
        thead {
            background-color: var(--color-secondary);
        }

        td.ONLINE, td.true {
            background-color: green;
        }

        td.ARCHIVED, td.false {
            background-color: darkred;
        }

        #overlay {
            background-color: rgba(60, 60, 60, 0.5);
        }

        #messageEdit button.btn-red {
            background-color: rgb(255, 87, 57);
            color: var(--color-text3);
            box-shadow: #360f01 0.5rem 0.5rem;
        }

        #messageEdit button.btn-orange {
            background-color: rgb(233, 177, 22);
            color: var(--color-text3);
            box-shadow: #362a01 0.5rem 0.5rem;
        }
    }
</style>