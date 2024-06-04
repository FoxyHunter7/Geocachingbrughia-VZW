import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import { LanguageProvider } from '@/services/LanguageService.js'
import NotFoundView from '@/views/NotFoundView.vue';

await initLanguageProvider();

async function initLanguageProvider() {
  const languageProvider = new LanguageProvider()
  await languageProvider.init()
}

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
        return `/${to.params.lang}/${LanguageProvider.DICTIONARY.NavHome[to.params.lang.toLocaleUpperCase()]}`
      }
    },
    {
      path: LanguageProvider.ROUTES.navHome.path,
      name: "home",
      props: true,
      component: HomeView,
      alias: LanguageProvider.ROUTES.navHome.aliases
    },
    {
      path: LanguageProvider.ROUTES.navEvents.path,
      name: "events",
      props: true,
      component: HomeView,
      alias: LanguageProvider.ROUTES.navEvents.aliases
    },
    {
      path: LanguageProvider.ROUTES.navGeocaches.path,
      name: "geocaches",
      props: true,
      component: HomeView,
      alias: LanguageProvider.ROUTES.navGeocaches.aliases
    },
    {
      path: LanguageProvider.ROUTES.navShop.path,
      name: "shop",
      props: true,
      component: HomeView,
      alias: LanguageProvider.ROUTES.navShop.aliases
    },
    {
      path: '/:pathMatch(.*)',
      name: "NotFound",
      component: NotFoundView,
    }
  ]
})

export default router
