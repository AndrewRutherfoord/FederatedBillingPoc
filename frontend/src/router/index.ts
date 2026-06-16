import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'Home',
      component: HomeView,
    },
    {
      path: '/billing-providers/onboarding-complete-callback',
      name: 'BillingOnboardingCompleteCallback',
      component: () => import('../views/billing-providers/OnboardingRedirectView.vue'),
    },
    {
      path: '/billing-providers/:id/account-overview',
      name: 'BillingAccountOverview',
      component: () => import('../views/billing-providers/AccountOverviewView.vue'),
      props: true,
    },
    {
      path: '/cloud-service-providers/onboarding-complete-callback',
      name: 'CspOnboardingCompleteCallback',
      component: () => import('../views/cloud-service-providers/OnboardingRedirectView.vue'),
    },
    {
      path: '/about',
      name: 'About',
      component: () => import('../views/AboutView.vue'),
    },
  ],
})

export default router
