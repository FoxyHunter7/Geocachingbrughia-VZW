<script setup>
import uiTranslations from '../data/uiTranslations.json'
import routingTranslations from '../data/routingTranslations.json'

defineProps({
    lang: String,
    pathName: String
})

function goToLang(lang, path) {
    return window.location.href = `/${lang}${routingTranslations[path][lang]}`
}

</script>

<template>
    <p>{{ uiTranslations[lang].languageSelector }}</p>
    <li>
        <ul><button @click="goToLang('en', pathName)">-- EN --</button></ul>
        <ul><button @click="goToLang('nl', pathName)">-- NL --</button></ul>
        <ul><button @click="goToLang('fr', pathName)">-- FR --</button></ul>
        <ul><button @click="goToLang('de', pathName)">-- DE --</button></ul>
    </li>
</template>