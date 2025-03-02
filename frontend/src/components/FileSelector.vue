<script setup lang="ts">
import {ref} from 'vue'
import {Button} from '@/components/ui/button'
import {Label} from '@/components/ui/label'
import {SelectFolder} from "../../wailsjs/go/main/App";

import {FolderSearch} from "lucide-vue-next";
import FileTree from "@/components/FileTree.vue";

const selectedPathData = ref({path: '', items: {}})
const emit = defineEmits(['selected', "file-clicked"])

const handleFolderSelect = async () => {
  try {
    const res = await SelectFolder()
    if (res) {
      emit("selected", res.path)
      selectedPathData.value = res
    }
  } catch (error) {
    console.error('Error selecting file:', error)
  }
}

function truncateFilePath(filePath: string, depth: number): string {
  const parts = filePath.split(/[/\\]/); // Split by both forward and backslashes
  const truncatedParts = parts.slice(-depth); // Get the last 3 elements
  return parts.length > depth ? `.../${truncatedParts.join("/")}` : filePath;
}
</script>

<template>
  <div class="space-y-4">
    <div
        class="rounded-md"
    >
      <Button
          id="file-select-btn"
          @click="handleFolderSelect"
          variant="default"
      >
        <FolderSearch/>
      </Button>
      <Label
          for="file-select-btn"
          :class="['break-all ml-3', selectedPathData.path ? 'font-[Courier]': 'italic']"
      >
        {{ selectedPathData.path ? truncateFilePath(selectedPathData.path, 2) : 'No folder selected' }}
      </Label>
      <FileTree
          class="h-[70vh]"
          v-if="selectedPathData.path"
          :folder-data="selectedPathData"
          @file-clicked="(val) => $emit('file-clicked', val)"
      />
    </div>
  </div>
</template>