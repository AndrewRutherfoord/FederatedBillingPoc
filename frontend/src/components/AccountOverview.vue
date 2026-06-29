<template>
    <main v-if="billingAccount">

        <div class="mt-4">
            <div class="flex justify-between items-center">
                <h3>Cost Breakdown</h3>
                <div class="flex gap-2">
                    <Button asChild v-slot="slotProps" size="small">
                        <RouterLink :to="{ name: 'ChargeBatches', params: { id: props.id } }" :class="slotProps.class">
                        View Charge Batches
                        </RouterLink>
                    </Button>
                    <Button asChild v-slot="slotProps" size="small">
                        <RouterLink :to="{ name: 'ResourceCharges', params: { id: props.id } }" :class="slotProps.class">View Resource Charges</RouterLink>
                    </Button>
                    <Button asChild v-slot="slotProps" size="small">
                        <RouterLink :to="{ name: 'Invoices', params: { id: props.id } }" :class="slotProps.class">View Invoices</RouterLink>
                    </Button>
                </div>
            </div>
            <div v-if="resourceCharges && resourceCharges.length > 0" class="flex flex-col md:flex-row gap-4 mt-2">
                <div class="card flex-1 flex flex-col items-center">
                    <h4>By Cloud Provider</h4>
                    <Chart type="pie" :data="cspChartData" :options="chartOptions" class="w-full md:w-[24rem]" />
                </div>
                <div class="card flex-1 flex flex-col items-center">
                    <h4>By Service Category</h4>
                    <Chart type="pie" :data="categoryChartData" :options="chartOptions" class="w-full md:w-[24rem]" />
                </div>
            </div>
            <p v-else class="text-sm text-gray-500 mt-2">No cost data available yet.</p>
        </div>

        <div>
            <div class="flex justify-between items-center my-2">
                <h3>Linked Cloud Providers</h3>
                <Button asChild v-slot="slotProps" size="small">
                    <RouterLink :to="{ name: 'CloudProviderLinks', params: { id: props.id } }" :class="slotProps.class">
                        Manage Cloud Provider Links
                    </RouterLink>
                </Button>
            </div>
            <DataTable :value="cloudProviderLinks" class="mt-2">
                <template #empty>No linked cloud providers found.</template>
                <Column field="cloud_provider_name" header="Cloud Provider"></Column>
                <Column field="cloud_provider_id" header="Provider ID"></Column>
                <Column header="Unpaid Cost">
                    <template #body="slotProps">
                        <span v-if="slotProps.data.total_cost !== undefined">
                            {{ truncate3dp(slotProps.data.total_cost) }} {{ slotProps.data.billing_currency }}
                        </span>
                        <span v-else>
                            N/A
                        </span>
                    </template>
                </Column>
            </DataTable>
        </div>
    </main>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { DataTable, Column } from 'primevue';
import Chart from 'primevue/chart';
import { useAsyncState } from '@vueuse/core';
import client from '@/api/client';
import type { components } from '@/api/api-schema';

type ResourceCharge = components['schemas']['handlers.ResourceChargeEntry'];

const props = defineProps<{
    id: string
}>();

const { state: billingAccount } = useAsyncState(async () => {
    const { data } = await client.GET("/billing/accounts/{id}", {
        params: { path: { id: props.id } }
    });
    return data;
}, null);

const { state: cloudProviderLinks } = useAsyncState(async () => {
    const { data } = await client.GET("/billing/accounts/{id}/cloud-provider-accounts", {
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

const CHART_COLOR_NAMES = ['cyan', 'orange', 'gray', 'green', 'purple', 'red', 'blue', 'teal', 'yellow', 'pink', 'indigo', 'lime'];

const truncate3dp = (value: number) => (Math.trunc(value * 1000) / 1000).toFixed(3);

const formatLabel = (value: string) => value.split('_').map((word) => word.charAt(0).toUpperCase() + word.slice(1)).join(' ');

const aggregateCost = (charges: ResourceCharge[], keyOf: (charge: ResourceCharge) => string) => {
    const totals = new Map<string, number>();
    for (const charge of charges) {
        const key = keyOf(charge);
        totals.set(key, (totals.get(key) ?? 0) + (charge.total_billed_cost ?? 0));
    }
    return totals;
};

const buildPieData = (totals: Map<string, number>, labelOf: (key: string) => string) => {
    const documentStyle = getComputedStyle(document.body);
    const entries = [...totals.entries()];
    return {
        labels: entries.map(([key]) => labelOf(key)),
        datasets: [
            {
                data: entries.map(([, cost]) => cost),
                backgroundColor: entries.map((_, i) => documentStyle.getPropertyValue(`--p-${CHART_COLOR_NAMES[i % CHART_COLOR_NAMES.length]}-500`)),
                hoverBackgroundColor: entries.map((_, i) => documentStyle.getPropertyValue(`--p-${CHART_COLOR_NAMES[i % CHART_COLOR_NAMES.length]}-400`)),
            }
        ]
    };
};

const cspNameById = computed(() => {
    const lookup = new Map<string, string>();
    for (const link of cloudProviderLinks.value ?? []) {
        if (link.cloud_provider_id) {
            lookup.set(link.cloud_provider_id, link.cloud_provider_name ?? link.cloud_provider_id);
        }
    }
    return lookup;
});

const cspChartData = computed(() => {
    const totals = aggregateCost(resourceCharges.value ?? [], (c) => c.cloud_service_provider_id ?? 'Unknown');
    return buildPieData(totals, (id) => cspNameById.value.get(id) ?? id);
});

const categoryChartData = computed(() => {
    const totals = aggregateCost(resourceCharges.value ?? [], (c) => c.service_category ?? 'other');
    return buildPieData(totals, formatLabel);
});

const chartOptions = computed(() => {
    const documentStyle = getComputedStyle(document.documentElement);
    return {
        plugins: {
            legend: {
                labels: {
                    usePointStyle: true,
                    color: documentStyle.getPropertyValue('--p-text-color')
                }
            }
        }
    };
});
</script>

<style scoped></style>
