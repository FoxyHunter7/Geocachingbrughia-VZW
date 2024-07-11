<script setup>
    import { deleteEvent, fetchEvents, getProfileData, postEvent, updateEvent } from '@/services/AdminService';
    import { onMounted, ref, computed, watch, toRaw } from 'vue';
    import { useRouter } from 'vue-router';
    import config from '@/data/config.json';
    import TipTapEditor from '@/components/TipTapEditor.vue';
    import StaticContentProvider from '@/services/StaticContentService';

    const router = useRouter();

    async function verifyLogin() {
        const response = await getProfileData(currPage.value);

        if (response.status) {
            fetchData();
        } else {
            router.push({ name: "admin" });
        }
    }

    const events = ref([]);
    const currPage = ref(1);
    const lastPage = ref(1);
    const search = ref("");

    async function fetchData() {
        const pagedResponse = await fetchEvents();

        if (pagedResponse.data) {
            events.value = pagedResponse.data;
            currPage.value = pagedResponse.current_page;
            lastPage.value = pagedResponse.last_page;
        } else if (pagedResponse.access_denied) {
            window.alert(response.access_denied);
            router.push({ name: "admin" });
        } else {
            router.push({ name: "admin" });
        }
    }

    const filteredEvents = computed(() => {
        if (!search.value) {
            return events.value;
        }

        return events.value.filter(event =>
            event.title.toLowerCase().includes(search.value.toLowerCase())
        );
    });

    function nextPage() {
        if (currPage.value < lastPage.value) {
            currPage.value++;
        }
    }

    function prevPage() {
        if (currPage.value > 1) {
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
            hour: "2-digit",
            minute: "2-digit"
        });
    }

    const currentlyEditing = ref(-1);
    const currentlyEditingData = ref({});

    const fileInput = ref(null);
    const filePreviewURL = ref(null);
    const editors = ref([]);

    function previewFile() {
        if (fileInput.value.files.length > 0) {
            const reader = new FileReader();
            reader.onload = (e) => {
                filePreviewURL.value = e.target.result;
            };
            reader.readAsDataURL(fileInput.value.files[0]);
        }
    }

    function editEvent(event) {
        currentlyEditing.value = event.id;
        currentlyEditingData.value = structuredClone(toRaw(event));
        delete currentlyEditingData.value.id;
        filePreviewURL.value = `${config.apiUrl}images/${currentlyEditingData.value.imageUrl}`;
    }

    function createEvent() {
        const translations = [];

        StaticContentProvider.LANGUAGES.forEach(language => {
            translations.push({lang_code: language.code, description: '{"type":"doc","content":[]}'});
        });

        currentlyEditingData.value = {
            state: "DRAFT",
            on_home: false,
            title: "",
            geolink: "",
            type: "REGULAR",
            location: "",
            start_date: "",
            end_date: "",
            translations: translations
        }

        currentlyEditing.value = -2;
    }

    function stopEditing() {
        editors.value = [];
        filePreviewURL.value = "";
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

        const response = await postEvent(formData);
        handleResponse(response);
    }

    async function update(state) {
        if (!verifyInputs(false)) {
            return;
        }
        saving.value = true;

        currentlyEditingData.value.state = state;
        const formData = setFormData();

        const response = await updateEvent(currentlyEditing.value, formData);
        handleResponse(response);
    }

    async function remove() {
        if (window.confirm(`Zeker dat je event: ${currentlyEditingData.value.title} wilt verwijderen?`)) {
            saving.value = true;
            
            const response = await deleteEvent(currentlyEditing.value);

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
        const form = document.querySelector("#eventEdit");

        validateImage.value = imageValidation;
        if (!form.checkValidity()) {
            form.reportValidity();
            return false;
        }

        if (StaticContentProvider.LANGUAGES.length !== currentlyEditingData.value.translations.length) {
            const translations = [];
            StaticContentProvider.LANGUAGES.forEach(language => {
                translations.push({lang_code: language.lang_code, description: '{"type":"doc","content":[]}'});
            });
            currentlyEditingData.value.translations = translations;

            window.alert("Ontbrekende vertaling, talen zonder vertalingen gevonden, gelieve dit aan te vullen.");
            return false;
        }

        if (imageValidation && fileInput.value.files[0].size / 1024 > 4096) {
            window.alert("Afbeelding te groot: " + fileInput.value.files[0].size / (1024 * 1024) + "MB, max: 4MB");
            return false;
        }

        return true;
    }

    function setFormData() {
        const translations = [];

        editors.value.forEach(editor => {
            translations.push(editor.getContent());
        });

        const [startDate, endDate] = [currentlyEditingData.value.start_date, currentlyEditingData.value.end_date];
        const formattedStartDate = (startDate.split("T").length > 1) ? `${startDate.split("T")[0]} ${startDate.split("T")[1]}:00` : startDate;
        const formattedEndDate = (endDate.split("T").length > 1) ? `${endDate.split("T")[0]} ${endDate.split("T")[1]}:00` : endDate;

        const formData = new FormData();

        formData.append("translations", JSON.stringify(translations));
        formData.append("start_date", formattedStartDate);
        formData.append("end_date", formattedEndDate);

        for (let property in currentlyEditingData.value) {
            if (property !== "translations" && property !== "start_date" && property !== "end_date") {
                if (currentlyEditingData.value[property] || property === "on_home") {
                    formData.append(property, currentlyEditingData.value[property]);
                }
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
    watch(currPage, async () => await fetchData());
</script>

<template>
    <header>
        <div>
            <button @click="() => (currentlyEditing === -1) ? router.push({ name: 'admin'}) : stopEditing()" type="button">Terug</button>
            <button v-show="currentlyEditing === -1" type="button" @click="createEvent">Nieuw Event</button>
        </div>
        <form v-show="events.length > 0 && events[0] !== 'loading' && currentlyEditing === -1" method="post" @submit.prevent="" id="search">
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
                <th>Op Start</th>
                <th>Titel</th>
                <th>Start</th>
                <th>Eind</th>
                <th></th>
            </thead>
            <tbody>
                <tr v-for="event in filteredEvents">
                    <td>{{ event.id }}</td>
                    <td><div><img :src="`/assets/media/eventtypes/${event.type}.png`" :alt="event.type" :title="event.type"></div></td>
                    <td :class="event.state">{{ event.state }}</td>
                    <td :class="event.on_home">{{ event.on_home }}</td>
                    <td>{{ event.title }}</td>
                    <td>{{ dateTimeFormatter(event.start_date) }}</td>
                    <td>{{ dateTimeFormatter(event.end_date) }}</td>
                    <td><div @click="() => editEvent(event)" class="icon-edit"></div></td>
                </tr>
            </tbody>
        </table>
        <div id="pager" v-show="events.length > 0 && events[0] !== 'loading' && currentlyEditing === -1">
            <div class="prev pagerNavBtn" :class="{ disabled: currPage === 1 }" @click="prevPage"></div>
            <p>{{ currPage }} / {{ lastPage }}</p>
            <div class="next pagerNavBtn" :class="{ disabled: currPage === lastPage}" @click="nextPage"></div>
        </div>
        <form @submit.prevent="" method="post" id="eventEdit" v-show="currentlyEditing !== -1">
            <section class="general">
                <label for="onHome">Toon event op de homepagina<span>*</span></label>
                <select v-model="currentlyEditingData.on_home" id="onHome" name="onHome">
                    <option value="false">Neen</option>
                    <option value="true">Ja</option>
                </select>
                <label for="type">Type<span>*</span></label>
                <select v-model="currentlyEditingData.type" id="type" name="type" required>
                    <option value="REGULAR">Regular</option>
                    <option value="CITO">CITO</option>
                    <option value="MEGA">Mega</option>
                    <option value="GIGA">Giga</option>
                    <option value="BLOCK">Block</option>
                </select>
                <label for="title">Titel<span>*</span></label>
                <input v-model="currentlyEditingData.title" type="text" max="100" id="title" name="title" required>
                <label for="geolink">Geocaching Link<span>*</span><br><i>Moet beginnen met: https://www.geocaching.com/geocache/</i></label>
                <input v-model="currentlyEditingData.geolink" type="text" id="geolink" title="Moet beginnen met: https://www.geocaching.com/geocache/" name="geolink" pattern="https://www\.geocaching\.com/geocache/.+" required>
                <label for="location">Locatie<br><i>In GMD notatie, bv: N 34° 56.789 E 123° 45.678</i></label>
                <input v-model="currentlyEditingData.location" type="text" id="location" name="location" title="In GMD notatie, bv: N 34° 56.789 E 123° 45.678" pattern="^[NS]\s\d+°\s\d+\.\d+\s[EW]\s\d+°\s\d+\.\d+$">
                <label for="startDate">Start Datum & Tijd<span>*</span></label>
                <input v-model="currentlyEditingData.start_date" type="datetime-local" id="startDate" name="startDate" required>
                <label for="endDate">Eind Datum & Tijd<span>*</span></label>
                <input v-model="currentlyEditingData.end_date" type="datetime-local" id="endDate" name="endDate" required>
            </section>
            <section class="image-upload">
                <label for="imgUpload">Poster Foto</label>
                <input type="file" ref="fileInput" accept="image/*" @change="previewFile" id="imgUpload" name="imgUpload" :required="validateImage">
                <img :src="filePreviewURL">
            </section>
            <div>
                <section class="translations">
                    <div v-for="translation in currentlyEditingData.translations">
                        <p>{{ translation.lang_code }}</p>
                        <TipTapEditor :content="translation.description" :langCode="translation.lang_code" :editable="true" ref="editors"/>
                    </div>
                </section>
                <section class="buttons">
                    <button class="btn-red" @click="remove()" type="button" v-if="currentlyEditing !== -2">Verwijderen</button>
                    <button class="btn-orange" @click="update('ARCHIVED')" type="button" v-if="currentlyEditingData.state === 'ONLINE'">Archiveren</button>
                    <button class="btn-orange" @click="() => (currentlyEditing === -2) ? save('DRAFT') : update('DRAFT')" type="button" v-if="currentlyEditingData.state !== 'ONLINE'">Opslaan Als Concept</button>
                    <button @click="() => (currentlyEditing === -2) ? save('ONLINE') : update('ONLINE')" type="button">Publiceren</button>
                </section>
            </div>
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
        mask: url(../assets/media/search.svg);
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

    th:nth-child(5), td:nth-child(5) {
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

    th:nth-child(3), th:nth-child(4) {
        width: 6rem;
    }

    td:nth-child(4) {
        text-transform: uppercase;
    }

    th:nth-child(6), th:nth-child(7) {
        width: 12rem;
    }

    th:nth-child(8) {
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
        mask: url(../assets/media/edit.svg);
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
        mask: url(../assets/media/chevron-right.svg);
        mask-size: contain;
        mask-repeat: no-repeat;
        mask-position: center;
    }

    #pager .pagerNavBtn.prev {
        mask: url(../assets/media/chevron-left.svg);
        mask-size: contain;
        mask-repeat: no-repeat;
        mask-position: center;
    }

    #pager .pagerNavBtn.disabled {
        filter: opacity(30%);
        cursor: auto;
    }

    #eventEdit {
        display: grid;
        grid-template-columns: 1fr 1fr;
        padding: 1rem;
    }

    #eventEdit > div {
        grid-column: span 2;
    }

    #eventEdit section.general, #eventEdit section.image-upload {
        display: flex;
        height: max-content;
        flex-direction: column;
    }

    #eventEdit section.image-upload {
        align-items: center;
    }

    #eventEdit label {
        margin-bottom: 0.3rem;
    }

    #eventEdit label:not(:first-child) {
        margin-top: 1.5rem;
    }

    #eventEdit label span {
        color: red;
    }

    #eventEdit label i {
        opacity: 60%;
    }

    #eventEdit input {
        width: 80%;
        font-size: 0.85rem;
    }

    #eventEdit section.general input:nth-of-type(4), #eventEdit section.general input:nth-of-type(5) {
        text-align: center;
        width: 12rem;
    }

    #eventEdit select {
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

    #eventEdit section.image-upload input[type="file"] {
        width: 18rem;
        padding-left: 0.5rem;
    }

    #eventEdit section.image-upload img {
        margin-top: 3rem;
        height: 30rem;
        max-width: 80%;
        object-fit: contain;
        object-position: 0;
    }

    #eventEdit section.translations {
        margin-top: 3rem;
        display: flex;
        justify-content: center;
        gap: 3rem;
        flex-wrap: wrap;
    }

    #eventEdit section.translations div {
        width: 40rem;
        border-radius: 0.3rem;
    }

    #eventEdit section.buttons {
        margin: 1rem;
        display: flex;
        gap: 1rem;
        justify-content: center;
    }

    #eventEdit button.btn-red {
        background-color: rgb(255, 66, 66);
        box-shadow: #e0b8aa 0.5rem 0.5rem;
    }

    #eventEdit button.btn-orange {
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

        #eventEdit button.btn-red {
            background-color: rgb(255, 87, 57);
            color: var(--color-text3);
            box-shadow: #360f01 0.5rem 0.5rem;
        }

        #eventEdit button.btn-orange {
            background-color: rgb(233, 177, 22);
            color: var(--color-text3);
            box-shadow: #362a01 0.5rem 0.5rem;
        }
    }
</style>