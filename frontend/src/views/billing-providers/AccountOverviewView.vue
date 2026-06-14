<template>
    <main v-if="state">
        <h2>{{ state?.alias }} - Billing Account Overview</h2>
        <hr>
        <div>
            <div class="flex justify-between items-center my-2">
                <h2>Linked Cloud Providers</h2>
                <Button label="Link Cloud Provider" icon="pi pi-plus" size="small" @click="formDialogVisible = true" />
            </div>
            <!-- <DataTable :value="state" class="mt-2">
                <template #empty>No linked cloud providers found.</template>
<Column field="cloud_provider_name" header="Cloud Provider"></Column>
</DataTable> -->
        </div>

        <Dialog v-model:visible="formDialogVisible" modal header="Link Cloud Provider" :style="{ width: '25rem' }">
            <Form @submit="onSubmit">
                <SelectField name="cloud_service_provider_id" label="Cloud Service Provider" :options="cloudProviderOptions" class="mb-2"></SelectField>
                <Button label="Link Provider" type="submit" class="mt-2" />
            </Form>
        </Dialog>
    </main>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import { DataTable, Column, Dialog } from 'primevue';
import { useAsyncState } from '@vueuse/core';
import client from '@/api/client';
import TextField from '@/components/form/TextField.vue';
import { Form } from 'vee-validate';
import SelectField from '@/components/form/SelectField.vue';

const formDialogVisible = ref(false);

const props = defineProps<{
    id: string
}>();


const { state, execute: refreshCsps } = useAsyncState(async () => {
    const { data } = await client.GET("/billing/accounts/{id}", {
        params: {
            path: {
                id: props.id
            }
        }
    });
    return data;
}, null);

// @ts-ignore - Unecessarily pedandic...
const cloudProviderOptions = computed<{ name: string; id: string }[]>(() => {
    if (!state.value) return [];
    return state.value.supported_cloud_providers?.map(csp => ({
        id: csp.id,
        name: csp.name,
    })) ?? []
})

const onSubmit = async (values: { cloud_service_provider_id: string }) => {


}
</script>

<style scoped></style>