<script setup lang="ts">

import {Button} from "@/components/ui/button"
import {Triangle} from "lucide-vue-next"
</script>

<template>
  <header class="sticky top-0 z-40 flex h-[57px] items-center gap-1 border-b bg-background px-4">
    <router-link to="/" class="text-xl font-semibold">
      <Button variant="outline" size="icon" aria-label="Home">
        <Triangle class="size-5 fill-foreground"/>
      </Button>
      Wintercloak
    </router-link>

    <!-- divider -->
    <span
        class="ml-auto"
    />

  </header>
</template>

<style scoped>

</style>