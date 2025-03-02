<script setup lang="ts">
import {ref} from "vue";
import FileTreeNode from "@/components/FileTreeNode.vue";
import {ScrollArea} from "@/components/ui/scroll-area"

defineProps({
  folderData: {
    type: Object,
    required: true,
  },
});

const expandedFolders = ref({});

</script>

<template>
  <ScrollArea :class="['p-4 w-full max-w-lg h-full']">
    <FileTreeNode
        v-for="item in folderData.items"
        :key="item.name"
        :item="item"
        v-model:expandedFolders="expandedFolders"
        @file-clicked="(val) => $emit('file-clicked', val)"
    />
  </ScrollArea>

</template>
