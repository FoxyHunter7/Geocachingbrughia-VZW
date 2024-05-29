import { createRouter, createWebHistory } from 'vue-router'
import { defaultLanguage } from '../services/LanguageService'
import HomeView from '@/views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: () => {
        return { path: `/${defaultLanguage()}` }
      }
    },
    {
      path: '/:lang([A-Za-z]{2})',
      redirect: to => {
        switch(to.params.lang) {
          case 'en':
            return { path: '/en/home' }
          case 'nl':
            return { path: '/nl/startpagina' }
          case 'fr':
            return { path: '/fr/accueil' }
          case 'de':
            return { path: '/de/startseite' }
        }
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
