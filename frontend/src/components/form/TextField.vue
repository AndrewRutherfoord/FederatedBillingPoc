<template>
    <div class="flex flex-col gap-2 mb-2">
        <label :for="name">{{ label }}</label>
        <InputText :id="name" type="text" v-model="value" :disabled="disabled"></InputText>
        <small class="text-muted" :id="`${name}-helper`">{{ helpText }}</small>
        <Message v-if="errorMessage" severity="error">{{ errorMessage }}</Message>
    </div>

</template>

<script setup lang="ts">
import { ref } from 'vue';
import InputText from 'primevue/inputtext';
import { useField } from 'vee-validate';
import Message from 'primevue/message';

const props = defineProps({
    name: { type: String, required: true },
    label: String,
    helpText: {
        required: false,
        type: String
    },
    disabled: {
        type: Boolean,
        default: false
    }
})

const { value, errorMessage } = useField<string>(() => props.name)

</script>

<style scoped></style>