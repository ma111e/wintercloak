<script setup lang="ts">
import Header from "@/views/Header.vue"
import LeftSidebar from "@/views/LeftSidebar.vue";
import {engines} from "../wailsjs/go/models.ts";
import {onBeforeMount, reactive} from "vue";
import {ListAvailableEngines} from "../wailsjs/go/main/App";
import Toaster from '@/components/ui/toast/Toaster.vue'
import { useColorMode} from "@vueuse/core";

onBeforeMount(() => {
  refreshEngines().then(() => {
    // if (data.availableEngines.length > 0) {
    //   router.push(`/playground/${data.availableEngines[0].name.toLowerCase()}`)
    // }
  })
})

const data = reactive({
  availableEngines: [] as engines.EngineMeta[],
})

function refreshEngines(): Promise<void> {
  return ListAvailableEngines().then((res: engines.EngineMeta[]) => {
    data.availableEngines = res
  })
}

useColorMode({
  modes: {
    lightOrange: 'lightOrange',
    darkYellow: 'darkYellow',
    lightYellow: 'lightYellow',
    lightRose: 'lightRose',
    darkRose: 'darkRose',
    lightSlate: "lightSlate",
    darkSlate: "darkSlate",
    darkOrange: "darkOrange"
  },
})

</script>

<template>
  <Header/>
  <div class="flex bg-background max-h-[calc(100vh_-_57px)]">
    <LeftSidebar
        :engines="data.availableEngines"
    />
    <router-view
        :engines="data.availableEngines"
    />
  </div>
  <Toaster/>
</template>