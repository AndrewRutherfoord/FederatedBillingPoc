<template>
    <h2>Billing Account Overview</h2>
    <hr>
    <div>
        <div class="flex justify-between items-center my-2">
            <h2>Linked Cloud Providers</h2>
            <Button label="Link Cloud Provider" icon="pi pi-plus" size="small" />
        </div>
        <DataTable :value="state" class="mt-2">
            <template #empty>No linked cloud providers found.</template>
            <Column field="cloud_provider_name" header="Cloud Provider"></Column>
        </DataTable>
    </div>

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
<script setup lang="ts">
import { ref } from 'vue';
import { DataTable, Column, Dialog } from 'primevue';
import { useAsyncState } from '@vueuse/core';
import client from '@/api/client';
import TextField from '@/components/form/TextField.vue';
import { Form } from 'vee-validate';

const formDialogVisible = ref(false);

const { state, execute: refreshCsps } = useAsyncState(async () => {
    const { data } = await client.GET("/billing/accounts");
    return data;
}, null);

const onSubmit = async (values: Record<string, unknown>) => {

}
</script>

<style scoped></style>