<script setup lang="ts">
import { useAsyncState } from '@vueuse/core';
import client from '@/api/client';
import { DataTable, Column, Dialog } from 'primevue';
import { ref } from 'vue';
import TextField from '@/components/form/TextField.vue';
import { Form } from 'vee-validate';
import type { components as ApiComponents } from '@/api/api-schema';

type RegisterRequest = ApiComponents["schemas"]["handlers.RegisterAcccountRequest"];

const formDialogVisible = ref(false);

const { state, error, execute: refreshAccounts } = useAsyncState(async () => {
    const { data } = await client.GET("/billing/accounts");
    return data;
}, null);

const onSubmit = async (values: Record<string, unknown>) => {
    const { data, error } = await client.POST("/billing/accounts/register", {
        body: {
            ...values,
            return_url: `${window.location.origin}/billing-providers/onboarding-complete-callback`,
        } as RegisterRequest
    });

    if (error) {
        console.error("Error creating billing account:", error);
        return;
    }

    formDialogVisible.value = false;

    if (data?.redirect_url) {
        window.location.href = data.redirect_url;
    } else {
        refreshAccounts();
    }
};
</script>

<template>
    <div class="flex justify-between items-center my-2">
        <h2>Billing Accounts</h2>
        <Button label="Create Billing Account" icon="pi pi-plus" @click="formDialogVisible = true" size="small" />
    </div>
    <DataTable :value="state" class="mt-2">
        <template #empty>No billing accounts found.</template>
        <Column field="alias" header="Account Alias"></Column>
        <Column field="billing_provider_name" header="Billing Provider"></Column>
        <Column>
            <template #body="{ data }">
                <RouterLink :to="{ name: 'BillingAccountOverview', params: { id: data.id } }" as="Button">View</RouterLink>
            </template>
        </Column>
    </DataTable>

    <Dialog v-model:visible="formDialogVisible" modal header="Create Billing Account" :style="{ width: '25rem' }">
        <Form @submit="onSubmit">
            <TextField name="account_alias" label="Account Alias"
                helpText="The name you will use to identify the billing account." />
            <TextField name="billing_provider_base_url" label="Billing Provider URL"
                helpText="The URL of the billing provider to create the account with." />
            <Button label="Create" type="submit" class="mt-2" />
        </Form>
    </Dialog>
</template>
