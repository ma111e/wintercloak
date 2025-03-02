<script setup lang="ts">
import {Card} from '@/components/ui/card'
import {type BasicColorSchema, useColorMode} from '@vueuse/core'
import {computed} from 'vue'
import {Icon} from '@iconify/vue'
import {Switch} from '@/components/ui/switch'

const isDarkColorMode = computed(() => {
  return colorMode.value.includes("dark")
})

const colorMode = useColorMode({
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

const themes = ['Slate', 'Orange', 'Yellow', 'Rose'] as const
type Theme = typeof themes[number]

const toggleDark = () => {
  if (isDarkColorMode.value) {
    colorMode.value = colorMode.value.replace("dark", "light") as BasicColorSchema
  } else {
    colorMode.value = colorMode.value.replace("light", "dark") as BasicColorSchema
  }
}

const currentTheme = computed(() => {
  // Extract the theme name from the color mode
  // e.g., "lightOrange" -> "Orange", "darkSlate" -> "Slate"
  return colorMode.value.replace(/^(light|dark)/, '') as Theme
})

const isCurrentTheme = (theme: Theme) => theme === currentTheme.value

const selectTheme = (theme: Theme) => {
  if (theme === currentTheme.value) return
  colorMode.value = `${isDarkColorMode.value ? 'dark' : 'light'}${theme}`
}
</script>

<template>
  <main class="grid flex-1 gap-4 overflow-auto p-4 pb-0">
    <div class="container mx-auto p-6 w-full">
      <h1 class="text-3xl font-bold mb-6 w-full">Theme Settings
        <Switch
            class="align-middle ml-2"
            :checked="isDarkColorMode"
            @click="toggleDark"
        >
          <template #thumb>
            <span
                class="flex justify-center align-middle h-full w-full size-3"
            >
              <Icon v-if="isDarkColorMode" icon="radix-icons:moon" class="self-center size-3"/>
              <Icon v-else icon="radix-icons:sun" class="self-center size-3"/>
            </span>
          </template>
        </Switch>
      </h1>
      <!--      <div class="grid gap-6 grid-cols-2">-->
      <!--          Dark mode-->

      <!--      </div>-->
      <div class="grid gap-6 md:grid-cols-3 lg:grid-cols-4">
        <Card
            v-for="theme in themes"
            :key="theme"
            class="p-4 cursor-pointer hover:bg-muted/50"
            :class="{ 'ring-2 ring-primary': isCurrentTheme(theme) }"
            @click="selectTheme(theme)"
        >
          <h2 class="text-xl font-semibold mb-4">{{ theme }}</h2>
          <div
              class="flex gap-2"
              :class="[
              isDarkColorMode ? `dark${theme}` : `light${theme}`,
            ]"
          >
            <div class="w-8 h-8 rounded-sm shadow-sm" style="background-color: var(--primary)"/>
            <div class="w-8 h-8 rounded-sm shadow-sm" style="background-color: var(--secondary)"/>
            <div class="w-8 h-8 rounded-sm shadow-sm" style="background-color: var(--muted)"/>
          </div>
        </Card>
      </div>
    </div>
  </main>
</template>