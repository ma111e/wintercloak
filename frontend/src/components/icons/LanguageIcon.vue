<script setup lang="ts">
import LanguagePowerShell from "@/components/icons/LanguagePowerShell.vue";
import LanguageGo from "@/components/icons/LanguageGo.vue";
import LanguagePython from "@/components/icons/LanguagePython.vue";
import FileQuestion from "@/components/icons/FileQuestion.vue";
import CodeBlockBraces from "@/components/icons/CodeBlockBraces.vue";

interface LanguageIconProps {
  lang: string
}

const props = defineProps<LanguageIconProps>()
</script>

<template>
  <LanguagePowerShell
      v-if="props.lang === 'powershell'"
  />
  <LanguageGo
      v-else-if="props.lang === 'golang'"
  />
  <LanguagePython
      v-else-if="props.lang === 'python'"
  />
  <CodeBlockBraces
      v-else-if="props.lang === 'text'"
  />
<!--                          TODO: add compiled linux 32 + 64bit and windows compiled 32 + 64bit -->
  <FileQuestion
      v-else
  />
</template>
