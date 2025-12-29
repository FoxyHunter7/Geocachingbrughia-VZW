<script setup>
  import LanguageProvider from '@/services/LanguageService';
  import StaticContentProvider from '@/services/StaticContentService';
  import { computed, onMounted, ref } from 'vue';
  import { useRouter } from 'vue-router';

  const router = useRouter();

  const lang = computed(() => LanguageProvider.CURR_LANG.value);
  const dictionary = StaticContentProvider.DICTIONARY;

  const redirectTimer = ref(10);

  async function countDownToRedirect() {
    while (redirectTimer.value > 0) {
      await new Promise(r => setTimeout(r, 1000));
      redirectTimer.value--;
    }

    router.push('/');
  }

  onMounted(countDownToRedirect);
</script>

<template>
  <main>
    <h1>{{ dictionary.UIPageNotFound[lang] }}</h1>
    <p>{{ dictionary.UIPageNotFoundSubTxt[lang].split("///t///")[0] }} {{ redirectTimer }} {{ dictionary.UIPageNotFoundSubTxt[lang].split("///t///")[1] }}</p>
  </main>
</template>

<style scoped>
  h1 {
    font-size: 1.2rem;
    font-weight: bold;
    margin-top: 3rem;
  }

  h1, p {
    text-align: center;
  }
</style>
