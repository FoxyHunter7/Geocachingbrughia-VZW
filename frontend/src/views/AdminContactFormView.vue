<script setup>
    import { fetchContactFormResponses, getProfileData } from '@/services/AdminService';
    import { onMounted, ref, computed, watch } from 'vue';
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

    const formResponses = ref([]);
    const currPage = ref(1);
    const lastPage = ref(1);
    const search = ref("");

    async function fetchData() {
        const pagedResponse = await fetchContactFormResponses(currPage.value);

        if (pagedResponse.data) {
            formResponses.value = pagedResponse.data;
            currPage.value = pagedResponse.current_page;
            lastPage.value = pagedResponse.last_page;
        } else if (pagedResponse.access_denied) {
            window.alert(response.access_denied);
            router.push({ name: "admin" });
        } else {
            router.push({ name: "admin" });
        }
    }

    const filteredFormResponses = computed(() => {
        if (!search.value) {
            return formResponses.value;
        }

        return formResponses.value.filter(formResponse =>
            formResponse.subject.toLowerCase().includes(search.value.toLowerCase())
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

    onMounted(verifyLogin);
    watch(currPage, async () => await fetchData());
</script>


<template>
    <header>
        <div>
            <button @click="router.push({ name: 'admin'})" type="button">Terug</button>
        </div>
        <form v-show="formResponses.length > 0 && formResponses[0] !== 'loading'" method="post" @submit.prevent="" id="search">
            <div>
                <input v-model="search" type="search" id="search" name="search" autocomplete="search" required placeholder="Zoeken">
            </div>
        </form>
    </header>
    <main>
        <table>
            <thead>
                <th>ID</th>
                <th>Datum Verzonden</th>
                <th>Onderwerp</th>
                <th>Bericht</th>
                <th>E-Mail Addres</th>
            </thead>
            <tbody>
                <tr v-for="formResponse in filteredFormResponses">
                    <td>{{ formResponse.id }}</td>
                    <td>{{ dateTimeFormatter(dateTimeFormatter(formResponse.created_at)) }}</td>
                    <td>{{ formResponse.subject }}</td>
                    <td>{{ formResponse.message }}</td>
                    <td>{{ formResponse.email }}</td>
                </tr>
            </tbody>
        </table>
        <div id="pager" v-show="formResponses.length > 0 && formResponses[0] !== 'loading'">
            <div class="prev pagerNavBtn" :class="{ disabled: currPage === 1 }" @click="prevPage"></div>
            <p>{{ currPage }} / {{ lastPage }}</p>
            <div class="next pagerNavBtn" :class="{ disabled: currPage === lastPage}" @click="nextPage"></div>
        </div>
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
        box-shadow: var(--color-background-2) 0.5rem 0.5rem;
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

    th:nth-child(3), th:nth-child(4), td:nth-child(3), td:nth-child(4) {
        text-align: start;
    }

    td {
        height: 1.8rem;
        padding: 0.5rem;
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
        width: 10rem;
    }

    th:nth-child(3) {
        width: 15rem
    }

    th:nth-child(4) {
        width: max-content;
    }

    th:nth-child(5) {
        width: 16rem;
    }

    #pager {
        display: flex;
        justify-content: center;
        align-items: center;
        gap: 2rem;
        margin: 1rem 0;
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
        background-color: var(--color-primary);
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
</style>