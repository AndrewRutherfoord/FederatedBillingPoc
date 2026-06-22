<template>
    <main>
        <div class="flex justify-between items-center">
            <h2>{{ billingAccount?.alias }} - Resource Charges</h2>
            <Button as="router-link" :to="{ name: 'Home', query: { account: props.id } }" text size="small">
                Back to Overview
            </Button>
        </div>
        <p class="text-sm text-gray-500">
            Billed cost aggregated per resource, sourced entirely from the FOCUS line items fetched directly from
            the cloud service provider.
        </p>
        <DataTable :value="resourceCharges" class="mt-2">
            <template #empty>No resource charges found.</template>
            <Column header="Resource">
                <template #body="{ data }">
                    <span>{{ data.resource_name ?? data.resource_id ?? 'Unattributed' }}</span>
                </template>
            </Column>
            <Column field="resource_type" header="Resource Type"></Column>
            <Column field="service_name" header="Service"></Column>
            <Column field="service_category" header="Service Category"></Column>
            <Column header="Cloud Provider">
                <template #body="{ data }">
                    {{ data.cloud_service_provider_id }}
                </template>
            </Column>
            <Column header="Total Billed Cost">
                <template #body="{ data }">
                    {{ data.total_billed_cost }} {{ data.billing_currency }}
                </template>
            </Column>
            <Column field="line_item_count" header="Line Items"></Column>
        </DataTable>
    </main>
</template>

<script setup lang="ts">
import { DataTable, Column } from 'primevue';
import { useAsyncState } from '@vueuse/core';
import client from '@/api/client';

const props = defineProps<{
    id: string
}>();

const { state: billingAccount } = useAsyncState(async () => {
    const { data } = await client.GET("/billing/accounts/{id}", {
        params: { path: { id: props.id } }
    });
    return data;
}, null);

const { state: resourceCharges } = useAsyncState(async () => {
    const { data } = await client.GET("/billing/accounts/{id}/resource-charges", {
        params: { path: { id: props.id } }
    });
    return data ?? [];
}, []);
</script>

<style scoped></style>
