<script setup lang="ts">
import { useAsyncState } from '@vueuse/core';
import { Tabs, TabList, Tab, TabPanels, TabPanel } from 'primevue';
import { computed, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import client from '@/api/client';
import AccountOverview from '@/components/AccountOverview.vue';

const route = useRoute();
const router = useRouter();

const { state: accounts } = useAsyncState(async () => {
    const { data } = await client.GET("/billing/accounts");
    return data ?? [];
}, []);

const firstAccountId = computed(() => accounts.value?.[0]?.id);

const pageTitle = computed(() => {
    const list = accounts.value ?? [];
    return list.length === 1 ? `${list[0]?.alias} - Overview` : 'Billing Accounts Overview';
});

const selectedAccountId = ref<string>('');

watch(accounts, (list) => {
    if (!list || list.length === 0) return;
    const queryId = route.query.account as string | undefined;
    selectedAccountId.value = list.find(a => a.id === queryId)?.id ?? firstAccountId.value ?? '';
}, { immediate: true });

watch(selectedAccountId, (id) => {
    if (id && route.query.account !== id) {
        router.replace({ query: { ...route.query, account: id } });
    }
});
</script>

<template>
    <div class="flex justify-between items-center my-2">
        <h2>{{ pageTitle }}</h2>
        <Button as="router-link" :to="{ name: 'ManageAccounts' }" text size="small">
            Manage Billing Accounts
        </Button>
    </div>

    <div v-if="!accounts || accounts.length === 0" class="text-center py-8">
        <p class="text-gray-500 mb-4">No billing accounts found.</p>
        <Button as="router-link" :to="{ name: 'ManageAccounts' }">Create a Billing Account</Button>
    </div>

    <AccountOverview v-else-if="accounts.length === 1 && firstAccountId" :id="firstAccountId" />

    <Tabs v-else v-model:value="selectedAccountId">
        <TabList>
            <Tab v-for="account in accounts" :key="account.id" :value="account.id!">{{ account.alias }}</Tab>
        </TabList>
        <TabPanels>
            <TabPanel v-for="account in accounts" :key="account.id" :value="account.id!">
                <AccountOverview :id="account.id!" />
            </TabPanel>
        </TabPanels>
    </Tabs>
</template>

<style scoped></style>
