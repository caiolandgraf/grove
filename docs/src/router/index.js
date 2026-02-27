import { createRouter, createWebHashHistory } from 'vue-router'
import HomePage from '@/pages/HomePage.vue'
import DocsPage from '@/pages/DocsPage.vue'
import ContributorsPage from '@/pages/ContributorsPage.vue'

const routes = [
  { path: '/', name: 'home', component: HomePage, meta: { title: 'Grove — Go Foundation CLI' } },
  { path: '/docs', name: 'docs', component: DocsPage, meta: { title: 'Docs — Grove' } },
  { path: '/contributors', name: 'contributors', component: ContributorsPage, meta: { title: 'Contributors — Grove' } },
  { path: '/:pathMatch(.*)*', redirect: '/' }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
  scrollBehavior(to, _from, savedPosition) {
    if (savedPosition) return savedPosition
    if (to.hash) return { el: to.hash, top: 80, behavior: 'smooth' }
    return { top: 0, behavior: 'smooth' }
  }
})

router.afterEach(to => { document.title = to.meta?.title ?? 'Grove' })

export default router
