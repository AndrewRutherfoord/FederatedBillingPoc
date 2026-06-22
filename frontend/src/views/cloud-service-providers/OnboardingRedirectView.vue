<script setup lang="ts">
import { useAsyncState } from '@vueuse/core';
import { useRoute, useRouter } from 'vue-router';
import client from '@/api/client';

const route = useRoute();
const router = useRouter();

const accountId = route.query.account_id as string | undefined;
const billingAccountId = route.query.billing_account_id as string | undefined;
const cspProviderId = route.query.csp_provider_id as string | undefined;

const { isLoading, isReady, error } = useAsyncState(async () => {
    if (!billingAccountId || !cspProviderId) return null;

    const { error } = await (client as any).POST(
        `/billing/accounts/${billingAccountId}/cloud-provider-accounts/complete`,
        { body: { csp_provider_id: cspProviderId } }
    );

    if (error) throw new Error('Failed to complete onboarding');
    return true;
}, null);
</script>

<template>
    <div class="flex items-center justify-center min-h-[60vh]">
        <div class="text-center max-w-md">
            <template v-if="!billingAccountId || !cspProviderId">
                <p class="text-red-500">Missing onboarding parameters.</p>
            </template>
            <template v-else-if="isLoading">
                <p class="text-gray-500">Completing setup…</p>
            </template>
            <template v-else-if="error">
                <p class="text-red-500">Failed to complete setup. Please try again.</p>
            </template>
            <template v-else-if="isReady">
                <div class="text-green-600 text-5xl mb-4">✓</div>
                <h1 class="text-2xl font-semibold mb-2">Cloud provider linked</h1>
                <p class="text-gray-600 mb-1">
                    Your cloud provider account has been linked to your billing account.
                </p>
                <p v-if="accountId" class="text-sm text-gray-400 mb-6">Account ID: {{ accountId }}</p>
                <button
                    class="px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700 text-sm"
                    @click="router.push({ name: 'Home', query: { account: billingAccountId } })"
                >
                    Go to dashboard
                </button>
            </template>
        </div>
    </div>
</template>
