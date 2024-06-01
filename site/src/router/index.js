import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import { LanguageProvider } from '@/services/LanguageService.js'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: () => {
        return { path: `/${LanguageProvider.CURR_LANG.toLocaleLowerCase()}` }
      }
    },
    {
      path: '/:lang([A-Za-z]{2})',
      redirect: to => {
        return `${to.params.lang}/${LanguageProvider.DICTIONARY.navHome[LanguageProvider.CURR_LANG]}`
      }
    },
    {
      path: '/en/home',
      name: "home",
      props: true,
      component: HomeView,
      alias: [
        '/nl/startpagina',
        '/fr/accueil',
        '/de/startseite'
      ]
    }
  ]
})

export default router
