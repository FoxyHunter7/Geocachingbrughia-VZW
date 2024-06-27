<script setup>
    import { computed, ref } from 'vue';
    import LanguageProvider from '@/services/LanguageService';
    import StaticContentProvider from '@/services/StaticContentService';
    import { postContact } from '@/services/ContactService';
    import config from '../data/config.json';

    const lang = computed(() => LanguageProvider.CURR_LANG.value);
    const dictionary = StaticContentProvider.DICTIONARY;

    const email = ref("");
    const subject = ref("");
    const message = ref("");
    function handleFormSubmit() {
        postContact(email.value, subject.value, message.value);
    }
</script>

<template>
    <section>
        <h2>{{ dictionary.ContactHelpQuestionTxt[lang] }}</h2>
        <p>{{ dictionary.ContactHelpTxt[lang] }}</p>
        <div>
            <form method="post" @submit.prevent="handleFormSubmit">
                <div>
                    <label for="email">{{ dictionary.FormMail[lang] }}</label>
                    <input v-model="email" type="email" id="email" name="email" autocomplete="email" required>
                </div>
                <div>
                    <label for="subject">{{ dictionary.FormSubject[lang] }}</label>
                    <input v-model="subject" type="text" id="subject" name="subject" autocomplete="off" required>
                </div>
                <div>
                    <label for="message">{{ dictionary.FormMessage[lang] }}</label>
                    <textarea v-model="message" id="message" name="message" autocomplete="off"></textarea>
                </div>
                <input type="submit" value="submit">
            </form>
            <ul>
                <li>
                    <p>{{ dictionary.ContactPostalMailTxt[lang] }}</p>
                    <a target="_blank" href="https://www.google.com/maps/place/Korte+Kwadeplasstraat+6,+8020+Oostkamp/@51.1593091,3.2333893,17z/data=!3m1!4b1!4m6!3m5!1s0x47c3506bbbd8d363:0xbed866c1f1258d86!8m2!3d51.1593091!4d3.2359696!16s%2Fg%2F11rcy98rq8?entry=ttu">Korte Kwadeplasstraat 6, 8020 Oostkamp</a>
                </li>
                <li>
                    <p>{{ dictionary.ContactCallTxt[lang] }}</p>
                    <a href="tel:+32 50 841 331">+32 50 841 331</a>
                    <a href="tel:+32 487 906 431">+32 487 906 431</a>
                </li>
                <li>
                    <p>{{ dictionary.ContactMailTxt[lang] }}</p>
                    <a href="mailto:info@geocachingbrughia.be">info@geocachingbrughia.be</a>
                </li>
            </ul>
        </div>
    </section>
</template>

<style scoped>
    section {
        padding: 3rem;
        max-height: 40rem;
    }

    h2 {
        text-align: center;
        text-transform: uppercase;
        font-weight: bold;
    }

    section > p {
        text-align: center;
        line-height: 200%;
        font-size: 1rem;
        margin-bottom: 3rem;
    }

    section > div {
        display: flex;
        justify-content: space-evenly;
        flex-wrap: wrap;
        gap: 5rem;
        margin: auto;
    }

    form {
        display: flex;
        flex-direction: column;
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
        width: 30rem;
        max-width: 90vw;
    }

    form input, form textarea {
        height: 2rem;
        border-radius: 0.3rem;
        border: solid 0.1rem var(--color-text);
        outline: none;
        font-family: inherit;
        font-size: 0.85rem;
        box-sizing: border-box;
    }

    form textarea {
        height: 8rem;
        resize: vertical;
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
    }

    form input[type="submit"]:hover {
        cursor: pointer;
        scale: 103%;
        transition: scale 0.25s;
    }

    ul {
        width: 30rem;
        display: flex;
        flex-direction: column;
        align-items: center;
        padding-inline-start: none;
        list-style: none;
        text-transform: capitalize;
    }

    li {
        display: block;
        width: 20rem;
        margin-bottom: 1.5rem;
    }

    li a {
        display: block;
        margin-top: 0.3rem;
        text-decoration: none;
        color: var(--color-quaternary);
    }
</style>