<script setup lang="ts">
import Playground from "@/views/Playground.vue";
import Preprocessor from "@/views/Preprocessor.vue";
import {shallowRef, watch} from "vue";
import {useRoute} from "vue-router";
import {engines} from "../../wailsjs/go/models.ts";

const route = useRoute();
const currentComponent = shallowRef(Playground);

interface ConfigurationProps {
  engines: engines.EngineMeta[],
}

const props = defineProps<ConfigurationProps>()

// Watch for mode changes and update component
watch(() => route.query.mode, (newMode) => {
  currentComponent.value = newMode === "preprocessor" ? Preprocessor : Playground;
});
</script>

<template>
  <component :engines="props.engines" :is="currentComponent"/>
</template>