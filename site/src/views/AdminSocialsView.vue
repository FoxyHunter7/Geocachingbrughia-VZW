<script setup>
    import { deleteSocial, fetchSocials, getProfileData, postSocial, updateSocial } from '@/services/AdminService';
    import { onMounted, ref, computed, toRaw } from 'vue';
    import config from '@/data/config.js';
    import { useRouter } from 'vue-router';

    const router = useRouter();

    async function verifyLogin() {
        const response = await getProfileData();

        if (response.status) {
            fetchData();
        } else {
            router.push({ name: "admin" });
        }
    }

    const socials = ref([]);
    const search = ref("");

    async function fetchData() {
        const response = await fetchSocials();
        if (response.access_denied) {
            window.alert(response.access_denied);
        } else {
            socials.value = response;
        }
    }

    const filteredSocials = computed(() => {
        if (!search.value) {
            return socials.value;
        }

        return socials.value.filter(social =>
            social.name.toLowerCase().includes(search.value.toLowerCase())
        );
    });

    const currentlyEditing = ref(-1);
    const currentlyEditingData = ref({});

    const fileInput = ref(null);
    const filePreviewURL = ref(null);

    function previewFile() {
        if (fileInput.value.files.length > 0) {
            const reader = new FileReader();
            reader.onload = (e) => {
                filePreviewURL.value = e.target.result;
            };
            reader.readAsDataURL(fileInput.value.files[0]);
        }
    }

    function editSocial(social) {
        currentlyEditing.value = social.id;
        currentlyEditingData.value = structuredClone(toRaw(social));
        delete currentlyEditingData.value.id;
        filePreviewURL.value = `${config.apiUrl}images/${currentlyEditingData.value.imageUrl}`;
    }

    function createSocial() {
        currentlyEditingData.value = {
            name: "",
            url: ""
        }

        currentlyEditing.value = -2;
    }

    function stopEditing() {
        filePreviewURL.value = "";
        currentlyEditing.value = -1;
        currentlyEditingData.value = {};
    }

    const saving = ref(false);

    async function save() {
        if (!verifyInputs()) {
            return;
        }
        saving.value = true;
     
        const formData = setFormData();

        const response = await postSocial(formData);
        handleResponse(response);
    }

    async function update() {
        if (!verifyInputs(false)) {
            return;
        }
        saving.value = true;
     
        const formData = setFormData();

        const response = await updateSocial(formData);
        handleResponse(response);
    }

    async function remove() {
        if (window.confirm(`Zeker dat je social: ${currentlyEditingData.value.name} wilt verwijderen?`)) {
            saving.value = true;
            
            const response = await deleteSocial(currentlyEditing.value);

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

    const validateImage = ref(true);

    function verifyInputs(imageValidation = true) {
        const form = document.querySelector("#socialEdit");

        for (let element of form.elements) {
            if (!element.checkValidity()) {
                element.reportValidity();
                return false;
            }
        }

        if (imageValidation && fileInput.value.files[0].size / 1024 > 4096) {
            window.alert("Afbeelding te groot: " + fileInput.value.files[0].size / (1024 * 1024) + "MB, max: 4MB");
            return false;
        }

        return true;
    }

    function setFormData() {
        const formData = new FormData();

        for (let property in currentlyEditingData.value) {
            if (currentlyEditingData.value[property]) {
                formData.append(property, currentlyEditingData.value[property]);
            }
        }

        if (fileInput.value.files[0]) {
            formData.append("image", fileInput.value.files[0]);
        }

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
            <button v-show="currentlyEditing === -1" type="button" @click="createSocial">Nieuwe Social</button>
        </div>
        <form v-show="socials.length > 0 && socials[0] !== 'loading' && currentlyEditing === -1" method="post" @submit.prevent="" id="search">
            <div>
                <input v-model="search" type="search" id="search" name="search" autocomplete="search" required placeholder="Zoeken">
            </div>
        </form>
    </header>
    <main>
        <table v-show="currentlyEditing === -1">
            <thead>
                <th>ID</th>
                <th>Name</th>
                <th></th>
            </thead>
            <tbody>
                <tr v-for="social in filteredSocials">
                    <td>{{ social.id }}</td>
                    <td>{{ social.name }}</td>
                    <td><div @click="() => editSocial(social)" class="icon-edit"></div></td>
                </tr>
            </tbody>
        </table>
        <form @submit.prevent="" method="post" id="socialEdit" v-show="currentlyEditing !== -1">
            <section class="general">
                <label for="name">Naam<span>*</span></label>
                <input v-model="currentlyEditingData.name" type="text" max="70" id="name" name="name" required>
                <label for="url">Link<span>*</span><br><i>Moet beginnen met https://</i></label>
                <input v-model="currentlyEditingData.url" type="text" id="url" name="url" pattern="https://.+" required>
            </section>
            <section class="image-upload">
                <label for="imgUpload">Social Media Icon</label>
                <input type="file" ref="fileInput" accept="image/*" @change="previewFile" id="imgUpload" name="imgUpload" :required="validateImage">
                <img :src="filePreviewURL">
            </section>
            <section class="buttons">
                    <button class="btn-red" @click="remove()" type="button" v-if="currentlyEditing !== -2">Verwijderen</button>
                    <button @click="() => (currentlyEditing === -2) ? save() : update()" type="button">Publiceren</button>
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
        background-color: var(--color-primary);
        border: none;
        width: 9rem;
        height: 2rem;
        font-family: inherit;
        border-radius: 0.4rem;
        box-shadow: var(--color-accent-dark) 0.5rem 0.5rem;
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
        background-color: var(--color-background-2);
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
        background-color: var(--color-background-2);
    }

    th:nth-child(2), td:nth-child(2) {
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

    th:nth-child(3) {
        width: 3rem;
    }

    td.ONLINE, td.true {
        background-color: rgb(125, 209, 0);
    }

    td.DRAFT {
        background-color: rgb(255, 136, 0);
    }

    td.ARCHIVED, td.false {
        background-color: red;
        color: var(--color-background-2);
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

    #socialEdit {
        display: grid;
        grid-template-columns: 1fr 1fr;
        padding: 1rem;
    }

    #socialEdit > div {
        grid-column: span 2;
    }

    #socialEdit section.general, #socialEdit section.image-upload {
        display: flex;
        height: max-content;
        flex-direction: column;
    }

    #socialEdit section.image-upload {
        align-items: center;
    }

    #socialEdit label {
        margin-bottom: 0.3rem;
    }

    #socialEdit label:not(:first-child) {
        margin-top: 1.5rem;
    }

    #socialEdit label span {
        color: red;
    }

    #socialEdit label i {
        opacity: 60%;
    }

    #socialEdit input {
        width: 80%;
        font-size: 0.85rem;
    }

    #socialEdit section.general input:nth-of-type(4), #socialEdit section.general input:nth-of-type(5) {
        text-align: center;
        width: 12rem;
    }

    #socialEdit select {
        width: 12rem;
        height: 2rem;
        font: inherit;
        font-size: 0.9rem;
        border-radius: 0.3rem;
        border: solid 0.1rem var(--color-text);
        padding: 0.2rem 0;
        outline: none;
        background-color: var(--color-background-2);
        color: var(--color-text);
    }

    #socialEdit section.image-upload input[type="file"] {
        width: 18rem;
        padding-left: 0.5rem;
    }

    #socialEdit section.image-upload img {
        margin-top: 3rem;
        height: 10rem;
        max-width: 80%;
        object-fit: contain;
        object-position: 0;
    }

    #socialEdit section.buttons {
        margin: 3rem 1rem 1rem 1rem;
        display: flex;
        gap: 1rem;
        justify-content: center;
        grid-column: span 2;
    }

    #socialEdit button.btn-red {
        background-color: rgb(255, 66, 66);
        box-shadow: #e0b8aa 0.5rem 0.5rem;
    }

    #socialEdit button.btn-orange {
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

        #socialEdit button.btn-red {
            background-color: rgb(255, 87, 57);
            color: var(--color-text3);
            box-shadow: #360f01 0.5rem 0.5rem;
        }

        #socialEdit button.btn-orange {
            background-color: rgb(233, 177, 22);
            color: var(--color-text3);
            box-shadow: #362a01 0.5rem 0.5rem;
        }
    }
</style>