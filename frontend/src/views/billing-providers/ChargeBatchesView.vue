<template>
    <main>
        <div class="flex justify-between items-center">
            <h2>{{ billingAccount?.alias }} - Charge Batches</h2>
            <Button as="router-link" :to="{ name: 'Home', query: { account: props.id } }" text size="small">
                Back to Overview
            </Button>
        </div>
        <p class="text-sm text-gray-500">
            Each batch is reported independently by the billing provider and by the cloud service provider.
            Rows are flagged when the two reports disagree or when one side hasn't reported the batch at all.
        </p>
        <DataTable :value="chargeBatches" class="mt-2">
            <template #empty>No charge batches found.</template>
            <Column field="batch_id" header="Batch ID">
                <template #body="{ data }">
                    <span class="font-mono text-xs" :title="data.batch_id">{{ shortId(data.batch_id) }}</span>
                </template>
            </Column>
            <Column field="cloud_service_provider_id" header="Cloud Provider"></Column>
            <Column header="Status">
                <template #body="{ data }">
                    <Tag :value="statusLabel(data.status)" :severity="statusSeverity(data.status)" />
                </template>
            </Column>
            <Column header="Billing Provider Report">
                <template #body="{ data }">
                    <span v-if="data.billing_provider_report">
                        {{ data.billing_provider_report.total_items }} items,
                        {{ data.billing_provider_report.total_cost }}
                        {{ data.billing_provider_report.billed_currency }}
                    </span>
                    <span v-else class="text-gray-400">Not reported</span>
                </template>
            </Column>
            <Column header="Cloud Provider Report">
                <template #body="{ data }">
                    <span v-if="data.cloud_provider_report">
                        {{ data.cloud_provider_report.total_items }} items,
                        {{ data.cloud_provider_report.total_cost }}
                        {{ data.cloud_provider_report.billed_currency }}
                    </span>
                    <span v-else class="text-gray-400">Not reported</span>
                </template>
            </Column>
            <Column header="Created">
                <template #body="{ data }">
                    {{ formatDate(data.billing_provider_report?.created_at ?? data.cloud_provider_report?.created_at) }}
                </template>
            </Column>
        </DataTable>
    </main>
</template>

<script setup lang="ts">
import { DataTable, Column, Tag } from 'primevue';
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

const { state: chargeBatches } = useAsyncState(async () => {
    const { data } = await client.GET("/billing/accounts/{id}/charge-batches", {
        params: { path: { id: props.id } }
    });
    return data ?? [];
}, []);

const shortId = (id?: string) => id ? `${id.slice(0, 8)}…` : '';

const formatDate = (value?: string) => value ? new Date(value).toLocaleString() : 'N/A';

const statusLabel = (status?: string) => ({
    matched: 'Matched',
    mismatched: 'Mismatched',
    missing_from_csp: 'Missing CSP report',
    missing_from_bp: 'Missing BP report',
}[status ?? ''] ?? status ?? 'Unknown');

const statusSeverity = (status?: string) => ({
    matched: 'success',
    mismatched: 'danger',
    missing_from_csp: 'warn',
    missing_from_bp: 'warn',
}[status ?? ''] ?? 'secondary');
</script>

<style scoped></style>
