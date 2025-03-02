<script setup lang="ts">
import {ChevronDownIcon, ChevronRightIcon, FileTextIcon, FolderIcon,} from "lucide-vue-next";

// Props: Single file or folder item + expanded state
const props = defineProps({
  item: {
    type: Object,
    required: true,
  },
  expandedFolders: {
    type: Object,
    required: true,
  },
});

// Toggle folder open/close
const toggleFolder = () => {
  if (props.item.type === "folder") {
    props.expandedFolders[props.item.name] = !props.expandedFolders[props.item.name];
  }
};

</script>

<template>
  <ul>
    <li>
      <!-- Folder or File -->
      <div
        class="flex items-center space-x-2 cursor-pointer hover:bg-muted p-2 rounded-md"
        @click="item.type === 'folder' ? toggleFolder() : $emit('file-clicked', item)"
      >
        <!-- Icon -->
        <component
          :is="item.type === 'folder' ? FolderIcon : FileTextIcon"
          class="w-5 h-5 text-primary"
        />
        <!-- Name -->
        <span class="text-sm">{{ item.name }}</span>

        <!-- Expand/Collapse Icon for Folders -->
        <component
          v-if="item.type === 'folder'"
          :is="expandedFolders[item.name] ? ChevronDownIcon : ChevronRightIcon"
          class="w-4 h-4 text-muted-foreground ml-auto"
        />
      </div>

      <!-- Render Nested Items -->
      <ul
        v-if="item.type === 'folder' && expandedFolders[item.name]"
        class="ml-6 border-l border-gray-300 pl-3"
      >
<!--        TODO: click on files not working ?-->
        <FileTreeNode
          v-for="subItem in item.items"
          :key="subItem.name"
          :item="subItem"
          :expandedFolders="expandedFolders"
        />
      </ul>
    </li>

  </ul>
</template>