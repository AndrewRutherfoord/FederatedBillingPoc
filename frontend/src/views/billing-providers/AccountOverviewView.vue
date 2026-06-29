<template>
    <main v-if="billingAccount">
        <div class="flex justify-between items-center">
            <h2>{{ billingAccount?.alias }} - Billing Account Overview</h2>
            <div class="flex gap-2">
                <RouterLink :to="{ name: 'ChargeBatches', params: { id: props.id } }" class="p-button p-button-text">
                    View Charge Batches
                </RouterLink>
                <RouterLink :to="{ name: 'ResourceCharges', params: { id: props.id } }" class="p-button p-button-text">
                    View Resource Charges
                </RouterLink>
                <RouterLink :to="{ name: 'Invoices', params: { id: props.id } }" class="p-button p-button-text">
                    View Invoices
                </RouterLink>
            </div>
        </div>
        <hr>
        <div>
            <div class="flex justify-between items-center my-2">
                <h3>Linked Cloud Providers</h3>
                <Button label="Link Cloud Provider" icon="pi pi-plus" size="small" @click="formDialogVisible = true" />
            </div>
            <DataTable :value="cloudProviderLinks" class="mt-2">
                <template #empty>No linked cloud providers found.</template>
                <Column field="cloud_provider_name" header="Cloud Provider"></Column>
                <Column field="cloud_provider_id" header="Provider ID"></Column>
                <Column header="Total Cost">
                    <template #body="slotProps">
                        <span v-if="slotProps.data.total_cost !== undefined">
                            {{ slotProps.data.total_cost }} {{ slotProps.data.billing_currency }}
                        </span>
                        <span v-else>
                            N/A
                        </span>
                    </template>
                </Column>
            </DataTable>
        </div>

        <div class="mt-4">
            <h3>Cost Breakdown</h3>
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

        <Dialog v-model:visible="formDialogVisible" modal header="Link Cloud Provider" :style="{ width: '25rem' }">
            <Form @submit="onSubmit">
                <SelectField name="cloud_service_provider_id" label="Cloud Service Provider"
                    :options="cloudProviderOptions" class="mb-2"></SelectField>
                <p v-if="linkError" class="text-red-500 text-sm mt-1">{{ linkError }}</p>
                <Button label="Link Provider" type="submit" class="mt-2" :loading="linking" />
            </Form>
        </Dialog>
    </main>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import { DataTable, Column, Dialog } from 'primevue';
import Chart from 'primevue/chart';
import { useAsyncState } from '@vueuse/core';
import client from '@/api/client';
import { Form } from 'vee-validate';
import SelectField from '@/components/form/SelectField.vue';
import type { components } from '@/api/api-schema';

type ResourceCharge = components['schemas']['handlers.ResourceChargeEntry'];

const formDialogVisible = ref(false);
const linking = ref(false);
const linkError = ref<string | null>(null);

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

const cloudProviderOptions = computed<{ name: string; id: string }[]>(() => {
    if (!billingAccount.value) return [];
    // @ts-ignore
    return billingAccount.value.supported_cloud_providers?.map((csp: any) => ({
        id: csp.id,
        name: csp.name,
    })) ?? [];
});

const onSubmit = async (values: { cloud_service_provider_id: string }) => {
    if (!billingAccount.value) return;
    linkError.value = null;
    linking.value = true;

    const returnURL = `${window.location.origin}/cloud-service-providers/onboarding-complete-callback`;

    try {
        const { data, error } = await client.POST("/billing/accounts/{id}/cloud-provider-accounts/register", {
            params: { path: { id: props.id } },
            body: {
                account_id: props.id,
                cloud_provider_id: values.cloud_service_provider_id,
                return_url: returnURL,
            },
        });

        if (error || !data) {
            linkError.value = 'Failed to initiate cloud provider onboarding.';
            return;
        }

        window.location.href = (data as any).redirect_url;
    } catch {
        linkError.value = 'An unexpected error occurred.';
    } finally {
        linking.value = false;
    }
};
</script>

<style scoped></style>
