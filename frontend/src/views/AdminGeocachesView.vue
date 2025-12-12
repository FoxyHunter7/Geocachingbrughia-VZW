<script setup>
    import { deleteGeocache, fetchGeocaches, getProfileData, postGeocache, updateGeocache } from '@/services/AdminService';
    import { onMounted, ref, computed, watch, toRaw } from 'vue';
    import { useRouter } from 'vue-router';

    const router = useRouter();

    async function verifyLogin() {
        const response = await getProfileData(currPage.value);

        if (response.status) {
            fetchData();
        } else {
            router.push({ name: "admin" });
        }
    }

    const geocaches = ref([]);
    const currPage = ref(1);
    const lastPage = ref(1);
    const search = ref("");

    async function fetchData() {
        const pagedResponse = await fetchGeocaches();

        if (pagedResponse.data) {
            geocaches.value = pagedResponse.data;
            currPage.value = pagedResponse.current_page;
            lastPage.value = pagedResponse.last_page;
        } else if (pagedResponse.access_denied) {
            window.alert(response.access_denied);
            router.push({ name: "admin" });
        } else {
            router.push({ name: "admin" });
        }
    }

    const filteredGeocaches = computed(() => {
        if (!search.value) {
            return geocaches.value;
        }

        return geocaches.value.filter(geocaches =>
            geocaches.title.toLowerCase().includes(search.value.toLowerCase())
        );
    });

    function nextPage() {
        if (currPage.value < lastPage.value) {
            currPage.value++;
        }
    }

    function prevPage() {
        if (currPage > 1) {
            currPage.value--;
        }
    }

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
        });
    }

    const currentlyEditing = ref(-1);
    const currentlyEditingData = ref({});

    function editGeocache(geocache) {
        currentlyEditing.value = geocache.id;
        currentlyEditingData.value = structuredClone(toRaw(geocache));
        delete currentlyEditingData.value.id;
    }

    function createGeocache() {
        currentlyEditingData.value = {
            state: "DRAFT",
            title: "",
            geolink: "",
            type: "TRADITIONAL",
            difficulty: null,
            terrain: null,
            placed_on: ""
        }

        currentlyEditing.value = -2;
    }

    function stopEditing() {
        currentlyEditing.value = -1;
        currentlyEditingData.value = {};
    }

    const saving = ref (false);

    async function save(state) {
        if (!verifyInputs()) {
            return;
        }
        saving.value = true;

        currentlyEditingData.value.state = state;
        const formData = setFormData();

        const response = await postGeocache(formData);
        handleResponse(response);
    }

    async function update(state) {
        if (!verifyInputs()) {
            return;
        }
        saving.value = true;

        currentlyEditingData.value.state = state;
        const formData = setFormData();

        const response = await updateGeocache(currentlyEditing.value, formData);
        handleResponse(response);
    }

    async function remove() {
        if (window.confirm(`Zeker dat je geocache: ${currentlyEditingData.title} wilt verwijderen?`)) {
            saving.value = true;

            const response = await deleteGeocache(currentlyEditing.value);

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
        const form = document.querySelector("#geocacheEdit");

        for (let element of form.elements) {
            if (!element.checkValidity()) {
                element.reportValidity();
                return false;
            }
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
    watch(currPage, async () => await fetchData());
</script>

<template>
    <header>
        <div>
            <button @click="() => (currentlyEditing === -1) ? router.push({ name: 'admin'}) : stopEditing()" type="button">Terug</button>
            <button v-show="currentlyEditing === -1" type="button" @click="createGeocache">Nieuwe Geocache</button>
        </div>
        <form v-show="geocaches.length > 0 && geocaches[0] !== 'loading' && currentlyEditing === -1" method="post" @submit.prevent="" id="search">
            <div>
                <input v-model="search" type="search" id="search" name="search" autocomplete="search" required placeholder="Zoeken">
            </div>
        </form>
    </header>
    <main>
        <table v-show="currentlyEditing === -1">
            <thead>
                <th>ID</th>
                <th>Type</th>
                <th>Status</th>
                <th>Titel</th>
                <th>Moeilijkheid</th>
                <th>Terrein</th>
                <th>Geplaatst Op</th>
                <th></th>
            </thead>
            <tbody>
                <tr v-for="geocache in filteredGeocaches">
                    <td>{{ geocache.id }}</td>
                    <td><div><img :src="`/assets/media/cachetypes/${geocache.type}.png`" :alt="geocache.type" :title="geocache.type"></div></td>
                    <td :class="geocache.state">{{ geocache.state }}</td>
                    <td>{{ geocache.title }}</td>
                    <td>{{ geocache.difficulty }}</td>
                    <td>{{ geocache.terrain }}</td>
                    <td>{{ dateTimeFormatter(geocache.placed_on) }}</td>
                    <td><div @click="() => editGeocache(geocache)" class="icon-edit"></div></td>
                </tr>
            </tbody>
        </table>
        <div id="pager" v-show="geocaches.length > 0 && geocaches[0] !== 'loading' && currentlyEditing === -1">
            <div class="prev pagerNavBtn" :class="{ disabled: currPage === 1 }" @click="prevPage"></div>
            <p>{{ currPage }} / {{ lastPage }}</p>
            <div class="next pagerNavBtn" :class="{ disabled: currPage === lastPage}" @click="nextPage"></div>
        </div>
        <form @submit.prevent="" method="post"id="geocacheEdit" v-show="currentlyEditing !== -1">
            <section class="general">
                <label for="title">Titel<span>*</span></label>
                <input v-model="currentlyEditingData.title" type="text" max="100" id="title" name="title" required>
                <label for="geolink">Geocaching Link<span>*</span><br><i>Moet beginnen met: https://www.geocaching.com/geocache/</i></label>
                <input v-model="currentlyEditingData.geolink" type="text" id="geolink" title="Moet beginnen met: https://www.geocaching.com/geocache/" name="geolink" pattern="https://www\.geocaching\.com/geocache/.+" required>
                <label for="type">Type<span>*</span></label>
                <select v-model="currentlyEditingData.type" id="type" name="type" required>
                    <option value="TRADITIONAL">Traditional</option>
                    <option value="MULTI">Multi</option>
                    <option value="MYSTERY">Mystery</option>
                    <option value="EARTH">Earth</option>
                    <option value="LETTERBOX">LetterBox</option>
                    <option value="WHEREIGO">WhereIGo</option>
                    <option value="VIRTUAL">Virtual</option>
                    <option value="LAB">Lab</option>
                    <option value="WEBCAM">Webcam</option>
                </select>
            </section>
            <section class="general">
                <label for="difficulty">Moeilijkheidsgraad<span>*</span></label>
                <input v-model="currentlyEditingData.difficulty" type="number" id="difficulty" name="difficulty" min="1" max="5" required>
                <label for="terrain">Terrein<span>*</span></label>
                <input v-model="currentlyEditingData.terrain" type="number" id="terrain" name="terrain" min="1" max="5" required>
                <label for="placedOn">Datum Geplaatst</label>
                <input v-model="currentlyEditingData.placed_on" type="date" id="placedOn" name="placedOn">
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
        background-color: var(--color-background-2);
    }

    th:nth-child(4), td:nth-child(4) {
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

    th:nth-child(3) {
        width: 6rem;
    }

    th:nth-child(5), th:nth-child(6) {
        width: 7rem;
    }

    th:nth-child(7) {
        width: 8.5rem;
    }

    th:nth-child(8) {
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

    #pager {
        display: flex;
        justify-content: center;
        align-items: center;
        gap: 2rem;
        margin-top: 1rem;
    }

    #pager p {
        text-align: center;
        font-size: 1rem;
        font-weight: bold;
        letter-spacing: 0.5rem;
    }

    #pager .pagerNavBtn {
        width: 2rem;
        height: 2rem;
        background-color: var(--color-text);
        cursor: pointer;
    }

    #pager .pagerNavBtn.next {
        mask: url(@/assets/media/chevron-right.svg);
        mask-size: contain;
        mask-repeat: no-repeat;
        mask-position: center;
    }

    #pager .pagerNavBtn.prev {
        mask: url(@/assets/media/chevron-left.svg);
        mask-size: contain;
        mask-repeat: no-repeat;
        mask-position: center;
    }

    #pager .pagerNavBtn.disabled {
        filter: opacity(30%);
        cursor: auto;
    }

    #geocacheEdit {
        display: grid;
        grid-template-columns: 1fr 1fr;
        padding: 1rem;
    }

    #geocacheEdit > div {
        grid-column: span 2;
    }

    #geocacheEdit section.general{
        display: flex;
        height: max-content;
        flex-direction: column;
    }

    #geocacheEdit label {
        margin-bottom: 0.3rem;
    }

    #geocacheEdit label:not(:first-child) {
        margin-top: 1.5rem;
    }

    #geocacheEdit label span {
        color: red;
    }

    #geocacheEdit label i {
        opacity: 60%;
    }

    #geocacheEdit input {
        width: 80%;
        font-size: 0.85rem;
    }

    #geocacheEdit section.general input:nth-of-type(4) {
        text-align: center;
        width: 12rem;
    }

    #geocacheEdit select {
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

    #geocacheEdit section.buttons {
        margin: 3rem 1rem 1rem 1rem;
        display: flex;
        gap: 1rem;
        justify-content: center;
        grid-column: span 2;
    }

    #geocacheEdit button.btn-red {
        background-color: rgb(255, 66, 66);
        box-shadow: #e0b8aa 0.5rem 0.5rem;
    }

    #geocacheEdit button.btn-orange {
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

        #geocacheEdit button.btn-red {
            background-color: rgb(255, 87, 57);
            color: var(--color-text3);
            box-shadow: #360f01 0.5rem 0.5rem;
        }

        #geocacheEdit button.btn-orange {
            background-color: rgb(233, 177, 22);
            color: var(--color-text3);
            box-shadow: #362a01 0.5rem 0.5rem;
        }
    }
</style>