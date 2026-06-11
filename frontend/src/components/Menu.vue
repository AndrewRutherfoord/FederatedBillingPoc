<template>
    <Menubar :model="items">
        <template #item="{ item, props, hasSubmenu, root }">
            <RouterLink v-if="item.route" v-ripple class="flex items-center" :to="{ name: item.route }" v-bind="props.action">
                <span>{{ item.label }}</span>
                <Badge v-if="item.badge" :class="{ 'ml-auto': !root, 'ml-2': root }" :value="item.badge" />
                <span v-if="item.shortcut"
                    class="ml-auto border border-surface rounded bg-emphasis text-muted-color text-xs p-1">{{
                    item.shortcut }}</span>
                <i v-if="hasSubmenu"
                    :class="['pi pi-angle-down ml-auto', { 'pi-angle-down': root, 'pi-angle-right': !root }]"></i>
            </RouterLink>
        </template>
        <template #end>
            <div class="flex items-center gap-2">
                <InputText placeholder="Search" type="text" class="w-32 sm:w-auto" />
                <Avatar image="/images/avatar/amyelsner.png" shape="circle" />
            </div>
        </template>
    </Menubar>

</template>
<script setup lang="ts">
import { ref } from 'vue';
import { RouterLink } from 'vue-router';
import { Menubar } from "primevue";

const items = ref([
    {
        label: 'Home',
        icon: 'pi pi-home',
        route: 'home'
    },
    {
        label: 'About',
        icon: 'pi pi-info-circle',
        route: 'about'
    },
]);


</script>

<style scoped></style>