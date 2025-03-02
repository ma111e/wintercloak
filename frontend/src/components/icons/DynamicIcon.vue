<script setup lang="ts">
import { ref, watch, defineProps } from 'vue';

const props = defineProps<{ name: string }>();
const iconComponent = ref<any>(null);

const loadIcon = async () => {
  try {
    const icon = await import(`lucide-vue-next/dist/esm/icons/${props.name}`);
    iconComponent.value = icon.default;
  } catch (error) {
    console.error(`Failed to load icon: ${props.name}`, error);
    iconComponent.value = null;
  }
};

watch(() => props.name, loadIcon, { immediate: true });
</script>

<template>
  <component :is="iconComponent" v-if="iconComponent" v-bind="$attrs" />
</template>