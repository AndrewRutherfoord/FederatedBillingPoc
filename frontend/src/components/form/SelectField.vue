<template>
  <div class="flex flex-col">
    <label :for="name">{{ label }}</label>

    <Select v-if="!multiple" :id="name" v-model="value" :placeholder="placeholder" :multiple="multiple"
      :disabled="disabled" :options="options" :optionLabel="optionLabel" :optionValue="optionValue">
      <template #value="slotProps">
        <slot name="value" v-bind="slotProps"></slot>
      </template>
      <template #option="slotProps">
        <slot name="option" v-bind="slotProps"></slot>
      </template>
    </Select>

    <MultiSelect v-else :id="name" v-model="value" :disabled="disabled" :options="options" :optionLabel="optionLabel"
      :optionValue="optionValue" display="chip">
      <template #option="slotProps">
        <slot name="option" :slotProps="slotProps"></slot>
      </template>
    </MultiSelect>
    <small class="text-muted" :id="`${name}-helper`">{{ helpText }}</small>
    <Message v-if="errorMessage" severity="error" class="mt-2">{{ errorMessage }}</Message>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import Select from "primevue/select";
import MultiSelect from "primevue/multiselect";
import { useField } from "vee-validate";
import Message from "primevue/message";

const props = defineProps({
  name: {
    type: String,
    required: true
  },
  label: String,
  helpText: {
    required: false,
    type: String
  },
  placeholder: {
    type: String,
    default: "Select an option..."
  },
  disabled: {
    type: Boolean,
    default: false
  },
  multiple: {
    type: Boolean,
    default: false
  },
  options: {
    type: Array,
    required: true
  },
  optionLabel: {
    type: String,
    default: "name"
  },
  optionValue: {
    type: String,
    default: "id"
  },
  modelValue: {
    type: [String, Number, Array],
    default: null
  }
});

const emit = defineEmits(["update:modelValue"]);

const { value, errorMessage } = useField(() => props.name);

watch(() => props.modelValue, (newValue) => {
  value.value = newValue;
});

watch(() => value.value, (newValue) => {
  emit("update:modelValue", newValue);
});
</script>

<style scoped></style>
