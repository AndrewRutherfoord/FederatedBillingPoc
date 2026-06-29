<template>
    <main>
        <div class="flex justify-between items-center">
            <h2>{{ billingAccount?.alias }} - Invoices</h2>
            <Button as="router-link" :to="{ name: 'Home', query: { account: props.id } }" text size="small">
                Back to Overview
            </Button>
        </div>
        <p class="text-sm text-gray-500">
            Expand a row to see each cloud provider's share of the invoice and the merkle root proving
            its charge batches weren't tampered with.
        </p>
        <DataTable :value="invoices" v-model:expandedRows="expandedRows" dataKey="id" class="mt-2">
            <template #empty>No invoices found.</template>
            <Column expander style="width: 3rem" />
            <Column header="Invoice ID">
                <template #body="{ data }">
                    <span class="font-mono text-xs" :title="data.id">{{ shortId(data.id) }}</span>
                </template>
            </Column>
            <Column header="Billing Period">
                <template #body="{ data }">
                    <span class="font-mono text-xs" :title="data.billing_period_id">{{ shortId(data.billing_period_id) }}</span>
                </template>
            </Column>
            <Column header="Amount">
                <template #body="{ data }">
                    {{ data.amount }} {{ data.currency }}
                </template>
            </Column>
            <Column header="Status">
                <template #body="{ data }">
                    <Tag :value="data.status" :severity="statusSeverity(data.status)" />
                </template>
            </Column>
            <Column header="Issued">
                <template #body="{ data }">
                    {{ formatDate(data.issued_at) }}
                </template>
            </Column>
            <Column header="Due">
                <template #body="{ data }">
                    {{ formatDate(data.due_at) }}
                </template>
            </Column>
            <template #expansion="{ data }">
                <DataTable :value="data.provider_line_items">
                    <Column field="cloud_service_provider_id" header="Cloud Provider"></Column>
                    <Column header="Amount">
                        <template #body="{ data: lineItem }">
                            {{ lineItem.amount }} {{ data.currency }}
                        </template>
                    </Column>
                    <Column field="batch_count" header="Batches"></Column>
                    <Column header="Merkle Root">
                        <template #body="{ data: lineItem }">
                            <span class="font-mono text-xs">{{ lineItem.merkle_root }}</span>
                        </template>
                    </Column>
                </DataTable>
            </template>
        </DataTable>
    </main>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { DataTable, Column, Tag } from 'primevue';
import { useAsyncState } from '@vueuse/core';
import client from '@/api/client';

const props = defineProps<{
    id: string
}>();

const expandedRows = ref({});

const { state: billingAccount } = useAsyncState(async () => {
    const { data } = await client.GET("/billing/accounts/{id}", {
        params: { path: { id: props.id } }
    });
    return data;
}, null);

const { state: invoices } = useAsyncState(async () => {
    const { data } = await client.GET("/billing/accounts/{id}/invoices", {
        params: { path: { id: props.id } }
    });
    return data ?? [];
}, []);

const shortId = (id?: string) => id ? `${id.slice(0, 8)}…` : '';

const formatDate = (value?: string) => value ? new Date(value).toLocaleString() : 'N/A';

const statusSeverity = (status?: string) => ({
    issued: 'info',
    paid: 'success',
    overdue: 'danger',
}[status ?? ''] ?? 'secondary');
</script>

<style scoped></style>
