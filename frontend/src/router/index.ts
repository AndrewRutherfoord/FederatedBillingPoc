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
      path: '/billing-providers/manage',
      name: 'ManageAccounts',
      component: () => import('../views/billing-providers/ManageAccountsView.vue'),
    },
    {
      path: '/billing-providers/:id/charge-batches',
      name: 'ChargeBatches',
      component: () => import('../views/billing-providers/ChargeBatchesView.vue'),
      props: true,
    },
    {
      path: '/billing-providers/:id/resource-charges',
      name: 'ResourceCharges',
      component: () => import('../views/billing-providers/ResourceChargesView.vue'),
      props: true,
    },
    {
      path: '/billing-providers/:id/invoices',
      name: 'Invoices',
      component: () => import('../views/billing-providers/InvoicesView.vue'),
      props: true,
    },
    {
      path: '/billing-providers/:id/cloud-provider-links',
      name: 'CloudProviderLinks',
      component: () => import('../views/billing-providers/CloudProviderLinksView.vue'),
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
