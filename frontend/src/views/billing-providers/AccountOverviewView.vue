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
import { useAsyncState } from '@vueuse/core';
import client from '@/api/client';
import { Form } from 'vee-validate';
import SelectField from '@/components/form/SelectField.vue';

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
