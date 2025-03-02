<script setup lang="ts">
import {Card, CardDescription, CardTitle,} from '@/components/ui/card'
import {engines} from "../../wailsjs/go/models.ts";
import LanguageIcon from "@/components/icons/LanguageIcon.vue";

interface SidebarProps {
  engines: engines.EngineMeta[]
}

defineProps<SidebarProps>()
</script>

<template>
  <main class="grid flex-1 gap-4 p-4 pb-0 h-[80vh]">
    <div
        class="container relative flex flex-col h-full w-full max-w-6xl px-10 mx-auto xl:px-0 mt-5 justify-center align-middle">
      <div class="grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
        <div class=""
             v-for="eng in $props.engines"
             :key="eng.name"
        >
          <div class="relative h-full">
            <Card
                class="p-6 pb-4 cursor-pointer bg-secondary h-full hover:bg-primary/10"
                @click="$router.push(`/engines/${eng.name.toLowerCase()}`)"
            >
              <CardTitle
                  class="flex justify-start items-center mb-4"
              >
                <span class="w-full">
                  {{ eng.name }}
                </span>
                <!--                <span-->
                <!--                class="ml-2"-->
                <!--                >-->

                <!--                </span>-->

                <span
                    class="w-full flex justify-end"
                >
<!--                <ExternalLink/>-->
                  <LanguageIcon :lang="eng.language"/>
                </span>
              </CardTitle>
              <CardDescription
                  class="text-pretty hyphens-auto"
              >
                {{ eng.shortDescription }}
              </CardDescription>
            </Card>
          </div>
        </div>
      </div>
    </div>
  </main>

</template>

<style scoped>

</style>