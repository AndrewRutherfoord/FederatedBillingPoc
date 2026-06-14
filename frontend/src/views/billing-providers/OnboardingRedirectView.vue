<script setup lang="ts">
import { useAsyncState } from '@vueuse/core';
import { useRoute, useRouter } from 'vue-router';
import client from '@/api/client';
import type { components as ApiComponents } from '@/api/api-schema';

type BillingAccount = ApiComponents["schemas"]["handlers.billingAccountResponse"];

const route = useRoute();
const router = useRouter();

const accountId = route.query.account_id as string | undefined;

const { state, isLoading, isReady } = useAsyncState(async () => {
    if (!accountId) return null;
    const { data } = await client.GET("/billing/accounts");
    const accounts = (data as unknown as BillingAccount[]) ?? [];
    return accounts.find(a => a.id === accountId) ?? null;
}, null);
</script>

<template>
    <div class="flex items-center justify-center min-h-[60vh]">
        <div class="text-center max-w-md">
            <template v-if="!accountId">
                <p class="text-red-500">No account ID provided.</p>
            </template>
            <template v-else-if="isLoading">
                <p class="text-gray-500">Confirming your account…</p>
            </template>
            <template v-else-if="isReady && state">
                <div class="text-green-600 text-5xl mb-4">✓</div>
                <h1 class="text-2xl font-semibold mb-2">Billing account set up</h1>
                <p class="text-gray-600 mb-1">
                    <span class="font-medium">{{ state.alias }}</span> is now linked to
                    <span class="font-medium">{{ state.billing_provider_name }}</span>.
                </p>
                <p class="text-sm text-gray-400 mb-6">Account ID: {{ accountId }}</p>
                <button
                    class="px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700 text-sm"
                    @click="router.push('/')"
                >
                    Go to dashboard
                </button>
            </template>
            <template v-else>
                <p class="text-gray-500">Account not found. <button class="text-indigo-600 underline" @click="router.push('/')">Return home</button></p>
            </template>
        </div>
    </div>
</template>
