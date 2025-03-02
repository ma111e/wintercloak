<script setup lang="ts">
import {Button} from '@/components/ui/button'
import {DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger} from '@/components/ui/dropdown-menu'
import {Icon} from '@iconify/vue'
import {type BasicColorSchema, useColorMode, useDark} from '@vueuse/core'

const isDark = useDark()
// const toggleDark = useToggle(isDark)

// Pass { disableTransition: false } to enable transitions
// const mode = useColorMode()
const mode = useColorMode({
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

const setTheme = (theme: string) => {
  mode.value = theme as BasicColorSchema
  // isDark.value = theme.includes("dark")
}
</script>

<template>
  <div>
    <DropdownMenu>
      <DropdownMenuTrigger as-child>
        <Button variant="outline">
          <Icon icon="radix-icons:moon"
                class="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0"/>
          <Icon icon="radix-icons:sun"
                class="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100"/>
          <span class="sr-only">Toggle theme</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
<!--        <DropdownMenuItem-->
<!--            class="cursor-pointer"-->
<!--        >-->
<!--&lt;!&ndash;          <Switch&ndash;&gt;-->
<!--&lt;!&ndash;              v-model:checked="isDark"&ndash;&gt;-->
<!--&lt;!&ndash;              id="dark-mode">&ndash;&gt;-->
<!--&lt;!&ndash;            <template #thumb>&ndash;&gt;-->
<!--&lt;!&ndash;              <span&ndash;&gt;-->
<!--&lt;!&ndash;                class="flex justify-center align-middle"&ndash;&gt;-->
<!--&lt;!&ndash;              >&ndash;&gt;-->
<!--&lt;!&ndash;              <Icon v-if="isDark" icon="radix-icons:moon" class="size-3"/>&ndash;&gt;-->
<!--&lt;!&ndash;              <Icon v-else icon="radix-icons:sun" class="size-3"/>&ndash;&gt;-->
<!--&lt;!&ndash;              </span>&ndash;&gt;-->
<!--&lt;!&ndash;            </template>&ndash;&gt;-->
<!--&lt;!&ndash;          </Switch>&ndash;&gt;-->
<!--          <Label for="dark-mode"-->
<!--                 class="cursor-pointer"-->
<!--          >-->
<!--            Dark mode-->
<!--          </Label>-->
<!--        </DropdownMenuItem>-->
        <DropdownMenuItem v-if="isDark" @click="setTheme('dark')">
          Dark
        </DropdownMenuItem>
        <DropdownMenuItem v-if="isDark" @click="setTheme('darkSlate')">
          Dark Slate
        </DropdownMenuItem>
        <DropdownMenuItem v-if="isDark" @click="setTheme('darkYellow')">
          Dark Yellow
        </DropdownMenuItem>
        <DropdownMenuItem v-if="isDark" @click="setTheme('darkRose')">
          Dark Rose
        </DropdownMenuItem>
        <DropdownMenuItem v-if="isDark" @click="setTheme('darkOrange')">
          Dark Orange
        </DropdownMenuItem>

        <DropdownMenuItem v-if="!isDark" @click="setTheme('light')">
          Light
        </DropdownMenuItem>
        <DropdownMenuItem v-if="!isDark" @click="setTheme('lightSlate')">
          Light Slate
        </DropdownMenuItem>
        <DropdownMenuItem v-if="!isDark" @click="setTheme('lightYellow')">
          Light Yellow
        </DropdownMenuItem>
        <DropdownMenuItem v-if="!isDark" @click="setTheme('lightRose')">
          Light Rose
        </DropdownMenuItem>
        <DropdownMenuItem v-if="!isDark" @click="setTheme('lightOrange')">
          Light Orange
        </DropdownMenuItem>

        <DropdownMenuItem @click="mode = 'auto'">
          System
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  </div>
</template>